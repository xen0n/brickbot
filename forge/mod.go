// SPDX-License-Identifier: GPL-3.0-or-later

package forge

import "net/http"

// IForgeHook is the interface that all forge hooks implement.
type IForgeHook interface {
	// HookRequest hooks an incoming webhook request to trigger actions.
	HookRequest(req *http.Request)
}

// IEvent is abstraction for forge-generated events.
type IEvent interface {
	// Raw returns the raw payload from forges.
	//
	// TODO: This is for debugging purposes, and is very likely to be removed
	// before initial release.
	Raw() interface{}
}
