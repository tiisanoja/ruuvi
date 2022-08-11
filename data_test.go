package main

import "testing"
import "math"

//Test Absolutely Humidity calculation logic
func TestAbsHumidity(t *testing.T){
    var got float64
    var want float64

    //Test with 20C, 100% Humidity
    got = math.Round(calculateAbsHumidity(float64(20), float64(100)) *100) / 100
    want = 17.28
    if got != want {
        t.Errorf("T=20C H=100: got %f, wanted %f", got, want)
    }

    //Test with -14C, 25% Humidity
    got = math.Round(calculateAbsHumidity(float64(-14), float64(25)) *100) / 100
    want = 0.43
    if got != want {
        t.Errorf("T=-14 H=25: got %f, wanted %f", got, want)
    }

    //Test with -10C, 99% Humidity
    got = math.Round(calculateAbsHumidity(float64(-10), float64(99)) *100) / 100
    want = 2.34
    if got != want {
        t.Errorf("T=-10C H=99: got %f, wanted %f", got, want)
    }

    //Test with 40C, 95% Humidity
    got = math.Round(calculateAbsHumidity(float64(40), float64(95)) *10) / 10
    want = 48.6
    if got != want {
        t.Errorf("T=40 H=95: got %f, wanted %f", got, want)
    }

    //Test with 30C, 19% Humidity
    got = math.Round(calculateAbsHumidity(float64(30), float64(19)) *10) / 10
    want = 5.8
    if got != want {
        t.Errorf("T=30C H=19: got %f, wanted %f", got, want)
    }


    //Test with -40C, 70% Humidity
    got = math.Round(calculateAbsHumidity(float64(-40), float64(70)) *100) / 100
    want = 0.12
    if got != want {
        t.Errorf("T=-40C H=70: got %f, wanted %f", got, want)
    }

}


//Test dew point calculation logic
func TestDewPoint(t *testing.T){
    var got float64
    var want float64

    //Rounded values taken from Ilmatieteenlaitos.fi
    got = math.Round(calculateDewPoint(float64(16.4), float64(74)) * 10) / 10
    want = 11.8

    if got != want {
        t.Errorf("Testing dew point: got %f, wanted %f", got, want)
    }

    //Rounded values taken from calculator.net
    got = math.Round(calculateDewPoint(float64(0), float64(100)) * 10) / 10
    want = 0

    if got != want {
        t.Errorf("Testing dew point: got %f, wanted %f", got, want)
    }


    //Rounded values taken from calculator.net
    got = math.Round(calculateDewPoint(float64(-35), float64(20)) * 10) / 10
    want = -50.1

    if got != want {
        t.Errorf("Testing dew point: got %f, wanted %f", got, want)
    }

    //Rounded values taken from calculator.net
    got = math.Round(calculateDewPoint(float64(40), float64(90)) * 10) / 10
    want = 38.0

    if got != want {
        t.Errorf("Testing dew point: got %f, wanted %f", got, want)
    }

}

func TestLockingLogic(t *testing.T){
    //Try to lock
    got := lockSensor("1")
    want := true
    if got != want {
        t.Error("Locking failed.")
    }

    //Try to lock when there is a lock
    got = lockSensor("1")
    want = false
    if got != want {
        t.Error("Locking not working.")
    }

    StoreDelay = 0
    //Release
    got = releaseSensor("1")
    want = true
    if got != want {
        t.Error("Release failed.")
    }
}

//Test Parsing with some example data
func TestParsingWithRealData(t *testing.T){
    data := []byte{0x04, 0x99, 0x05, 0x12, 0xFC, 0x53, 0x94, 0xC3, 0x7C, 0x00, 0x04, 0xFF, 0xFC, 0x04, 0x0C, 0xAC, 0x36, 0x42, 0x00, 0xCD, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F}

    //
    //Expected values
    //
    Temp := 24.3
    Humidity := 53.49
    Pressure := uint32(100044)
    TXPower := int16(4)
    Battery := uint16(2977)
    MAC := "cb:b8:33:4c:88:4f"
    Seq := uint16(205)
    AccelerationX := int16(4)
    AccelerationY := int16(-4)
    AccelerationZ := int16(1036)

    parsingHelper(data, Temp, Humidity, Pressure, TXPower, Battery, MAC, Seq, AccelerationX, AccelerationY, AccelerationZ, t)
}


