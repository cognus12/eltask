package counter

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Counter struct {
	urls  []string
	rgxp  *regexp.Regexp
	data  map[string]int
	total int
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

func (c *Counter) processUrl(u string) {
	html := c.getHtml(u)
	entries := c.findAllSubsritngs(html)
	count := len(entries)
	c.data[u] = count

	c.total += count
}

func NewCounter(urls []string, pattern string) (*Counter, error) {
	rgxp, err := regexp.Compile(pattern)

	if err != nil {
		return nil, err
	}

	return &Counter{
		urls: urls,
		rgxp: rgxp,
		data: make(map[string]int),
	}, nil
}

func (c *Counter) Start() {
	for _, u := range c.urls {
		c.processUrl(u)
	}
}

func (c *Counter) Print() {
	for k, v := range c.data {
		fmt.Printf("Count for %v: %v \n", k, v)
	}
	fmt.Println("Total: ", c.total)
}
