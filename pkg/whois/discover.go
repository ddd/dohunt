package whois

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func FetchSrvList(client *fasthttp.Client) (map[string]string, error) {

	// Create a new request object
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// Set the request URL
	req.SetRequestURI("https://raw.githubusercontent.com/rfc1036/whois/next/tld_serv_list")

	// Send the request and get the response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	data := resp.Body()

	// Parse the data into a map
	whoisSrvMap := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			// Skip empty lines and lines starting with '#'
			continue
		}

		fields := strings.Fields(line)
		if fields[1] == "WEB" || fields[1] == "NONE" {
			continue
		}

		// Remove preceding "."
		tld := fields[0][1:]

		if isUpper(fields[1]) && len(fields) > 2 {
			whoisSrvMap[tld] = fields[2]
		} else {
			whoisSrvMap[tld] = fields[1]
		}
	}

	return whoisSrvMap, nil
}
