package actions

import (
	"encoding/json"
	"fmt"

	"github.com/gebv/sayto/lib"
	"github.com/gobuffalo/buffalo"
)

type selectOptions struct {
	Options []selectOption `json:"options"`
}

type selectOption struct {
	Text  string `json:"text"`
	Value string `json:"value"`
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

	options := selectOptions{
		Options: []selectOption{
			{Text: "Label 1", Value: "Value 1"},
			{Text: "Label 2", Value: "Value 2"},
			{Text: "Label 3", Value: "Value 3"},
			{Text: "Label 4", Value: "Value 4"},
		},
	}

	fmt.Println("Query string", encodedCallback.Value)
	_, err = YouTrackAPI.SearchIssues(encodedCallback.Value)
	if err != nil {
		fmt.Println("Search Error ", err)
		return nil
	}
	// TODO: add search task by https://www.jetbrains.com/help/youtrack/standalone/Intellisense-for-issue-search.html

	return c.Render(200, r.JSON(options))
}
