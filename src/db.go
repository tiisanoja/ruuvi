package main

import (
	"context"
	"log"

	qdb "github.com/questdb/go-questdb-client/v4"
)

// Stores ruuvitag measurement related data
func StoreMeasurement(sensorData SensorData) {

	insert(sensorData)
}

// Stores ruuvitag Hardware related data
func StoreHWmeasurement(sensorData SensorData) {

	insertHW(sensorData)
}

// Create Tags to be given for each measurement
// Tags to be added:
// Address: Street Address for Ruuvitag. If used in multiple addresses. For example if you move you can use this to filter measurements
// MAC: Ruuvitag's MAC address
// Location: Location in the Adress. For example livingroom
func createTags(sensorData SensorData) (map[string]string, bool) {
	tags := map[string]string{"Address": Address}

	bFound := false
	//Add adress
	for i := 0; i != len(sensorAddresses); i++ {
		if sensorData.MAC == sensorAddresses[i] {
			tags["Location"] = aLocations[i]
			bFound = true
		}
	}

	tags["Device"] = sensorData.MAC

	return tags, bFound
}

// Context and client for QuestDB connection
var ctx context.Context
var client qdb.LineSender = nil

// Databse connect
// Open it only once
func dbConnect() {
	ctx = context.TODO()

	var err error

	log.Println("Opening Databse connection.")

	//Creating Connection string
	connectionString := ConnectionString + "auto_flush_interval=1000;"

	client, err = qdb.LineSenderFromConf(ctx, connectionString)
	if err != nil {
		log.Panic("Failed to connect to database.")
	}
	log.Println("Database connection opened.")

}

// Close the database connection
// Close it only once
func dbClose() {
	if client == nil {
		log.Println("WARNING: Database connection was already closed or not opened at all.")
		return
	}

	err := client.Flush(ctx)
	if err != nil {
		log.Printf("ERROR: Database Flush failed: %s\n", err.Error())
	}

	err = client.Close(ctx)
	if err != nil {
		log.Printf("ERROR: Database Close failed: %s\n", err.Error())
	}

	client = nil
}

// Insert measurement points to database
// Uses: Measurement table
func insert(sensorData SensorData) {

	//Create tags/Symbols for each row
	tags, bFound := createTags(sensorData)

	/*Check if MAC was found from the list to be stored to DB*/
	if bFound == false {
		/*Not found. not storing*/
		log.Printf("%s Ruuvitag is not listed in config.yml to be stored to database.\n", sensorData.MAC)
		return
	}

	//Write measurement to database
	err := client.Table("measurements").
		Symbol("Address", tags["Address"]).
		Symbol("Device", tags["Device"]).
		Symbol("Location", tags["Location"]).
		Float64Column("Temperature", sensorData.Temp).
		Int64Column("Pressure", int64(sensorData.Pressure)).
		Float64Column("Humidity", sensorData.Humidity).
		Float64Column("Dewpoint", sensorData.DewPoint).
		Float64Column("AbsoluteHumidity", sensorData.AbsHumidity).
		Int64Column("AccelerationX", int64(sensorData.AccelerationX)).
		Int64Column("AccelerationY", int64(sensorData.AccelerationY)).
		Int64Column("AccelerationZ", int64(sensorData.AccelerationZ)).
		AtNow(ctx)

	if err != nil {
		log.Panicf("Failed to insert data to Database: %s\n", err.Error())
	}

	//Flush data
	err = client.Flush(ctx)
	if err != nil {
		log.Panicf("Failed to flush data: %s\n", err.Error())
	}

}

// Insert RuuviTag hardware related points to database
// Uses Hardware table
func insertHW(sensorData SensorData) {

	//Create tags/Symbols for each row
	tags, bFound := createTags(sensorData)

	/*Check if MAC was found from the list to be stored to DB*/
	if bFound == false {
		/*Not found, not storing*/
		log.Printf("%s Ruuvitag is not listed in config.yml to be stored to database.\n", sensorData.MAC)
		return
	}

	//Write measurement to database
	err := client.Table("HWmeasurements").
		Symbol("Address", tags["Address"]).
		Symbol("Device", tags["Device"]).
		Symbol("Location", tags["Location"]).
		Int64Column("Battery", int64(sensorData.Battery)).
		Int64Column("TxPower", int64(sensorData.TXPower)).
		AtNow(ctx)

	if err != nil {
		log.Panicf("Failed to insert data to Database: %s\n", err.Error())
	}

	//Flush data
	err = client.Flush(ctx)
	if err != nil {
		log.Panicf("Failed to flush data: %s\n", err.Error())
	}

}
