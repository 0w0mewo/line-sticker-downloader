package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

func HttpGetWithProcessor(ctx context.Context, client *http.Client, url string, processor func(r io.Reader) error) error {
	errch := make(chan error, 1)
	var wg sync.WaitGroup
	pr, pw := io.Pipe()

	// get
	resp, err := get(ctx, client, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read resp body stream from pipe
	// let the callback function prcoess the stream
	wg.Add(1)
	go func() {
		defer wg.Done()
		errch <- processor(pr)

		pr.Close()
	}()

	// write resp body stream to pipe
	io.Copy(pw, resp.Body)
	pw.Close()

	wg.Wait()

	return <-errch
}

func HttpGetJson(ctx context.Context, client *http.Client, url string, res interface{}) error {
	return HttpGetWithProcessor(ctx, client, url, func(r io.Reader) error {
		return json.NewDecoder(r).Decode(res)
	})
}

func header(req *http.Request) {
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0")

}

func get(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	header(req)

	return client.Do(req)
}
