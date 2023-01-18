package webhook_test

import (
	"log"
	"net/http"

	"github.com/mapleque/gitlab-hook-sdk/webhook"
)

func Example_Option() {
	secret := "my-gitlab-hook-secret"
	opts := []webhook.Option{
		webhook.Options.Secret(secret),
		webhook.Handlers.Default(func(payload webhook.Payload) {
			log.Printf("recieve webhook hook payload: %v", payload)
		}),
		webhook.Handlers.IssueEvent(func(payload webhook.IssueEventPayload) {
			log.Printf("here is an issue event payload: %v", payload)
		}),
	}
	hook, _ := webhook.New(opts...)

	http.Handle("/mywebhook", hook)
}
