package main

import (
	"eltask/pkg/counter"
	"log"
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
	counter, err := counter.NewCounter(urls, `Go`)

	if err != nil {
		log.Fatal(err)
	}

	counter.Start()
}
