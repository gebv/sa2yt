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

		switch encodedCallback.CallbackID {
		case "new_task":
			sendDialogWindow(&encodedCallback)
		case "create_task":
			createIssueAndSendAnswer(&encodedCallback)
		}
	}()

	return c.Render(200, r.Plain(""))
}

func sendDialogWindow(encodedCallback *lib.SlackActionCallback) {
	var projectOptions []lib.SlackDialogElementOption
	for _, project := range YouTrackAPI.CachedProjects {
		projectOptions = append(projectOptions, lib.SlackDialogElementOption{
			Label: project.ID,
			Value: project.ID,
		})
	}

	originalMessage := fmt.Sprintf("[%s](дата сообщения из слак)", encodedCallback.Message.Ts)

	lib.OpenDialogInSlack(
		&lib.SlackDialogResponse{
			TriggerID: encodedCallback.TriggerID,
			Dialog: lib.SlackDialog{
				CallbackID:  "create_task",
				State:       originalMessage,
				Title:       "Create new Task",
				SubmitLabel: "Request",
				Elements: []lib.SlackDialogResponseElement{
					{
						Type:    "select",
						Label:   "Project ID",
						Name:    "projectID",
						Options: projectOptions,
					},
					{
						Type:        "text",
						Label:       "Summary",
						Name:        "summary",
						Placeholder: "Task Summary",
					},
					{
						Type:  "textarea",
						Label: "Description",
						Name:  "description",
						Hint:  "Explaint your task",
					},
				},
			},
		},
	)
}

func createIssueAndSendAnswer(encodedCallback *lib.SlackActionCallback) {
	urlToTask, err := YouTrackAPI.CreateIssue(
		encodedCallback.Submission.ProjectID,
		encodedCallback.Submission.Summary,
		encodedCallback.Submission.Description+"\n"+encodedCallback.State)
	if err != nil {
		lib.SendAnswerToSlack(encodedCallback.ResponseURL, &lib.SlackResponse{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Error create issue in YouTrack: %v", err),
		})
	}

	lib.SendAnswerToSlack(encodedCallback.ResponseURL, &lib.SlackResponse{
		Text: "Task was created",
		Attachments: []lib.SlackAttachment{
			{
				Fallback: fmt.Sprintf("View Task In YouTrack %s.", urlToTask),
				Actions: []lib.SlackAction{
					{
						Type: "button",
						Text: "View Task In YouTrack",
						URL:  urlToTask,
					},
				},
			},
		},
	})
}

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
