// SPDX-License-Identifier: GPL-3.0-or-later

package forge

import (
	"net/http"

	"github.com/xen0n/brickbot/bot/v1alpha1"
)

// IForgeHook is the interface that all forge hooks implement.
type IForgeHook interface {
	// HookRequest hooks an incoming webhook request to trigger actions.
	HookRequest(req *http.Request) (*v1alpha1.Event, error)
}
