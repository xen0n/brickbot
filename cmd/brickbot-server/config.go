// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Server serverConfig `toml:"server"`
	GitHub githubConfig `toml:"github"`
	GitLab gitlabConfig `toml:"gitlab"`
	WeCom  wecomConfig  `toml:"wecom"`
	Bot    botConfig    `toml:"bot"`
}

type serverConfig struct {
	// ListenAddr is the address brickbot-server listens at.
	//
	// The format is just what net.Listen accepts for TCP.
	ListenAddr string `toml:"listen_addr"`
}

type githubConfig struct {
	Enabled bool   `toml:"enabled"`
	Secret  string `toml:"secret"`
}

type gitlabConfig struct {
	Enabled bool   `toml:"enabled"`
	Secret  string `toml:"secret"`
}

type wecomConfig struct {
	Enabled    bool   `toml:"enabled"`
	CorpID     string `toml:"corpid"`
	CorpSecret string `toml:"corpsecret"`
	AgentID    int64  `toml:"agentid"`
}

type botConfig struct {
	PluginPath string `toml:"plugin_path"`
	ConfigPath string `toml:"config_path"`
}

func parseConfig(path string) (config, error) {
	var result config
	_, err := toml.DecodeFile(path, &result)
	if err != nil {
		return config{}, err
	}
	return result, nil
}
