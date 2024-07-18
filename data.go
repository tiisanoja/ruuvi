package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

var (
	httpURL = os.Getenv("HTTP_URL")
)

type ValidDataDef struct {
	Temp          bool
	Humidity      bool
	Pressure      bool
	Battery       bool
	Address       bool
	AccelerationX bool
	AccelerationY bool
	AccelerationZ bool
	TimeStamp     bool
	Seq           bool
	TXPower       bool
	MAC           bool
}

// SensorData to be posted
type SensorData struct {
	Temp          float64
	Humidity      float64
	Pressure      uint32
	Battery       uint16
	Address       string
	AccelerationX int16
	AccelerationY int16
	AccelerationZ int16
	TimeStamp     time.Time
	Seq           uint16
	TXPower       int16
	MAC           string

	//Calculated values

	//Absolutely Humidity
	AbsHumidity float64
	DewPoint    float64
	SSI         float64
	Humidex     float64

	//Struct to mark if read data is invalid
	ValidData ValidDataDef
}

type SensorFormat5 struct {
	ManufacturerID   uint16
	DataFormat       uint8
	Temperature      int16
	Humidity         uint16
	Pressure         uint16
	AccelerationX    int16
	AccelerationY    int16
	AccelerationZ    int16
	BatteryVoltageMv uint16
	Movement         uint8
	Sequence         uint16
	MAC1             uint8
	MAC2             uint8
	MAC3             uint8
	MAC4             uint8
	MAC5             uint8
	MAC6             uint8
}

func parseTemperature(t uint8, f uint8) float64 {
	var mask uint8
	mask = (1 << 7)
	isNegative := (t & mask) > 0
	temp := float64(t&^mask) + float64(f)/100.0
	if isNegative {
		temp *= -1
	}
	return temp
}

var mWait map[string]bool = make(map[string]bool)
var mutex = &sync.Mutex{}

// Parses sensorData from binary data sent by RuuviTag over BLE
// Calculates also few values
func parseSensorFormat5(data []byte) (error, SensorData) {
	reader := bytes.NewReader(data)
	result := SensorFormat5{}
	sensorData := SensorData{}

	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return err, sensorData
	}

	sensorData.Temp = float64(result.Temperature) * 0.005
	sensorData.Humidity = float64(result.Humidity) * 0.0025
	sensorData.Pressure = uint32(result.Pressure) + uint32(50000+PressureCorrection)
	bat := result.BatteryVoltageMv & 0xffe0
	bat = ((bat >> 5) + 1600)
	sensorData.Battery = bat
	sensorData.TXPower = int16(0x001F&result.BatteryVoltageMv)*2 - 40
	sensorData.AccelerationX = result.AccelerationX
	sensorData.AccelerationY = result.AccelerationY
	sensorData.AccelerationZ = result.AccelerationZ
	sensorData.Seq = result.Sequence
	sensorData.MAC = fmt.Sprintf("%x:%x:%x:%x:%x:%x", result.MAC1, result.MAC2, result.MAC3, result.MAC4, result.MAC5, result.MAC6)

	//Calculate values
	sensorData.AbsHumidity = calculateAbsHumidity(sensorData.Temp, sensorData.Humidity)
	sensorData.DewPoint = calculateDewPoint(sensorData.Temp, sensorData.Humidity)
	sensorData.SSI = calculateSSI(sensorData.Temp, sensorData.Humidity)
	sensorData.Humidex = calculateHumidex(sensorData.Temp, sensorData.DewPoint)

	sensorData = checkNAN(sensorData)

	return err, sensorData

}

