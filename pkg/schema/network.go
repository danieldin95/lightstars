package schema

type Range struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Network struct {
	UUID    string  `json:"uuid"`
	Name    string  `json:"name"`
	Bridge  string  `json:"bridge"`
	State   string  `json:"state"`
	Address string  `json:"address"`
	Netmask string  `json:"netmask,omitempty"`
	Prefix  string  `json:"prefix,omitempty"`
	Range   []Range `json:"range"`
	Mode    string  `json:"mode"`           // nat, router.
	Type    string  `json:"type,omitempty"` // linux bridge or openvswitch
}

type ListNetwork struct {
	List
	Items []Network `json:"items"`
}
