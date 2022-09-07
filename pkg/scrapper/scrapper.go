package scrapper

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

type Scrapper struct {
	rgxp    *regexp.Regexp
	results []string
	total   int
	m       sync.Mutex
}

func NewScrapper(rgxp *regexp.Regexp) *Scrapper {
	return &Scrapper{
		rgxp: rgxp,
	}
}

func (s *Scrapper) getContent(url string) string {
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

func (s *Scrapper) findSubsritngs(content string) []string {
	return s.rgxp.FindAllString(content, -1)
}

func (s *Scrapper) process(url string) {
	content := s.getContent(url)
	entries := s.findSubsritngs(content)
	count := len(entries)

	s.m.Lock()
	s.results = append(s.results, fmt.Sprintf("Count for %v: %v", url, count))
	s.total += count
	s.m.Unlock()
}

func (s *Scrapper) print() {
	for _, v := range s.results {
		fmt.Println(v)
	}
	fmt.Println("Total: ", s.total)
}

func (s *Scrapper) Run(urls *[]string, maxPoolSize int) {
	wg := sync.WaitGroup{}

	urlsCount := len(*urls)

	var workerPoolSize int

	if urlsCount < maxPoolSize {
		workerPoolSize = urlsCount
	} else {
		workerPoolSize = maxPoolSize
	}

	dataCh := make(chan string, workerPoolSize)

	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range dataCh {
				s.process(data)
			}
		}()
	}

	for _, u := range *urls {
		dataCh <- u
	}

	close(dataCh)

	wg.Wait()
	s.print()
}
