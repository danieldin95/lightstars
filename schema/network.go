package schema

type Range struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type Network struct {
	UUID    string  `json:"uuid"`
	Name    string  `json:"name"`
	State   string  `json:"state"`
	Address string  `json:"address"`
	Netmask string  `json:"netmask,omitempty"`
	Prefix  string  `json:"prefix,omitempty"`
	Range   []Range `json:"range"`
	DHCP    string  `json:"dhcp,omitempty"`
	Mode    string  `json:"mode"` // nat, router.
}
