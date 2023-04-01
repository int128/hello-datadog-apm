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
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	tracer.Start(
		tracer.WithService("hello-datadog-apm"),
		tracer.WithEnv("github-actions"),
	)
	defer tracer.Stop()
	httptrace.WrapClient(http.DefaultClient)
	waitForDatadogAgent()

	var err error
	ctx := context.Background()
	span, ctx := tracer.StartSpanFromContext(ctx, "run")
	defer func() {
		span.Finish(tracer.WithError(err))
	}()
	err = do(ctx)
	if err != nil {
		log.Printf("error: %s", err)
		return 1
	}
	return 0
}

func waitForDatadogAgent() {
	for i := 0; i < 120; i++ {
		resp, err := http.Get("http://localhost:8126/info")
		if resp != nil {
			resp.Body.Close()
		}
		if err == nil && resp.StatusCode == 200 {
			log.Printf("datadog-agent is ready")
			return
		}
		log.Printf("Waiting for datadog-agent: %s", err)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	os.Exit(run())
}
