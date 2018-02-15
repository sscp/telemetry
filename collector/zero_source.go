package collector

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	sscpproto "github.com/sscp/telemetry/proto"
	"golang.org/x/time/rate"
)

// ZeroPacketSource is a PacketSource that returns only zeroed out DataMessages
// at a given rate
type ZeroPacketSource struct {
	outChan  chan *ContextPacket
	doneChan chan bool
	limiter  *rate.Limiter
}

// Packets is the stream of zeroed binary packets
// It is simply a reference to outChan
func (zps *ZeroPacketSource) Packets() <-chan *ContextPacket {
	return zps.outChan
}

// Listen begins sending zeroed packets to the Packets channel.
// It launches a gorountine that sen
func (zps *ZeroPacketSource) Listen() {
	go func() {
		for {
			select {
			case <-zps.doneChan:
				return
			default:
				err := zps.limiter.Wait(context.TODO())
				if err == nil {
					zPacket, _ := CreateZeroPacket()
					zps.outChan <- &ContextPacket{
						ctx:    context.TODO(),
						packet: zPacket,
					}
				} else {
					fmt.Println("too fast")
				}
			}
		}
	}()
}

func (zps *ZeroPacketSource) Close() {
	zps.doneChan <- true
	close(zps.doneChan)
	close(zps.outChan)
}

func NewZeroPacketSource(packetsPerSecond int) PacketSource {
	return &ZeroPacketSource{
		outChan:  make(chan *ContextPacket),
		doneChan: make(chan bool),
		// Only allow one packet out at a time
		limiter: rate.NewLimiter(rate.Limit(packetsPerSecond), 1),
	}
}

func CreateZeroPacket() ([]byte, error) {
	zdm := CreateZeroDataMessage()
	return proto.Marshal(zdm)
}

