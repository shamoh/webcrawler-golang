package gcrawler

import (
	"github.com/petrbouda/golang-http-client"
	"net/http"
	"sync"
	"github.com/Sirupsen/logrus"
	"os"
	"strings"
)

func init() {
	logger = logrus.New()
	logger.Out = os.Stdout
}

var logger *logrus.Logger

func NewCrawler() *Crawler {
	return &Crawler{
		libraryParser:  LibraryParser{},
		searchClient:   NewGoogleSearchClient(),
		httpClient:    	http_client.NewHttpClient(false),
		done:           make(chan struct{}),
		maxRequest: 	make(chan struct{}, 20),
	}
}

type Crawler struct {
	libraryParser Parser
	searchClient  *GoogleSearchClient
	httpClient    *http.Client
	done     	  chan struct{}
	maxRequest	  chan struct{}
}

func (crawler Crawler) Crawl(term string) ([]string, error) {
	pages, err := crawler.searchClient.Search(term)
	if err != nil {
		return nil, err
	}

	var wait sync.WaitGroup
	libraries := make(chan []string, len(pages))
	for _, page := range pages {
		wait.Add(1)
		go crawler.processPage(page, &wait, libraries)
	}

	// Wait for the all processed pages and close the channel for incoming results
	go func() {
		wait.Wait()
		close(libraries)
	}()

	storedLibraries := make([]string, len(pages))
	printer:
	for {
		select {
		case <- crawler.done:
			// Drain libraries to allow existing goroutines to finish.
			for range libraries {
				// Do nothing.
			}
			break printer
		case library, ok := <- libraries:
			if !ok {
				// libraries channel has been already closed
				break printer
			}
			logrus.Infof("Found library: `%s`", library)
			storedLibraries = append(storedLibraries, library...)
		}
	}

	return storedLibraries, nil
}

func (crawler Crawler) processPage(url string, wait *sync.WaitGroup, libraries chan <- []string) {
	defer wait.Done()
	if crawler.cancelled() {
		return
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Warnf("GET request cannot be created for url `%s`: %v\n", url)
		return
	}

	resp, err := invoke(crawler.httpClient, req)
	if err != nil {
		logger.Warnf("Error occurred during HTTP processing, URL:`%s`, Error: %v\n", url)
		return
	}

	libraryName := crawler.libraryParser.Parse(string(resp[:]))
	libraries <- libraryName
}

// Http invoker with inner throttling for concurrent http requests.
// Number of concurrent messages can be changed as a size of `maxRequest` channel
func (crawler Crawler) httpInvoke(client *http.Client, request *http.Request) ([]byte, error) {
	select {
	// Add a record to the channel indicating that one HTTP request will be processed.
	case crawler.maxRequest <- struct{}{}:
	// Crawler has been already cancelled therefore none HTTP request should be created.
	case <- crawler.done:
		return nil, nil
	}
	// Add the end of the function remove the record indicating HTTP processing.
	defer func() {
		<- crawler.maxRequest
	}()

	return invoke(client, request)
}

// Method recognizes whether the Crawler has been already cancelled or still running.
func (crawler Crawler) cancelled() bool {
	select {
	case <- crawler.done:
		return true
	default:
		return false
	}
}