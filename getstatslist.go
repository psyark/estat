package estat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// 政府統計の総合窓口（e-Stat）で提供している統計表の情報を取得します。リクエストパラメータの指定により条件を絞った情報の取得も可能です。
func GetStatsList(ctx context.Context, query url.Values) (*GetStatsListContent, error) {
	container := GetStatsListContainer{}

	// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0#api_2_1
	return &container.GetStatsList, callAPI(
		ctx,
		http.MethodGet,
		"http://api.e-stat.go.jp/rest/3.0/app/json/getStatsList?"+query.Encode(),
		WithDataHandler(func(data []byte) error { return json.Unmarshal(data, &container) }),
	)
}

// 3.2. 統計表情報取得

// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0#api_4_4
type GetStatsListContainer struct {
	GetStatsList GetStatsListContent `json:"GET_STATS_LIST"`
}

type GetStatsListContent struct {
	DATALIST_INF any    `json:"DATALIST_INF"`
	PARAMETER    any    `json:"PARAMETER"`
	RESULT       Result `json:"RESULT"`
}
