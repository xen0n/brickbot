// SPDX-License-Identifier: GPL-3.0-or-later

package forge

import "net/http"

// IForgeHook is the interface that all forge hooks implement.
type IForgeHook interface {
	// HookRequest hooks an incoming webhook request to trigger actions.
	HookRequest(req *http.Request)
}
