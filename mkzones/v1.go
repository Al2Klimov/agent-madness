package main

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var word = regexp.MustCompile(`\A\w+\z`)

func v1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := query.Get("name")
	if !word.MatchString(name) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pending.Lock()
	pending.names[name] = struct{}{}

	pending.cancel()
	ctx, cancel := context.WithCancel(context.Background())
	go signal(ctx)
	pending.cancel = cancel

	pending.Unlock()

	w.WriteHeader(http.StatusAccepted)
}

func signal(ctx context.Context) {
	select {
	case <-time.After(10 * time.Second):
		pending.cond.Signal()
	case <-ctx.Done():
	}
}
