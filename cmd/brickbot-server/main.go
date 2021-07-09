// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
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

	return http.ListenAndServe(conf.Server.ListenAddr, mux)
}

func dummyHealthzHandler(rw http.ResponseWriter, r *http.Request) {
	// Does nothing for now, just report healthy.
	rw.WriteHeader(http.StatusOK)
}
