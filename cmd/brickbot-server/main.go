// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/xen0n/brickbot/forge"
	forgeGH "github.com/xen0n/brickbot/forge/github"
	forgeGL "github.com/xen0n/brickbot/forge/gitlab"
	"github.com/xen0n/brickbot/im"
	imWeCom "github.com/xen0n/brickbot/im/wecom"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var configPath string
	flag.StringVar(&configPath, "c", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.Debug().Str("path", configPath).Msg("using this config")

	conf, err := parseConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse config file")
		os.Exit(1)
	}

	err = runServer(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
		os.Exit(1)
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
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to initialize WeCom integration")
			return err
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
				log.Error().Err(err).Msg("failed to initialize GitHub integration")
				return err
			}

			mux.HandleFunc("/github", makeForgeHookHandler(fh, wecom))
		}

		if conf.GitLab.Enabled {
			fh, err := forgeGL.New(conf.GitLab.Secret)
			if err != nil {
				log.Error().Err(err).Msg("failed to initialize GitLab integration")
				return err
			}

			mux.HandleFunc("/gitlab", makeForgeHookHandler(fh, wecom))
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
		botEvent, err := fh.HookRequest(r)
		if err != nil {
			log.Error().Err(err).Msg("failed to process incoming webhook event")

			// TODO: is returning failure the best thing to do in this case?
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if botEvent == nil {
			// Event is boring, do nothing.
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		log.Debug().Str("event", fmt.Sprintf("%+v", botEvent)).Msg("parsed incoming event")

		// TODO: bot logic

		// Most webhooks ignore the response body, but might retry in case of
		// failed deliveries, so send 204.
		rw.WriteHeader(http.StatusNoContent)
	}
}
