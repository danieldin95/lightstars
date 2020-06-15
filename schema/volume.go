package schema

type VolumeInfo struct {
	Pool       string `json:"pool"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Capacity   uint64 `json:"capacity"`
	Allocation uint64 `json:"allocation"`
}

type VolumeInfos map[string]VolumeInfo
