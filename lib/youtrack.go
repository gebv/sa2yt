package lib

import (
	"fmt"
	"net/http"
	"net/url"
)

// YouTrackAPI - struct for youtrack api
type YouTrackAPI struct {
	Token  string
	Domain string
}

// CreateIssue - create New Issue in YouTrack
func (api *YouTrackAPI) CreateIssue() error {
	api.sendRequest("GET", &url.URL{Path: "rest/admin/project"})
	return nil
}

func (api *YouTrackAPI) sendRequest(method string, path *url.URL) error {
	client := &http.Client{}
	baseURL, err := url.Parse(api.Domain)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, baseURL.ResolveReference(path).String(), nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Token))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}
