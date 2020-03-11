package schema

type User struct {
	Type     string `json:"type"` // admin, guest or other
	Name     string `json:"name"`
	Password string `json:"password"`
}
