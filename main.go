package main

import (
	"eltask/pkg/counter"
	"log"
)

func main() {
	urls := []string{"https://golang.org"}
	counter, err := counter.NewCounter(urls, `Go`)

	if err != nil {
		log.Fatal(err)
	}

	counter.Start()
	counter.Print()
}
