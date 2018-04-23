// Package log should be used for all logging. These log messages are
// context-aware so that they can log to OpenTracing and Jaeger. This also
// encourages the pattern of logging only when context is available, which
// encourages returning errors in other functions rather than logging.
//
// If you are looking to add some kind of instrumentation, this is the place to
// do it!
package log

import (
	"context"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

const eventKey = "event"
const errorTypeKey = "error_type"

// Error logs the error message to the console and to the active span if there
// is one so the error shows up in traces.
func Error(ctx context.Context, err error, msg string) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		// handle errors by recording them in the span
		span.SetTag(string(ext.Error), true)
		if msg != "" {
			span.LogFields(
				otlog.Error(err),
				otlog.String(errorTypeKey, msg),
			)
		} else {
			span.LogFields(
				otlog.Error(err),
				otlog.String(errorTypeKey, msg),
			)
		}

	}
	// Print error out
	log.Printf("%v: %v", msg, err)
}

// Event logs a message to the console and to the active span if there is one.
func Event(ctx context.Context, msg string) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		// record message in span
		span.LogFields(otlog.String(eventKey, msg))
	}
	// Print message out
	log.Print(msg)
}
