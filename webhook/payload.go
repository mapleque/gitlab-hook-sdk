package webhook

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Event Gitlab event enum type
type Event string

const (
	IssueEvent        Event = "Issue Hook"
	MergeRequestEvent Event = "Merge Request Hook"
	CommentEvent      Event = "Note Hook"
)

type NoteableType string

const (
	NoteableTypeCommit       NoteableType = "Commit"
	NoteableTypeIssue        NoteableType = "Issue"
	NoteableTypeMergeRequest NoteableType = "MergeRequest"
	NoteableTypeSnippet      NoteableType = "Snippet"
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

	ObjectKind string  `json:"object_kind"`
	EventType  string  `json:"event_type"`
	User       User    `json:"user"`
	Project    Project `json:"project"`
}

// Event return the current payload event name
func (p AbstractEventPayload) Event() Event {
	return p.event
}

// Body return the payload content body as bytes
func (p AbstractEventPayload) Body() []byte {
	return p.body
}

// IssueEventPayload contains the information for Gitlab's issue event
type IssueEventPayload struct {
	AbstractEventPayload

	ObjectAttributes IssueObjectAttributes `json:"object_attributes"`
	Labels           []Label               `json:"labels"`
	Changes          Changes               `json:"changes"`
	Assignees        []Assignee            `json:"assignees"`
}

// MergeRequestEventPayload contains the information
// for Gitlab's merge request event
type MergeRequestEventPayload struct {
	AbstractEventPayload

	ObjectAttributes MergeRequestObjectAttributes `json:"object_attributes"`
	Labels           []Label                      `json:"labels"`
	Changes          Changes                      `json:"changes"`
	Repository       Repository                   `json:"repository"`
	Assignees        []Assignee                   `json:"assignees"`
}

// CommentRequestPayload contains the information for Giltab's comment event
type CommentEventPayload struct {
	AbstractEventPayload

	ObjectAttributes CommentObjectAttributes      `json:"object_attributes"`
	Repository       Repository                   `json:"repository"`
	Issue            IssueObjectAttributes        `json:"issue"`
	Commit           Commit                       `json:"commit"`
	MergeRequest     MergeRequestObjectAttributes `json:"merge_request"`
}

