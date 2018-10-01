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

// SlackActionCallback - callback message after action from slack
type SlackActionCallback struct {
	Type     string `json:"type"`
	Token    string `json:"token"`
	ActionTs string `json:"action_ts"`
	Team     struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	CallbackID string `json:"callback_id"`
	TriggerID  string `json:"trigger_id"`
	MessageTs  string `json:"message_ts"`
	Message    struct {
		Type        string `json:"type"`
		User        string `json:"user"`
		Text        string `json:"text"`
		ClientMsgID string `json:"client_msg_id"`
		Ts          string `json:"ts"`
	} `json:"message"`
	ResponseURL string `json:"response_url"`
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
