package sundae

import (
	"context"
	"fmt"
	"math"

	"github.com/golang/protobuf/proto"
	//"github.com/opentracing/opentracing-go"

	// Fork of https://github.com/fatih/structs/ that adds an "indirect"
	// option to dereference pointers to get values, not pointers in map
	"github.com/jackbeasley/structs"

	"github.com/sscp/telemetry/events"
	"github.com/sscp/telemetry/log"
)

//go:generate protoc -I=. --go_out=. ./sundae.proto
//go:generate protoc-go-inject-tag -input=./sundae.pb.go

const tryToHandlePadding = false

func verifyFloat32(ctx context.Context, val *float32, fieldName string) *float32 {
	if val == nil {
		return nil
	}
	if math.IsNaN(float64(*val)) || math.IsInf(float64(*val), 0) {
		log.Event(ctx, fmt.Sprintf("%v found in field %v", *val, fieldName))
		return nil
	}
	return val
}

// VerifyFloatValues makes sure that all floats are not INF of NaN as those are
// invalid numbers. If there are invalid numbers, they are replaced with nil and
// thus removed from the struct. This is all done in-place and errors are
// logged, but not returned as these errors are not failures, processing can
// continue with the invalid values removed.
func VerifyFloatValues(ctx context.Context, dm *SundaeDataMessage) {
	dm.LowVoltPower = verifyFloat32(ctx, dm.LowVoltPower, "LowVoltPower")

	dm.BatteryTemp_1 = verifyFloat32(ctx, dm.BatteryTemp_1, "BatteryTemp_1")
	dm.BatteryTemp_2 = verifyFloat32(ctx, dm.BatteryTemp_2, "BatteryTemp_2")
	dm.BatteryTemp_3 = verifyFloat32(ctx, dm.BatteryTemp_3, "BatteryTemp_3")
	dm.BatteryTemp_4 = verifyFloat32(ctx, dm.BatteryTemp_4, "BatteryTemp_4")
	dm.BatteryTemp_5 = verifyFloat32(ctx, dm.BatteryTemp_5, "BatteryTemp_5")
	dm.BatteryTemp_6 = verifyFloat32(ctx, dm.BatteryTemp_6, "BatteryTemp_6")
	dm.BmsPrechargeBatteryAdc = verifyFloat32(ctx, dm.BmsPrechargeBatteryAdc, "BmsPrechargeBatteryAdc")
	dm.BmsPrechargeCarAdc = verifyFloat32(ctx, dm.BmsPrechargeCarAdc, "BmsPrechargeCarAdc")
	dm.LowVoltOutputCurrent = verifyFloat32(ctx, dm.LowVoltOutputCurrent, "LowVoltOutputCurrent")
	dm.BatteryCurrent = verifyFloat32(ctx, dm.BatteryCurrent, "BatteryCurrent")
	dm.RightMotorControllerPower = verifyFloat32(ctx, dm.RightMotorControllerPower, "RightMotorControllerPower")
	dm.AmpHours = verifyFloat32(ctx, dm.AmpHours, "AmpHours")

	dm.MotorControllerSpeed = verifyFloat32(ctx, dm.MotorControllerSpeed, "MotorControllerSpeed")
	dm.MotorControllerRpm = verifyFloat32(ctx, dm.MotorControllerRpm, "MotorControllerRpm")
	dm.AvgOdometer = verifyFloat32(ctx, dm.AvgOdometer, "AvgOdometer")
	dm.LeftMotorTemp = verifyFloat32(ctx, dm.LeftMotorTemp, "LeftMotorTemp")
	dm.RightMotorTemp = verifyFloat32(ctx, dm.RightMotorTemp, "RightMotorTemp")
	dm.LeftMotorControllerTemp = verifyFloat32(ctx, dm.LeftMotorControllerTemp, "LeftMotorControllerTemp")
	dm.RightMotorControllerTemp = verifyFloat32(ctx, dm.RightMotorControllerTemp, "RightMotorControllerTemp")
	dm.LeftMotorControllerAlive = verifyFloat32(ctx, dm.LeftMotorControllerAlive, "LeftMotorControllerAlive")
	dm.RightMotorControllerAlive = verifyFloat32(ctx, dm.RightMotorControllerAlive, "RightMotorControllerAlive")
	dm.LeftMotorControllerCurrent = verifyFloat32(ctx, dm.LeftMotorControllerCurrent, "LeftMotorControllerCurrent")
	dm.RightMotorControllerCurrent = verifyFloat32(ctx, dm.RightMotorControllerCurrent, "RightMotorControllerCurrent")
	dm.MotorControllerCurrentDiff = verifyFloat32(ctx, dm.MotorControllerCurrentDiff, "MotorControllerCurrentDiff")

	dm.LeftMotorControllerBusVoltage = verifyFloat32(ctx, dm.LeftMotorControllerBusVoltage, "LeftMotorControllerBusVoltage")
	dm.RightMotorControllerBusVoltage = verifyFloat32(ctx, dm.RightMotorControllerBusVoltage, "RightMotorControllerBusVoltage")
	dm.LeftMotorController_15VVoltage = verifyFloat32(ctx, dm.LeftMotorController_15VVoltage, "LeftMotorController_15VVoltage")
	dm.RightMotorController_15VVoltage = verifyFloat32(ctx, dm.RightMotorController_15VVoltage, "RightMotorController_15VVoltage")
	dm.LeftMotorController_3V3Voltage = verifyFloat32(ctx, dm.LeftMotorController_3V3Voltage, "LeftMotorController_3V3Voltage")
	dm.RightMotorController_3V3Voltage = verifyFloat32(ctx, dm.RightMotorController_3V3Voltage, "RightMotorController_3V3Voltage")
	dm.LeftMotorController_1V9Voltage = verifyFloat32(ctx, dm.LeftMotorController_1V9Voltage, "LeftMotorController_1V9Voltage")
	dm.RightMotorController_1V9Voltage = verifyFloat32(ctx, dm.RightMotorController_1V9Voltage, "RightMotorController_1V9Voltage")
	dm.LeftMotorControllerDspTemp = verifyFloat32(ctx, dm.LeftMotorControllerDspTemp, "LeftMotorControllerDspTemp")
	dm.RightMotorControllerDspTemp = verifyFloat32(ctx, dm.RightMotorControllerDspTemp, "RightMotorControllerDspTemp")
	dm.LeftMotorControllerPhaseCurrent = verifyFloat32(ctx, dm.LeftMotorControllerPhaseCurrent, "LeftMotorControllerPhaseCurrent")
	dm.RightMotorControllerPhaseCurrent = verifyFloat32(ctx, dm.RightMotorControllerPhaseCurrent, "RightMotorControllerPhaseCurrent")
	dm.LeftMotorRpmCommand = verifyFloat32(ctx, dm.LeftMotorRpmCommand, "LeftMotorRpmCommand")
	dm.RightMotorRpmCommand = verifyFloat32(ctx, dm.RightMotorRpmCommand, "RightMotorRpmCommand")
	dm.LeftMotorCurrentCommand = verifyFloat32(ctx, dm.LeftMotorCurrentCommand, "LeftMotorCurrentCommand")
	dm.RightMotorCurrentCommand = verifyFloat32(ctx, dm.RightMotorCurrentCommand, "RightMotorCurrentCommand")

	dm.GpsLatitude = verifyFloat32(ctx, dm.GpsLatitude, "GpsLatitude")
	dm.GpsLongitude = verifyFloat32(ctx, dm.GpsLongitude, "GpsLongitude")
	dm.GpsSpeed = verifyFloat32(ctx, dm.GpsSpeed, "GpsSpeed")
	dm.GpsAltitude = verifyFloat32(ctx, dm.GpsAltitude, "GpsAltitude")
	dm.GpsBearing = verifyFloat32(ctx, dm.GpsBearing, "GpsBearing")

	dm.MpptArrayPower = verifyFloat32(ctx, dm.MpptArrayPower, "MpptArrayPower")
	dm.Mppt_A0VoltIn = verifyFloat32(ctx, dm.Mppt_A0VoltIn, "Mppt_A0VoltIn")
	dm.Mppt_A0VoltOut = verifyFloat32(ctx, dm.Mppt_A0VoltOut, "Mppt_A0VoltOut")
	dm.Mppt_A0Current = verifyFloat32(ctx, dm.Mppt_A0Current, "Mppt_A0Current")
	dm.Mppt_A0Temp = verifyFloat32(ctx, dm.Mppt_A0Temp, "Mppt_A0Temp")
	dm.Mppt_A1VoltIn = verifyFloat32(ctx, dm.Mppt_A1VoltIn, "Mppt_A1VoltIn")
	dm.Mppt_A1VoltOut = verifyFloat32(ctx, dm.Mppt_A1VoltOut, "Mppt_A1VoltOut")
	dm.Mppt_A1Current = verifyFloat32(ctx, dm.Mppt_A1Current, "Mppt_A1Current")
	dm.Mppt_A1Temp = verifyFloat32(ctx, dm.Mppt_A1Temp, "Mppt_A1Temp")
	dm.Mppt_B0VoltIn = verifyFloat32(ctx, dm.Mppt_B0VoltIn, "Mppt_B0VoltIn")
	dm.Mppt_B0VoltOut = verifyFloat32(ctx, dm.Mppt_B0VoltOut, "Mppt_B0VoltOut")
	dm.Mppt_B0Current = verifyFloat32(ctx, dm.Mppt_B0Current, "Mppt_B0Current")
	dm.Mppt_B0Temp = verifyFloat32(ctx, dm.Mppt_B0Temp, "Mppt_B0Temp")
	dm.Mppt_B1VoltIn = verifyFloat32(ctx, dm.Mppt_B1VoltIn, "Mppt_B1VoltIn")
	dm.Mppt_B1VoltOut = verifyFloat32(ctx, dm.Mppt_B1VoltOut, "Mppt_B1VoltOut")
	dm.Mppt_B1Current = verifyFloat32(ctx, dm.Mppt_B1Current, "Mppt_B1Current")
	dm.Mppt_B1Temp = verifyFloat32(ctx, dm.Mppt_B1Temp, "Mppt_B1Temp")
	dm.Mppt_C0VoltIn = verifyFloat32(ctx, dm.Mppt_C0VoltIn, "Mppt_C0VoltIn")
	dm.Mppt_C0VoltOut = verifyFloat32(ctx, dm.Mppt_C0VoltOut, "Mppt_C0VoltOut")
	dm.Mppt_C0Current = verifyFloat32(ctx, dm.Mppt_C0Current, "Mppt_C0Current")
	dm.Mppt_C0Temp = verifyFloat32(ctx, dm.Mppt_C0Temp, "Mppt_C0Temp")
	dm.Mppt_C1VoltIn = verifyFloat32(ctx, dm.Mppt_C1VoltIn, "Mppt_C1VoltIn")
	dm.Mppt_C1VoltOut = verifyFloat32(ctx, dm.Mppt_C1VoltOut, "Mppt_C1VoltOut")
	dm.Mppt_C1Current = verifyFloat32(ctx, dm.Mppt_C1Current, "Mppt_C1Current")
	dm.Mppt_C1Temp = verifyFloat32(ctx, dm.Mppt_C1Temp, "Mppt_C1Temp")
	dm.Mppt_D0VoltIn = verifyFloat32(ctx, dm.Mppt_D0VoltIn, "Mppt_D0VoltIn")
	dm.Mppt_D0VoltOut = verifyFloat32(ctx, dm.Mppt_D0VoltOut, "Mppt_D0VoltOut")
	dm.Mppt_D0Current = verifyFloat32(ctx, dm.Mppt_D0Current, "Mppt_D0Current")
	dm.Mppt_D0Temp = verifyFloat32(ctx, dm.Mppt_D0Temp, "Mppt_D0Temp")
	dm.Mppt_D1VoltIn = verifyFloat32(ctx, dm.Mppt_D1VoltIn, "Mppt_D1VoltIn")
	dm.Mppt_D1VoltOut = verifyFloat32(ctx, dm.Mppt_D1VoltOut, "Mppt_D1VoltOut")
	dm.Mppt_D1Current = verifyFloat32(ctx, dm.Mppt_D1Current, "Mppt_D1Current")
	dm.Mppt_D1Temp = verifyFloat32(ctx, dm.Mppt_D1Temp, "Mppt_D1Temp")
	dm.CellVolt_1 = verifyFloat32(ctx, dm.CellVolt_1, "CellVolt_1")
	dm.CellVolt_2 = verifyFloat32(ctx, dm.CellVolt_2, "CellVolt_2")
	dm.CellVolt_3 = verifyFloat32(ctx, dm.CellVolt_3, "CellVolt_3")
	dm.CellVolt_4 = verifyFloat32(ctx, dm.CellVolt_4, "CellVolt_4")
	dm.CellVolt_5 = verifyFloat32(ctx, dm.CellVolt_5, "CellVolt_5")
	dm.CellVolt_6 = verifyFloat32(ctx, dm.CellVolt_6, "CellVolt_6")
	dm.CellVolt_7 = verifyFloat32(ctx, dm.CellVolt_7, "CellVolt_7")
	dm.CellVolt_8 = verifyFloat32(ctx, dm.CellVolt_8, "CellVolt_8")
	dm.CellVolt_9 = verifyFloat32(ctx, dm.CellVolt_9, "CellVolt_9")
	dm.CellVolt_10 = verifyFloat32(ctx, dm.CellVolt_10, "CellVolt_10")
	dm.CellVolt_11 = verifyFloat32(ctx, dm.CellVolt_11, "CellVolt_11")
	dm.CellVolt_12 = verifyFloat32(ctx, dm.CellVolt_12, "CellVolt_12")
	dm.CellVolt_13 = verifyFloat32(ctx, dm.CellVolt_13, "CellVolt_13")
	dm.CellVolt_14 = verifyFloat32(ctx, dm.CellVolt_14, "CellVolt_14")
	dm.CellVolt_15 = verifyFloat32(ctx, dm.CellVolt_15, "CellVolt_15")
	dm.CellVolt_16 = verifyFloat32(ctx, dm.CellVolt_16, "CellVolt_16")
	dm.CellVolt_17 = verifyFloat32(ctx, dm.CellVolt_17, "CellVolt_17")
	dm.CellVolt_18 = verifyFloat32(ctx, dm.CellVolt_18, "CellVolt_18")
	dm.CellVolt_19 = verifyFloat32(ctx, dm.CellVolt_19, "CellVolt_19")
	dm.CellVolt_20 = verifyFloat32(ctx, dm.CellVolt_20, "CellVolt_20")
	dm.CellVolt_21 = verifyFloat32(ctx, dm.CellVolt_21, "CellVolt_21")
	dm.CellVolt_22 = verifyFloat32(ctx, dm.CellVolt_22, "CellVolt_22")
	dm.CellVolt_23 = verifyFloat32(ctx, dm.CellVolt_23, "CellVolt_23")
	dm.CellVolt_24 = verifyFloat32(ctx, dm.CellVolt_24, "CellVolt_24")
	dm.CellVolt_25 = verifyFloat32(ctx, dm.CellVolt_25, "CellVolt_25")
	dm.CellVolt_26 = verifyFloat32(ctx, dm.CellVolt_26, "CellVolt_26")
	dm.CellVolt_27 = verifyFloat32(ctx, dm.CellVolt_27, "CellVolt_27")
	dm.CellVolt_28 = verifyFloat32(ctx, dm.CellVolt_28, "CellVolt_28")
	dm.CellVolt_29 = verifyFloat32(ctx, dm.CellVolt_29, "CellVolt_29")
	dm.CellVolt_30 = verifyFloat32(ctx, dm.CellVolt_30, "CellVolt_30")
	dm.CellVolt_31 = verifyFloat32(ctx, dm.CellVolt_31, "CellVolt_31")
	dm.CellVoltMin = verifyFloat32(ctx, dm.CellVoltMin, "CellVoltMin")
	dm.CellVoltMax = verifyFloat32(ctx, dm.CellVoltMax, "CellVoltMax")
	dm.CellVoltAvg = verifyFloat32(ctx, dm.CellVoltAvg, "CellVoltAvg")
	dm.CellVoltDiff = verifyFloat32(ctx, dm.CellVoltDiff, "CellVoltDiff")

	dm.BmsLeftMotorControllerCurrent = verifyFloat32(ctx, dm.BmsLeftMotorControllerCurrent, "BmsLeftMotorControllerCurrent")
	dm.BmsRightMotorControllerCurrent = verifyFloat32(ctx, dm.BmsRightMotorControllerCurrent, "BmsRightMotorControllerCurrent")
	dm.BmsMotorControllerCurrentSum = verifyFloat32(ctx, dm.BmsMotorControllerCurrentSum, "BmsMotorControllerCurrentSum")
	dm.PacketsPerSec = verifyFloat32(ctx, dm.PacketsPerSec, "PacketsPerSec")
}

