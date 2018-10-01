package actions

import (
	"fmt"

	"github.com/gebv/sayto/lib"
	"github.com/gobuffalo/buffalo"
)

// SlackActionsCreate default implementation.
func SlackActionsCreate(c buffalo.Context) error {
	fmt.Printf("Form: %v \n", c.Request().Form)

	go func() {
		responseURL := c.Request().Form.Get("payload")
		fmt.Println("responseURL:  ", responseURL)

		lib.SendAnswerToSlack(responseURL, &lib.SlackResponse{
			Text: "Task was created",
			Attachments: []lib.SlackAttachment{
				{
					Fallback: fmt.Sprintf("View Task In YouTrack %s.", "urlToTask"),
					Actions: []lib.SlackAction{
						{
							Type: "button",
							Text: "View Task In YouTrack",
							URL:  "urlToTask",
						},
					},
				},
			},
		})

		// map[payload:[{
		// 		"type":"message_action",
		// 		"token":"YmLuJsrTXziYJyioh17Kl7O7",
		// 		"action_ts":"1538416708.046393",
		// 		"team":{"id":"TCTMD8KE2","domain":"my-test-companytalk"},
		// 		"user":{"id":"UCTMD8L30","name":"gavrilov.ea"},
		// 		"channel":{"id":"CCUE4GAAC","name":"new-test-app"},
		// 		"callback_id":"create_task",
		// 		"trigger_id":"446610868420.435727291478.fd381ff42baaa14d5138081e0deabc8e",
		// 		"message_ts":"1537472211.000200",
		// 		"message":{"type":"message","user":"UCTMD8L30","text":"hi","client_msg_id":"7ca3327e-fde4-46de-8302-64d405ddcb39","ts":"1537472211.000200"},
		// 		"response_url":"https:\/\/hooks.slack.com\/app\/TCTMD8KE2\/446693190419\/i3Nr0YMuy90XC8LucmUxQjaz"}]]
	}()

	return c.Render(200, r.Plain(""))
}
