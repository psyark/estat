package estat

import (
	"encoding/json"
)

// https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0
type Response struct {
	GetStatsData GetStatsData `json:"GET_STATS_DATA"`
}

type GetStatsData struct {
	Parameter       json.RawMessage `json:"PARAMETER"`
	Result          json.RawMessage `json:"RESULT"`
	StatisticalData StatisticalData `json:"STATISTICAL_DATA"`
}

type StatisticalData struct {
	ResultInf ResultInf       `json:"RESULT_INF"`
	TableInf  TableInf        `json:"TABLE_INF"`
	ClassInf  json.RawMessage `json:"CLASS_INF"`
	DataInf   json.RawMessage `json:"DATA_INF"`
}

type ResultInf struct {
	TotalNumber int `json:"TOTAL_NUMBER"`
	FromNumber  int `json:"FROM_NUMBER"`
	ToNumber    int `json:"TO_NUMBER"`
}

type TableInf struct {
	ID                   string                `json:"@id"`
	COLLECT_AREA         string                `json:"COLLECT_AREA"`
	CYCLE                string                `json:"CYCLE"`
	DESCRIPTION          string                `json:"DESCRIPTION"`
	GOV_ORG              CodeWithDescription   `json:"GOV_ORG"`
	MAIN_CATEGORY        CodeWithDescription   `json:"MAIN_CATEGORY"`
	OPEN_DATE            string                `json:"OPEN_DATE"`
	OVERALL_TOTAL_NUMBER int                   `json:"OVERALL_TOTAL_NUMBER"`
	SMALL_AREA           int                   `json:"SMALL_AREA"`
	STATISTICS_NAME      string                `json:"STATISTICS_NAME"`
	STATISTICS_NAME_SPEC json.RawMessage       `json:"STATISTICS_NAME_SPEC"`
	STAT_NAME            CodeWithDescription   `json:"STAT_NAME"`
	SUB_CATEGORY         CodeWithDescription   `json:"SUB_CATEGORY"`
	SURVEY_DATE          json.RawMessage       `json:"SURVEY_DATE"`
	TITLE                NumberWithDescription `json:"TITLE"`
	TITLE_SPEC           json.RawMessage       `json:"TITLE_SPEC"`
	UPDATED_DATE         string                `json:"UPDATED_DATE"`
}

type CodeWithDescription struct {
	Code        string `json:"@code"`
	Description string `json:"$"`
}

type NumberWithDescription struct {
	No          string `json:"@no"`
	Description string `json:"$"`
}