func deserializeProto(ctx context.Context, packet []byte, handlePadding bool) (*SundaeDataMessage, error) {
	// Deserialize ProtoBuf using sundae proto
	dMsg := SundaeDataMessage{}
	err := proto.Unmarshal(packet, &dMsg)
	if err != nil {
		if handlePadding {
			log.Event(ctx, "SLOW: Trying to deserialize padded protobuf")
			for i := len(packet); i >= 10; i-- {
				err = proto.Unmarshal(packet[0:i], &dMsg)
				if err == nil {
					fmt.Println(i)
					return &dMsg, nil
				}
			}
		}
		return nil, err
	}
	return &dMsg, nil
}

// Deserialize unpacks a sundae protobuf, verifies that the fields are valid, then derives any values as needed
func Deserialize(ctx context.Context, rawEvent events.RawEvent) (events.DataEvent, error) {

	dMsg, err := deserializeProto(ctx, rawEvent.Data, tryToHandlePadding)
	if err != nil {
		log.Error(ctx, err, "Could not deserialize protobuf")
		return events.DataEvent{}, err
	}
	// Verify the the proto and clean up data
	VerifyFloatValues(ctx, dMsg)

	dataFields := structs.Map(dMsg)

	return events.DataEvent{
		EventMeta: rawEvent.EventMeta,
		Data:      dataFields,
	}, nil
}
