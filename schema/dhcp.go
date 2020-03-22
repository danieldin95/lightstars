package schema

type DHCPLease struct {
	Type     int    `json:"type"`
	Mac      string `json:"mac"`
	IPAddr   string `json:"ipAddr"`
	Prefix   uint   `json:"prefix"`
	Hostname string `json:"hostname"`
}
