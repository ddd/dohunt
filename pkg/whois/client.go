package whois

import (
	"dohunt/pkg/domains"
	"net"
	"time"
)

func Query(server string, domain string, tld string) (domains.Domain, error) {

	conn, err := net.Dial("tcp", server+":43")
	if err != nil {
		return domains.Domain{}, err
	}

	defer conn.Close()

	// Send the query
	query := domain + "\r\n"
	_, err = conn.Write([]byte(query))
	if err != nil {
		return domains.Domain{}, err
	}

	err = conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		return domains.Domain{}, err
	}

	// Receive the response
	buffer := make([]byte, 4096)
	var response []byte

	for {

		x, err := conn.Read(buffer)

		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return domains.Domain{}, err
			}
		}

		response = append(response, buffer[:x]...)

	}

	domainData, err := parseWhoisResponse(domain, tld, response)
	if err != nil {
		return domains.Domain{}, err
	}

	return domainData, nil
}
