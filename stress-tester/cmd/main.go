package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"stress-tester/internal/usecase"
)

func main() {

	url, requests, concurrency := handleFlags()
	usecase.RoutineGet(*url, *requests, *concurrency)
}

func handleFlags() (url *string, requests *int, concurrency *int) {

	url = flag.String("url", "http://localhost:8080", "Url to be tested.")
	requests = flag.Int("requests", 105, "Qt of requests.")
	concurrency = flag.Int("concurrency", 10, "Qt of concurrent requests.")

	flag.Parse()

	errors := []string{}

	if *url == "" {
		errors = append(errors, "url must not be empty")
	}
	if *requests <= 0 {
		errors = append(errors, "requests must be greater than 0")
	}
	if *concurrency <= 0 {
		errors = append(errors, "concurrency must be greater than 0")
	}
	if len(errors) == 0 {
		req, err := http.Get(*url)
		if err != nil {
			errors = append(errors, err.Error())
		}
		if req != nil && req.StatusCode != 200 {
			errors = append(errors, fmt.Sprintf("Status code should be 200, but is %d. Check the URL.", req.StatusCode))
		}
	}

	if len(errors) > 0 {
		fmt.Println("Invalid parameters:")
		for _, err := range errors {
			fmt.Println(err)
			slog.Error(err)
		}
		fmt.Println("Usage: go run main.go --url=http://localhost:8080 --requests=1000 --concurrency=10")
		os.Exit(1)
	}

	return
}
