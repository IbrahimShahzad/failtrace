package main

import (
	"context"

	"github.com/IbrahimShahzad/failtrace"
)

func main() {
	ctx := failtrace.WithLogger(context.Background())
	handle(ctx)
}

func handle(ctx context.Context) {
	log := failtrace.FromContext(ctx)
	log.Debug("handling request")
	a(ctx)
}

func a(ctx context.Context) {
	log := failtrace.FromContext(ctx)
	defer log.FlushIf(nil)
	log.Debug("inside a")
	b(ctx)
}

func b(ctx context.Context) {
	log := failtrace.FromContext(ctx)
	defer log.FlushIf(nil)
	log.Debug("inside b")
	var err error
	if err != nil {
		log.FlushIf(err) // This will write the logs to the output
		return
	}
	log.FlushIf(nil)
}
