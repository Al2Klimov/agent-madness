package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

//go:embed zones.conf
var zonesConf string

func deploy() {
	for {
		pending.Lock()
		for len(pending.names) < 1 {
			pending.cond.Wait()
		}

		names := pending.names
		pending.names = map[string]struct{}{}
		pending.Unlock()

		var zones struct {
			Results []struct {
				Name string `json:"name"`
			} `json:"results"`
		}

		for {
			err := i2req(http.MethodGet, "/v1/objects/zones", map[string][]string{"attrs": {"name"}}, &zones)
			if err == nil {
				break
			}

			_, _ = fmt.Fprintf(os.Stderr, "failed to fetch zones: %v\n", err)
			time.Sleep(3 * time.Second)
		}

		for _, zone := range zones.Results {
			delete(names, zone.Name)
		}

		if len(names) < 1 {
			continue
		}

		var pkg string
		for {
			pkg = strconv.FormatInt(time.Now().Unix(), 10)

			err := i2req(http.MethodPost, "/v1/config/packages/"+pkg, nil, nil)
			if err == nil {
				break
			}

			_, _ = fmt.Fprintf(os.Stderr, "failed to create package %q: %v\n", pkg, err)
			time.Sleep(3 * time.Second)
		}

		cfg := make(map[string]string, len(names))
		for name := range names {
			cfg["zones.d/master/"+name+".conf"] = fmt.Sprintf(zonesConf, name)
		}

		for {
			err := i2req(http.MethodPost, "/v1/config/stages/"+pkg, map[string]any{"files": cfg}, nil)
			if err == nil {
				break
			}

			_, _ = fmt.Fprintf(os.Stderr, "failed to deploy package %q: %v\n", pkg, err)
			time.Sleep(3 * time.Second)
		}

		for name := range names {
			for {
				err := i2req(http.MethodGet, "/v1/objects/zones/"+name, map[string][]string{"attrs": {"name"}}, nil)
				if err == nil {
					break
				}

				_, _ = fmt.Fprintf(os.Stderr, "failed to fetch zone %q: %v\n", name, err)
				time.Sleep(10 * time.Second)
			}

			break
		}
	}
}
