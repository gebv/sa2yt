package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
)

// TasksCreate create new task in YouTrack.
func TasksCreate(c buffalo.Context) error {
	// fmt.Printf("PARAMS: %v", c.Data())
	fmt.Printf("Form: %v \n", c.Request().Form)

	return c.Render(200, r.JSON(map[string]string{"status": "OK"}))
}
