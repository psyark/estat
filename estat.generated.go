package estat

import "encoding/json"

func (c *ClassHelper) UnmarshalJSON(d []byte) error {
	if d[0] == '{' {
		*c = make([]Class, 1)
		return json.Unmarshal(d, &(*c)[0])
	} else {
		return json.Unmarshal(d, (*[]Class)(c))
	}
}

func (c *NoteHelper) UnmarshalJSON(d []byte) error {
	if d[0] == '{' {
		*c = make([]Note, 1)
		return json.Unmarshal(d, &(*c)[0])
	} else {
		return json.Unmarshal(d, (*[]Note)(c))
	}
}
