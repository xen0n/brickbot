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
	"github.com/xen0n/brickbot/im"
	imWeCom "github.com/xen0n/brickbot/im/wecom"
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
	// IM integration.
	var wecom im.IProvider
	if conf.WeCom.Enabled {
		p, err := imWeCom.New(
			conf.WeCom.CorpID,
			conf.WeCom.CorpSecret,
			conf.WeCom.AgentID,
			conf.WeCom.ChatID,
		)
		if err != nil {
			panic(err)
		}

		wecom = p
	}

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

			mux.HandleFunc("/github", makeForgeHookHandler(fh, wecom))
		}
	}

	return http.ListenAndServe(conf.Server.ListenAddr, mux)
}

func dummyHealthzHandler(rw http.ResponseWriter, r *http.Request) {
	// Does nothing for now, just report healthy.
	rw.WriteHeader(http.StatusOK)
}

func makeForgeHookHandler(
	fh forge.IForgeHook,
	imProvider im.IProvider,
) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Invoke the forge-specific logic.
		hookResult, err := fh.HookRequest(r)
		if err != nil {
			// TODO: is returning failure the best thing to do in this case?
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if hookResult.IsInteresting {
			// Only process interesting events.
			// For now, directly report to IM for debugging.
			if imProvider != nil {
				err := imProvider.SendTeamMessage(hookResult.Event)
				if err != nil {
					// This error is not related to webhook request itself, so
					// don't return failure status code.
					fmt.Printf("XXX failed to send team message: %s\n", err.Error())
				}
			}
		}

		// Most webhooks ignore the response body, but might retry in case of
		// failed deliveries, so send 204.
		rw.WriteHeader(http.StatusNoContent)
	}
}
