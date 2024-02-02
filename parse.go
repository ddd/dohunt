package main

import (
	"errors"
	"fmt"
	"strings"

	"dohunt/pkg/dns"

	"github.com/valyala/fasthttp"
)

func parseDomain(client *fasthttp.Client, fullDomain string) (sld string, tld string, domain string, err error) {

	//fmt.Println("Recieved request to parse", fullDomain)

	if !strings.Contains(fullDomain, ".") {
		return "", "", "", errors.New("invalid domain")
	}

	x := strings.Split(strings.ToLower(fullDomain), ".")

	if len(x) >= 3 {

		ownedByRegistry, err := dns.IsSLDOwnedByRegistry(client, x[len(x)-2]+"."+x[len(x)-1])

		if err != nil {
			return "", "", "", err
		}

		if ownedByRegistry {
			return x[len(x)-2] + "." + x[len(x)-1], x[len(x)-1], x[len(x)-3] + "." + x[len(x)-2] + "." + x[len(x)-1], nil
		} else {
			//fmt.Printf("%v is not owned by the registry\n", x[len(x)-2]+"."+x[len(x)-1])
			return x[len(x)-1], x[len(x)-1], x[len(x)-2] + "." + x[len(x)-1], nil
		}

	} else if len(x) == 2 {

		return x[len(x)-1], x[len(x)-1], fullDomain, nil
	} else {
		return "", "", "", fmt.Errorf("invalid domain %v", fullDomain)
	}
}
