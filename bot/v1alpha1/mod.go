// SPDX-License-Identifier: GPL-3.0-or-later

package v1alpha1

const PluginAPIVersion = 1

type EventType int

// All event types.
const (
	EventTypeUnknown          EventType = 0
	EventTypeWebhookInstalled EventType = 1
	EventTypePROpened         EventType = 2
	EventTypePRClosed         EventType = 3
	EventTypePRMerged         EventType = 4
	EventTypePRRenamed        EventType = 5
	EventTypePRReviewed       EventType = 6
	EventTypePRReady          EventType = 7
	EventTypePRWithdrawn      EventType = 8
	EventTypeCIFinished       EventType = 9
	EventTypeReviewPing       EventType = 10
)

type IssueState int

// All issue states.
const (
	IssueStateUnknown IssueState = 0
	IssueStateOpen    IssueState = 1
	IssueStateClosed  IssueState = 2
	IssueStateMerged  IssueState = 3
)

type CIState int

// All CI states.
const (
	CIStateUnknown CIState = 0
	CIStatePassed  CIState = 1
	CIStateFailed  CIState = 2
	CIStateErrored CIState = 3
)

type ReviewType int

// All review types.
const (
	ReviewTypeUnknown        ReviewType = 0
	ReviewTypeComment        ReviewType = 1
	ReviewTypeApprove        ReviewType = 2
	ReviewTypeRequestChanges ReviewType = 3
)

type ForgeUser struct {
	Forge    string
	UserName string
}

type Repo struct {
	User     ForgeUser
	RepoName string
}

type PR struct {
	Repo   Repo
	Number int
	Title  string
	Author ForgeUser
	State  IssueState
}

type CIRun struct {
	Repo  Repo
	PR    PR
	State CIState
}

type WebhookInstalledParams struct {
	Repo Repo
}

func (x *WebhookInstalledParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type PROpenedParams struct {
	Actor ForgeUser
	PR    PR
}

func (x *PROpenedParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type PRClosedParams struct {
	Actor ForgeUser
	PR    PR
}

func (x *PRClosedParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type PRMergedParams struct {
	Actor ForgeUser
	PR    PR
}

func (x *PRMergedParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type PRRenamedParams struct {
	Actor ForgeUser
	PR    PR
}

func (x *PRRenamedParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type PRReviewedParams struct {
	Actor  ForgeUser
	PR     PR
	Review ReviewType
}

func (x *PRReviewedParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type PRReadyParams struct {
	Actor ForgeUser
	PR    PR
}

func (x *PRReadyParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type PRWithdrawnParams struct {
	PR PR
}

func (x *PRWithdrawnParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type CIFinishedParams struct {
	Run CIRun
}

func (x *CIFinishedParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type ReviewPingParams struct {
	Actor ForgeUser
	PR    PR
}

func (x *ReviewPingParams) IntoEvent() *Event {
	return &Event{
		inner: x,
	}
}

type Event struct {
	inner interface{}
}

func (e *Event) Type() EventType {
	switch e.inner.(type) {
	case *WebhookInstalledParams:
		return EventTypeWebhookInstalled
	case *PROpenedParams:
		return EventTypePROpened
	case *PRClosedParams:
		return EventTypePRClosed
	case *PRMergedParams:
		return EventTypePRMerged
	case *PRRenamedParams:
		return EventTypePRRenamed
	case *PRReviewedParams:
		return EventTypePRReviewed
	case *PRReadyParams:
		return EventTypePRReady
	case *PRWithdrawnParams:
		return EventTypePRWithdrawn
	case *CIFinishedParams:
		return EventTypeCIFinished
	case *ReviewPingParams:
		return EventTypeReviewPing
	default:
		return EventTypeUnknown
	}
}

func (e *Event) WebhookInstalled() (*WebhookInstalledParams, bool) {
	params, ok := e.inner.(*WebhookInstalledParams)
	return params, ok
}

func (e *Event) PROpened() (*PROpenedParams, bool) {
	params, ok := e.inner.(*PROpenedParams)
	return params, ok
}

func (e *Event) PRClosed() (*PRClosedParams, bool) {
	params, ok := e.inner.(*PRClosedParams)
	return params, ok
}

func (e *Event) PRMerged() (*PRMergedParams, bool) {
	params, ok := e.inner.(*PRMergedParams)
	return params, ok
}

func (e *Event) PRRenamed() (*PRRenamedParams, bool) {
	params, ok := e.inner.(*PRRenamedParams)
	return params, ok
}

func (e *Event) PRReviewed() (*PRReviewedParams, bool) {
	params, ok := e.inner.(*PRReviewedParams)
	return params, ok
}

func (e *Event) PRReady() (*PRReadyParams, bool) {
	params, ok := e.inner.(*PRReadyParams)
	return params, ok
}

func (e *Event) PRWithdrawn() (*PRWithdrawnParams, bool) {
	params, ok := e.inner.(*PRWithdrawnParams)
	return params, ok
}

func (e *Event) CIFinished() (*CIFinishedParams, bool) {
	params, ok := e.inner.(*CIFinishedParams)
	return params, ok
}

func (e *Event) ReviewPing() (*ReviewPingParams, bool) {
	params, ok := e.inner.(*ReviewPingParams)
	return params, ok
}

type IIMProvider interface {
	SendTextToPerson(userID string, text string) error
	SendTextToChat(chatID string, text string) error
	SendMarkdownToPerson(userID string, md string) error
	SendMarkdownToChat(chatID string, md string) error
}

// IPlugin is the interface all plugins must implement.
type IPlugin interface {
	Setup() error
	ProcessEvent(e *Event, im IIMProvider) error
	Teardown() error
}

// IPluginConfigFactoryFunc is signature for the plugin's exported
// "BrickbotPluginConfigFactory" function.
//
// You should return the zero value of your desired config struct.
type IPluginConfigFactoryFunc = func() interface{}

// IPluginFactoryFunc is signature for the plugin's exported
// "BrickbotPluginFactory" function.
type IPluginFactoryFunc = func(config interface{}) (IPlugin, error)
