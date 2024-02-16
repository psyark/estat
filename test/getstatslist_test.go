package test

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"testing"

	"github.com/psyark/estat"
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

	err = cyclicTest(data, func() ([]byte, error) {
		typed := estat.GetStatsListContainer{}
		if err := json.Unmarshal(data, &typed); err != nil {
			return nil, err
		}
		return json.Marshal(typed)
	})
	if err != nil {
		t.Fatal(err)
	}
}
