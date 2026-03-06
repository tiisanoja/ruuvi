package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
	"tinygo.org/x/bluetooth"
)

var isPoweredOn = false
var scanMutex = sync.Mutex{}

/*TODO: To be moved*/
var PressureCorrection = 0
var StoreDelay = 15
var ConnectionString = ""
var Address = "Home"
var sensorAddresses []string
var aLocations []string
var DBToken = ""
var DBOrg = ""
var DBBucket = ""

var adapter = bluetooth.DefaultAdapter

func scanHandler(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
	/*
		    data := device.ManufacturerData()
			//Verify that there is data
			if len(data) < 1 {
				return
			}
			byteData := data[0].Data
			ParseRuuviData(byteData)
	*/

	if len(device.ManufacturerData()) < 1 {
		return
	}

	ParseRuuviData(device.ManufacturerData()[0].Data, device.ManufacturerData()[0].CompanyID)
}

func initialize() {
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
	viper.SetDefault("Database.ConnectionString", "http://localhost:8086")
	viper.SetDefault("Address", "Home")
	viper.SetDefault("Measurements.StoreDelay", 15)
	viper.SetDefault("Database.Token", "")
	viper.SetDefault("Database.Org", "")
	viper.SetDefault("Database.Bucket", "weather")

	PressureCorrection = viper.GetInt("Pressure.Correction")
	ConnectionString = viper.GetString("Database.ConnectionString")
	Address = viper.GetString("Address")
	StoreDelay = viper.GetInt("Measurements.StoreDelay")
	DBToken = viper.GetString("Database.Token")
	DBOrg = viper.GetString("Database.Org")
	DBBucket = viper.GetString("Database.Bucket")
}

func main() {

	initialize()

	log.Printf("Pressure correction: %d\n", viper.GetInt("Pressure.Correction"))
	log.Printf("Database connection: %s\n", viper.GetString("Database.ConnectionString"))
	log.Printf("Address: %s\n", viper.GetString("Address"))

	//Open Database connection
	dbConnect()
	defer dbClose()

	log.Printf("Following Ruuvitags will be used:\n")

	iSensors := len(viper.GetStringSlice("Sensors"))
	sensorAddresses = viper.GetStringSlice("Sensors")
	for i := 0; i != iSensors; i++ {
		sTemp := fmt.Sprintf("%s.Location", sensorAddresses[i])
		aLocations = append(aLocations, viper.GetString(sTemp))
		log.Printf("%d. MAC: %s, Location: %s\n", i, sensorAddresses[i], aLocations[i])
	}

	log.Printf("Opening BluetoothLE device...\n")

	if err := adapter.Enable(); err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Start Scanning....")
	adapter.Scan(scanHandler)

	for {
		time.Sleep(time.Second)
		log.Println("Scanning againg...")
	}

}
