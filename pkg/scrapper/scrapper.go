package scrapper

import (
	"eltask/pkg/pool"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

type scrapper struct {
	rgxp  *regexp.Regexp
	total int
	m     sync.Mutex
}

type Scrapper interface {
	Run(urls *[]string, rgxp *regexp.Regexp, maxPoolSize int)
}

func NewScrapper() Scrapper {
	return &scrapper{}
}

func (s *scrapper) getContent(url string) string {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return ""
	}

	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)

		return ""
	}

	return string(content)
}

func (s *scrapper) findSubsritngs(content string) []string {
	return s.rgxp.FindAllString(content, -1)
}

func (s *scrapper) write(msg string) {
	fmt.Println(msg)
}

func (s *scrapper) parse(url string) {
	content := s.getContent(url)
	entries := s.findSubsritngs(content)
	count := len(entries)

	s.write(fmt.Sprintf("Count for %v: %v", url, count))

	s.m.Lock()
	s.total += count
	s.m.Unlock()
}

func (s *scrapper) clean() {
	s.rgxp = nil
	s.total = 0
}

func (s *scrapper) Run(urls *[]string, rgxp *regexp.Regexp, maxPoolSize int) {
	s.rgxp = rgxp
	pool.Process(urls, s.parse, maxPoolSize)
	s.write(fmt.Sprintf("Total: %v", s.total))
}
