package estat

import "encoding/json"

func unmarshalList[T any](list *[]T, data []byte) error {
	if data[0] == '{' {
		*list = make([]T, 1)
		return json.Unmarshal(data, &(*list)[0])
	} else {
		return json.Unmarshal(data, list)
	}
}

func marshalList[T any](list []T) ([]byte, error) {
	if len(list) == 1 {
		return json.Marshal(list[0])
	} else {
		return json.Marshal(list)
	}
}
