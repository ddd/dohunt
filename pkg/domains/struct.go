package domains

type Domain struct {
	FQDN    string   `json:"-"`
	Status  string   `json:"status"`  // taken, available
	Expires int64    `json:"expires"` // unix timestamp
	Flags   []string `json:"flags"`   // client hold
	Source  string   `json:"source"`  // rdap, whois, dig
	SLD     string   `json:"-"`       // .com.br // .co.uk // .com
	TLD     string   `json:"-"`       //
}
