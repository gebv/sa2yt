package actions

import (
	"fmt"

	"github.com/gebv/sayto/lib"
	"github.com/gobuffalo/buffalo"
)

type selectOptions struct {
	Options []lib.SlackDialogElementOption `json:"options"`
}

// SlackOptionsIndex default implementation.
func SlackOptionsIndex(c buffalo.Context) error {
	fmt.Printf("Form Options: %v \n", c.Request().Form)

	// payload := c.Request().Form.Get("payload")
	// var encodedCallback lib.SlackActionCallback

	// err := json.Unmarshal([]byte(payload), &encodedCallback)
	// if err != nil {
	// 	fmt.Printf("ERROR: Can't encode slack message: %v \n", err)
	// 	return nil
	// }

	// fmt.Println("Query string", encodedCallback.Value)
	// _, err = YouTrackAPI.SearchIssues(encodedCallback.Value)
	// if err != nil {
	// 	fmt.Println("Search Error ", err)
	// 	return nil
	// }

	options := selectOptions{
		Options: []lib.SlackDialogElementOption{
			{Label: "Label 1", Value: "Value 1"},
			{Label: "Label 2", Value: "Value 2"},
			{Label: "Label 3", Value: "Value 3"},
			{Label: "Label 4", Value: "Value 4"},
		},
	}

	return c.Render(200, r.JSON(options))
}
