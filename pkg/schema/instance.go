package schema

type Channel struct {
	Type          string `json:"type"`
	TargetType    string `json:"targetType"`
	TargetName    string `json:"targetName"`
	SourceChannel string `json:"sourceChannel"`
}

type Graphics struct {
	Type     string `json:"type"`
	Password string `json:"password"`
	Listen   string `json:"listen"`
	Port     string `json:"port"`
	AutoPort string `json:"autoport"`
}

type ListGraphics struct {
	List
	Items []Graphics `json:"items"`
}

type Processor struct {
	Cpu  string `json:"cpu"`  // configure
	Mode string `json:"mode"` // configure
	Time uint64 `json:"time"` // MicroSeconds
}

type Memory struct {
	Size string `json:"size"` // configure
	Unit string `json:"unit"` // configure
}

type Instance struct {
	UUID        string       `json:"uuid"`
	Name        string       `json:"name"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Family      string       `json:"family,omitempty"` // linux, windows or others
	State       string       `json:"state,omitempty"`
	Arch        string       `json:"arch,omitempty"` // x86_64 or i386
	Type        string       `json:"type,omitempty"`
	Boots       string       `json:"boots,omitempty"`
	DataStore   string       `json:"datastore,omitempty"`
	Start       string       `json:"start,omitempty"` // whether booting with created
	CpuMode     string       `json:"cpuMode,omitempty"`
	MaxCpu      uint         `json:"maxCpu,omitempty"`
	MaxMem      uint64       `json:"maxMem,omitempty"`  // KiB
	Memory      uint64       `json:"memory,omitempty"`  // KiB
	CpuTime     uint64       `json:"cpuTime,omitempty"` // micro seconds
	Disks       []Disk       `json:"disks,omitempty"`
	Interfaces  []Interface  `json:"interfaces,omitempty"`
	Controllers []Controller `json:"controllers,omitempty"`
	Graphics    []Graphics   `json:"graphics,omitempty"`
	Channels    []Channel    `json:"channels,omitempty"`
}

type ListInstance struct {
	List
	Items []Instance `json:"items"`
}

type Disk struct {
	Domain     Instance `json:"domain"`
	Seq        string   `json:"seq,omitempty"`       // configure
	Name       string   `json:"name,omitempty"`      // disk name
	UUID       string   `json:"uuid,omitempty"`      // disk UUID
	Store      string   `json:"datastore,omitempty"` // disk saved to datastore
	Size       string   `json:"size"`                // configure
	SizeUnit   string   `json:"sizeUnit,omitempty"`  //configure
	Format     string   `json:"format"`
	Source     string   `json:"source"`
	Device     string   `json:"device"`
	Bus        string   `json:"bus"`      //configre
	AddrType   string   `json:"addrType"` // pci, and drive
	AddrSlot   string   `json:"addrSlot"`
	AddrDomain string   `json:"addrDomain"`
	AddrBus    string   `json:"addrBus"`
	AddrFunc   string   `json:"addrFunc"`
	AddrCtl    string   `json:"addrCtl"`
	AddrTgt    string   `json:"addrTgt"`
	AddrUnit   string   `json:"addrUnit"`
	Volume     Volume   `json:"volume"`
}

type ListDisk struct {
	List
	Items []Disk `json:"items"`
}

type Interface struct {
	Domain     Instance `json:"domain"`
	Seq        string   `json:"seq,omitempty"` //configure
	Name       string   `json:"name,omitempty"`
	UUID       string   `json:"uuid,omitempty"`
	Type       string   `json:"type,omitempty"` //bridge or openvswitch
	IpAddr     string   `json:"ipaddr"`
	Address    string   `json:"address"`
	Network    string   `json:"network"`
	Source     string   `json:"source"`
	HostDev    string   `json:"hostDev"`
	Model      string   `json:"model"` // configure
	Device     string   `json:"device"`
	AddrType   string   `json:"addrType"` // now only pci.
	AddrSlot   string   `json:"addrSlot"`
	AddrDomain string   `json:"addrDomain"`
	AddrBus    string   `json:"addrBus"`
	AddrFunc   string   `json:"addrFunc"`
}

type ListInterface struct {
	List
	Items []Interface `json:"items"`
}

type Controller struct {
	Type    string  `json:"source"`
	Model   string  `json:"model"`
	Index   string  `json:"device"`
	Address Address `json:"address"`
}

type Address struct {
	Type     string `json:"type"`
	Domain   string `json:"domain"`
	Bus      string `json:"bus"`
	Slot     string `json:"slot"`
	Function string `json:"function"`
}

type InstancesStats struct {
	Running      uint   `json:"running"`
	Shutdown     uint   `json:"shutdown"`
	Others       uint   `json:"others"`
	AllocMem     uint64 `json:"allocMem"`
	OccupiedMem  uint64 `json:"occupiedMem"`
	AllocCpu     uint   `json:"allocCpu"`
	OccupiedCpu  uint   `json:"occupiedCpu"`
	AllocStorage uint64 `json:"allocStorage"`
}

type Snapshot struct {
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	Uptime    int64  `json:"uptime"`
	State     string `json:"state"`
	IsCurrent bool   `json:"isCurrent"`
}

type ListSnapshot struct {
	List
	Items []Snapshot `json:"items"`
}
