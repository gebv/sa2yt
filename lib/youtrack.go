package lib

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// YouTrackAPI - struct for youtrack api
type YouTrackAPI struct {
	Token  string
	Domain string
}

// CreateIssue - create New Issue in YouTrack
func (api *YouTrackAPI) CreateIssue(projectID, summary, description string) (string, error) {
	response, err := api.sendRequest("PUT", &url.URL{Path: "youtrack/rest/issue"}, map[string]string{
		"project":     projectID,
		"summary":     summary,
		"description": description,
	})

	if err != nil {
		return "", err
	}

	if response.StatusCode != 201 {
		return "", fmt.Errorf("Wrong response status from Youtrack is %d", response.StatusCode)
	}

	fmt.Println("Create Issue resp: ", response)

	restURL := response.Header.Get("Location")
	return strings.Replace(restURL, "/rest", "", 1), nil
}

func (api *YouTrackAPI) sendRequest(method string, path *url.URL, params map[string]string) (*http.Response, error) {
	client := &http.Client{}
	baseURL, err := url.Parse(api.Domain)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, baseURL.ResolveReference(path).String(), prepareParams(params))
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Token))

	fmt.Println("URL --- ", baseURL.ResolveReference(path).String())
	fmt.Println("TOKEN --- ", fmt.Sprintf("Bearer %s", api.Token))
	fmt.Printf("REQUEST --- %v \n", request)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("RESPPPPP", response)

	return response, nil
}

func prepareParams(params map[string]string) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	values := url.Values{}
	for param, value := range params {
		values.Set(param, value)
	}

	buffer.WriteString(values.Encode())
	return buffer
}
