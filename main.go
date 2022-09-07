package main

import (
	"eltask/pkg/scrapper"
	"log"
	"regexp"
)

func main() {
	urls := []string{"https://golang.org",
		"https://gobyexample.com/worker-pools",
		"https://gobyexample.com/worker-pools",
		"https://gobyexample.com/waitgroups",
		"https://gobyexample.com/rate-limiting",
		"https://gobyexample.com/atomic-counters",
		"https://gobyexample.com/mutexes",
	}

	rgxp, err := regexp.Compile(`Go`)

	if err != nil {
		log.Fatal(err)
	}

	s := scrapper.NewScrapper()

	s.Run(&urls, rgxp, 5)
}
