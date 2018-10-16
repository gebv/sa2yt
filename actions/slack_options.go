package actions

import (
	"encoding/json"
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

	payload := c.Request().Form.Get("payload")
	var encodedCallback lib.SlackActionCallback

	err := json.Unmarshal([]byte(payload), &encodedCallback)
	if err != nil {
		fmt.Printf("ERROR: Can't encode slack message: %v \n", err)
		return nil
	}

	fmt.Println("Query string", encodedCallback.Value)
	issues, err := YouTrackAPI.SearchIssues(encodedCallback.Value)
	if err != nil {
		fmt.Println("Search Error ", err)
		return nil
	}

	options := selectOptions{}
	for _, issue := range issues {
		options.Options = append(options.Options, lib.SlackDialogElementOption{
			Label: fmt.Sprintf("%s: %s", issue.EntityID(), issue.Summary),
			Value: issue.EntityID(),
		})
	}

	return c.Render(200, r.JSON(options))
}
