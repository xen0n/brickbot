// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/xen0n/brickbot/bot"
	"github.com/xen0n/brickbot/bot/v1alpha1"
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

	botPlugin, err := bot.LoadPlugin(conf.Bot.PluginPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load bot plugin")
		os.Exit(1)
	}

	bot, err := botPlugin.InitWithConfigTOML(conf.Bot.ConfigPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to construct bot plugin")
		os.Exit(2)
	}

	err = bot.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup bot plugin")
		os.Exit(2)
	}

	srv, err := makeServer(&conf, bot)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	exitcodeChan := make(chan int)

	go func() {
		// only SIGINT for now
		// gracefully quit on catching that signal
		<-signalChan

		log.Info().Msg("caught SIGINT")

		// TODO: timeout for graceful quit
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("error occurred during shutdown")
			exitcodeChan <- 10
			return
		}

		// shutdown successful
		exitcodeChan <- 0
	}()

	err = srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to serve")
			os.Exit(1)
		}

		// server closed, do nothing
	}

	// wait for shutdown to complete
	exitcode := <-exitcodeChan
	os.Exit(exitcode)
}

func makeServer(conf *config, bot v1alpha1.IPlugin) (*http.Server, error) {
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
			return nil, err
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
				return nil, err
			}

			mux.HandleFunc("/github", makeForgeHookHandler(fh, bot, wecom))
		}

		if conf.GitLab.Enabled {
			fh, err := forgeGL.New(conf.GitLab.Secret)
			if err != nil {
				log.Error().Err(err).Msg("failed to initialize GitLab integration")
				return nil, err
			}

			mux.HandleFunc("/gitlab", makeForgeHookHandler(fh, bot, wecom))
		}
	}

	return &http.Server{
		Addr:        conf.Server.ListenAddr,
		Handler:     mux,
		ReadTimeout: 1 * time.Minute,
	}, nil
}

func dummyHealthzHandler(rw http.ResponseWriter, r *http.Request) {
	// Does nothing for now, just report healthy.
	rw.WriteHeader(http.StatusOK)
}

func makeForgeHookHandler(
	fh forge.IForgeHook,
	bot v1alpha1.IPlugin,
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

		// Call bot plugin asynchronously.
		go func(e *v1alpha1.Event) {
			err := bot.ProcessEvent(e, imProvider)
			if err != nil {
				log.Error().Err(err).Msg("bot returned failure")
			}
		}(botEvent)

		// Most webhooks ignore the response body, but might retry in case of
		// failed deliveries, so send 204.
		rw.WriteHeader(http.StatusNoContent)
	}
}
