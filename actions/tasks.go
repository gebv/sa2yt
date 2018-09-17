package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
)

// TasksCreate create new task in YouTrack.
func TasksCreate(c buffalo.Context) error {
	// fmt.Printf("PARAMS: %v", c.Data())
	fmt.Printf("Form: %v \n", c.Request().Form)

	err := YouTrackAPI.CreateIssue()
	if err != nil {
		return c.Render(200, r.Plain(fmt.Sprintf("Error: %v", err)))
	}

	return c.Render(200, r.Plain("All OK"))
}
