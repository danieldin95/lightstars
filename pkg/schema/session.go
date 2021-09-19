package schema

type Session struct {
	Value   string `json:"value"`
	Date    string `json:"date"`
	Client  string `json:"client"`
	Expires string `json:"expires"`
	Uuid    string `json:"uuid"`
}

type ListSession struct {
	List
	Items []Session `json:"items"`
}
