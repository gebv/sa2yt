package actions

import (
	"fmt"
	"strings"

	"github.com/gebv/sayto/lib"

	"github.com/gobuffalo/buffalo"
)

// TasksCreate create new task in YouTrack.
func TasksCreate(c buffalo.Context) error {
	fmt.Printf("Form: %v \n", c.Request().Form)

	splitFunc := func(c rune) bool { return c == '#' }
	texts := strings.FieldsFunc(c.Request().Form.Get("text"), splitFunc)

	fmt.Printf("TEXTS:   %v, %d\n", texts, len(texts))
	if len(texts) < 2 {
		return c.Render(200, r.JSON(lib.SlackResponse{
			ResponseType: "ephemeral",
			Text:         "ProjectID and Title required params. Please use this format: \"#ProjectID#Title#Description\"",
		}))
	}

	projectID := texts[0]
	title := texts[1]
	description := ""
	if len(texts) == 3 {
		description = texts[2]
	}

	urlToTask, err := YouTrackAPI.CreateIssue(projectID, title, description)
	if err != nil {
		return c.Render(200, r.JSON(lib.SlackResponse{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Error create issue in YouTrack: %v", err),
		}))
	}

	return c.Render(200, r.JSON(lib.SlackResponse{
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
	}))
}
