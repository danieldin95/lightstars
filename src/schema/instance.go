package schema

type Graphics struct {
	Type     string `json:"type"`
	Password string `json:"password"`
	Listen   string `json:"listen"`
	Port     string `json:"port"`
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
	Action      string       `json:"action,omitempty"` // If is "", means not action.
	UUID        string       `json:"uuid"`
	Name        string       `json:"name"`
	Family      string       `json:"family"` // linux, windows or others
	State       string       `json:"state"`
	Arch        string       `json:"arch"` // x86_64 or i386
	Type        string       `json:"type"`
	Boots       string       `json:"boots,omitempty"`
	DataStore   string       `json:"datastore,omitempty"`
	Start       string       `json:"start,omitempty"` // whether booting with created
	CpuMode     string       `json:"cpuMode"`
	MaxCpu      uint         `json:"maxCpu"`
	MaxMem      uint64       `json:"maxMem"`  // KiB
	Memory      uint64       `json:"memory"`  // KiB
	CpuTime     uint64       `json:"cpuTime"` // micro seconds
	Disks       []Disk       `json:"disks,omitempty"`
	Interfaces  []Interface  `json:"interfaces,omitempty"`
	Controllers []Controller `json:"controllers,omitempty"`
	Graphics    []Graphics   `json:"graphics"`
}

type ListInstance struct {
	List
	Items []Instance `json:"items"`
}

type Disk struct {
	Action     string `json:"action,omitempty"`
	Seq        string `json:"seq,omitempty"`       // configure
	Name       string `json:"name,omitempty"`      // disk name
	UUID       string `json:"uuid,omitempty"`      // disk UUID
	Store      string `json:"datastore,omitempty"` // disk saved to datastore
	Size       string `json:"size"`                // configure
	SizeUnit   string `json:"sizeUnit,omitempty"`  //configure
	Format     string `json:"format"`
	Source     string `json:"source"`
	Device     string `json:"device"`
	Bus        string `json:"bus"`      //configre
	AddrType   string `json:"addrType"` // pci, and drive
	AddrSlot   uint16 `json:"addrSlot"`
	AddrDomain uint16 `json:"addrDomain"`
	AddrBus    uint16 `json:"addrBus"`
	AddrFunc   uint16 `json:"addrFunc"`
	AddrCtl    uint16 `json:"addrCtl"`
	AddrTgt    uint16 `json:"addrTgt"`
	AddrUnit   uint16 `json:"addrUnit"`
	Volume     Volume `json:"volume"`
}

type ListDisk struct {
	List
	Items []Disk `json:"items"`
}

type Interface struct {
	Action     string `json:"action,omitempty"` // If is "", means not action.
	Seq        string `json:"seq,omitempty"`    //configure
	Name       string `json:"name,omitempty"`
	UUID       string `json:"uuid,omitempty"`
	Type       string `json:"type,omitempty"` //bridge or openvswitch
	Address    string `json:"address"`
	Source     string `json:"source"`
	Model      string `json:"model"` // configure
	Device     string `json:"device"`
	AddrType   string `json:"addrType"` // now only pci.
	AddrSlot   uint16 `json:"addrSlot"`
	AddrDomain uint16 `json:"addrDomain"`
	AddrBus    uint16 `json:"addrBus"`
	AddrFunc   uint16 `json:"addrFunc"`
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
