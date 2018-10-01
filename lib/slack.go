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

// SlackDialogResponse - response for dialog end-point
type SlackDialogResponse struct {
	Token     string      `json:"token"`
	TriggerID string      `json:"trigger_id"`
	Dialog    SlackDialog `json:"dialog"`
}

// SlackDialog - struct for dialog
type SlackDialog struct {
	CallbackID  string                       `json:"callback_id"`
	Title       string                       `json:"title"`
	SubmitLabel string                       `json:"submit_label"`
	State       string                       `json:"state"`
	Elements    []SlackDialogResponseElement `json:"elements"`
}

// SlackDialogResponseElement - element for dialog form
type SlackDialogResponseElement struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Name  string `json:"name"`
}

// SlackDialogURL - url for dialogs in slack
const SlackDialogURL = "https://slack.com/api/dialog.open"

// SendAnswerToSlack - send answer to slack chat
func SendAnswerToSlack(url string, slackResponse *SlackResponse) error {
	buffer := new(bytes.Buffer)
	responseBody, err := json.Marshal(slackResponse)
	if err != nil {
		return err
	}
	buffer.WriteString(string(responseBody))

	response, err := sendRequestToSlack("POST", url, buffer)
	if err != nil {
		return err
	}

	fmt.Println("RESPPPPP", response)

	return nil
}

// OpenDialogInSlack - Open dialog window in slack
func OpenDialogInSlack(dialog *SlackDialogResponse) error {
	buffer := new(bytes.Buffer)
	responseBody, err := json.Marshal(dialog)
	if err != nil {
		return err
	}
	buffer.WriteString(string(responseBody))

	response, err := sendRequestToSlack("POST", SlackDialogURL, buffer)
	if err != nil {
		return err
	}

	fmt.Println("Dialog RESPPPPP", response)

	var respBody []byte
	response.Body.Read(respBody)

	fmt.Println("PARSED RESPPPPP", string(respBody))

	return nil
}

func sendRequestToSlack(method, url string, buffer *bytes.Buffer) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest(method, url, buffer)
	request.Header.Set("content-type", "application/json")
	request.Header.Set("Accept", "application/json")

	fmt.Printf("REQUEST TO SLACK --- %v \n", request)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
