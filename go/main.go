package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func do() error {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		return fmt.Errorf("failed to get: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read: %w", err)
	}
	log.Printf("got: %s", b)
	return nil
}

func run() int {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	tracer.Start(
		tracer.WithService("hello-datadog-apm"),
		tracer.WithEnv("github-actions"),
	)
	defer tracer.Stop()

	var err error
	span := tracer.StartSpan("run")
	defer span.Finish(tracer.WithError(err))
	err = do()
	if err != nil {
		log.Printf("error: %s", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
