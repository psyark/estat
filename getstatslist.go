package estat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// 政府統計の総合窓口（e-Stat）で提供している統計表の情報を取得します。リクエストパラメータの指定により条件を絞った情報の取得も可能です。
func GetStatsList(ctx context.Context, query url.Values) (*GetStatsListResult, error) {
	// 2.1. 統計表情報取得
	// http(s)://api.e-stat.go.jp/rest/<バージョン>/app/json/getStatsList?<パラメータ群>
	response := GetStatsListResponse{}

	// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0#api_2_1
	err := callAPI(
		ctx,
		http.MethodGet,
		"http://api.e-stat.go.jp/rest/3.0/app/json/getStatsList?"+query.Encode(),
		func(d *json.Decoder) error { return d.Decode(&response) },
	)

	if err != nil {
		return nil, err
	}

	return &response.GetStatsList, nil
}

// 3.2. 統計表情報取得

// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0#api_4_4
type GetStatsListResponse struct {
	GetStatsList GetStatsListResult `json:"GET_STATS_LIST"`
}

type GetStatsListResult struct {
}
