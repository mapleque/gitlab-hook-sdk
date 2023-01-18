package webhook

import "log"

// WebhookHandler is the handler func using in hook
type WebhookHandler func(Payload)

// Handlers is a var of HandlerOptions
var Handlers = HandlerOptions{}

// HandlerOptions is a helper for register hook handler as options
type HandlerOptions struct{}

// Default set the default handler for all hook events
func (HandlerOptions) Default(h WebhookHandler) Option {
	return func(hook *Webhook) error {
		hook.defaultHandler = h
		return nil
	}
}

func defaultHandler(payload Payload) {
	log.Printf("recieve webhook hook payload: %v", payload)
}

// IssueEvent register issueEventHandler
func (HandlerOptions) IssueEvent(h func(IssueEventPayload)) Option {
	return func(hook *Webhook) error {
		hook.issueEventHandler = func(p Payload) {
			h(p.(IssueEventPayload))
		}
		return nil
	}
}

// TagEvent register tagEventHandler
func (HandlerOptions) TagEvent(h func(TagEventPayload)) Option {
	return func(hook *Webhook) error {
		hook.tagEventHandler = func(p Payload) {
			h(p.(TagEventPayload))
		}
		return nil
	}
}

// MergeRequestEvent register mergeRequestEventHandler
func (HandlerOptions) MergeRequestEvent(
	h func(MergeRequestEventPayload),
) Option {
	return func(hook *Webhook) error {
		hook.mergeRequestEventHandler = func(p Payload) {
			h(p.(MergeRequestEventPayload))
		}
		return nil
	}
}

// CommentEvent register commentEventHandler
func (HandlerOptions) CommentEvent(h func(CommentEventPayload)) Option {
	return func(hook *Webhook) error {
		hook.commentEventHandler = func(p Payload) {
			h(p.(CommentEventPayload))
		}
		return nil
	}
}
