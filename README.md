
## dohunt

Fast and efficient tool for tracking domain expiry written in Go.


## Installation

Clone the repository to your machine
```
git clone https://github.com/ddd/dohunt
```

Compile the binary (requires Go)

```
cd dohunt; go build
```

Run the binary

```
./dohunt -p 8080 -d domains.txt
```

The local web interface should be available at `http://localhost:8080` and the local API is available at `http://localhost:8080/api/domains`

Add authentication **(optional)**

```
DOCHKEY=secret ./dohunt -p 8080 -d domains.txt
```

The local web interface will be available at `http://localhost:8080/#secret` and the local API is available at `http://localhost:8080/api/domains?key=secret`


## Known Issues:
- Some TLDs without whois/rdap servers such as `ph` respond to DNS queries of invalid domain names with `NOERROR` instead of `NXDOMAIN`, hence making it impossible to check domain availability

- In order to detect if a subdomain is owned by the registry (ex. `com.br` vs `rr.com`), the code checks if the subdomain has the same nameserver as the root TLD. This may sometimes be inaccurate.