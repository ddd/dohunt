package rdap

type Event struct {
	Action string `json:"eventAction"`
	Date   string `json:"eventDate"`
}

type RDAPResponse struct {
	Flags  []string `json:"status"`
	Events []Event  `json:"events"`
}
