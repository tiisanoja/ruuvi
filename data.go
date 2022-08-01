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
    Temp       bool
    Humidity   bool
    Pressure   bool
    Battery    bool
    Address    bool
    AccelerationX bool
    AccelerationY bool
    AccelerationZ bool
    TimeStamp  bool
    Seq        bool
    TXPower    bool
    MAC        bool
}


//SensorData to be posted
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
   AbsHumidity    float64
   DewPoint       float64

   //Struct to mark if read data is invalid
   ValidData    ValidDataDef
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

//Parses sensorData from binary data sent by RuuviTag over BLE
//Calculates also few values
func parseSensorFormat5(data []byte) (error, SensorData)  {
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

    sensorData = checkNAN(sensorData)

    return err, sensorData

}

//Check if RuuviTag returns error value (Not a Number - NAN) 
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

//Calculates absolutely humidity based on temperature and humidity%
func calculateAbsHumidity(temp float64, humidity float64) float64 {
    steamSaturationPressure := 6.1078 * math.Pow(10, (7.5*temp/(temp+237.3)))
    absHumidity := (216.679 * (humidity * steamSaturationPressure) / 100) / (temp + 273.15)
    return absHumidity
}

//Calculates dew point based on temperature and humidity%
func calculateDewPoint(temp float64, humidity float64) float64 {

    tTemp := ((17.27 * temp) / (237.7 + temp)) + math.Log(humidity/100)
    tDewpoint := 0.0
    if tTemp != 17.27 {
        tDewpoint = (237.7 * tTemp) / (17.27 - tTemp)
    } else {
        log.Printf("INFO: In tDewpoint: Dividing by zero. Adding 0.01 to temperature.\n")
        tDewpoint = (237.7 * (tTemp + 0.01)) / (17.27 - tTemp + 0.01)
    }
    return tDewpoint
}


//Tries to lock RuuviTag sensor
//If MAC address is already locked returns false
//If locking was successful returns true
func lockSensor( MAC string) bool{
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

//Release lock from the sensor
func releaseSensor ( MAC string ) bool {
    //Let's not store all measurements per device
    if StoreDelay > 0 {
        time.Sleep(StoreDelay)
    }

    mutex.Lock()
    mWait[MAC] = false
    mutex.Unlock()

    return true
}

//ParseRuuviData parses raw ruuvidata received from RuuviTag
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
