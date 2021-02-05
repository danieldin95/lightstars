package libvirtc

import (
	"github.com/beevik/etree"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/libvirt/libvirt-go"
)

var (
	DomainAll                    = libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE
	DomainDeviceModifyConfig     = libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	DomainDeviceModifyPersistent = libvirt.DOMAIN_DEVICE_MODIFY_LIVE | libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	DomainDestroyGraceful        = libvirt.DOMAIN_DESTROY_GRACEFUL
	DomainShutdownAcpi           = libvirt.DOMAIN_SHUTDOWN_ACPI_POWER_BTN
	DomainCpuMaximum             = libvirt.DOMAIN_VCPU_MAXIMUM | libvirt.DOMAIN_VCPU_CONFIG
	DomainCpuConfig              = libvirt.DOMAIN_VCPU_CONFIG
	DomainMemMaximum             = libvirt.DOMAIN_MEM_MAXIMUM | libvirt.DOMAIN_MEM_CONFIG
	DomainMemConfig              = libvirt.DOMAIN_MEM_CONFIG
)

type Domain struct {
	libvirt.Domain
	hyper *HyperVisor
}

func NewDomainFromVir(dom *libvirt.Domain, hyper *HyperVisor) *Domain {
	return &Domain{*dom, hyper}
}

func (d *Domain) document() (*etree.Document, error) {
	xml, err := d.GetXMLDesc(true)
	if err != nil {
		return nil, err
	}
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return nil, err
	}
	return doc, nil
}

func (d *Domain) SetCpu(max, mode string) error {
	if d.hyper == nil {
		return libstar.NewErr("hyper is nil")
	}

	d.hyper.Lock.Lock()
	defer d.hyper.Lock.Unlock()

	doc, err := d.document()
	if err != nil {
		return err
	}
	domEle := doc.FindElement("/domain")
	if domEle == nil {
		return libstar.NewErr("domain tag not found")
	}
	// edit vCpu
	vCpuEle := domEle.FindElement("./vcpu")
	if vCpuEle == nil {
		vCpuEle = domEle.CreateElement("vcpu")
	}
	vCpuEle.SetText(max)
	// edit cpu mode
	cpuEle := domEle.FindElement("./cpu")
	if cpuEle == nil {
		cpuEle = domEle.CreateElement("cpu")
	}
	cpuEle.CreateAttr("mode", mode)
	libstar.Debug("Domain.SetCpu %v", cpuEle)
	newXml, err := doc.WriteToString()
	if err != nil {
		return err
	}
	libstar.Debug("Domain.SetCpu %s", newXml)
	if err := d.reDefine(newXml); err != nil {
		return err
	}
	return nil
}

func (d *Domain) SetMemory(size, unit string) error {
	if d.hyper == nil {
		return libstar.NewErr("hyper is nil")
	}

	d.hyper.Lock.Lock()
	defer d.hyper.Lock.Unlock()

	doc, err := d.document()
	if err != nil {
		return err
	}
	domEle := doc.FindElement("/domain")
	if domEle == nil {
		return libstar.NewErr("domain tag not found")
	}
	// edit memory
	memEle := domEle.FindElement("./memory")
	if memEle == nil {
		memEle = domEle.CreateElement("memory")
	}
	memEle.CreateAttr("unit", unit)
	memEle.SetText(size)

	// edit current memory
	curMemEle := domEle.FindElement("./currentMemory")
	if curMemEle == nil {
		curMemEle = domEle.CreateElement("currentMemory")
	}
	curMemEle.CreateAttr("unit", unit)
	curMemEle.SetText(size)

	newXml, err := doc.WriteToString()
	if err != nil {
		return err
	}
	libstar.Debug("Domain.SetMemory %s", newXml)
	if err := d.reDefine(newXml); err != nil {
		return err
	}
	return nil
}

func (d *Domain) reDefine(newXml string) error {
	if d.hyper == nil {
		return libstar.NewErr("hyper is nil")
	}
	oldXml, err := d.GetXMLDesc(true)
	if err != nil {
		return err
	}
	if err := d.Undefine(); err != nil {
		return err
	}
	_, err = d.hyper.Conn.DomainDefineXML(newXml)
	if err != nil {
		_, _ = d.hyper.Conn.DomainDefineXML(oldXml) // define by old xml.
		return err
	}
	return nil
}

func (d *Domain) GetXMLDesc(secure bool) (string, error) {
	if secure {
		return d.Domain.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	} else {
		return d.Domain.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	}
}

func (d *Domain) GetMetadataTitle(title bool) (string, error) {
	tipus := libvirt.DOMAIN_METADATA_DESCRIPTION
	if title {
		tipus = libvirt.DOMAIN_METADATA_TITLE
	}
	return d.Domain.GetMetadata(tipus, "", libvirt.DOMAIN_AFFECT_CURRENT)
}

// virsh desc 3 --config --live --title your_title
func (d *Domain) SetMetadataTitle(value string, title bool) error {
	tipus := libvirt.DOMAIN_METADATA_DESCRIPTION
	if title {
		tipus = libvirt.DOMAIN_METADATA_TITLE
	}
	if active, _ := d.Domain.IsActive(); active {
		if err := d.Domain.SetMetadata(tipus, value, "", "", libvirt.DOMAIN_AFFECT_LIVE); err != nil {
			return err
		}
	}
	return d.Domain.SetMetadata(tipus, value, "", "", libvirt.DOMAIN_AFFECT_CONFIG)
}

func ListDomains() ([]Domain, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return hyper.ListAllDomains()
}

func DomainState2Str(state libvirt.DomainState) string {
	switch state {
	case libvirt.DOMAIN_NOSTATE:
		return "no-state"
	case libvirt.DOMAIN_RUNNING:
		return "running"
	case libvirt.DOMAIN_BLOCKED:
		return "blocked"
	case libvirt.DOMAIN_PAUSED:
		return "paused"
	case libvirt.DOMAIN_SHUTDOWN:
		return "shutdown"
	case libvirt.DOMAIN_CRASHED:
		return "crashed"
	case libvirt.DOMAIN_PMSUSPENDED:
		return "pm-suspended"
	case libvirt.DOMAIN_SHUTOFF:
		return "shutoff"
	default:
		return "unknown"
	}
}
