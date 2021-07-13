// SPDX-License-Identifier: GPL-3.0-or-later

package github

import (
	"github.com/go-playground/webhooks/v6/github"

	"github.com/xen0n/brickbot/bot/v1alpha1"
)

// Adapters for whole event param structs

func intoWebhookInstalledParams(x *github.PingPayload) v1alpha1.WebhookInstalledParams {
	return v1alpha1.WebhookInstalledParams{
		Repo: botModelFromRepoContainingPayload(x),
	}
}

func intoPROpenedParams(x *github.PullRequestPayload) v1alpha1.PROpenedParams {
	return v1alpha1.PROpenedParams{
		Actor: botModelFromSenderContainingPayload(x),
		PR:    botModelFromPRContainingPayload(x),
	}
}

func intoPRClosedParams(x *github.PullRequestPayload) v1alpha1.PRClosedParams {
	return v1alpha1.PRClosedParams{
		Actor: botModelFromSenderContainingPayload(x),
		PR:    botModelFromPRContainingPayload(x),
	}
}

func intoPRMergedParams(x *github.PullRequestPayload) v1alpha1.PRMergedParams {
	return v1alpha1.PRMergedParams{
		Actor: botModelFromSenderContainingPayload(x),
		PR:    botModelFromPRContainingPayload(x),
	}
}

func intoPRReadyParams(x *github.PullRequestPayload) v1alpha1.PRReadyParams {
	return v1alpha1.PRReadyParams{
		Actor: botModelFromSenderContainingPayload(x),
		PR:    botModelFromPRContainingPayload(x),
	}
}

// Adapters for component fields

func botModelFromSenderContainingPayload(x interface{}) v1alpha1.ForgeUser {
	switch x := x.(type) {
	case *github.PullRequestPayload:
		return v1alpha1.ForgeUser{
			Forge:    forgeType,
			UserName: x.Sender.Login,
		}

	default:
		panic("should never happen")
	}
}

func botModelFromRepoContainingPayload(x interface{}) v1alpha1.Repo {
	switch x := x.(type) {
	case *github.PingPayload:
		return v1alpha1.Repo{
			User: v1alpha1.ForgeUser{
				Forge:    forgeType,
				UserName: x.Repository.Owner.Login,
			},
			RepoName: x.Repository.Name,
		}

	case *github.PullRequestPayload:
		return v1alpha1.Repo{
			User: v1alpha1.ForgeUser{
				Forge:    forgeType,
				UserName: x.Repository.Owner.Login,
			},
			RepoName: x.Repository.Name,
		}

	default:
		panic("should never happen")
	}
}

func botModelFromPRContainingPayload(x interface{}) v1alpha1.PR {
	switch x := x.(type) {
	case *github.PullRequestPayload:
		return v1alpha1.PR{
			Repo:   botModelFromRepoContainingPayload(x),
			Number: int(x.PullRequest.Number),
			Title:  x.PullRequest.Title,
			Author: v1alpha1.ForgeUser{
				Forge:    forgeType,
				UserName: x.PullRequest.User.Login,
			},
			State: issueStateFromPRState(x.PullRequest.State, x.PullRequest.Merged),
		}

	case *github.PullRequestReviewPayload:
		return v1alpha1.PR{
			Repo:   botModelFromRepoContainingPayload(x),
			Number: int(x.PullRequest.Number),
			Title:  x.PullRequest.Title,
			Author: v1alpha1.ForgeUser{
				Forge:    forgeType,
				UserName: x.PullRequest.User.Login,
			},
			State: issueStateFromPRState(x.PullRequest.State, false), // no "merged" field
		}

	default:
		panic("should never happen")
	}
}

func issueStateFromPRState(state string, merged bool) v1alpha1.IssueState {
	switch state {
	case "open":
		return v1alpha1.IssueStateOpen
	case "closed":
		if merged {
			return v1alpha1.IssueStateMerged
		}
		return v1alpha1.IssueStateClosed
	default:
		panic("should never happen")
	}
}

//nolint:deadcode,unused // used by upcoming logic
func isDraftPRFromPRContainingPayload(x interface{}) bool {
	switch x := x.(type) {
	case *github.PullRequestPayload:
		return x.PullRequest.Draft

	default:
		panic("should never happen")
	}
}
