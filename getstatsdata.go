package estat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// 指定した統計表ID又はデータセットIDに対応する統計データ（数値データ）を取得します。
func GetStatsData(ctx context.Context, query url.Values, options ...Option) (*GetStatsDataContent, error) {
	container := GetStatsDataContainer{}
	options = append(options, WithDataHandler(func(data []byte) error { return json.Unmarshal(data, &container) }))

	// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0#api_2_3
	return &container.GetStatsData, callAPI(
		ctx,
		http.MethodGet,
		"http://api.e-stat.go.jp/rest/3.0/app/json/getStatsData?"+query.Encode(),
		options...,
	)
}

// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0#api_4_4
type GetStatsDataContainer struct {
	GetStatsData GetStatsDataContent `json:"GET_STATS_DATA"`
}

type GetStatsDataContent struct {
	Result          Result                `json:"RESULT"`
	Parameter       GetStatsDataParameter `json:"PARAMETER"`
	StatisticalData *StatisticalData      `json:"STATISTICAL_DATA,omitempty"`
}

// GetStatsDataParameter is 4.4.1. PARAMETER タグ
// リクエスト時に指定されたパラメータを出力します。パラメータ名を間違えた場合や別のAPIのパラメータを指定した場合は出力されません。
type GetStatsDataParameter struct {
	AnnotationGetFlg  string `json:"ANNOTATION_GET_FLG,omitempty"`
	CntGetFlg         string `json:"CNT_GET_FLG,omitempty"`
	DataFormat        string `json:"DATA_FORMAT"`
	ExplanationGetFlg string `json:"EXPLANATION_GET_FLG,omitempty"`
	Lang              string `json:"LANG"`
	MetagetFlg        string `json:"METAGET_FLG"`
	ReplaceSpChars    int    `json:"REPLACE_SP_CHARS,omitempty"`
	SectionHeaderFlg  int    `json:"SECTION_HEADER_FLG,omitempty"`
	StartPosition     int    `json:"START_POSITION"`
	StatsDataID       string `json:"STATS_DATA_ID"`
	Limit             int    `json:"LIMIT,omitempty"`
}

// StatisticalData is 4.4.2. STATISTICAL_DATA タグ
// 統計データの情報を出力します。エラーがあった場合はこのタグ自体出力されません。
type StatisticalData struct {
	ResultInf GetStatsDataResultInf `json:"RESULT_INF"`
	TableInf  TableInf              `json:"TABLE_INF"`
	ClassInf  ClassInf              `json:"CLASS_INF"`
	DataInf   DataInf               `json:"DATA_INF"`
}

type GetStatsDataResultInf struct {
	TotalNumber int `json:"TOTAL_NUMBER"`
	FromNumber  int `json:"FROM_NUMBER"`
	ToNumber    int `json:"TO_NUMBER"`
	NextKey     int `json:"NEXT_KEY,omitempty"`
}

