// SPDX-License-Identifier: CC0-1.0

//nolint:goheader // Examples are put in public domain.

//+build example

package main

import (
	"errors"
	"fmt"

	"github.com/xen0n/brickbot/bot/v1alpha1"
)

var BrickbotPluginAPIVersion = v1alpha1.PluginAPIVersion

type pluginConfig struct {
	TeamChatID string `toml:"team_chatid"`
}

func BrickbotPluginConfigFactory() interface{} {
	return pluginConfig{}
}

type plugin struct {
	teamChatID string
}

var _ v1alpha1.IPlugin = (*plugin)(nil)

func BrickbotPluginFactory(config interface{}) (v1alpha1.IPlugin, error) {
	c, ok := config.(pluginConfig)
	if !ok {
		return nil, errors.New("wrong config type; should never happen")
	}

	return &plugin{
		teamChatID: c.TeamChatID,
	}, nil
}

func (p *plugin) Setup() error {
	return nil
}

func (p *plugin) ProcessEvent(e *v1alpha1.Event, im v1alpha1.IIMProvider) error {
	switch e.Type() {
	case v1alpha1.EventTypeWebhookInstalled:
		ee, _ := e.WebhookInstalled()
		err := im.SendTextToChat(
			p.teamChatID,
			fmt.Sprintf("🎉 搬砖 Bot 在 %s/%s 安装成功！", ee.Repo.User.UserName, ee.Repo.RepoName),
		)
		if err != nil {
			return err
		}

	case v1alpha1.EventTypePROpened:
		ee, _ := e.PROpened()
		err := im.SendTextToChat(
			p.teamChatID,
			fmt.Sprintf(
				"%s 提交了 %s/%s #%d\n\n%s",
				ee.Actor.UserName,
				ee.PR.Repo.User.UserName,
				ee.PR.Repo.RepoName,
				ee.PR.Number,
				ee.PR.Title,
			),
		)
		if err != nil {
			return err
		}

	case v1alpha1.EventTypePRMerged:
		ee, _ := e.PRMerged()
		err := im.SendTextToChat(
			p.teamChatID,
			fmt.Sprintf(
				"%s 合并了 %s/%s #%d\n\n%s",
				ee.Actor.UserName,
				ee.PR.Repo.User.UserName,
				ee.PR.Repo.RepoName,
				ee.PR.Number,
				ee.PR.Title,
			),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *plugin) Teardown() error {
	return nil
}
