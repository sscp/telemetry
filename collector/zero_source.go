package collector

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	sundaeproto "github.com/sscp/telemetry/collector/sundae"

	"github.com/golang/protobuf/proto"
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
					recievedTime := time.Now()
					// Create context with time of receiving packet
					ctx := ContextWithRecievedTime(context.Background(), recievedTime)

					zps.outChan <- &ContextPacket{
						ctx:    ctx,
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

func zeroUInt32() *uint32 {
	z := uint32(0)
	return &z
}

func zeroUInt64() *uint64 {
	z := uint64(0)
	return &z
}

func zeroInt32() *int32 {
	z := int32(0)
	return &z
}

func zeroInt64() *int64 {
	z := int64(0)
	return &z
}

func zeroFloat32() *float32 {
	z := float32(0)
	return &z
}

func randUInt32() *uint32 {
	z := uint32(rand.Intn(100))
	return &z
}

func randUInt64() *uint64 {
	z := uint64(rand.Intn(100))
	return &z
}

func randInt32() *int32 {
	z := int32(rand.Intn(100))
	return &z
}

func randInt64() *int64 {
	z := int64(rand.Intn(100))
	return &z
}

func randFloat32() *float32 {
	z := rand.Float32() * 100
	return &z
}

func CreateZeroDataMessage() *sundaeproto.DataMessage {
	zdm := sundaeproto.DataMessage{
		RegenEnabled:                     zeroUInt32(),
		RegenCommand:                     zeroUInt32(),
		BatteryPower:                     zeroUInt32(),
		ArrayPower:                       zeroFloat32(),
		ReverseOn:                        zeroUInt32(),
		LowVoltPower:                     zeroFloat32(),
		HazardLightsOn:                   zeroInt32(),
		BatteryVoltage:                   zeroFloat32(),
		Ltc6804Badpec:                    zeroUInt32(),
		BmsState:                         zeroUInt32(),
		ChargeEnabled:                    zeroInt32(),
		DischargeEnabled:                 zeroInt32(),
		HighsideContactorOn:              zeroInt32(),
		LowsideContactorOn:               zeroInt32(),
		PrechargeOn:                      zeroInt32(),
		LowVoltBusOn:                     zeroInt32(),
		BatteryTemp_1:                    zeroFloat32(),
		BatteryTemp_2:                    zeroFloat32(),
		BatteryTemp_3:                    zeroFloat32(),
		BatteryTemp_4:                    zeroFloat32(),
		BatteryTemp_5:                    zeroFloat32(),
		BatteryTemp_6:                    zeroFloat32(),
		BmsPrechargeBatteryAdc:           zeroFloat32(),
		BmsPrechargeCarAdc:               zeroFloat32(),
		LowVoltOutputCurrent:             zeroFloat32(),
		BatteryCurrent:                   zeroFloat32(),
		RightMotorControllerPower:        zeroFloat32(),
		AmpHours:                         zeroFloat32(),
		HeadLightsOn:                     zeroInt32(),
		BrakeLightsOn:                    zeroInt32(),
		RightBlinkerOn:                   zeroInt32(),
		LeftBlinkerOn:                    zeroInt32(),
		BrakePressed:                     zeroInt32(),
		ThrottlePressed:                  zeroInt32(),
		DriveMode:                        zeroUInt32(),
		MotorControllerEnabled:           zeroInt32(),
		MotorControllerSpeed:             zeroFloat32(),
		MotorControllerRpm:               zeroFloat32(),
		AvgOdometer:                      zeroFloat32(),
		LeftMotorTemp:                    zeroFloat32(),
		RightMotorTemp:                   zeroFloat32(),
		LeftMotorControllerTemp:          zeroFloat32(),
		RightMotorControllerTemp:         zeroFloat32(),
		LeftMotorControllerAlive:         zeroFloat32(),
		RightMotorControllerAlive:        zeroFloat32(),
		LeftMotorControllerCurrent:       zeroFloat32(),
		RightMotorControllerCurrent:      zeroFloat32(),
		MotorControllerCurrentDiff:       zeroFloat32(),
		LeftMotorControllerError:         zeroUInt32(),
		RightMotorControllerError:        zeroUInt32(),
		LeftMotorControllerLimit:         zeroUInt32(),
		RightMotorControllerLimit:        zeroUInt32(),
		LeftMotorControllerRxErrorCount:  zeroUInt32(),
		RightMotorControllerRxErrorCount: zeroUInt32(),
		LeftMotorControllerTxErrorCount:  zeroUInt32(),
		RightMotorControllerTxErrorCount: zeroUInt32(),
		LeftMotorControllerBusVoltage:    zeroFloat32(),
		RightMotorControllerBusVoltage:   zeroFloat32(),
		LeftMotorController_15VVoltage:   zeroFloat32(),
		RightMotorController_15VVoltage:  zeroFloat32(),
		LeftMotorController_3V3Voltage:   zeroFloat32(),
		RightMotorController_3V3Voltage:  zeroFloat32(),
		LeftMotorController_1V9Voltage:   zeroFloat32(),
		RightMotorController_1V9Voltage:  zeroFloat32(),
		LeftMotorControllerDspTemp:       zeroFloat32(),
		RightMotorControllerDspTemp:      zeroFloat32(),
		LeftMotorControllerPhaseCurrent:  zeroFloat32(),
		RightMotorControllerPhaseCurrent: zeroFloat32(),
		LeftMotorRpmCommand:              zeroFloat32(),
		RightMotorRpmCommand:             zeroFloat32(),
		LeftMotorCurrentCommand:          zeroFloat32(),
		RightMotorCurrentCommand:         zeroFloat32(),
		GpsTime:                          zeroInt64(),
		GpsLatitude:                      zeroFloat32(),
		GpsLongitude:                     zeroFloat32(),
		GpsSpeed:                         zeroFloat32(),
		GpsAltitude:                      zeroFloat32(),
		GpsBearing:                       zeroFloat32(),
		LedState:                         zeroInt32(),
		MpptArrayPower:                   zeroFloat32(),
		Mppt_A0VoltIn:                    zeroFloat32(),
		Mppt_A0VoltOut:                   zeroFloat32(),
		Mppt_A0Current:                   zeroFloat32(),
		Mppt_A0Temp:                      zeroFloat32(),
		Mppt_A1VoltIn:                    zeroFloat32(),
		Mppt_A1VoltOut:                   zeroFloat32(),
		Mppt_A1Current:                   zeroFloat32(),
		Mppt_A1Temp:                      zeroFloat32(),
		Mppt_B0VoltIn:                    zeroFloat32(),
		Mppt_B0VoltOut:                   zeroFloat32(),
		Mppt_B0Current:                   zeroFloat32(),
		Mppt_B0Temp:                      zeroFloat32(),
		Mppt_B1VoltIn:                    zeroFloat32(),
		Mppt_B1VoltOut:                   zeroFloat32(),
		Mppt_B1Current:                   zeroFloat32(),
		Mppt_B1Temp:                      zeroFloat32(),
		Mppt_C0VoltIn:                    zeroFloat32(),
		Mppt_C0VoltOut:                   zeroFloat32(),
		Mppt_C0Current:                   zeroFloat32(),
		Mppt_C0Temp:                      zeroFloat32(),
		Mppt_C1VoltIn:                    zeroFloat32(),
		Mppt_C1VoltOut:                   zeroFloat32(),
		Mppt_C1Current:                   zeroFloat32(),
		Mppt_C1Temp:                      zeroFloat32(),
		Mppt_D0VoltIn:                    zeroFloat32(),
		Mppt_D0VoltOut:                   zeroFloat32(),
		Mppt_D0Current:                   zeroFloat32(),
		Mppt_D0Temp:                      zeroFloat32(),
		Mppt_D1VoltIn:                    zeroFloat32(),
		Mppt_D1VoltOut:                   zeroFloat32(),
		Mppt_D1Current:                   zeroFloat32(),
		Mppt_D1Temp:                      zeroFloat32(),
		CellVolt_1:                       zeroFloat32(),
		CellVolt_2:                       zeroFloat32(),
		CellVolt_3:                       zeroFloat32(),
		CellVolt_4:                       zeroFloat32(),
		CellVolt_5:                       zeroFloat32(),
		CellVolt_6:                       zeroFloat32(),
		CellVolt_7:                       zeroFloat32(),
		CellVolt_8:                       zeroFloat32(),
		CellVolt_9:                       zeroFloat32(),
		CellVolt_10:                      zeroFloat32(),
		CellVolt_11:                      zeroFloat32(),
		CellVolt_12:                      zeroFloat32(),
		CellVolt_13:                      zeroFloat32(),
		CellVolt_14:                      zeroFloat32(),
		CellVolt_15:                      zeroFloat32(),
		CellVolt_16:                      zeroFloat32(),
		CellVolt_17:                      zeroFloat32(),
		CellVolt_18:                      zeroFloat32(),
		CellVolt_19:                      zeroFloat32(),
		CellVolt_20:                      zeroFloat32(),
		CellVolt_21:                      zeroFloat32(),
		CellVolt_22:                      zeroFloat32(),
		CellVolt_23:                      zeroFloat32(),
		CellVolt_24:                      zeroFloat32(),
		CellVolt_25:                      zeroFloat32(),
		CellVolt_26:                      zeroFloat32(),
		CellVolt_27:                      zeroFloat32(),
		CellVolt_28:                      zeroFloat32(),
		CellVolt_29:                      zeroFloat32(),
		CellVolt_30:                      zeroFloat32(),
		CellVolt_31:                      zeroFloat32(),
		CellVoltMin:                      zeroFloat32(),
		CellVoltMax:                      zeroFloat32(),
		CellVoltAvg:                      zeroFloat32(),
		CellVoltDiff:                     zeroFloat32(),
		PowerSaveOn:                      zeroInt32(),
		RearviewOn:                       zeroInt32(),
		MicEnabled:                       zeroInt32(),
		ImuTemp:                          zeroInt32(),
		ImuMagnetX:                       zeroInt32(),
		ImuMagnetY:                       zeroInt32(),
		ImuMagnetZ:                       zeroInt32(),
		ImuGyroX:                         zeroInt32(),
		ImuGyroY:                         zeroInt32(),
		ImuGyroZ:                         zeroInt32(),
		ImuAccelX:                        zeroInt32(),
		ImuAccelY:                        zeroInt32(),
		ImuAccelZ:                        zeroInt32(),
		BmsLeftMotorControllerCurrent:    zeroFloat32(),
		BmsRightMotorControllerCurrent:   zeroFloat32(),
		BmsMotorControllerCurrentSum:     zeroFloat32(),
		PacketsPerSec:                    zeroFloat32(),
	}
	return &zdm
}

func CreateRandomDataMessage() *sundaeproto.DataMessage {
	zdm := sundaeproto.DataMessage{
		RegenEnabled:                     randUInt32(),
		RegenCommand:                     randUInt32(),
		BatteryPower:                     randUInt32(),
		ArrayPower:                       randFloat32(),
		ReverseOn:                        randUInt32(),
		LowVoltPower:                     randFloat32(),
		HazardLightsOn:                   randInt32(),
		BatteryVoltage:                   randFloat32(),
		Ltc6804Badpec:                    randUInt32(),
		BmsState:                         randUInt32(),
		ChargeEnabled:                    randInt32(),
		DischargeEnabled:                 randInt32(),
		HighsideContactorOn:              randInt32(),
		LowsideContactorOn:               randInt32(),
		PrechargeOn:                      randInt32(),
		LowVoltBusOn:                     randInt32(),
		BatteryTemp_1:                    randFloat32(),
		BatteryTemp_2:                    randFloat32(),
		BatteryTemp_3:                    randFloat32(),
		BatteryTemp_4:                    randFloat32(),
		BatteryTemp_5:                    randFloat32(),
		BatteryTemp_6:                    randFloat32(),
		BmsPrechargeBatteryAdc:           randFloat32(),
		BmsPrechargeCarAdc:               randFloat32(),
		LowVoltOutputCurrent:             randFloat32(),
		BatteryCurrent:                   randFloat32(),
		RightMotorControllerPower:        randFloat32(),
		AmpHours:                         randFloat32(),
		HeadLightsOn:                     randInt32(),
		BrakeLightsOn:                    randInt32(),
		RightBlinkerOn:                   randInt32(),
		LeftBlinkerOn:                    randInt32(),
		BrakePressed:                     randInt32(),
		ThrottlePressed:                  randInt32(),
		DriveMode:                        randUInt32(),
		MotorControllerEnabled:           randInt32(),
		MotorControllerSpeed:             randFloat32(),
		MotorControllerRpm:               randFloat32(),
		AvgOdometer:                      randFloat32(),
		LeftMotorTemp:                    randFloat32(),
		RightMotorTemp:                   randFloat32(),
		LeftMotorControllerTemp:          randFloat32(),
		RightMotorControllerTemp:         randFloat32(),
		LeftMotorControllerAlive:         randFloat32(),
		RightMotorControllerAlive:        randFloat32(),
		LeftMotorControllerCurrent:       randFloat32(),
		RightMotorControllerCurrent:      randFloat32(),
		MotorControllerCurrentDiff:       randFloat32(),
		LeftMotorControllerError:         randUInt32(),
		RightMotorControllerError:        randUInt32(),
		LeftMotorControllerLimit:         randUInt32(),
		RightMotorControllerLimit:        randUInt32(),
		LeftMotorControllerRxErrorCount:  randUInt32(),
		RightMotorControllerRxErrorCount: randUInt32(),
		LeftMotorControllerTxErrorCount:  randUInt32(),
		RightMotorControllerTxErrorCount: randUInt32(),
		LeftMotorControllerBusVoltage:    randFloat32(),
		RightMotorControllerBusVoltage:   randFloat32(),
		LeftMotorController_15VVoltage:   randFloat32(),
		RightMotorController_15VVoltage:  randFloat32(),
		LeftMotorController_3V3Voltage:   randFloat32(),
		RightMotorController_3V3Voltage:  randFloat32(),
		LeftMotorController_1V9Voltage:   randFloat32(),
		RightMotorController_1V9Voltage:  randFloat32(),
		LeftMotorControllerDspTemp:       randFloat32(),
		RightMotorControllerDspTemp:      randFloat32(),
		LeftMotorControllerPhaseCurrent:  randFloat32(),
		RightMotorControllerPhaseCurrent: randFloat32(),
		LeftMotorRpmCommand:              randFloat32(),
		RightMotorRpmCommand:             randFloat32(),
		LeftMotorCurrentCommand:          randFloat32(),
		RightMotorCurrentCommand:         randFloat32(),
		GpsTime:                          randInt64(),
		GpsLatitude:                      randFloat32(),
		GpsLongitude:                     randFloat32(),
		GpsSpeed:                         randFloat32(),
		GpsAltitude:                      randFloat32(),
		GpsBearing:                       randFloat32(),
		LedState:                         randInt32(),
		MpptArrayPower:                   randFloat32(),
		Mppt_A0VoltIn:                    randFloat32(),
		Mppt_A0VoltOut:                   randFloat32(),
		Mppt_A0Current:                   randFloat32(),
		Mppt_A0Temp:                      randFloat32(),
		Mppt_A1VoltIn:                    randFloat32(),
		Mppt_A1VoltOut:                   randFloat32(),
		Mppt_A1Current:                   randFloat32(),
		Mppt_A1Temp:                      randFloat32(),
		Mppt_B0VoltIn:                    randFloat32(),
		Mppt_B0VoltOut:                   randFloat32(),
		Mppt_B0Current:                   randFloat32(),
		Mppt_B0Temp:                      randFloat32(),
		Mppt_B1VoltIn:                    randFloat32(),
		Mppt_B1VoltOut:                   randFloat32(),
		Mppt_B1Current:                   randFloat32(),
		Mppt_B1Temp:                      randFloat32(),
		Mppt_C0VoltIn:                    randFloat32(),
		Mppt_C0VoltOut:                   randFloat32(),
		Mppt_C0Current:                   randFloat32(),
		Mppt_C0Temp:                      randFloat32(),
		Mppt_C1VoltIn:                    randFloat32(),
		Mppt_C1VoltOut:                   randFloat32(),
		Mppt_C1Current:                   randFloat32(),
		Mppt_C1Temp:                      randFloat32(),
		Mppt_D0VoltIn:                    randFloat32(),
		Mppt_D0VoltOut:                   randFloat32(),
		Mppt_D0Current:                   randFloat32(),
		Mppt_D0Temp:                      randFloat32(),
		Mppt_D1VoltIn:                    randFloat32(),
		Mppt_D1VoltOut:                   randFloat32(),
		Mppt_D1Current:                   randFloat32(),
		Mppt_D1Temp:                      randFloat32(),
		CellVolt_1:                       randFloat32(),
		CellVolt_2:                       randFloat32(),
		CellVolt_3:                       randFloat32(),
		CellVolt_4:                       randFloat32(),
		CellVolt_5:                       randFloat32(),
		CellVolt_6:                       randFloat32(),
		CellVolt_7:                       randFloat32(),
		CellVolt_8:                       randFloat32(),
		CellVolt_9:                       randFloat32(),
		CellVolt_10:                      randFloat32(),
		CellVolt_11:                      randFloat32(),
		CellVolt_12:                      randFloat32(),
		CellVolt_13:                      randFloat32(),
		CellVolt_14:                      randFloat32(),
		CellVolt_15:                      randFloat32(),
		CellVolt_16:                      randFloat32(),
		CellVolt_17:                      randFloat32(),
		CellVolt_18:                      randFloat32(),
		CellVolt_19:                      randFloat32(),
		CellVolt_20:                      randFloat32(),
		CellVolt_21:                      randFloat32(),
		CellVolt_22:                      randFloat32(),
		CellVolt_23:                      randFloat32(),
		CellVolt_24:                      randFloat32(),
		CellVolt_25:                      randFloat32(),
		CellVolt_26:                      randFloat32(),
		CellVolt_27:                      randFloat32(),
		CellVolt_28:                      randFloat32(),
		CellVolt_29:                      randFloat32(),
		CellVolt_30:                      randFloat32(),
		CellVolt_31:                      randFloat32(),
		CellVoltMin:                      randFloat32(),
		CellVoltMax:                      randFloat32(),
		CellVoltAvg:                      randFloat32(),
		CellVoltDiff:                     randFloat32(),
		PowerSaveOn:                      randInt32(),
		RearviewOn:                       randInt32(),
		MicEnabled:                       randInt32(),
		ImuTemp:                          randInt32(),
		ImuMagnetX:                       randInt32(),
		ImuMagnetY:                       randInt32(),
		ImuMagnetZ:                       randInt32(),
		ImuGyroX:                         randInt32(),
		ImuGyroY:                         randInt32(),
		ImuGyroZ:                         randInt32(),
		ImuAccelX:                        randInt32(),
		ImuAccelY:                        randInt32(),
		ImuAccelZ:                        randInt32(),
		BmsLeftMotorControllerCurrent:    randFloat32(),
		BmsRightMotorControllerCurrent:   randFloat32(),
		BmsMotorControllerCurrentSum:     randFloat32(),
		PacketsPerSec:                    randFloat32(),
	}
	return &zdm
}