type TableInf struct {
	ID                 string             `json:"@id"`
	CollectArea        string             `json:"COLLECT_AREA"`
	Cycle              string             `json:"CYCLE"`
	Description        string             `json:"DESCRIPTION"`
	GovOrg             AnnotatedCode      `json:"GOV_ORG"`
	MainCategory       AnnotatedCode      `json:"MAIN_CATEGORY"`
	OpenDate           string             `json:"OPEN_DATE"`
	OverallTotalNumber int                `json:"OVERALL_TOTAL_NUMBER"`
	SmallArea          int                `json:"SMALL_AREA"`
	StatisticsName     string             `json:"STATISTICS_NAME"`
	StatisticsNameSpec StatisticsNameSpec `json:"STATISTICS_NAME_SPEC"`
	StatName           AnnotatedCode      `json:"STAT_NAME"`
	SubCategory        AnnotatedCode      `json:"SUB_CATEGORY"`
	SurveyDate         any                `json:"SURVEY_DATE"` // "200104-200203" or 0
	Title              any                `json:"TITLE"`       // "第１表　月（12区分）、施設所在地(47区分及び運輸局等)、従業者数(4区分)、宿泊目的割合(2区分)別施設数" or {"$":"新規就農者調査 就農形態別新規就農者数", "@no":1}
	TitleSpec          TitleSpec          `json:"TITLE_SPEC"`
	UpdatedDate        string             `json:"UPDATED_DATE"`
}
type StatisticsNameSpec struct {
	TabulationCategory     string `json:"TABULATION_CATEGORY"`
	TabulationSubCategory1 string `json:"TABULATION_SUB_CATEGORY1,omitempty"`
	TabulationSubCategory2 string `json:"TABULATION_SUB_CATEGORY2,omitempty"`
	TabulationSubCategory3 string `json:"TABULATION_SUB_CATEGORY3,omitempty"`
	TabulationSubCategory4 string `json:"TABULATION_SUB_CATEGORY4,omitempty"`
}
type TitleSpec struct {
	TableCategory     string `json:"TABLE_CATEGORY,omitempty"`
	TableName         string `json:"TABLE_NAME,omitempty"`
	TableExplanation  string `json:"TABLE_EXPLANATION,omitempty"`
	TableSubCategory1 any    `json:"TABLE_SUB_CATEGORY1,omitempty"` // 0003090650 で 数値が入る
	TableSubCategory2 any    `json:"TABLE_SUB_CATEGORY2,omitempty"`
	TableSubCategory3 any    `json:"TABLE_SUB_CATEGORY3,omitempty"`
}

type ClassInf struct {
	ClassObj []ClassObj `json:"CLASS_OBJ"`
}

func (i *ClassInf) GetClassObj(id string) *ClassObj {
	for _, o := range i.ClassObj {
		if o.ID == id {
			return &o
		}
	}
	return nil
}

type ClassObj struct {
	ID          string        `json:"@id"`
	Name        string        `json:"@name"`
	Description string        `json:"@description,omitempty"`
	Class       ClassList     `json:"CLASS"`
	Explanation []AnnotatedId `json:"EXPLANATION,omitempty"`
}

func (o *ClassObj) GetClass(code string) *Class {
	for _, c := range o.Class {
		if c.Code == code {
			return &c
		}
	}
	return nil
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
	Note  NoteList  `json:"NOTE,omitempty"`
	Value ValueList `json:"VALUE"`
}

type Note struct {
	Annotation string `json:"$"`
	Char       string `json:"@char"`
}

// VALUE	統計数値(セル)の情報です。データ件数分だけ出力されます。
// 属性として表章事項コード(tab)、分類事項コード(cat01 ～ cat15)、地域事項コード(area)、時間軸事項コード(time)、単位(unit)、注釈記号(anotation)を保持します。全ての属性はデータがある場合のみ出力されます。
type Value struct {
	Value string `json:"$"` // 二重引用符付きの数値が入ったり、特殊文字 "-" 等が入ったりする
	Tab   string `json:"@tab,omitempty"`
	Cat01 string `json:"@cat01,omitempty"`
	Cat02 string `json:"@cat02,omitempty"`
	Cat03 string `json:"@cat03,omitempty"`
	Cat04 string `json:"@cat04,omitempty"`
	Cat05 string `json:"@cat05,omitempty"`
	Cat06 string `json:"@cat06,omitempty"`
	Cat07 string `json:"@cat07,omitempty"`
	Cat08 string `json:"@cat08,omitempty"`
	Cat09 string `json:"@cat09,omitempty"`
	Cat10 string `json:"@cat10,omitempty"`
	Cat11 string `json:"@cat11,omitempty"`
	Cat12 string `json:"@cat12,omitempty"`
	Cat13 string `json:"@cat13,omitempty"`
	Cat14 string `json:"@cat14,omitempty"`
	Cat15 string `json:"@cat15,omitempty"`
	Area  string `json:"@area,omitempty"`
	Time  string `json:"@time,omitempty"`
	Unit  string `json:"@unit,omitempty"`
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
