package events

import (
	"context"
	"time"
)

type EventMeta struct {
	CollectedTimeNanos int64
}

func (e EventMeta) GetCollectedTime() time.Time {
	return time.Unix(0, e.CollectedTimeNanos)
}

type DataEvent struct {
	EventMeta
	Data map[string]interface{}
}

// ContextDataEvent adds context to a RawEvent to hold request-scopped info
type ContextDataEvent struct {
	context.Context
	DataEvent
}

// ContextRawEvent adds context to a RawEvent to hold request-scopped info
type ContextRawEvent struct {
	context.Context
	RawEvent
}

type RawEvent struct {
	EventMeta
	Data []byte
}

func NewRawEventNow(packet []byte) RawEvent {
	return RawEvent{
		EventMeta: EventMeta{
			CollectedTimeNanos: time.Now().UnixNano(),
		},
		Data: packet,
	}
}
