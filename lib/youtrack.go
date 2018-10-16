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

// YouTrackIssue - issue from YouTrack
type YouTrackIssue struct {
	ID       string      `json:"id"`
	EntityID string      `json:"entityId"`
	JiraID   interface{} `json:"jiraId"`
	Field    []struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"field"`
}

// Summary - get summary field from issue
func (issue *YouTrackIssue) Summary() string {
	var summary string

	for _, field := range issue.Field {
		if field.Name == "summary" {
			summary = field.Value.(string)
			break
		}
	}

	return summary
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
func (api *YouTrackAPI) SearchIssues(query string) ([]YouTrackIssue, error) {
	// old end point  - rest/issue
	// path := url.URL{Path: "youtrack/api/issues"}
	// params := path.Query()
	// params.Set("query", query)
	// params.Set("fields", "project(id,name),id,numberInProject,summary")
	// path.RawQuery = params.Encode()
	jsonStr, _ := json.Marshal(map[string]string{
		"query":  query,
		"fields": "project(id,name),id,numberInProject,summary,description",
	})

	response, err := api.sendJSONRequest("GET", &url.URL{Path: "youtrack/api/issues"}, jsonStr)

	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("SEARCH ISSUES RESP BODY %v\n", string(respBody))

	type searchAnswer struct {
		Issues []YouTrackIssue `json:"issue"`
	}
	var searchResp searchAnswer

	err = json.Unmarshal(respBody, &searchResp)
	if err != nil {
		return nil, err
	}

	return searchResp.Issues, nil
}

// CreateComment - add comment to specified Issue
func (api *YouTrackAPI) CreateComment(issueID, comment string) error {
	path := fmt.Sprintf("youtrack/api/issues/%s/comments", issueID)
	response, err := api.sendJSONRequest("POST", &url.URL{Path: path}, []byte(`{"text": "test comment"}`))
	//[]byte(fmt.Sprintf(`{"text": "%s"}`, comment))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Wrong response status from Youtrack is %d", response.StatusCode)
	}

	fmt.Println("Create Comment resp: ", response)

	return nil
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
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Token))

	fmt.Printf("REQUEST --- %v \n", request)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("RESPPPPP", response)

	return response, nil
}

// TODO: extract common parties
func (api *YouTrackAPI) sendJSONRequest(method string, path *url.URL, JSONstr []byte) (*http.Response, error) {
	client := &http.Client{}
	baseURL, err := url.Parse(api.Domain)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, baseURL.ResolveReference(path).String(), bytes.NewBuffer(JSONstr))
	request.Header.Set("content-type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Token))

	fmt.Printf("JSON REQUEST --- %v \n", request)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("JSON RESPPPPP", response)

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
