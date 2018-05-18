package contextkeys

import (
	"context"
	"time"
)

type contextKey string

func (c contextKey) String() string {
	return "packetsource" + string(c)
}

var (
	contextKeyRecievedTime = contextKey("recievedTime")
)

// ContextWithRecievedTime Adds time to the ctx and returns it, used to collect
// the time the packet was received before it is deserialized
func ContextWithRecievedTime(ctx context.Context, time time.Time) context.Context {
	return context.WithValue(ctx, contextKeyRecievedTime, time)
}

// RecievedTimeFromContext returns the recievedTime recorded by packetSource as
// well as a bool that is true only if there is a time in the context
func RecievedTimeFromContext(ctx context.Context) (time.Time, bool) {
	t, ok := ctx.Value(contextKeyRecievedTime).(time.Time)
	return t, ok
}
