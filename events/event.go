package events

import (
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

type RawEvent struct {
	EventMeta
	Data []byte
}

func NewRawDataEvent(packet []byte) RawEvent {
	return RawEvent{
		EventMeta: EventMeta{
			CollectedTimeNanos: time.Now().UnixNano(),
		},
		Data: packet,
	}
}
