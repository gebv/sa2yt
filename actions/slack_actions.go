package actions

import (
	"fmt"

	"encoding/json"

	"github.com/gebv/sayto/lib"
	"github.com/gobuffalo/buffalo"
)

// SlackActionsCreate default implementation.
func SlackActionsCreate(c buffalo.Context) error {
	fmt.Printf("Form: %v \n", c.Request().Form)

	go func() {
		payload := c.Request().Form.Get("payload")
		var encodedCallback lib.SlackActionCallback

		err := json.Unmarshal([]byte(payload), &encodedCallback)
		if err != nil {
			fmt.Printf("ERROR: Can't encode slack message: %v \n", err)
			return
		}

		fmt.Println("responseURL:  ", encodedCallback.ResponseURL)
		lib.OpenDialogInSlack(
			&lib.SlackDialogResponse{
				TriggerID: encodedCallback.TriggerID,
				Dialog: lib.SlackDialog{
					CallbackID:  encodedCallback.CallbackID,
					State:       "Limo",
					Title:       "Request a Ride",
					SubmitLabel: "Request",
					//TODO: change dialog buttons https://api.slack.com/methods/dialog.open
					Elements: []lib.SlackDialogResponseElement{
						{
							Type:  "text",
							Label: "Pickup Location",
							Name:  "loc_origin",
						},
						{
							Type:  "text",
							Label: "Dropoff Location",
							Name:  "loc_destination",
						},
					},
				},
			},
		)

		// Simple Answer
		// lib.SendAnswerToSlack(encodedCallback.ResponseURL, &lib.SlackResponse{
		// 	Text: "Task was created",
		// 	Attachments: []lib.SlackAttachment{
		// 		{
		// 			Fallback: fmt.Sprintf("View Task In YouTrack %s.", "urlToTask"),
		// 			Actions: []lib.SlackAction{
		// 				{
		// 					Type: "button",
		// 					Text: "View Task In YouTrack",
		// 					URL:  "urlToTask",
		// 				},
		// 			},
		// 		},
		// 	},
		// })
	}()

	return c.Render(200, r.Plain(""))
}
