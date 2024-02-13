package estat

import (
	"context"
	"encoding/json"
	"net/http"
)

func callAPI(ctx context.Context, method string, url string, cb func(*json.Decoder) error) error {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return cb(json.NewDecoder(resp.Body))
}
