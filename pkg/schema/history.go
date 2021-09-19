package schema

type History struct {
	User   string `json:"user"`
	Date   string `json:"date"`
	Method string `json:"method"`
	Client string `json:"client"`
	Url    string `json:"url"`
}

type ListHistory struct {
	List
	Items []History `json:"items"`
}
