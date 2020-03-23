package schema

type NFS struct {
	Host   string `json:"host"`
	Path   string `json:"path"`
	Format string `json:"format"`
}

type DataStore struct {
	UUID       string `json:"uuid"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Format     string `json:"format"`
	State      string `json:"state"`
	Capacity   uint64 `json:"capacity"`   // bytes
	Allocation uint64 `json:"allocation"` // bytes
	Available  uint64 `json:"available"`  // Bytes
	Source     string `json:"source"`
	NFS        *NFS   `json:"nfs"`
}

type ListDataStore struct {
	List
	Items []DataStore `json:"items"`
}
