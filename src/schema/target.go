package schema

type Target struct {
	Name   string `json:"name"`
	Target string `json:"target"`
	Host   string `json:"host"`
}

func (tgt Target) ID() string {
	return tgt.Host + ":" + tgt.Target
}
