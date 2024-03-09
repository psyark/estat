package estat

import "encoding/json"

type List[T any] []T

func (list *List[T]) UnmarshalJSON(data []byte) error {
	if data[0] == '{' {
		*list = make([]T, 1)
		return json.Unmarshal(data, &(*list)[0])
	} else {
		return json.Unmarshal(data, (*[]T)(list))
	}
}

func (list List[T]) MarshalJSON() ([]byte, error) {
	if len(list) == 1 {
		return json.Marshal(list[0])
	} else {
		return json.Marshal(([]T)(list))
	}
}
