package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/psyark/estat"
	"github.com/wI2L/jsondiff"
)

func TestXxx2(t *testing.T) {

	query := url.Values{}
	query.Set("appId", os.Getenv("appId"))
	query.Set("limit", "10")

	resp, err := http.Get("http://api.e-stat.go.jp/rest/3.0/app/json/getStatsList?" + query.Encode())
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

	typed := estat.GetStatsListContainer{}
	if err := json.Unmarshal(data, &typed); err != nil {
		t.Fatal(err)
	}

	data1, _ := json.MarshalIndent(untyped, "", "  ")

	os.WriteFile("testdata/list.json", data1, 0666)

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

}
