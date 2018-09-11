package cars

import (
	"context"

	"github.com/sscp/telemetry/cars/sundae"

	"github.com/sscp/telemetry/events"
)

//go:generate enumer -type=Car -json -text -output=gen_cars.go

// Car refers to a car-specific deserialization routine.
type Car int

const (
	// Sundae is sscp's 2017 WSC car
	Sundae Car = iota
)

func stubCarSupport(ctx context.Context, event events.RawEvent) (events.DataEvent, error) {
	return events.DataEvent{
		EventMeta: event.EventMeta,
		Data:      make(map[string]interface{}, 0),
	}, nil
}

func GetCarDeserializer(car Car) func(ctx context.Context, rawEvent events.RawEvent) (events.DataEvent, error) {
	switch car {
	case Sundae:
		return sundae.Deserialize
	default:
		return stubCarSupport
	}
}
