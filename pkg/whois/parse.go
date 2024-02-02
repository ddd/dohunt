package whois

import (
	"dohunt/pkg/domains"
	"fmt"
	"time"
)

func parseWhoisResponse(domain string, tld string, resp []byte) (domains.Domain, error) {

	result := WhoisExpiryMap[tld].Re.FindSubmatch(resp)

	if len(result) != 2 {
		return domains.Domain{}, fmt.Errorf("failed to parse whois response for %v (%v)", domain, tld)
	}

	//fmt.Println(string(result[1]))
	dateString := string(result[1])
	t, err := time.Parse(WhoisExpiryMap[tld].Format, dateString)
	if err != nil {
		return domains.Domain{}, fmt.Errorf("failed to parse date %v for %v (%v)", dateString, domain, tld)
	}

	return domains.Domain{FQDN: domain, Status: "taken", Source: "whois", Expires: t.Unix()}, nil

}
