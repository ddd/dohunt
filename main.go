package main

import (
	"bufio"
	"dohunt/pkg/domains"
	"dohunt/pkg/rdap"
	"dohunt/pkg/whois"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	domainResultMap = make(map[string]domains.Domain)
	rdapSrvMap      = make(map[string]string)
	whoisSrvMap     = make(map[string]string)
)

func manager(client *fasthttp.Client, dnsCh chan<- domains.Domain, domainFile string) {

	var domainsToCheck []domains.Domain

	file, err := os.Open(domainFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	tldList, err := rdap.FetchTLDList(client)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		x := strings.ToLower(scanner.Text())
		sld, tld, domain, err := parseDomain(client, x)
		if err != nil {
			log.Printf("An unknown error occurred while parsing domain %v: %v\n", x, err)
			continue
		}

		if _, ok := tldList[tld]; !ok {
			log.Printf("Invalid TLD for domain %v: %v\n", domain, tld)
			continue
		}

		if _, ok := rdapSrvMap[tld]; !ok {
			if _, ok := whoisSrvMap[tld]; !ok {
				fmt.Printf("\033[90mExpiry for %v can't be tracked as TLD %v does not have a known rdap/whois server\033[0m\n", domain, tld)
			}
		}

		domainData := domains.Domain{FQDN: domain, SLD: sld, TLD: tld}
		dnsCh <- domainData
		domainsToCheck = append(domainsToCheck, domainData)
	}

	for {
		for _, domain := range domainsToCheck {
			dnsCh <- domain
		}
		time.Sleep(time.Minute * 20)
	}

}

func main() {

	// Define command line flags
	port := flag.Int("p", 8080, "port")
	threads := flag.Int("t", 1, "thread count")
	domainFile := flag.String("d", "domains.txt", "file containing list of domains to track expiry of")

	// Parse the command line flags
	flag.Parse()

	client := &fasthttp.Client{}

	var err error
	rdapSrvMap, err = rdap.FetchSrvList(client)
	if err != nil {
		log.Fatal(err)
	}

	whoisSrvMap, err = whois.FetchSrvList(client)
	if err != nil {
		log.Fatal(err)
	}

	dnsCh := make(chan domains.Domain, 1000)
	whoisCh := make(chan domains.Domain, 1000)
	rdapCh := make(chan domains.Domain, 1000)

	go manager(client, dnsCh, *domainFile)

	for i := 0; i < *threads; i++ {
		go dnsWorker(client, dnsCh, rdapCh, whoisCh)
		go rdapWorker(client, rdapCh, whoisCh)
		go whoisWorker(whoisCh)
	}

	if os.Getenv("DOCHKEY") == "" {
		fmt.Println("\033[33mWarning: DOCHKEY environment variable not set, no key required to access local API\033[0m")
	}

	startAPI(&domainResultMap, *port)

}
