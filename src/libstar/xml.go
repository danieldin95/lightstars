package libstar

import "encoding/xml"

type XMLBase struct {
}

func (x *XMLBase) Decode(data string) error {
	if err := xml.Unmarshal([]byte(data), x); err != nil {
		Error("XMLBase.Decode %s", err)
		return err
	}
	return nil
}

func (x *XMLBase) Encode() string {
	data, err := xml.Marshal(x)
	if err != nil {
		Error("XMLBase.Encode %s", err)
		return ""
	}
	return string(data)
}
