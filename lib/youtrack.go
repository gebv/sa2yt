package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// api.sendRequest("GET", &url.URL{Path: "rest/issue/NTA-1"})

	// api.sendRequest("PUT", &url.URL{Path: "/rest/issue?project=NTA&summary=new&description=description"})
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

	fmt.Println("URL --- ", baseURL.ResolveReference(path).String())
	fmt.Println("TOKEN --- ", fmt.Sprintf("Bearer %s", api.Token))
	fmt.Printf("REQUEST --- %v \n", req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("RESPPPPP", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jsonBody interface{}
	json.Unmarshal(body, &jsonBody)
	fmt.Printf("RESP BODY %v \n", jsonBody)

	return nil
}
