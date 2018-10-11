package actions

import (
	"github.com/gebv/sayto/lib"
	"github.com/gobuffalo/buffalo"
)

// SlackOptionsIndex default implementation.
func SlackOptionsIndex(c buffalo.Context) error {
	options := []lib.SlackDialogElementOption{
		{Label: "Label 1", Value: "Value 1"},
		{Label: "Label 2", Value: "Value 2"},
		{Label: "Label 3", Value: "Value 3"},
		{Label: "Label 4", Value: "Value 4"},
	}

	return c.Render(200, r.JSON(options))
}
