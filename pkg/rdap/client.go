package rdap

import (
	"dohunt/pkg/domains"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

func Query(client *fasthttp.Client, server string, domain string, tld string) (domains.Domain, error) {

	// Create a new request object
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// Set the request URL
	req.SetRequestURI(fmt.Sprintf("%vdomain/%v", server, domain))

	// Send the request and get the response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return domains.Domain{}, err
	}

	status := resp.StatusCode()
	switch status {
	case 200:
	case 404:
		return domains.Domain{FQDN: domain, Status: "available"}, nil
	default:
		return domains.Domain{}, fmt.Errorf("unknown status code %v", status)
	}

	// Unmarshal the response JSON
	var response RDAPResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return domains.Domain{}, err
	}

	domainData := domains.Domain{FQDN: domain, Status: "taken", Flags: response.Flags, Source: "rdap"}

	for _, event := range response.Events {
		if event.Action == "expiration" || event.Action == "soft expiration" {

			var t time.Time
			var err error

			c, ok := RDAPExpiryMap[tld]
			if ok {
				t, err = time.Parse(c.(string), event.Date)
			} else {
				t, err = time.Parse(time.RFC3339, event.Date)
			}
			if err != nil {
				return domains.Domain{}, err
			}

			domainData.Expires = t.Unix()
			break
		}
	}

	return domainData, nil
}
