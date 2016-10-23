package gcrawler

import (
	"net/http"
	"io/ioutil"
	"github.com/petrbouda/golang-http-client"
)

const (
	googleSearchUrl = "http://www.google.com/search"
	userAgentHeader = "User-Agent"
	userAgentValue = "Mozilla/5.0"
)

type GoogleSearchClient struct{
	httpClient *http.Client
	linkParser *GoogleLinkParser
}

// Search a provided term using Google Search engine
// returns the links found on search page.
func (c *GoogleSearchClient) Search(term string) ([]string, error) {
	req, err := http.NewRequest("GET", googleSearchUrl, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", term)
	req.URL.RawQuery = q.Encode()
	req.Header.Add(userAgentHeader, userAgentValue)

	resp, err := invoke(c.httpClient, req)
	if err != nil {
		return nil, err
	}

	searchContent := string(resp)
	return c.linkParser.Parse(searchContent), nil
}

func invoke(client *http.Client, request *http.Request) ([]byte, error) {
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// return entire response body in bytes
	return ioutil.ReadAll(resp.Body)
}

func NewGoogleSearchClient() *GoogleSearchClient {
	return &GoogleSearchClient{
		httpClient: http_client.NewHttpClient(false),
		linkParser: &GoogleLinkParser{},
	}
}