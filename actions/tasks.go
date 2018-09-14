package actions

import "github.com/gobuffalo/buffalo"

// TasksCreate create new task in YouTrack.
func TasksCreate(c buffalo.Context) error {
	return c.Render(200, r.JSON(map[string]string{"status": "OK"}))
}
