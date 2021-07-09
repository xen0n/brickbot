// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/xen0n/brickbot/forge"
	forgeGH "github.com/xen0n/brickbot/forge/github"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("configPath", configPath)

	conf, err := parseConfig(configPath)
	if err != nil {
		panic(err)
	}

	err = runServer(&conf)
	if err != nil {
		panic(err)
	}
}

func runServer(conf *config) error {
	mux := http.NewServeMux()

	// Health check endpoints.
	mux.HandleFunc("/healthz", dummyHealthzHandler)
	mux.HandleFunc("/livez", dummyHealthzHandler)
	mux.HandleFunc("/readyz", dummyHealthzHandler)

	// Metrics endpoints.
	mux.Handle("/metrics", promhttp.Handler())

	// Webhook endpoints.
	{
		if conf.GitHub.Enabled {
			fh, err := forgeGH.New(conf.GitHub.Secret)
			if err != nil {
				panic(err)
			}

			mux.HandleFunc("/github", makeForgeHookHandler(fh))
		}
	}

	return http.ListenAndServe(conf.Server.ListenAddr, mux)
}

func dummyHealthzHandler(rw http.ResponseWriter, r *http.Request) {
	// Does nothing for now, just report healthy.
	rw.WriteHeader(http.StatusOK)
}

func makeForgeHookHandler(fh forge.IForgeHook) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Invoke the forge-specific logic.
		fh.HookRequest(r)

		// Most webhooks ignore the response body, but might retry in case of
		// failed deliveries, so send 204.
		rw.WriteHeader(http.StatusNoContent)
	}
}
