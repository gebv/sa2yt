package lib

import (
	"bytes"
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
func (api *YouTrackAPI) CreateIssue(projectID, title, description string) error {
	api.sendRequest("PUT", &url.URL{Path: "youtrack/rest/issue"}, map[string]string{
		"project":     "NTA",
		"summary":     "New Issue from API",
		"description": "Full issue description",
	})

	// response &{201 Created 201 HTTP/2.0 2 0 map[Set-Cookie:[JSESSIONID=1uwuvkwsou61c6jo4b4jhutpx;Path=/youtrack;Secure;HttpOnly] Cache-Control:[no-cache, no-store, no-transform, must-revalidate] X-Frame-Options:[SAMEORIGIN] Content-Length:[0] Access-Control-Expose-Headers:[Location] Server:[nginx] X-Content-Type-Options:[nosniff] Referrer-Policy:[strict-origin-when-cross-origin] Vary:[Accept-Encoding, User-Agent] Location:[https://syato-test-app.myjetbrains.com/youtrack/rest/issue/NTA-5] X-Xss-Protection:[1; mode=block] Expires:[Thu, 01 Jan 1970 00:00:00 GMT] Strict-Transport-Security:[max-age=31536000; includeSubdomains;] Date:[Wed, 19 Sep 2018 18:36:30 GMT]] {0xc4201e6c80} 0 [] false false map[] 0xc42013a400 0xc4201d5ad0}
	return nil
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
