package webhook

// Option is a configuration option
type Option func(*Webhook) error

// Options is a var of WebhookOptions
var Options = WebhookOptions{}

// WebhookOptions is a helper for creating configration options
type WebhookOptions struct{}

// Secret setting secret for hook
func (WebhookOptions) Secret(secret string) Option {
	return func(hook *Webhook) error {
		hook.secret = secret
		return nil
	}
}
