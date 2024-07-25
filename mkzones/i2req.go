package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var i2trust = http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

func i2req(method string, uri string, body any, resp any) error {
	var bdy io.ReadCloser
	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			panic(err)
		}

		bdy = io.NopCloser(buf)
	}

	response, err := i2trust.Do(&http.Request{
		Method: method,
		URL: &url.URL{
			Scheme: "https",
			User:   url.UserPassword("root", "123456"),
			Host:   "master1:5665",
			Path:   uri,
		},
		Header: http.Header{"Accept": []string{"application/json"}},
		Body:   bdy,
	})
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	if resp != nil {
		return json.NewDecoder(response.Body).Decode(resp)
	}

	return nil
}
