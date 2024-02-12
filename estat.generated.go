package estat

import "encoding/json"

// Code generated by github.com/psyark/estat/codegen; DO NOT EDIT.

// ClassHelper は、Class の単一の値または配列を透過的にUnmarshal/Marshalするスライスです
type ClassHelper []Class

func (c *ClassHelper) UnmarshalJSON(d []byte) error {
	if d[0] == '{' {
		*c = make([]Class, 1)
		return json.Unmarshal(d, &(*c)[0])
	} else {
		return json.Unmarshal(d, (*[]Class)(c))
	}
}

func (c ClassHelper) MarshalJSON() ([]byte, error) {
	if len(c) == 1 {
		return json.Marshal(c[0])
	} else {
		return json.Marshal([]Class(c))
	}
}

// NoteHelper は、Note の単一の値または配列を透過的にUnmarshal/Marshalするスライスです
type NoteHelper []Note

func (c *NoteHelper) UnmarshalJSON(d []byte) error {
	if d[0] == '{' {
		*c = make([]Note, 1)
		return json.Unmarshal(d, &(*c)[0])
	} else {
		return json.Unmarshal(d, (*[]Note)(c))
	}
}

func (c NoteHelper) MarshalJSON() ([]byte, error) {
	if len(c) == 1 {
		return json.Marshal(c[0])
	} else {
		return json.Marshal([]Note(c))
	}
}
