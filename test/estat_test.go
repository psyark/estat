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

	query, err := url.ParseQuery("lang=J&statsDataId=0004009602&metaGetFlg=Y&cntGetFlg=N&explanationGetFlg=Y&annotationGetFlg=Y&sectionHeaderFlg=1&replaceSpChars=0")
	if err != nil {
		t.Fatal(err)
	}

	query.Set("appId", os.Getenv("appId"))

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
	data2, _ := json.MarshalIndent(typed, "", "  ")
	patch, err := jsondiff.CompareJSON(data2, data1)
	if err != nil {
		t.Fatal(err)
	}

	for _, op := range patch {
		fmt.Printf("%s\n", op)
	}

	if len(patch) != 0 {
		t.Fatal("unmatch")
	}
}
