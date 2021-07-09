// SPDX-License-Identifier: GPL-3.0-or-later

package github

import (
	"net/http"

	"github.com/go-playground/webhooks/v6/github"

	"github.com/xen0n/brickbot/forge"
)

type githubEvent struct {
	inner interface{}
}

var _ forge.IEvent = (*githubEvent)(nil)

// Raw returns the raw payload from forges.
//
// TODO: This is for debugging purposes, and is very likely to be removed
// before initial release.
func (e *githubEvent) Raw() interface{} {
	return e.inner
}

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
func (f *githubForge) HookRequest(req *http.Request) (*forge.HookResult, error) {
	payload, err := f.hook.Parse(
		req,
		// XXX This is everything for now, I don't know exactly what GitHub is
		// going to give out yet for our existing workflow...
		github.CheckRunEvent,
		github.CheckSuiteEvent,
		github.CommitCommentEvent,
		github.CreateEvent,
		github.DeleteEvent,
		github.DeploymentEvent,
		github.DeploymentStatusEvent,
		github.ForkEvent,
		github.GollumEvent,
		github.InstallationEvent,
		github.InstallationRepositoriesEvent,
		github.IntegrationInstallationEvent,
		github.IntegrationInstallationRepositoriesEvent,
		github.IssueCommentEvent,
		github.IssuesEvent,
		github.LabelEvent,
		github.MemberEvent,
		github.MembershipEvent,
		github.MilestoneEvent,
		github.MetaEvent,
		github.OrganizationEvent,
		github.OrgBlockEvent,
		github.PageBuildEvent,
		github.PingEvent,
		github.ProjectCardEvent,
		github.ProjectColumnEvent,
		github.ProjectEvent,
		github.PublicEvent,
		github.PullRequestEvent,
		github.PullRequestReviewEvent,
		github.PullRequestReviewCommentEvent,
		github.PushEvent,
		github.ReleaseEvent,
		github.RepositoryEvent,
		github.RepositoryVulnerabilityAlertEvent,
		github.SecurityAdvisoryEvent,
		github.StatusEvent,
		github.TeamEvent,
		github.TeamAddEvent,
		github.WatchEvent,
	)
	if err != nil {
		return nil, err
	}

	return &forge.HookResult{
		IsInteresting: true,
		Event: &githubEvent{
			inner: payload,
		},
	}, nil
}
