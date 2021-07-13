// SPDX-License-Identifier: GPL-3.0-or-later

package gitlab

import (
	"net/http"

	"github.com/go-playground/webhooks/v6/gitlab"

	"github.com/xen0n/brickbot/bot/v1alpha1"
	"github.com/xen0n/brickbot/forge"
)

//nolint:deadcode,unused,varcheck // pending rework
const forgeType = "gitlab"

type gitlabForge struct {
	hook *gitlab.Webhook
}

var _ forge.IForgeHook = (*gitlabForge)(nil)

// New returns a new GitLab forge hook instance.
func New(secret string) (forge.IForgeHook, error) {
	hook, err := gitlab.New(
		gitlab.Options.Secret(secret),
	)
	if err != nil {
		return nil, err
	}

	return &gitlabForge{
		hook: hook,
	}, nil
}

// HookRequest hooks an incoming webhook request to trigger actions.
func (f *gitlabForge) HookRequest(req *http.Request) (*v1alpha1.Event, error) {
	payload, err := f.hook.Parse(
		req,
		// XXX This is everything for now, I don't know exactly what GitLab is
		// going to give out yet for our existing workflow...
		gitlab.PushEvents,
		gitlab.TagEvents,
		gitlab.IssuesEvents,
		gitlab.ConfidentialIssuesEvents,
		gitlab.CommentEvents,
		gitlab.MergeRequestEvents,
		gitlab.WikiPageEvents,
		gitlab.PipelineEvents,
		gitlab.BuildEvents,
		gitlab.JobEvents,
		gitlab.SystemHookEvents,
	)
	if err != nil {
		return nil, err
	}

	// TODO
	_ = payload
	return nil, nil
}
