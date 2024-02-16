package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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
		query.Set("limit", "10000")

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

	err = cyclicTest(data, func() ([]byte, error) {
		if err := json.Unmarshal(data, &typed); err != nil {
			return nil, err
		}
		return json.Marshal(typed)
	})
	if err != nil {
		t.Fatal(err)
	}

	uncached := []estat.TableInf{}
	for _, t := range typed.GetStatsList.DatalistInf.TableInf {
		if _, err := os.Stat(fmt.Sprintf("testdata/%s.json", t.ID)); os.IsNotExist(err) {
			uncached = append(uncached, t)
		}
	}

	{
		t := uncached[rand.Intn(len(uncached))]
		fmt.Println(t.ID, t.GovOrg.Annotation, t.StatisticsName, t.Title)
	}
}
