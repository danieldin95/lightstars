package schema

type MetaData struct {
	Total  int `json:"total"`
	Size   int `json:"size"`
	Offset int `json:"offset"`
}

type List struct {
	Items []interface{} `json:"items"`
	Metadata  MetaData   `json:"metadata"`
}