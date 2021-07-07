// SPDX-License-Identifier: GPL-3.0-or-later

package github

import (
	"net/http"

	"github.com/go-playground/webhooks/v6/github"

	"github.com/xen0n/brickbot/forge"
)

type githubForge struct {
	hook *github.Webhook
}

var _ forge.IForgeHook = (*githubForge)(nil)

// New returns a new GitHub forge hook instance.
func New(secret string) (forge.IForgeHook, error) {
	hook, err := github.New(
		github.Options.Secret(secret),
	)
	if err != nil {
		return nil, err
	}

	return &githubForge{
		hook: hook,
	}, nil
}

// HookRequest hooks an incoming webhook request to trigger actions.
func (f *githubForge) HookRequest(req *http.Request) {
	panic("TODO")
}
