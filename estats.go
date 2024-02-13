package estat

import (
	"context"
	"io"
	"net/http"
)

func callAPI(ctx context.Context, method string, url string, cb func([]byte) error) error {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return cb(data)
}