func CreateZeroDataMessage() *sscpproto.DataMessage {
	Uint32 := uint32(0)
	Int32 := int32(0)
	Int64 := int64(0)
	Float32 := float32(0.0)

	zdm := sscpproto.DataMessage{
		RegenEnabled:                     &Uint32,
		RegenCommand:                     &Uint32,
		BatteryPower:                     &Float32,
		ArrayPower:                       &Float32,
		ReverseOn:                        &Uint32,
		LowVoltPower:                     &Float32,
		HazardLightsOn:                   &Int32,
		BatteryVoltage:                   &Float32,
		Ltc6804Badpec:                    &Uint32,
		BmsState:                         &Uint32,
		ChargeEnabled:                    &Int32,
		DischargeEnabled:                 &Int32,
		HighsideContactorOn:              &Int32,
		LowsideContactorOn:               &Int32,
		PrechargeOn:                      &Int32,
		LowVoltBusOn:                     &Int32,
		BatteryTemp_1:                    &Float32,
		BatteryTemp_2:                    &Float32,
		BatteryTemp_3:                    &Float32,
		BatteryTemp_4:                    &Float32,
		BatteryTemp_5:                    &Float32,
		BatteryTemp_6:                    &Float32,
		BmsPrechargeBatteryAdc:           &Float32,
		BmsPrechargeCarAdc:               &Float32,
		LowVoltOutputCurrent:             &Float32,
		BatteryCurrent:                   &Float32,
		RightMotorControllerPower:        &Float32,
		AmpHours:                         &Float32,
		HeadLightsOn:                     &Int32,
		BrakeLightsOn:                    &Int32,
		RightBlinkerOn:                   &Int32,
		LeftBlinkerOn:                    &Int32,
		BrakePressed:                     &Int32,
		ThrottlePressed:                  &Int32,
		DriveMode:                        &Uint32,
		MotorControllerEnabled:           &Int32,
		MotorControllerSpeed:             &Float32,
		MotorControllerRpm:               &Float32,
		AvgOdometer:                      &Float32,
		LeftMotorTemp:                    &Float32,
		RightMotorTemp:                   &Float32,
		LeftMotorControllerTemp:          &Float32,
		RightMotorControllerTemp:         &Float32,
		LeftMotorControllerAlive:         &Float32,
		RightMotorControllerAlive:        &Float32,
		LeftMotorControllerCurrent:       &Float32,
		RightMotorControllerCurrent:      &Float32,
		MotorControllerCurrentDiff:       &Float32,
		LeftMotorControllerError:         &Uint32,
		RightMotorControllerError:        &Uint32,
		LeftMotorControllerLimit:         &Uint32,
		RightMotorControllerLimit:        &Uint32,
		LeftMotorControllerRxErrorCount:  &Uint32,
		RightMotorControllerRxErrorCount: &Uint32,
		LeftMotorControllerTxErrorCount:  &Uint32,
		RightMotorControllerTxErrorCount: &Uint32,
		LeftMotorControllerBusVoltage:    &Float32,
		RightMotorControllerBusVoltage:   &Float32,
		LeftMotorController_15VVoltage:   &Float32,
		RightMotorController_15VVoltage:  &Float32,
		LeftMotorController_3V3Voltage:   &Float32,
		RightMotorController_3V3Voltage:  &Float32,
		LeftMotorController_1V9Voltage:   &Float32,
		RightMotorController_1V9Voltage:  &Float32,
		LeftMotorControllerDspTemp:       &Float32,
		RightMotorControllerDspTemp:      &Float32,
		LeftMotorControllerPhaseCurrent:  &Float32,
		RightMotorControllerPhaseCurrent: &Float32,
		LeftMotorRpmCommand:              &Float32,
		RightMotorRpmCommand:             &Float32,
		LeftMotorCurrentCommand:          &Float32,
		RightMotorCurrentCommand:         &Float32,
		GpsTime:                          &Int64,
		GpsLatitude:                      &Float32,
		GpsLongitude:                     &Float32,
		GpsSpeed:                         &Float32,
		GpsAltitude:                      &Float32,
		GpsBearing:                       &Float32,
		LedState:                         &Int32,
		MpptArrayPower:                   &Float32,
		Mppt_A0VoltIn:                    &Float32,
		Mppt_A0VoltOut:                   &Float32,
		Mppt_A0Current:                   &Float32,
		Mppt_A0Temp:                      &Float32,
		Mppt_A1VoltIn:                    &Float32,
		Mppt_A1VoltOut:                   &Float32,
		Mppt_A1Current:                   &Float32,
		Mppt_A1Temp:                      &Float32,
		Mppt_B0VoltIn:                    &Float32,
		Mppt_B0VoltOut:                   &Float32,
		Mppt_B0Current:                   &Float32,
		Mppt_B0Temp:                      &Float32,
		Mppt_B1VoltIn:                    &Float32,
		Mppt_B1VoltOut:                   &Float32,
		Mppt_B1Current:                   &Float32,
		Mppt_B1Temp:                      &Float32,
		Mppt_C0VoltIn:                    &Float32,
		Mppt_C0VoltOut:                   &Float32,
		Mppt_C0Current:                   &Float32,
		Mppt_C0Temp:                      &Float32,
		Mppt_C1VoltIn:                    &Float32,
		Mppt_C1VoltOut:                   &Float32,
		Mppt_C1Current:                   &Float32,
		Mppt_C1Temp:                      &Float32,
		Mppt_D0VoltIn:                    &Float32,
		Mppt_D0VoltOut:                   &Float32,
		Mppt_D0Current:                   &Float32,
		Mppt_D0Temp:                      &Float32,
		Mppt_D1VoltIn:                    &Float32,
		Mppt_D1VoltOut:                   &Float32,
		Mppt_D1Current:                   &Float32,
		Mppt_D1Temp:                      &Float32,
		CellVolt_1:                       &Float32,
		CellVolt_2:                       &Float32,
		CellVolt_3:                       &Float32,
		CellVolt_4:                       &Float32,
		CellVolt_5:                       &Float32,
		CellVolt_6:                       &Float32,
		CellVolt_7:                       &Float32,
		CellVolt_8:                       &Float32,
		CellVolt_9:                       &Float32,
		CellVolt_10:                      &Float32,
		CellVolt_11:                      &Float32,
		CellVolt_12:                      &Float32,
		CellVolt_13:                      &Float32,
		CellVolt_14:                      &Float32,
		CellVolt_15:                      &Float32,
		CellVolt_16:                      &Float32,
		CellVolt_17:                      &Float32,
		CellVolt_18:                      &Float32,
		CellVolt_19:                      &Float32,
		CellVolt_20:                      &Float32,
		CellVolt_21:                      &Float32,
		CellVolt_22:                      &Float32,
		CellVolt_23:                      &Float32,
		CellVolt_24:                      &Float32,
		CellVolt_25:                      &Float32,
		CellVolt_26:                      &Float32,
		CellVolt_27:                      &Float32,
		CellVolt_28:                      &Float32,
		CellVolt_29:                      &Float32,
		CellVolt_30:                      &Float32,
		CellVolt_31:                      &Float32,
		CellVoltMin:                      &Float32,
		CellVoltMax:                      &Float32,
		CellVoltAvg:                      &Float32,
		CellVoltDiff:                     &Float32,
		PowerSaveOn:                      &Int32,
		RearviewOn:                       &Int32,
		MicEnabled:                       &Int32,
		ImuTemp:                          &Int32,
		ImuMagnetX:                       &Int32,
		ImuMagnetY:                       &Int32,
		ImuMagnetZ:                       &Int32,
		ImuGyroX:                         &Int32,
		ImuGyroY:                         &Int32,
		ImuGyroZ:                         &Int32,
		ImuAccelX:                        &Int32,
		ImuAccelY:                        &Int32,
		ImuAccelZ:                        &Int32,
		BmsLeftMotorControllerCurrent:    &Float32,
		BmsRightMotorControllerCurrent:   &Float32,
		BmsMotorControllerCurrentSum:     &Float32,
		PacketsPerSec:                    &Float32,
	}
	return &zdm
}
