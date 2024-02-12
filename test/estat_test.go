package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/psyark/estat"
	"github.com/wI2L/jsondiff"
)

func TestXxx(t *testing.T) {
	if err := godotenv.Load("testdata/secret.env"); err != nil {
		t.Fatal(err)
	}

	for _, statsDataId := range []string{"0003354197", "0004009602", "0003313482", "0002019042"} {
		statsDataId := statsDataId
		t.Run(statsDataId, func(t *testing.T) {
			query, err := url.ParseQuery("lang=J&&explanationGetFlg=Y&annotationGetFlg=Y&sectionHeaderFlg=1&replaceSpChars=0")
			if err != nil {
				t.Fatal(err)
			}

			query.Set("appId", os.Getenv("appId"))
			query.Set("statsDataId", statsDataId)

			resp, err := http.Get("http://api.e-stat.go.jp/rest/3.0/app/json/getStatsData?" + query.Encode())
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			untyped := map[string]any{}
			if err := json.Unmarshal(data, &untyped); err != nil {
				t.Fatal(err)
			}

			typed := estat.Response{}
			if err := json.Unmarshal(data, &typed); err != nil {
				t.Fatal(err)
			}

			data1, _ := json.MarshalIndent(untyped, "", "  ")

			os.WriteFile(fmt.Sprintf("testdata/%s.json", statsDataId), data1, 0666)

			data2, _ := json.MarshalIndent(typed, "", "  ")

			patch, err := jsondiff.CompareJSON(data1, data2)
			if err != nil {
				t.Fatal(err)
			}

			for _, op := range patch {
				fmt.Printf("%v %v: value=%v, oldValue=%v\n", op.Type, op.Path, op.Value, op.OldValue)
			}

			if len(patch) != 0 {
				t.Fatal("unmatch")
			}

			// fmt.Println(typed.GetStatsData.StatisticalData.TableInf.GovOrg.Annotation)
			// fmt.Println(typed.GetStatsData.StatisticalData.TableInf.StatisticsName)
			// fmt.Println(typed.GetStatsData.StatisticalData.TableInf.Title)

			// fmt.Println(string(data2))
			for _, classObj := range typed.GetStatsData.StatisticalData.ClassInf.ClassObj {
				fmt.Println(classObj.ID, classObj.Name)
				for _, class := range classObj.Class {
					fmt.Println("    ", class)
				}
			}
		})
	}
}
