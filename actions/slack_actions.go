package actions

import (
	"fmt"
	"net/url"
	"path"

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
			sendNewTaskWindow(&encodedCallback)
		case "create_task":
			createIssueAndSendAnswer(&encodedCallback)
		case "new_comment":
			sendNewCommentWindow(&encodedCallback)
		case "create_comment":
			createNewCommentAndSendAnswer(&encodedCallback)
		}
	}()

	return c.Render(200, r.Plain(""))
}

func sendNewTaskWindow(encodedCallback *lib.SlackActionCallback) {
	var projectOptions []lib.SlackDialogElementOption
	for _, project := range YouTrackAPI.CachedProjects {
		projectOptions = append(projectOptions, lib.SlackDialogElementOption{
			Label: project.ID,
			Value: project.ID,
		})
	}

	lib.OpenDialogInSlack(
		&lib.SlackDialogResponse{
			TriggerID: encodedCallback.TriggerID,
			Dialog: lib.SlackDialog{
				CallbackID:  "create_task",
				State:       encodedCallback.MessageLink(),
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

func sendNewCommentWindow(encodedCallback *lib.SlackActionCallback) {
	lib.OpenDialogInSlack(
		&lib.SlackDialogResponse{
			TriggerID: encodedCallback.TriggerID,
			Dialog: lib.SlackDialog{
				CallbackID:  "create_comment",
				State:       encodedCallback.MessageLink(),
				Title:       "Add new comment to Task",
				SubmitLabel: "Add",
				Elements: []lib.SlackDialogResponseElement{
					{
						Type:           "select",
						Label:          "Task ID",
						Name:           "taskID",
						DataSource:     "external",
						MinQueryLength: 2,
					},
					{
						Type:        "text",
						Label:       "Comment",
						Name:        "summary",
						Optional:    true,
						Placeholder: "Body Comment",
					},
				},
			},
		},
	)
}

func createIssueAndSendAnswer(encodedCallback *lib.SlackActionCallback) {
	parsedState := encodedCallback.ParseState()
	urlToTask, err := YouTrackAPI.CreateIssue(
		encodedCallback.Submission.ProjectID,
		encodedCallback.Submission.Summary,
		encodedCallback.Submission.Description+parsedState.FormattedLink())
	if err != nil {
		lib.SendAnswerToSlack(encodedCallback.ResponseURL, &lib.SlackResponse{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Error create issue in YouTrack: %v", err),
		})
		return
	}

	lib.SendAnswerToSlack(encodedCallback.ResponseURL, &lib.SlackResponse{
		Text: "Task was created",
		Attachments: []lib.SlackAttachment{
			{
				Fallback: fmt.Sprintf("View Task In YouTrack %s.", urlToTask),
				Color:    "good",
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

func createNewCommentAndSendAnswer(encodedCallback *lib.SlackActionCallback) {
	parsedState := encodedCallback.ParseState()
	err := YouTrackAPI.CreateComment(encodedCallback.Submission.TaskID,
		encodedCallback.Submission.Summary+parsedState.FormattedLink())

	if err != nil {
		lib.SendAnswerToSlack(encodedCallback.ResponseURL, &lib.SlackResponse{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Error create comment in YouTrack: %v", err),
		})
		return
	}

	pathToTask, _ := url.Parse(path.Join("/youtrack/issue/", encodedCallback.Submission.TaskID))
	domainURL, _ := url.Parse(YouTrackAPI.Domain)
	urlToTask := domainURL.ResolveReference(pathToTask)

	lib.SendAnswerToSlack(encodedCallback.ResponseURL, &lib.SlackResponse{
		ResponseType: "ephemeral",
		Text:         "Comment was added",
		Attachments: []lib.SlackAttachment{
			{
				Fallback: fmt.Sprintf("View Task In YouTrack %s.", urlToTask),
				Color:    "good",
				Actions: []lib.SlackAction{
					{
						Type: "button",
						Text: "View Task In YouTrack",
						URL:  urlToTask.String(),
					},
				},
			},
		},
	})
}
