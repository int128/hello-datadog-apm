package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func do(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/get", nil)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
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
	time.Sleep(30 * time.Second)
	tracer.Start(
		tracer.WithService("hello-datadog-apm"),
		tracer.WithEnv("github-actions"),
	)
	defer func() {
		tracer.Stop()
		time.Sleep(60 * time.Second)
	}()
	httptrace.WrapClient(http.DefaultClient)

	ctx := context.Background()
	span, ctx := tracer.StartSpanFromContext(ctx, "run")
	err := do(ctx)
	span.Finish(tracer.WithError(err))
	if err != nil {
		log.Printf("error: %s", err)
		return 1
	}
	return 0
}

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	code := run()
	log.Printf("Exiting with code %d", code)
	os.Exit(code)
}
