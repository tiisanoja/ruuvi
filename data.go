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

func parseSensorFormat5(data []byte) *SensorData {
    reader := bytes.NewReader(data)
    result := SensorFormat5{}
    sensorData := SensorData{}

    err := binary.Read(reader, binary.BigEndian, &result)
    if err != nil {
        log.Printf("ERROR: %s\n", err)
        return &sensorData
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

    mutex.Lock()
    if mWait[sensorData.MAC] == true {
        mutex.Unlock()
        log.Printf("INFO: Skipping %s\n", sensorData.MAC)
        return &sensorData
    }
    mWait[sensorData.MAC] = true
    mutex.Unlock()

    fmt.Printf("Seq: %d, Temp: %f, Pres: %d, Hum: %f, Battery: %d, TXPower: %d, MAC: %s\n",
        sensorData.Seq, sensorData.Temp, sensorData.Pressure, sensorData.Humidity, sensorData.Battery, sensorData.TXPower, sensorData.MAC)

    measurement := make(map[string]interface{})
    measurement["Temperature"] = sensorData.Temp
    if sensorData.Humidity < 150.0 {
        measurement["Humidity"] = sensorData.Humidity
    }

    steamSaturationPressure := 6.1078 * math.Pow(10, (7.5*sensorData.Temp/(sensorData.Temp+237.3)))
    absHumidity := (216.679 * (sensorData.Humidity * steamSaturationPressure) / 100) / (sensorData.Temp + 273.15)
    measurement["AbsoluteHumidity"] = absHumidity

    //Dew Point
    tTemp := ((17.27 * sensorData.Temp) / (237.7 + sensorData.Temp)) + math.Log(sensorData.Humidity/100)
    tDewpoint := 0.0
    if tTemp != 17.27 {
        tDewpoint = (237.7 * tTemp) / (17.27 - tTemp)
    } else {
        fmt.Printf("INFO: In tDewpoint: Dividing by zero. Adding 0.01 to temperature.\n")
        tDewpoint = (237.7 * (tTemp + 0.01)) / (17.27 - tTemp + 0.01)
    }

    measurement["Dewpoint"] = tDewpoint
    measurement["Pressure"] = int(sensorData.Pressure)

    Insert(measurement, sensorData.MAC)

    HWmeasurement := make(map[string]interface{})
    HWmeasurement["Battery"] = int(sensorData.Battery)
    HWmeasurement["TxPower"] = int(sensorData.TXPower)

    InsertHW(HWmeasurement, sensorData.MAC)

    //Let's take only every 15s measurement per device
    time.Sleep(StoreDelay)
    mutex.Lock()
    mWait[sensorData.MAC] = false
    mutex.Unlock()

    return &sensorData
}

//ParseRuuviData parses ruuvidata
func ParseRuuviData(data []byte, a string) {

    if len(data) > 2 && binary.LittleEndian.Uint16(data[0:2]) == 0x0499 {
        if len(data) > 25 {
            sensorFormat := data[2]
            fmt.Printf("Ruuvi data with sensor format %d\n", sensorFormat)
            switch sensorFormat {
            case 3:
                fmt.Printf("RuuviTag version 3 not supported. Please upgrade RuuviTag to version 5.")
            case 5:
                parseSensorFormat5(data)
            default:
                fmt.Printf("Unknown sensor format %d\n", sensorFormat)
            }
        } else {
            fmt.Printf("Incorrect lenght in Ruuvi device data: %d. Expected to be 26.\n", len(data))
        }

    } else {
        //fmt.Printf("Not a ruuvi device \n")
    }

}
