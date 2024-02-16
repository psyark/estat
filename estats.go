package estat

import (
	"context"
	"io"
	"net/http"
)

type optionStruct struct {
	dataHandlers []func([]byte) error
}
type Option func(o *optionStruct)

func callAPI(ctx context.Context, method string, url string, options ...Option) error {
	os := optionStruct{}
	for _, o := range options {
		o(&os)
	}

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

	for _, dh := range os.dataHandlers {
		if err := dh(data); err != nil {
			return err
		}
	}

	return nil
}

func WithDataHandler(handler func(b []byte) error) Option {
	return func(o *optionStruct) {
		o.dataHandlers = append(o.dataHandlers, handler)
	}
}
