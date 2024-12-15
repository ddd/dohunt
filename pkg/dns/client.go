package dns

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type Answer struct {
	Data string `json:"data"`
}

type DNSResponse struct {
	Status    int      `json:"status"`
	Answers   []Answer `json:"Answer"`
	Authority []Answer `json:"Authority"`
}

var ownedByRegistryCache = map[string]bool{
	"co.za": true,
	"ac.uk": true,
}

func QuerySOA(client *fasthttp.Client, query string) ([]string, error) {

	// Create a new request object
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// Set the request URL
	req.SetRequestURI(fmt.Sprintf("https://dns.google/resolve?name=%v&type=SOA", query))

	// Send the request and get the response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	status := resp.StatusCode()
	if status != 200 {
		return nil, fmt.Errorf("unknown status code %v when searching: %v", status, query)
	}

	var response DNSResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	if response.Status == 2 {
		return nil, fmt.Errorf("srvfail when fetching soa for %v", query)
	}

	var respData Answer
	if len(response.Answers) == 0 {
		respData = response.Authority[0]
	} else {
		respData = response.Answers[0]
	}

	x := strings.Split(respData.Data, " ")

	var ns []string
	for _, i := range x {
		if strings.Contains(i, ".") {
			ns = append(ns, i)
		}
	}

	return ns, nil

}

// This simply returns true if the response from Google DNS is NXDOMAIN
func CheckAvailable(client *fasthttp.Client, query string) (bool, error) {

	// Create a new request object
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// Set the request URL
	req.SetRequestURI(fmt.Sprintf("https://dns.google/resolve?name=%v&type=A", query))

	// Send the request and get the response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return false, err
	}

	status := resp.StatusCode()
	if status != 200 {
		return false, fmt.Errorf("unknown status code %v when searching: %v", status, query)
	}

	var response map[string]interface{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return false, err
	}

	if response["Status"].(float64) == 3 {
		return true, err
	} else {
		return false, err
	}

}

// This may be inaccurate, simply matches nameservers of the subdomain with that of the root tld
func IsSLDOwnedByRegistry(client *fasthttp.Client, potentialTLD string) (bool, error) {

	c, ok := ownedByRegistryCache[potentialTLD]
	if ok {
		return c, nil
	}

	// Create a new request object
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	x := strings.Split(potentialTLD, ".")

	rootNs, err := QuerySOA(client, x[1])
	if err != nil {
		return false, err
	}

	ns, err := QuerySOA(client, potentialTLD)
	if err != nil {
		return false, err
	}

	for _, x := range ns {
		for _, y := range rootNs {
			if x == y {
				ownedByRegistryCache[potentialTLD] = true
				return true, nil
			}
		}
	}

	ownedByRegistryCache[potentialTLD] = false
	return false, nil

}
