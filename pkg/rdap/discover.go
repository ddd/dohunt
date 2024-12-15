package rdap

import (
	"encoding/json"
	"strings"

	"github.com/valyala/fasthttp"
)

// Returns a map of every TLD which has a RDAP server and the respective server address
func FetchSrvList(client *fasthttp.Client) (map[string]string, error) {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("https://data.iana.org/rdap/dns.json")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	var response map[string]any
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	rdapSrvMap := make(map[string]string)
	services := response["services"].([]any)

	for _, service := range services {
		tlds := service.([]any)[0].([]any)

		for _, tld := range tlds {
			rdapSrvMap[tld.(string)] = service.([]any)[1].([]any)[0].(string)
		}
	}

	return rdapSrvMap, nil
}

// Returns list of valid TLDs
func FetchTLDList(client *fasthttp.Client) (map[string]int, error) {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(resp.Body()), "\n")

	// Create an empty array to store the TLDs
	tlds := make(map[string]int)

	// Iterate over the lines and add each TLD to the array
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			tlds[strings.ToLower(line)] = 0
		}
	}

	return tlds, nil
}
