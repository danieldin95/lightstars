package hoster

type SystemInfo struct {
	HostName  string `json:"hostname"`
	Cpu       int    `json:"cpu"`
	CpuVendor string `json:"cpuVendor"`
	Mem       int    `json:"memory"`
	MemFree   int    `json:"memoryFree"`
}

type IpLink struct {
	Alias  string `json:"alias"`
	Name   string `json:"name"`
	Master string `json:"master"`
	State  string `json:"state"`
}

type IpRoute struct {
	Prefix  string `json:"prefix"`
	NextHop string `json:"gateway"`
	Source  string `json:"source"`
	Link    string `json:"link"`
}
