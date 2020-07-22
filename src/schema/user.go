package schema

type User struct {
	Type      string   `json:"type"` // admin, guest or other
	Name      string   `json:"name"`
	Password  string   `json:"password,omitempty"`
	Language  string   `json:"language"`
	Instances []string `json:"instances,omitempty"`
}

type ListUser struct {
	List
	Items []User `json:"items"`
}