func eventParsing(event Event, payload []byte) (Payload, error) {
	abp := AbstractEventPayload{
		event: event,
		body:  payload,
	}
	switch event {
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

// User contains Gitlab event user information
type User struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

// Project contains Gitlab event project information
type Project struct {
	Id                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

// IssueObjectAttributes contains Gitlab event issue object attributes
// information
type IssueObjectAttributes struct {
	AuthorId         int64      `json:"author_id"`
	ClosedAt         customTime `json:"closed_at"`
	Confidential     bool       `json:"confidential"`
	CreatedAt        customTime `json:"created_at"`
	Description      string     `json:"description"`
	DiscussionLocked bool       `json:"discussion_locked"`
	DueDate          customTime `json:"due_date"`
	Id               int64      `json:"id"`
	IId              int64      `json:"iid"`
	LastEditedAt     customTime `json:"last_edited_at"`
	LastEditedById   int64      `json:"last_edited_by_id"`
	MilestoneId      int64      `json:"milestone_id"`
	Title            string     `json:"title"`
	UpdatedAt        customTime `json:"updated_at"`
	UpdatedById      int64      `json:"updated_by_id"`
	URL              string     `json:"url"`

	AssigneeIds []int64 `json:"assignee_ids"`
	Labels      []Label `json:"labels"`
	State       string  `json:"state"`
	Action      string  `json:"action"`
}

// MergeRequestObjectAttributes contains Gitlab event merge request object
// attributes information
type MergeRequestObjectAttributes struct {
	Id              int64      `json:"id"`
	TargetBranch    string     `json:"target_branch"`
	SourceBranch    string     `json:"source_branch"`
	AuthorId        int64      `json:"author_id"`
	AssigneeId      int64      `json:"assignee_id"`
	Title           string     `json:"title"`
	CreatedAt       customTime `json:"created_at"`
	UpdatedAt       customTime `json:"updated_at"`
	MilestoneId     int64      `json:"milestone_id"`
	State           string     `json:"state"`
	MergeStatus     string     `json:"unckecked"`
	TargetProjectId int64      `json:"target_project_id"`
	IId             int64      `json:"iid"`
	Description     string     `json:"description"`
	Source          Project    `json:"source"`
	Target          Project    `json:"target"`
	LastCommit      Commit     `json:"last_commit"`
	WorkInProgress  bool       `json:"work_in_progress"`
	URL             string     `json:"url"`
	Action          string     `json:"action"`
	AssigneeIds     []string   `json:"assignee_ids"`
}

// CommentObjectAttributes contains Gitlab event note object
// attributes information
type CommentObjectAttributes struct {
	AuthorId int64 `json:"author_id"`

	CreatedAt    customTime   `json:"created_at"`
	DiscussionId string       `json:"discussion_id"`
	Id           int64        `json:"id"`
	Note         string       `json:"note"`
	NoteableId   int64        `json:"noteable_id"`
	NoteableType NoteableType `json:"noteable_type"`
	ProjectId    int64        `json:"project_id"`
	UpdatedAt    customTime   `json:"updated_at"`
	UpdatedById  int64        `json:"updated_by_id"`
	URL          string       `json:"url"`
}

// Label contains Gitlab event label information
type Label struct {
	Id            int64      `json:"id"`
	Title         string     `json:"title"`
	Color         string     `json:"color"`
	ProjectId     int64      `json:"project_id"`
	CreatedAt     customTime `json:"created_at"`
	UpdatedAt     customTime `json:"updated_at"`
	Template      bool       `json:"template"`
	Description   string     `json:"description"`
	Type          string     `json:"type"`
	GroupId       int64      `json:"group_id"`
	RemoveOnClose bool       `json:"remove_on_close"`
}

// Changes contains Gitlab event changes information
type Changes struct {
	AuthorId    int64Change      `json:"author_id"`
	CreatedAt   customTimeChange `json:"created_at"`
	Description stringChange     `json:"description"`
	Id          int64Change      `json:"id"`
	IId         int64Change      `json:"iid"`
	ProjectId   int64Change      `json:"project_id"`
	Title       stringChange     `json:"title"`
	UpdatedAt   customTimeChange `json:"updated_at"`
}

// Assignee contains Gitlab event assignee information
type Assignee struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

// Repository contains Gitlab event repository information
type Repository struct {
	Name        string       `json:"name"`
	URL         string       `json:"url"`
	Description stringChange `json:"description"`
	Homepage    string       `json:"homepage"`
}

// Commit contains Gitlab event commit information
type Commit struct {
	Id        int64      `json:"id"`
	Message   string     `json:"message"`
	Title     string     `json:"title"`
	Timestamp customTime `json:"timestamp"`
	URL       string     `json:"url"`
	Author    Author     `json:"author"`
	Added     []string   `json:"added"`
	Modified  []string   `json:"modified"`
	Removed   []string   `json:"removed"`
}

// Author contains Gitlab event author information
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type int64Change struct {
	Previous int64 `json:"previous"`
	Current  int64 `json:"current"`
}

type stringChange struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type customTimeChange struct {
	Previous customTime `json:"previous"`
	Current  customTime `json:"current"`
}

type customTime struct {
	time.Time
}

func (t *customTime) UnmarshalJSON(b []byte) (err error) {
	layout := []string{
		"2006-01-02 15:04:05 MST",
		"2006-01-02 15:04:05 Z07:00",
		"2006-01-02 15:04:05 Z0700",
		time.RFC3339,
	}
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}
	for _, l := range layout {
		t.Time, err = time.Parse(l, s)
		if err == nil {
			return
		}
	}
	return
}
