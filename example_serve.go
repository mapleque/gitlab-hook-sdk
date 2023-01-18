package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mapleque/gitlab-hook-sdk/webhook"
)

func main() {
	var opts []webhook.Option
	secret := os.Getenv("GITLAB_WEBHOOK_SECRET")
	if secret != "" {
		opts = append(opts, webhook.Options.Secret(secret))
	}
	opts = append(
		opts,
		webhook.Handlers.Default(func(payload webhook.Payload) {
			log.Printf("recieve webhook hook payload: %v", payload)
		}),
		webhook.Handlers.IssueEvent(func(payload webhook.IssueEventPayload) {
			log.Printf("here is an issue event payload: %v", payload)
		}),
	)
	hook, _ := webhook.New(opts...)

	http.Handle("/mywebhook", hook)

	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatalf("webhook hook serve failed with error: %v", err)
	}
}
