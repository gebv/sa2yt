package lib

// SlackResponse - struct for response to Slack API
type SlackResponse struct {
	ResponseType string            `json:"response_type,omitempty"`
	Text         string            `json:"text"`
	Attachments  []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment - struct for attachments in slack message
type SlackAttachment struct {
	Text     string        `json:"text"`
	Fallback string        `json:"fallback"`
	Actions  []SlackAction `json:"actions"`
}

// SlackAction - interactive buttons for message
type SlackAction struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	URL   string `json:"url"`
	Style string `json:"style"`
}
