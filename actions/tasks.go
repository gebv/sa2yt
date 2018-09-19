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

	// texts := strings.Split(c.Request().Form.Get("text"), "#")
	splitFunc := func(c rune) bool { return c == '#' }
	texts := strings.FieldsFunc(c.Request().Form.Get("text"), splitFunc)

	fmt.Printf("TEXTS:   %v, %d\n", texts, len(texts))
	if len(texts) < 2 {
		return c.Render(200, r.JSON(lib.SlackResponse{
			ResponseType: "ephemeral",
			Text:         "ProjectID and Title required params. Please use this format: \"#ProjectID#Title#Description\"",
		}))
	}
	// err := YouTrackAPI.CreateIssue()
	// if err != nil {

	// 	return c.Render(200, r.String(fmt.Sprintf("Error: %v", err)))
	// }

	return c.Render(200, r.JSON(lib.SlackResponse{
		Text: "Task was created",
		Attachments: []lib.SlackAttachment{
			{
				Text: "Link to task", // TODO: return link from Location header from response
			},
		},
	}))
}
