package estat

import (
	"context"
	"net/url"
)

// 政府統計の総合窓口（e-Stat）で提供している統計表の情報を取得します。リクエストパラメータの指定により条件を絞った情報の取得も可能です。
func GetStatsList(ctx context.Context, query url.Values) {
	// 2.1. 統計表情報取得
	// http(s)://api.e-stat.go.jp/rest/<バージョン>/app/json/getStatsList?<パラメータ群>
}

// 3.2. 統計表情報取得