// Check if RuuviTag returns error value (Not a Number - NAN)
func checkNAN(sensorData SensorData) SensorData {
	//Error values (NAN)
	TempNAN := -163.84
	HumidityNAN := 163.8375
	PressureNAN := uint32(115535)
	TXPowerNAN := int16(22)
	BatteryNAN := uint16(3647)
	SeqNAN := uint16(65535)
	AccelerationXNAN := int16(-32768)
	AccelerationYNAN := int16(-32768)
	AccelerationZNAN := int16(-32768)

	sensorData.ValidData.MAC = true

	//Check if Temperature is having Error value
	if sensorData.Temp == TempNAN {
		sensorData.ValidData.Temp = false
		log.Printf("ERROR: RuuviTag's %s temperature is having errornous value\n", sensorData.MAC)
	} else {
		sensorData.ValidData.Temp = true
	}

	//Check if humidity is having Error value
	if sensorData.Humidity == HumidityNAN {
		sensorData.ValidData.Humidity = false
		log.Printf("ERROR: RuuviTag's %s humidity is having errornous value\n", sensorData.MAC)
	} else {
		sensorData.ValidData.Humidity = true
	}

	//Check if pressure is having Error value
	if sensorData.Pressure == PressureNAN {
		sensorData.ValidData.Pressure = false
		log.Printf("ERROR: RuuviTag's %s pressure is having errornous value\n", sensorData.MAC)
	} else {
		sensorData.ValidData.Pressure = true
	}

	//Check if Transmit Power is having Error value
	if sensorData.TXPower == TXPowerNAN {
		sensorData.ValidData.TXPower = false
		log.Printf("ERROR: RuuviTag's %s transmit power is having errornous value\n", sensorData.MAC)
	} else {
		sensorData.ValidData.TXPower = true
	}

	//Check if battery is having Error value
	if sensorData.Battery == BatteryNAN {
		sensorData.ValidData.Battery = false
		log.Printf("ERROR: RuuviTag's %s battery voltage is having errornous value\n", sensorData.MAC)
	} else {
		sensorData.ValidData.Battery = true
	}

	//Check if sequence is having Error value
	if sensorData.Seq == SeqNAN {
		sensorData.ValidData.Seq = false
		log.Printf("ERROR: RuuviTag's %s sequence is having errornous value\n", sensorData.MAC)
	} else {
		sensorData.ValidData.Seq = true
	}

	//Check if AccelerationX is having Error value
	if sensorData.AccelerationX == AccelerationXNAN {
		sensorData.ValidData.AccelerationX = false
		log.Printf("ERROR: RuuviTag's %s accelerationX is having errornous value\n", sensorData.MAC)
	} else {
		sensorData.ValidData.AccelerationX = true
	}

	//Check if AccelerationY is having Error value
	if sensorData.AccelerationY == AccelerationYNAN {
		sensorData.ValidData.AccelerationY = false
		log.Printf("ERROR: RuuviTag's %s accelerationY is having errornous value\n", sensorData.MAC)
	} else {
		sensorData.ValidData.AccelerationY = true
	}

	//Check if acclerationZ is having Error value
	if sensorData.AccelerationZ == AccelerationZNAN {
		sensorData.ValidData.AccelerationZ = false
		log.Printf("ERROR: RuuviTag's %s accelerationZ is having errornous value\n", sensorData.MAC)
	} else {
		sensorData.ValidData.AccelerationZ = true
	}

	return sensorData
}

// Calculates absolutely humidity approximation based on temperature and humidity%
//
// temp    : Temperature needs to be given in °C
// humidity: Humidity is humidity in percent (0-100)
//
// return value: Approximation of absolute Humidity (g/m3)
//
// Abolutely humidity is calculated using Bolton formula for steam saturated pressure
// Note! Returned value is approximation and has error. See links for detail for error.
// Also measurements has error which are effecting to result of approximation of absolutely humidity
//
// Calculated absolutely humidity is WITHOUT ANY WARRANTY!
func calculateAbsHumidity(temp float64, humidity float64) float64 {
	//Change Temperature to Kelvin
	tempInK := temp + 273.16

	//Calculates steam saturation pressure using Bolton (1980) formula
	//source for the formula: http://meteorologia.uib.eu/ROMU/formal/relative_humidity/relative_humidity.pdf
	//See from the link the precion for the formula. Result is not exact correct value
	steamSaturationPressure := 6.112 * math.Pow(math.E, ((17.67*temp)/(temp+243.5)))

	//absolutely humidity is calculated using formula found:
	//https://carnotcycle.wordpress.com/2012/08/04/how-to-convert-relative-humidity-to-absolute-humidity/comment-page-1/
	absHumidity := (humidity * steamSaturationPressure * 18.02) / (tempInK * 100 * 0.08314)
	return absHumidity
}

// Calculates approximation of dew point based on temperature and humidity%
//
// temp    : Temperature needs to be given in °C
// humidity: Humidity is humidity percent
//
// return value: Approximation of dew point in °C
//
// See good information source for dew point: https://en.wikipedia.org/wiki/Dew_point
// A well-known approximation formula is used to calculate the dew point. Formula can be found from mentioned wikipedia page.
//
// Below text is taken from https://en.wikipedia.org/wiki/Dew_point
//
// These valuations provide a maximum error of 0.1%, for −30 °C ≤ T ≤ 35°C and 1% < RH < 100%. Also noteworthy is the Sonntag1990,[17]
//
//	a = 6.112 mbar, b = 17.62, c = 243.12 °C; for −45 °C ≤ T ≤ 60 °C (error ±0.35 °C).
//
// Please note! Temperature and humidity has error in measurements so the total error for dew point is higher than mentioned above
//
// Calculated dew point is WITHOUT ANY WARRANTY!
func calculateDewPoint(temp float64, humidity float64) float64 {

	//Change on 20220806: b and c value have been changed from b=17.27, c=237.7
	//to work better with below zero temperatures
	b := 17.62
	c := 243.12
	tTemp := ((b * temp) / (c + temp)) + math.Log(humidity/100)
	tDewpoint := 0.0
	if tTemp != b {
		tDewpoint = (c * tTemp) / (b - tTemp)
	} else {
		log.Printf("INFO: In tDewpoint: Dividing by zero. Adding 0.01 to temperature.\n")
		tDewpoint = (c * (tTemp + 0.01)) / (b - tTemp + 0.01)
	}
	return tDewpoint
}

