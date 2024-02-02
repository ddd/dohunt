package main

import (
	"dohunt/pkg/dns"
	"dohunt/pkg/domains"
	"dohunt/pkg/rdap"
	"dohunt/pkg/whois"
	"log"

	"github.com/valyala/fasthttp"
)

func dnsWorker(client *fasthttp.Client, dnsCh <-chan domains.Domain, rdapCh chan<- domains.Domain, whoisCh chan<- domains.Domain) {

	for domain := range dnsCh {

		available, err := dns.CheckAvailable(client, domain.FQDN)
		if err != nil {
			log.Printf("An unknown error occcured when querying DNS for %v: %v\n", domain.FQDN, err)
		}

		// If the domain is taken, we query RDAP as first priority, if not available, whois for expiry if possible
		if !available {
			_, ok := rdapSrvMap[domain.TLD]
			if ok {
				rdapCh <- domain
				continue
			}

			whoisCh <- domain

		} else {
			// If domain is available, add to domainResultMap accordingly

			domain.Status = "available"
			domain.Source = "dns"
			domainResultMap[domain.FQDN] = domain
		}
	}
}

func rdapWorker(client *fasthttp.Client, rdapCh <-chan domains.Domain, whoisCh chan<- domains.Domain) {

	for domain := range rdapCh {

		result, err := rdap.Query(client, rdapSrvMap[domain.TLD], domain.FQDN, domain.TLD)
		if err != nil {
			log.Printf("An unknown error occurred while querying rdap for %v: %v\n", domain.FQDN, err)
			whoisCh <- domain
			continue
		}

		// rdap server didn't provide expiry, so we fallback to whois
		if result.Expires == 0 {
			whoisCh <- domain
			continue
		}

		domainResultMap[result.FQDN] = result
	}
}

func whoisWorker(whoisCh <-chan domains.Domain) {

	for domain := range whoisCh {

		// checking if tld has expiry regex for it. if not, skip
		_, ok := whois.WhoisExpiryMap[domain.TLD]
		if !ok {
			continue
		}

		srv, ok := whoisSrvMap[domain.SLD]
		if !ok {
			srv, ok = whoisSrvMap[domain.TLD]
			if !ok {
				// we skip domain if there is no whois server for sld or tld
				continue
			}
		}

		// expiry not shown for jp subdomain registrations, so skip.
		if domain.TLD == "jp" && (domain.SLD != domain.TLD) {
			continue
		}

		result, err := whois.Query(srv, domain.FQDN, domain.TLD)
		if err != nil {
			log.Printf("An unknown error occurred while querying whois for %v: %v\n", domain.FQDN, err)
			continue
		}

		domainResultMap[domain.FQDN] = result
	}

}
