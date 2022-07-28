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
        t.Errorf("got %f, wanted %f", got, want)
    }

    //Test with -14C, 25% Humidity
    got = math.Round(calculateAbsHumidity(float64(-14), float64(25)) *100) / 100
    want = 0.43
    if got != want {
        t.Errorf("got %f, wanted %f", got, want)
    }

    //Test with -10C, 99% Humidity
    got = math.Round(calculateAbsHumidity(float64(-10), float64(99)) *100) / 100
    want = 2.33
    if got != want {
        t.Errorf("got %f, wanted %f", got, want)
    }

    //Test with 40C, 95% Humidity
    got = math.Round(calculateAbsHumidity(float64(40), float64(95)) *10) / 10
    want = 48.5
    if got != want {
        t.Errorf("got %f, wanted %f", got, want)
    }

    //Test with 30C, 19% Humidity
    got = math.Round(calculateAbsHumidity(float64(30), float64(19)) *10) / 10
    want = 5.8
    if got != want {
        t.Errorf("got %f, wanted %f", got, want)
    }


    //Test with -40C, 70% Humidity
    got = math.Round(calculateAbsHumidity(float64(-40), float64(70)) *100) / 100
    want = 0.12
    if got != want {
        t.Errorf("got %f, wanted %f", got, want)
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
    want = -49.9

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

