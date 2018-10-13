package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// YouTrackAPI - struct for youtrack api
type YouTrackAPI struct {
	Token          string
	Domain         string
	CachedProjects []YouTrackProject
}

// YouTrackProject - projects on YouTrack
type YouTrackProject struct {
	ID  string `json:"id"`
	URL string `json:"url"`
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

// SearchIssues - search Issues in YouTrack
func (api *YouTrackAPI) SearchIssues(query string) (string, error) {
	response, err := api.sendRequest("GET", &url.URL{Path: "youtrack/rest/issue"}, map[string]string{
		"filter": query,
	})

	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Search Issue resp: ", respBody)

	return "", nil
}

// RefreshProjectsCache - get available projects from YouTrack
func (api *YouTrackAPI) RefreshProjectsCache() error {
	projects := []YouTrackProject{}
	response, err := api.getAllProjects()
	if err != nil {
		return err
	}

	fmt.Println("Get projects resp: ", response)
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, &projects)
	if err != nil {
		return err
	}

	api.CachedProjects = projects

	return nil
}

func (api *YouTrackAPI) getAllProjects() (*http.Response, error) {
	response, err := api.sendRequest("GET", &url.URL{Path: "youtrack/rest/admin/project"}, map[string]string{})

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Wrong response status from Youtrack is %d", response.StatusCode)
	}

	return response, nil
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
