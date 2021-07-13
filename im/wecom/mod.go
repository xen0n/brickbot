// SPDX-License-Identifier: GPL-3.0-or-later

package wecom

import (
	"errors"

	"github.com/xen0n/go-workwx"

	"github.com/xen0n/brickbot/im"
)

type wecomProvider struct {
	app *workwx.WorkwxApp
}

var _ im.IProvider = (*wecomProvider)(nil)

// New returns a new 企业微信 (WeCom) provider instance.
func New(
	corpID string,
	corpSecret string,
	agentID int64,
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

	cl := workwx.New(corpID)
	app := cl.WithApp(corpSecret, agentID)

	return &wecomProvider{
		app: app,
	}, nil
}

func (p *wecomProvider) SendTextToPerson(userID string, text string) error {
	rcpt := workwx.Recipient{
		UserIDs: []string{userID},
	}

	err := p.app.SendTextMessage(&rcpt, text, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *wecomProvider) SendTextToChat(chatID string, text string) error {
	rcpt := workwx.Recipient{
		ChatID: chatID,
	}

	err := p.app.SendTextMessage(&rcpt, text, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *wecomProvider) SendMarkdownToPerson(userID string, md string) error {
	rcpt := workwx.Recipient{
		UserIDs: []string{userID},
	}

	err := p.app.SendMarkdownMessage(&rcpt, md, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *wecomProvider) SendMarkdownToChat(chatID string, md string) error {
	rcpt := workwx.Recipient{
		ChatID: chatID,
	}

	err := p.app.SendMarkdownMessage(&rcpt, md, false)
	if err != nil {
		return err
	}

	return nil
}
