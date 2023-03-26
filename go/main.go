package main

import (
	"log"
	"os"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func run() int {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	tracer.Start(tracer.WithService("hello-datadog-apm"))
	defer tracer.Stop()

	span := tracer.StartSpan("run")
	defer span.Finish()

	log.Printf("Hello world")
	return 0
}

func main() {
	os.Exit(run())
}
