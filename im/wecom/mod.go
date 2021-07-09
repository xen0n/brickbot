// SPDX-License-Identifier: GPL-3.0-or-later

package wecom

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/xen0n/go-workwx"

	"github.com/xen0n/brickbot/im"
)

type wecomProvider struct {
	app *workwx.WorkwxApp

	teamChatID string
}

var _ im.IProvider = (*wecomProvider)(nil)

// New returns a new 企业微信 (WeCom) provider instance.
func New(
	corpID string,
	corpSecret string,
	agentID int64,
	teamChatID string,
) (im.IProvider, error) {
	if corpID == "" {
		return nil, errors.New("empty CorpID")
	}
	if corpSecret == "" {
		return nil, errors.New("empty CorpSecret")
	}
	if agentID == 0 {
		return nil, errors.New("empty AgentID")
	}
	if teamChatID == "" {
		return nil, errors.New("empty team ChatID")
	}

	cl := workwx.New(corpID)
	app := cl.WithApp(corpSecret, agentID)

	return &wecomProvider{
		app:        app,
		teamChatID: teamChatID,
	}, nil
}

// SendTeamMessage sends a message to team scope.
func (p *wecomProvider) SendTeamMessage(m im.IOutgoingMessage) error {
	raw := m.Raw()

	// DEBUG
	ty := reflect.TypeOf(raw)
	content := fmt.Sprintf("event %s\nrepr %+v", ty.Name(), raw)

	rcpt := workwx.Recipient{
		ChatID: p.teamChatID,
	}

	err := p.app.SendTextMessage(&rcpt, content, false)
	if err != nil {
		return err
	}

	return nil
}