// Calculate Summer Simmer Index (SSI)
//
// SSI is heat index used in Finland by FMI.
// Formula is taken from https://github.com/fmidev/smartmet-library-newbase/blob/master/newbase/NFmiMetMath.cpp#L418
func calculateSSI(temp float64, humidity float64) float64 {

	// SSI is calculated only to temperatures over 14.5C
	if temp < 14.5 {
		return temp
	}

	// When it is > 14.5 degrees in Finland, 60% is approximately
	// the minimum mean monthly humidity. However, it seems that
	// most humans feel most comfortable either at 45%, or
	// alternatively somewhere between 50-60%. Hence we choose
	// the middle ground 50%
	// Source: https://github.com/fmidev/smartmet-library-newbase/blob/master/newbase/NFmiMetMath.cpp#L418

	rhRef := 50.0 / 100.0
	rh := humidity / 100.0

	ssi := (1.8*temp - 0.55*(1-rh)*(1.8*temp-26) - 0.55*(1-rhRef)*26) / (1.8 * (1 - 0.55*(1-rhRef)))
	return ssi
}

//
// Calculate Humidex
//
// https://en.wikipedia.org/wiki/Humidex
//
// The humidex (short for humidity index) is an index number used by Canadian meteorologists
// to describe how hot the weather feels to the average person,
// by combining the effect of heat and humidity.

func calculateHumidex(temp float64, dewPoint float64) float64 {
	a := (1 / 273.15) - (1 / (273.15 + dewPoint))
	b := a * 5417.7530
	humidex := temp + 0.5555*(6.11*math.Exp(b)-10)
	return humidex
}

// Tries to lock RuuviTag sensor
// If MAC address is already locked returns false
// If locking was successful returns true
func lockSensor(MAC string) bool {
	mutex.Lock()
	if mWait[MAC] == true {
		mutex.Unlock()
		log.Printf("INFO: Skipping %s\n", MAC)
		return false
	}
	mWait[MAC] = true
	mutex.Unlock()
	return true
}

// Release lock from the sensor
func releaseSensor(MAC string) bool {
	//Let's not store all measurements per device
	if StoreDelay > 0 {
		time.Sleep(StoreDelay)
	}

	mutex.Lock()
	mWait[MAC] = false
	mutex.Unlock()

	return true
}

// ParseRuuviData parses raw ruuvidata received from RuuviTag
func ParseRuuviData(data []byte, a string) {

	if len(data) > 2 && binary.LittleEndian.Uint16(data[0:2]) == 0x0499 {
		if len(data) > 25 {
			sensorFormat := data[2]
			log.Printf("Ruuvi data with sensor format %d\n", sensorFormat)
			switch sensorFormat {
			case 3:
				log.Printf("RuuviTag version 3 not supported. Please upgrade RuuviTag to version 5.")
			case 5:
				err, sensorData := parseSensorFormat5(data)
				if err == nil {
					if lockSensor(sensorData.MAC) == false {
						return
					}

					log.Printf("Seq: %d, Temp: %f, Pres: %d, Hum: %f, Battery: %d, TXPower: %d, MAC: %s\n",
						sensorData.Seq,
						sensorData.Temp,
						sensorData.Pressure,
						sensorData.Humidity,
						sensorData.Battery,
						sensorData.TXPower,
						sensorData.MAC)

					StoreMeasurement(sensorData)
					StoreHWmeasurement(sensorData)
					releaseSensor(sensorData.MAC)
				}
			case 8:
				log.Printf("RuuviTag version 8 is having encrypted data which is not supported.")
			default:
				log.Printf("Unknown sensor format %d\n", sensorFormat)
			}
		} else {
			log.Printf("Incorrect lenght in Ruuvi device data: %d. Expected to be 26.\n", len(data))
		}

	} else {
		//fmt.Printf("Not a RuuviTag device \n")
	}

}
