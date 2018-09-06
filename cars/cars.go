package cars

import (
	"context"
	"github.com/sscp/telemetry/cars/sundae"
)

//go:generate enumer -type=Car -json -text -output=gen_cars.go

// A car
type Car int

const (
	// Sundae
	Sundae Car = iota
)

func stubCarSupport(ctx context.Context, packet []byte) (map[string]interface{}, error) {
	return make(map[string]interface{}, 0), nil
}

func GetCarSupport(car Car) func(ctx context.Context, packet []byte) (map[string]interface{}, error) {

	switch car {
	case Sundae:
		return sundae.Deserialize
	default:
		return stubCarSupport
	}
}
