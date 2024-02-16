package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/psyark/estat"
	"github.com/wI2L/jsondiff"
)

func TestGetStatsList(t *testing.T) {
	ctx := context.Background()

	if _, err := os.Stat("testdata/list.json"); os.IsNotExist(err) {
		query := url.Values{}
		query.Set("appId", os.Getenv("appId"))
		query.Set("limit", "100")

		_, err := estat.GetStatsList(ctx, query, estat.WithDataHandler(func(data []byte) error {
			return os.WriteFile("testdata/list.json", data, 0666)
		}))
		if err != nil {
			t.Fatal(err)
		}
	}

	data, err := os.ReadFile("testdata/list.json")
	if err != nil {
		t.Fatal(err)
	}

	typed := estat.GetStatsListContainer{}
	untyped := map[string]any{}
	if err := json.Unmarshal(data, &untyped); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &typed); err != nil {
		t.Fatal(err)
	}

	data1, _ := json.MarshalIndent(untyped, "", "  ")
	data2, _ := json.MarshalIndent(typed, "", "  ")

	patch, err := jsondiff.CompareJSON(data1, data2)
	if err != nil {
		t.Fatal(err)
	}

	for _, op := range patch {
		_, _ = fmt.Printf("%v %v: value=%v, oldValue=%v\n", op.Type, op.Path, op.Value, op.OldValue)
	}

	if len(patch) != 0 {
		t.Fatal("unmatch")
	}
}
