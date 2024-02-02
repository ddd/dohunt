package rdap

// This has to unfortunately exist as some ccTLDs implementing RDAP refuse to follow RFC
var RDAPExpiryMap = map[string]any{
	"id": "2006-01-02 15:04:05",
}
