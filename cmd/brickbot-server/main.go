// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
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

	fmt.Println(conf)
}
