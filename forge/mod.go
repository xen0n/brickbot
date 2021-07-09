// SPDX-License-Identifier: GPL-3.0-or-later

package forge

import "net/http"

// IForgeHook is the interface that all forge hooks implement.
type IForgeHook interface {
	// HookRequest hooks an incoming webhook request to trigger actions.
	HookRequest(req *http.Request) (*HookResult, error)
}

// HookResult represents a parsed event from a hook invocation.
type HookResult struct {
	// IsInteresting specifies whether the event should be processed by bot
	// instances.
	IsInteresting bool
	// Event contains the parsed event if IsInteresting is true.
	Event IEvent
}

// IEvent is abstraction for forge-generated events.
type IEvent interface {
	// Raw returns the raw payload from forges.
	//
	// TODO: This is for debugging purposes, and is very likely to be removed
	// before initial release.
	Raw() interface{}
}
