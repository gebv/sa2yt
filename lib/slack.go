package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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

// SendAnswerToSlack - send answer to slack chat
func SendAnswerToSlack(url string, slackResponse *SlackResponse) error {
	client := &http.Client{}

	buffer := new(bytes.Buffer)
	responseBody, err := json.Marshal(slackResponse)
	if err != nil {
		return err
	}
	buffer.WriteString(string(responseBody))

	request, err := http.NewRequest("POST", url, buffer)
	request.Header.Set("content-type", "application/json")
	request.Header.Set("Accept", "application/json")

	fmt.Printf("REQUEST TO SLACK --- %v \n", request)

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	fmt.Println("RESPPPPP", response)

	return nil
}