//Test Parsing with max values
func TestParsingWithMaxValues(t *testing.T){
    data := []byte{0x04, 0x99, 0x05, 0x7F, 0xFF, 0xFF, 0xFE, 0xFF, 0xFE, 0x7F, 0xFF, 0x7F, 0xFF, 0x7F, 0xFF, 0xFF, 0xDE, 0xFE, 0xFF, 0xFE, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F}

    //
    //Expected values
    //
    Temp := 163.835
    Humidity := 163.8350
    Pressure := uint32(115534)
    TXPower := int16(20)
    Battery := uint16(3646)
    MAC := "cb:b8:33:4c:88:4f"
    Seq := uint16(65534)
    AccelerationX := int16(32767)
    AccelerationY := int16(32767)
    AccelerationZ := int16(32767)

    parsingHelper(data, Temp, Humidity, Pressure, TXPower, Battery, MAC, Seq, AccelerationX, AccelerationY, AccelerationZ, t)
}


//Test Parsing with min values
func TestParsingWithMinValues(t *testing.T){
    data := []byte{0x04, 0x99, 0x05, 0x80, 0x01, 0x00, 0x00, 0x00, 0x00, 0x80, 0x01, 0x80, 0x01, 0x80, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F}

    //
    //Expected values
    //
    Temp := -163.835
    Humidity := 0.000
    Pressure := uint32(50000)
    TXPower := int16(-40)
    Battery := uint16(1600)
    MAC := "cb:b8:33:4c:88:4f"
    Seq := uint16(0)
    AccelerationX := int16(-32767)
    AccelerationY := int16(-32767)
    AccelerationZ := int16(-32767)

    sensorData := parsingHelper(data, Temp, Humidity, Pressure, TXPower, Battery, MAC, Seq, AccelerationX, AccelerationY, AccelerationZ, t)

    expectedValue := true
    //Verify that ValidData has correct values in case of valid input
    if sensorData.ValidData.Temp == !expectedValue {
        t.Errorf("SensorData.ValidData.Temp is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.Temp)
    }
    if sensorData.ValidData.Humidity == !expectedValue {
        t.Errorf("SensorData.ValidData.Humidity is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.Humidity)
    }
    if sensorData.ValidData.Pressure == !expectedValue {
        t.Errorf("SensorData.ValidData.Pressure is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.Pressure)
    }
    if sensorData.ValidData.TXPower == !expectedValue {
        t.Errorf("SensorData.ValidData.TXPower is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.TXPower)
    }
    if sensorData.ValidData.Battery == !expectedValue {
        t.Errorf("SensorData.ValidData.Battery is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.Battery)
    }
    if sensorData.ValidData.Seq == !expectedValue {
        t.Errorf("SensorData.ValidData.Seq is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.Seq)
    }
    if sensorData.ValidData.AccelerationX == !expectedValue {
        t.Errorf("SensorData.ValidData.AccelerationX is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.AccelerationX)
    }
    if sensorData.ValidData.AccelerationY == !expectedValue {
        t.Errorf("SensorData.ValidData.AccelerationY is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.AccelerationY)
    }
    if sensorData.ValidData.AccelerationZ == !expectedValue {
        t.Errorf("SensorData.ValidData.AccelerationZ is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.AccelerationZ)
    }
    if sensorData.ValidData.MAC == !expectedValue {
        t.Errorf("SensorData.ValidData.MAC is false. Expected: %v Got: %v", expectedValue, sensorData.ValidData.AccelerationZ)
    }


}

//Test Parsing with NAN values
func TestParsingWithNANValues(t *testing.T){
    data := []byte{0x04, 0x99, 0x05, 0x80, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00, 0x80, 0x00, 0x80, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

    //
    //Expected NAN values
    //
    Temp := -163.84
    Humidity := 163.8375
    Pressure := uint32(115535)
    TXPower := int16(22)
    Battery := uint16(3647)
    MAC := "ff:ff:ff:ff:ff:ff"
    Seq := uint16(65535)
    AccelerationX := int16(-32768)
    AccelerationY := int16(-32768)
    AccelerationZ := int16(-32768)

    sensorData := parsingHelper(data, Temp, Humidity, Pressure, TXPower, Battery, MAC, Seq, AccelerationX, AccelerationY, AccelerationZ, t)

    expectedValue := false
    //Verify that ValidData has correct values
    if sensorData.ValidData.Temp == !expectedValue {
        t.Errorf("SensorData.ValidData.Temp: Expected: %v Got: %v", expectedValue, sensorData.ValidData.Temp)
    }
    if sensorData.ValidData.Humidity == !expectedValue {
        t.Errorf("SensorData.ValidData.Humidity: Expected: %v Got: %v", expectedValue, sensorData.ValidData.Humidity)
    }
    if sensorData.ValidData.Pressure == !expectedValue {
        t.Errorf("SensorData.ValidData.Pressure: Expected: %v Got: %v", expectedValue, sensorData.ValidData.Pressure)
    }
    if sensorData.ValidData.TXPower == !expectedValue {
        t.Errorf("SensorData.ValidData.TXPower: Expected: %v Got: %v", expectedValue, sensorData.ValidData.TXPower)
    }
    if sensorData.ValidData.Battery == !expectedValue {
        t.Errorf("SensorData.ValidData.Battery: Expected: %v Got: %v", expectedValue, sensorData.ValidData.Battery)
    }
    if sensorData.ValidData.Seq == !expectedValue {
        t.Errorf("SensorData.ValidData.Seq: Expected: %v Got: %v", expectedValue, sensorData.ValidData.Seq)
    }
    if sensorData.ValidData.AccelerationX == !expectedValue {
        t.Errorf("SensorData.ValidData.AccelerationX: Expected: %v Got: %v", expectedValue, sensorData.ValidData.AccelerationX)
    }
    if sensorData.ValidData.AccelerationY == !expectedValue {
        t.Errorf("SensorData.ValidData.AccelerationY: Expected: %v Got: %v", expectedValue, sensorData.ValidData.AccelerationY)
    }
    if sensorData.ValidData.AccelerationZ == !expectedValue {
        t.Errorf("SensorData.ValidData.AccelerationZ: Expected: %v Got: %v", expectedValue, sensorData.ValidData.AccelerationZ)
    }


}


//
// Internal Functions
//

//Test Parsing
func parsingHelper(data []byte, Temp float64, Humidity float64, Pressure uint32, TXPower int16, Battery uint16, MAC string, Seq uint16,
    AccelerationX int16, AccelerationY int16, AccelerationZ int16, t *testing.T) SensorData {

    err, sensorData := parseSensorFormat5(data)
    if err == nil {

       //Verify Temperature
       if sensorData.Temp != Temp {
            t.Errorf("SensorData is having incorrect Temperature. Expected: %f Got: %f", Temp, sensorData.Temp)
       }

       //Verify Humidity
       if sensorData.Humidity!= Humidity {
            t.Errorf("SensorData is having incorrect Humidity. Expected: %f Got: %f", Humidity, sensorData.Humidity)
       }

       //Verify Pressure
       if sensorData.Pressure != Pressure {
            t.Errorf("SensorData is having incorrect Pressure. Expected: %d Got: %d",Pressure, sensorData.Pressure)
       }

       //Verify Transmit Power
       if sensorData.TXPower != TXPower {
            t.Errorf("SensorData is having incorrect Transmit Power. Expected: %d Got: %d", TXPower, sensorData.TXPower)
       }

       //Verify Voltage
       if sensorData.Battery != Battery {
            t.Errorf("SensorData is having incorrect Voltage. Expected: %d Got: %d",  Battery, sensorData.Battery)
       }

       //Verify Sequence
       if sensorData.Seq != Seq {
            t.Errorf("SensorData is having incorrect Sequence. Expected: %d Got: %d",  Seq, sensorData.Seq)
       }


       //Verify AccelerationX
       if sensorData.AccelerationX != AccelerationX {
            t.Errorf("SensorData is having incorrect AccelerationX. Expected: %d Got: %d",  AccelerationX, sensorData.AccelerationX)
       }

       //Verify AccelerationY
       if sensorData.AccelerationY != AccelerationY {
            t.Errorf("SensorData is having incorrect AccelerationY. Expected: %d Got: %d",  AccelerationY, sensorData.AccelerationY)
       }

       //Verify AccelerationZ
       if sensorData.AccelerationZ != AccelerationZ {
            t.Errorf("SensorData is having incorrect AccelerationZ. Expected: %d Got: %d",  AccelerationZ, sensorData.AccelerationZ)
       }

       //Verify MAC
       if sensorData.MAC != MAC {
            t.Errorf("SensorData is having incorrect MAC. Expected: %s Got: %s", MAC, sensorData.MAC)
       }

    } else {
        t.Errorf("Sample data parsing failed. Error: %s",err)
    }

    return sensorData
}
