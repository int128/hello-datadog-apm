package main

import (
	"log"
	"os"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func do() {
	span := tracer.StartSpan("run")
	defer span.Finish()
	log.Printf("Hello world")
}

func run() int {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	tracer.Start(
		tracer.WithService("hello-datadog-apm"),
		tracer.WithEnv("github-actions"),
	)
	defer tracer.Stop()

	for i := 0; i < 50; i++ {
		do()
	}
	return 0
}

func main() {
	os.Exit(run())
}
