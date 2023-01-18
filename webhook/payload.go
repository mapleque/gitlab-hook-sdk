package webhook

import (
	"encoding/json"
	"fmt"
)

// Event gitlab event enum type
type Event string

const (
	TagEvent          Event = "Tag Push Hook"
	IssueEvent        Event = "Issue Hook"
	MergeRequestEvent Event = "Merge Request Hook"
	CommentEvent      Event = "Note Hook"
)

// Payload is an event payload interface
type Payload interface {
	// Event return the current payload event name
	Event() Event

	// Body return the payload content body as bytes
	Body() []byte
}

// AbstractEventPayload contains Gitlab's event raw message
type AbstractEventPayload struct {
	event Event
	body  []byte
}

// Event return the current payload event name
func (p AbstractEventPayload) Event() Event {
	return p.event
}

// Body return the payload content body as bytes
func (p AbstractEventPayload) Body() []byte {
	return p.body
}

// TagEventPayload contains the information for Gitlab's tag event
type TagEventPayload struct {
	AbstractEventPayload
}

// IssueEventPayload contains the information for Gitlab's issue event
type IssueEventPayload struct {
	AbstractEventPayload
}

// MergeRequestEventPayload contains the information
// for Gitlab's merge request event
type MergeRequestEventPayload struct {
	AbstractEventPayload
}

// CommentRequestPayload contains the information for Giltab's comment event
type CommentEventPayload struct {
	AbstractEventPayload
}

func eventParsing(event Event, payload []byte) (Payload, error) {
	abp := AbstractEventPayload{
		event: event,
		body:  payload,
	}
	switch event {
	case TagEvent:
		p := TagEventPayload{AbstractEventPayload: abp}
		if err := json.Unmarshal(payload, &p); err != nil {
			return abp, fmt.Errorf("parse event %s payload error: %w", event, err)
		}
		return p, nil
	case IssueEvent:
		p := IssueEventPayload{AbstractEventPayload: abp}
		if err := json.Unmarshal(payload, &p); err != nil {
			return abp, fmt.Errorf("parse event %s payload error: %w", event, err)
		}
		return p, nil
	case MergeRequestEvent:
		p := MergeRequestEventPayload{AbstractEventPayload: abp}
		if err := json.Unmarshal(payload, &p); err != nil {
			return abp, fmt.Errorf("parse event %s payload error: %w", event, err)
		}
		return p, nil
	case CommentEvent:
		p := CommentEventPayload{AbstractEventPayload: abp}
		if err := json.Unmarshal(payload, &p); err != nil {
			return abp, fmt.Errorf("parse event %s payload error: %w", event, err)
		}
		return p, nil
	default:
		return abp, nil
	}
}
