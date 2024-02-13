package estat

import (
	"encoding/json"
	"time"
)

// Result is 4.1.1. RESULT タグ
// すべてのAPI共通で、以下の要素を出力します。
type Result struct {
	Date     time.Time `json:"DATE"`
	ErrorMsg string    `json:"ERROR_MSG"`
	Status   int       `json:"STATUS"`
}

func (r Result) MarshalJSON() ([]byte, error) {
	type Alias Result
	t := struct {
		Alias
		Date string `json:"DATE"`
	}{
		Alias: Alias(r),
	}

	if r.Date.UnixMilli() == 0 {
		t.Date = r.Date.Format("2006-01-02T15:04:05-07:00")
	} else {
		t.Date = r.Date.Format("2006-01-02T15:04:05.000-07:00")
	}
	return json.Marshal(t)
}
