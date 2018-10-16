package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gobuffalo/envy"
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
	Color    string        `json:"color"`
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
	State    string `json:"state"`
	ActionTs string `json:"action_ts"`
	Name     string `json:"name"`
	Value    string `json:"value"`
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
	Submission  struct {
		ProjectID   string `json:"projectID"`
		TaskID      string `json:"taskID"`
		Summary     string `json:"summary"`
		Description string `json:"description"`
	} `json:"submission"`
}

// SlackDialogResponse - response for dialog end-point
type SlackDialogResponse struct {
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
	Type           string                     `json:"type"`
	Label          string                     `json:"label"`
	Name           string                     `json:"name"`
	Placeholder    string                     `json:"placeholder"`
	Hint           string                     `json:"hint"`
	DataSource     string                     `json:"data_source"`
	MinQueryLength int                        `json:"min_query_length"`
	Options        []SlackDialogElementOption `json:"options"`
}

// SlackDialogElementOption - options for "select" input
type SlackDialogElementOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// StateData - message data for state attr
type StateData struct {
	Link    string `json:"link"`
	Message string `json:"message"`
}

// SlackDialogURL - url for dialogs in slack
const SlackDialogURL = "https://slack.com/api/dialog.open"

// SlackAccessToken - access token for slack app
var SlackAccessToken = envy.Get("SLACK_ACCESS_TOKEN", "")

// SlackDomain - slack domain
var SlackDomain = envy.Get("SLACK_DOMAIN", "")

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
	fmt.Println("Buffer for Response", string(responseBody))

	buffer.WriteString(string(responseBody))

	response, err := sendRequestToSlack("POST", SlackDialogURL, buffer)
	if err != nil {
		return err
	}

	fmt.Println("Dialog RESPPPPP", response)

	respBody, _ := ioutil.ReadAll(response.Body)
	fmt.Println("PARSED RESPPPPP", string(respBody))

	return nil
}

func sendRequestToSlack(method, url string, buffer *bytes.Buffer) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest(method, url, buffer)
	request.Header.Set("content-type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", SlackAccessToken))
	request.Header.Set("Accept", "application/json")

	fmt.Printf("REQUEST TO SLACK --- %v \n", request)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// MessageLink - link on message to Slack
func (callback *SlackActionCallback) MessageLink() string {
	message := StateData{
		Message: callback.Message.Text,
		Link: fmt.Sprintf("https://%s/archives/%s/p%s\n",
			SlackDomain,
			callback.Channel.ID,
			strings.Replace(callback.Message.Ts, ".", "", 1),
		),
	}

	encodedData, _ := json.Marshal(message)

	return string(encodedData)
}

// ParseState - parse state when callback returs
func (callback *SlackActionCallback) ParseState() StateData {
	message := StateData{}

	json.Unmarshal([]byte(callback.State), &message)

	return message
}
