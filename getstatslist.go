package estat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// 政府統計の総合窓口（e-Stat）で提供している統計表の情報を取得します。リクエストパラメータの指定により条件を絞った情報の取得も可能です。
func GetStatsList(ctx context.Context, query url.Values, options ...Option) (*GetStatsListContent, error) {
	container := GetStatsListContainer{}
	options = append(options, WithDataHandler(func(data []byte) error { return json.Unmarshal(data, &container) }))

	// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0#api_2_1
	return &container.GetStatsList, callAPI(
		ctx,
		http.MethodGet,
		"http://api.e-stat.go.jp/rest/3.0/app/json/getStatsList?"+query.Encode(),
		options...,
	)
}

// 3.2. 統計表情報取得

// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0#api_4_2
type GetStatsListContainer struct {
	GetStatsList GetStatsListContent `json:"GET_STATS_LIST"`
}

type GetStatsListContent struct {
	Parameter   any         `json:"PARAMETER"`
	DatalistInf DatalistInf `json:"DATALIST_INF"`
	Result      Result      `json:"RESULT"`
}

type DatalistInf struct {
	Number    int                   `json:"NUMBER"`
	ResultInf GetStatsListResultInf `json:"RESULT_INF"`
	TableInf  List[TableInf]        `json:"TABLE_INF"`
}

type GetStatsListResultInf struct {
	FromNumber int `json:"FROM_NUMBER"`
	ToNumber   int `json:"TO_NUMBER"`
	NextKey    int `json:"NEXT_KEY"`
}
