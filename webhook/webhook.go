package webhook

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	ErrInvalidHTTPMethod        = errors.New("invalid HTTP Method")
	ErrInvalidGitlabToken       = errors.New("X-Gitlab-Token validation failed")
	ErrMissingGitlabEventHeader = errors.New("missing X-Gitlab-Event Header")
	ErrInvalidPayload           = errors.New("parsing payload error")
)

// Webhook ...
type Webhook struct {
	secret string

	defaultHandler WebhookHandler

	tagEventHandler          WebhookHandler
	issueEventHandler        WebhookHandler
	mergeRequestEventHandler WebhookHandler
	commentEventHandler      WebhookHandler
}

// New creates Webhook instance
func New(opts ...Option) (*Webhook, error) {
	hook := &Webhook{
		secret:                   "",
		defaultHandler:           nil,
		tagEventHandler:          nil,
		issueEventHandler:        nil,
		mergeRequestEventHandler: nil,
		commentEventHandler:      nil,
	}
	for _, opt := range opts {
		if err := opt(hook); err != nil {
			return nil, fmt.Errorf("apply Option error: %v", err)
		}
	}
	// deal empty handler
	if hook.defaultHandler == nil {
		hook.defaultHandler = defaultHandler
	}

	if hook.tagEventHandler == nil {
		hook.tagEventHandler = hook.defaultHandler
	}

	if hook.issueEventHandler == nil {
		hook.issueEventHandler = hook.defaultHandler
	}

	if hook.mergeRequestEventHandler == nil {
		hook.mergeRequestEventHandler = hook.defaultHandler
	}

	if hook.commentEventHandler == nil {
		hook.commentEventHandler = hook.defaultHandler
	}
	return hook, nil
}

// Serve implement http.Handler
func (hook *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := hook.parseEvent(r)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	switch p := payload.(type) {
	case TagEventPayload:
		hook.tagEventHandler(p)
	case IssueEventPayload:
		hook.issueEventHandler(p)
	case MergeRequestEventPayload:
		hook.mergeRequestEventHandler(p)
	case CommentEventPayload:
		hook.commentEventHandler(p)
	default:
		err := fmt.Errorf("unsupported event: %v", payload.Event())
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}

func (hook *Webhook) parseEvent(r *http.Request) (Payload, error) {
	defer func() {
		_, _ = io.Copy(ioutil.Discard, r.Body)
		_ = r.Body.Close()
	}()

	if r.Method != http.MethodPost {
		return nil, ErrInvalidHTTPMethod
	}

	if len(hook.secret) > 0 {
		token := r.Header.Get("X-Gitlab-Token")
		if token != hook.secret {
			return nil, ErrInvalidGitlabToken
		}
	}

	event := r.Header.Get("X-Gitlab-Event")
	if len(event) == 0 {
		return nil, ErrMissingGitlabEventHeader
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return nil, ErrInvalidPayload
	}

	return eventParsing(Event(event), payload)
}
