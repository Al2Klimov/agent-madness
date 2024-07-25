package main

import (
	"net/http"
	"sync"
)

var pending struct {
	sync.Mutex

	cond   sync.Cond
	cancel func()
	names  map[string]struct{}
}

func main() {
	pending.cond.L = &pending
	pending.cancel = doNothing
	pending.names = map[string]struct{}{}

	go deploy()

	http.HandleFunc("/v1", v1)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func doNothing() {
}
