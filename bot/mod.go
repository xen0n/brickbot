// SPDX-License-Identifier: GPL-3.0-or-later

package bot

import "github.com/xen0n/brickbot/forge"

// IBot is the interface that all bots implement.
type IBot interface {
	// ConsumeForgeEvent consumes a forge event.
	ConsumeForgeEvent(e forge.IEvent)
}
