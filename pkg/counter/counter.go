package counter

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

type Counter struct {
	urls    []string
	rgxp    *regexp.Regexp
	results []string
	total   int
	m       sync.Mutex
}

func (c *Counter) getHtml(url string) string {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return ""
	}

	defer res.Body.Close()

	html, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)

		return ""
	}

	return string(html)
}

func (c *Counter) findAllSubsritngs(s string) []string {
	return c.rgxp.FindAllString(s, -1)
}

func (c *Counter) process(url string) {
	html := c.getHtml(url)
	entries := c.findAllSubsritngs(html)
	count := len(entries)

	c.m.Lock()
	c.results = append(c.results, fmt.Sprintf("Count for %v: %v", url, count))
	c.total += count
	c.m.Unlock()
}

func (c *Counter) print() {
	for _, v := range c.results {
		fmt.Println(v)
	}
	fmt.Println("Total: ", c.total)
}

func NewCounter(urls []string, pattern string) (*Counter, error) {
	rgxp, err := regexp.Compile(pattern)

	if err != nil {
		return nil, err
	}

	return &Counter{
		urls: urls,
		rgxp: rgxp,
	}, nil
}

func (c *Counter) Start() {
	wg := sync.WaitGroup{}

	urlsCount := len(c.urls)

	var workerPoolSize int

	if urlsCount < 5 {
		workerPoolSize = urlsCount
	} else {
		workerPoolSize = 5
	}

	dataCh := make(chan string, workerPoolSize)

	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range dataCh {
				c.process(data)
			}
		}()
	}

	for i := range c.urls {
		dataCh <- c.urls[i]
	}

	close(dataCh)

	wg.Wait()
	c.print()
}
