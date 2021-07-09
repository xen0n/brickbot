// SPDX-License-Identifier: GPL-3.0-or-later

package im

import "github.com/xen0n/brickbot/forge"

// IOutgoingMessage is abstraction for messages to be sent to IM providers.
type IOutgoingMessage = forge.IEvent

// IProvider is abstraction for IM backends.
type IProvider interface {
	// SendTeamMessage sends a message to team scope.
	SendTeamMessage(m IOutgoingMessage) error
}
