// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Server serverConfig
	GitHub githubConfig `toml:"github"`
	GitLab gitlabConfig `toml:"gitlab"`
	WeCom  wecomConfig  `toml:"wecom"`
}

type serverConfig struct {
	// ListenAddr is the address brickbot-server listens at.
	//
	// The format is just what net.Listen accepts for TCP.
	ListenAddr string `toml:"listen_addr"`
}

type githubConfig struct {
	Enabled bool
	Secret  string
}

type gitlabConfig struct {
	Enabled bool
	Secret  string
}

type wecomConfig struct {
	Enabled    bool
	CorpID     string `toml:"corpid"`
	CorpSecret string `toml:"corpsecret"`
	AgentID    int64  `toml:"agentid"`
}

func parseConfig(path string) (config, error) {
	var result config
	_, err := toml.DecodeFile(path, &result)
	if err != nil {
		return config{}, err
	}
	return result, nil
}
