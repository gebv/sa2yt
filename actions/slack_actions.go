package actions

import (
	"fmt"

	"github.com/gebv/sayto/lib"
	"github.com/gobuffalo/buffalo"
)

// SlackActionsCreate default implementation.
func SlackActionsCreate(c buffalo.Context) error {
	fmt.Printf("Form: %v \n", c.Request().Form)

	return c.Render(200, r.JSON(lib.SlackResponse{
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
	}))
}
