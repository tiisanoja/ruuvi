package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

var isPoweredOn = false
var scanMutex = sync.Mutex{}

/*TODO: To be moved*/
var PressureCorrection = 0
var StoreDelay = 15 * time.Second
var ConnectionString = ""
var Address = "Home"
var aSensors []string 
var aLocations []string

func beginScan(d gatt.Device) {
	scanMutex.Lock()
	for isPoweredOn {
		d.Scan(nil, true) //Scan for five seconds and then restart
		time.Sleep(5 * time.Second)
		d.StopScanning()
	}
	scanMutex.Unlock()
}

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("scanning...")
		isPoweredOn = true
		go beginScan(d)
		return
	case gatt.StatePoweredOff:
		log.Println("REINIT ON POWER OFF")
		isPoweredOn = false
		d.Init(onStateChanged)
	default:
		log.Println("WARN: unhandled state: ", string(s))
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	ParseRuuviData(a.ManufacturerData, p.ID())
}

func main() {
        // Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

        // Set undefined variables
	viper.SetDefault("Pressure.Correction", 0)
        viper.SetDefault("Database.ConnectionString","http://localhost:8086")
	viper.SetDefault("Address","Home")
	viper.SetDefault("Measurements.StoreDelay",15)

        PressureCorrection = viper.GetInt("Pressure.Correction")
        ConnectionString = viper.GetString("Database.ConnectionString")
	Address = viper.GetString("Address")
	StoreDelay = time.Duration(viper.GetInt("Measurements.StoreDelay")) * time.Second

        log.Printf("Pressure correction: %d\n",viper.GetInt("Pressure.Correction"))
        log.Printf("Database connection: %s\n", viper.GetString("Database.ConnectionString"))
        log.Printf("Address: %s\n", viper.GetString("Address"))

	log.Printf("Following Ruuvitags will be used:\n")
	
	iSensors := len(viper.GetStringSlice("Sensors"))
	aSensors = viper.GetStringSlice("Sensors")
	for i := 0; i != iSensors; i++ {
	    sTemp := fmt.Sprintf("%s.Location",aSensors[i])
	    aLocations = append(aLocations, viper.GetString(sTemp))
	    log.Printf("%d. MAC: %s, Location: %s\n", i, aSensors[i], aLocations[i])
	}

	log.Printf("Opening BluetoothLE device...\n")

	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))
	d.Init(onStateChanged)
	select {}
}
