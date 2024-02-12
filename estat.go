package estat

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

func GetStats(query url.Values) (*GetStatsData, error) {
	resp, err := http.Get("http://api.e-stat.go.jp/rest/3.0/app/json/getStatsData?" + query.Encode())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := Response{}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	return &response.GetStatsData, nil
}

// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0
type Response struct {
	GetStatsData GetStatsData `json:"GET_STATS_DATA"`
}

type GetStatsData struct {
	Result          Result          `json:"RESULT"`
	Parameter       Parameter       `json:"PARAMETER"`
	StatisticalData StatisticalData `json:"STATISTICAL_DATA"`
}

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
		Alias(r),
		r.Date.Format("2006-01-02T15:04:05.000-07:00"),
	}
	return json.Marshal(t)
}

// Parameter is 4.4.1. PARAMETER タグ
// リクエスト時に指定されたパラメータを出力します。パラメータ名を間違えた場合や別のAPIのパラメータを指定した場合は出力されません。
type Parameter struct {
	AnnotationGetFlg  string `json:"ANNOTATION_GET_FLG"`
	CntGetFlg         string `json:"CNT_GET_FLG"`
	DataFormat        string `json:"DATA_FORMAT"`
	ExplanationGetFlg string `json:"EXPLANATION_GET_FLG"`
	Lang              string `json:"LANG"`
	MetagetFlg        string `json:"METAGET_FLG"`
	ReplaceSpChars    int    `json:"REPLACE_SP_CHARS"`
	SectionHeaderFlg  int    `json:"SECTION_HEADER_FLG"`
	StartPosition     int    `json:"START_POSITION"`
	StatsDataID       string `json:"STATS_DATA_ID"`
}

// StatisticalData is 4.4.2. STATISTICAL_DATA タグ
// 統計データの情報を出力します。エラーがあった場合はこのタグ自体出力されません。
type StatisticalData struct {
	ResultInf ResultInf `json:"RESULT_INF"`
	TableInf  TableInf  `json:"TABLE_INF"`
	ClassInf  ClassInf  `json:"CLASS_INF"`
	DataInf   DataInf   `json:"DATA_INF"`
}

type ResultInf struct {
	TotalNumber int `json:"TOTAL_NUMBER"`
	FromNumber  int `json:"FROM_NUMBER"`
	ToNumber    int `json:"TO_NUMBER"`
}

type TableInf struct {
	ID                   string          `json:"@id"`
	COLLECT_AREA         string          `json:"COLLECT_AREA"`
	CYCLE                string          `json:"CYCLE"`
	DESCRIPTION          string          `json:"DESCRIPTION"`
	GOV_ORG              AnnotatedCode   `json:"GOV_ORG"`
	MAIN_CATEGORY        AnnotatedCode   `json:"MAIN_CATEGORY"`
	OPEN_DATE            string          `json:"OPEN_DATE"`
	OVERALL_TOTAL_NUMBER int             `json:"OVERALL_TOTAL_NUMBER"`
	SMALL_AREA           int             `json:"SMALL_AREA"`
	STATISTICS_NAME      string          `json:"STATISTICS_NAME"`
	STATISTICS_NAME_SPEC json.RawMessage `json:"STATISTICS_NAME_SPEC"`
	STAT_NAME            AnnotatedCode   `json:"STAT_NAME"`
	SUB_CATEGORY         AnnotatedCode   `json:"SUB_CATEGORY"`
	SURVEY_DATE          json.RawMessage `json:"SURVEY_DATE"`
	TITLE                json.RawMessage `json:"TITLE"`
	TITLE_SPEC           json.RawMessage `json:"TITLE_SPEC"`
	UPDATED_DATE         string          `json:"UPDATED_DATE"`
}

type ClassInf struct {
	ClassObj []ClassObj `json:"CLASS_OBJ"`
}
type ClassObj struct {
	ID          string        `json:"@id"`
	Name        string        `json:"@name"`
	Class       ClassHelper   `json:"CLASS"`
	Explanation []AnnotatedId `json:"EXPLANATION,omitempty"`
}

type Class struct {
	Code       string `json:"@code"`
	Level      string `json:"@level"`
	Name       string `json:"@name"`
	Unit       string `json:"@unit,omitempty"`
	ParentCode string `json:"@parentCode,omitempty"`
}

// DATA_INF	統計データの数値情報を出力します。
// 指定した絞り込み条件又はデータセットの条件又はその両方の条件によって抽出されるデータ件数が 0 の場合、このタグは出力されません。
// また、件数取得フラグ(cntGetFlg)に”Y”(件数のみ取得する)を指定した場合も出力されません。
type DataInf struct {
	Note  NoteHelper `json:"NOTE"`
	Value []Value    `json:"VALUE"`
}

type Note struct {
	Annotation string `json:"$"`
	Char       string `json:"@char"`
}

// VALUE	統計数値(セル)の情報です。データ件数分だけ出力されます。
// 属性として表章事項コード(tab)、分類事項コード(cat01 ～ cat15)、地域事項コード(area)、時間軸事項コード(time)、単位(unit)、注釈記号(anotation)を保持します。全ての属性はデータがある場合のみ出力されます。
type Value struct {
	Value any    `json:"$"`
	Tab   string `json:"@tab,omitempty"`
	Cat01 string `json:"@cat01,omitempty"`
	Cat02 string `json:"@cat02,omitempty"`
	Cat03 string `json:"@cat03,omitempty"`
	Area  string `json:"@area,omitempty"`
	Time  string `json:"@time"`
	Unit  string `json:"@unit"`
}

type AnnotatedCode struct {
	Code       string `json:"@code"`
	Annotation string `json:"$"`
}

type AnnotatedNo struct {
	No         string `json:"@no"`
	Annotation string `json:"$"`
}

type AnnotatedId struct {
	Id         string `json:"@id"`
	Annotation string `json:"$"`
}
