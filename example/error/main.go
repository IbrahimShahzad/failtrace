package main

import (
	"context"
	"errors"

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
	err := errors.New("an error occurred in b")
	if err != nil {
		log.FlushIf(err) // This will write the logs to the output
		return
	}
	log.FlushIf(nil)
}
