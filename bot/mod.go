// SPDX-License-Identifier: GPL-3.0-or-later

package bot

import (
	"errors"
	"plugin"

	"github.com/rs/zerolog/log"
	"github.com/xen0n/brickbot/bot/v1alpha1"
	"github.com/xen0n/brickbot/im"
)

// IBot is the interface that all bots implement.
type IBot interface {
	ProcessEvent(e *v1alpha1.Event, im im.IProvider) error
}

func InitWithPlugin(pluginPath string) (IBot, error) {
	pl, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	apiVersionSym, err := pl.Lookup("PluginAPIVersion")
	if err != nil {
		return nil, err
	}

	apiVersion, ok := apiVersionSym.(int)
	if !ok {
		return nil, err
	}

	if apiVersion != v1alpha1.PluginAPIVersion {
		return nil, errors.New("plugin API version mismatch")
	}

	processEventFnSym, err := pl.Lookup("ProcessEvent")
	if err != nil {
		return nil, err
	}

	processEventFn, ok := processEventFnSym.(v1alpha1.IProcessEventFunc)
	if !ok {
		return nil, errors.New("wrong type of ProcessEvent symbol")
	}

	return &pluginHostV1alpha1{
		processEventFn: processEventFn,
	}, nil
}

type pluginHostV1alpha1 struct {
	processEventFn v1alpha1.IProcessEventFunc
}

var _ IBot = (*pluginHostV1alpha1)(nil)

func (h *pluginHostV1alpha1) ProcessEvent(
	e *v1alpha1.Event,
	im im.IProvider,
) error {
	return h.processEventFn(e, makeIMProviderShim(im))
}

type imProviderShim struct {
	inner im.IProvider
}

var _ v1alpha1.IIMProvider = (*imProviderShim)(nil)

func makeIMProviderShim(p im.IProvider) *imProviderShim {
	return &imProviderShim{
		inner: p,
	}
}

func (p *imProviderShim) SendTextToPerson(userID string, text string) error {
	log.Debug().
		Str("userID", userID).
		Str("text", text).
		Msg("SendTextToPerson stub")

	return nil
}

func (p *imProviderShim) SendTextToChat(chatID string, text string) error {
	log.Debug().
		Str("chatID", chatID).
		Str("text", text).
		Msg("SendTextToChat stub")

	return nil
}

func (p *imProviderShim) SendMarkdownToPerson(userID string, md string) error {
	log.Debug().
		Str("userID", userID).
		Str("md", md).
		Msg("SendMarkdownToPerson stub")

	return nil
}

func (p *imProviderShim) SendMarkdownToChat(chatID string, md string) error {
	log.Debug().
		Str("chatID", chatID).
		Str("md", md).
		Msg("SendMarkdownToChat stub")

	return nil
}
