package main

import (
	"context"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	qdb "github.com/questdb/go-questdb-client/v4"
)

const (
	// Specifies name of bucket where data is stored
	DBERROR = "DB Error: "
)

func StoreMeasurement(sensorData SensorData) {
	measurement := make(map[string]interface{})

	//Store only valid values
	if sensorData.ValidData.Temp == true {
		measurement["Temperature"] = sensorData.Temp
	}

	//Sometimes Humidity is reported incorrectly
	//Let's not store it if it is way too high number
	if sensorData.ValidData.Humidity == true {
		measurement["Humidity"] = sensorData.Humidity
	}

	//Store only valid values
	if sensorData.ValidData.Pressure == true {
		measurement["Pressure"] = int(sensorData.Pressure)
	}

	//Accelarations
	if sensorData.ValidData.AccelerationX == true {
		measurement["AccelerationX"] = int(sensorData.AccelerationX)
	}
	if sensorData.ValidData.AccelerationY == true {
		measurement["AccelerationY"] = int(sensorData.AccelerationY)
	}
	if sensorData.ValidData.AccelerationZ == true {
		measurement["AccelerationZ"] = int(sensorData.AccelerationZ)
	}

	//Calculated values
	if sensorData.ValidData.Temp == true && sensorData.ValidData.Humidity == true {
		measurement["AbsoluteHumidity"] = sensorData.AbsHumidity
	}

	if sensorData.ValidData.Temp == true && sensorData.ValidData.Humidity == true {
		measurement["Dewpoint"] = sensorData.DewPoint
	}

	insert(sensorData, sensorData.MAC)
}

func StoreHWmeasurement(sensorData SensorData) {
	HWmeasurement := make(map[string]interface{})

	if sensorData.ValidData.Battery == true {
		HWmeasurement["Battery"] = int(sensorData.Battery)
	}

	if sensorData.ValidData.TXPower == true {
		HWmeasurement["TxPower"] = int(sensorData.TXPower)
	}

	insertHW(HWmeasurement, sensorData.MAC)
}

//
//
//

var ctx context.Context
var client qdb.LineSender = nil

// Databse connect
func dbConnect() {
	ctx = context.TODO()

	var err error

	log.Println("Opening Databse connection.")

	//Creating Connection string
	connectionString := "http::addr=" + "localhost:9000" + //ConnectionString +
		";token=" + DBToken + ";auto_flush_interval=1000;"

	client, err = qdb.LineSenderFromConf(ctx, connectionString)
	if err != nil {
		log.Panic("Failed to connect to database.")
	}
	log.Println("Database connection opened.")

}

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

// Insert points to database
// Uses: Measurement table
// func insert(measurement map[string]interface{}, MAC string) {
func insert(sensorData SensorData, MAC string) {

	tags := map[string]string{"Address": Address}

	bFound := false
	// Create a point and add to batch
	for i := 0; i != len(sensorAddresses); i++ {
		if MAC == sensorAddresses[i] {
			tags["Location"] = aLocations[i]
			bFound = true
		}
	}

	/*Check if MAC was found from the list to be stored to DB*/
	if bFound == false {
		/*Not found. not storing*/
		log.Printf("%s Ruuvitag is not listed in config.yml to be stored to database.\n", MAC)
		return
	}

	tags["Device"] = MAC

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

// Insert points to database
// Uses Hardware table
func insertHW(measurement map[string]interface{}, MAC string) {
	// Create client and set batch size to 2
	c := influxdb2.NewClientWithOptions(ConnectionString, DBToken, influxdb2.DefaultOptions().SetBatchSize(2))
	defer c.Close()

	// user blocking write client for writes to desired bucket
	writeAPI := c.WriteAPI(DBOrg, DBBucket)

	// Create a point and add to batch
	tags := map[string]string{"Address": Address}

	bFound := false
	// Create a point and add to batch
	for i := 0; i != len(sensorAddresses); i++ {
		if MAC == sensorAddresses[i] {
			tags["Location"] = aLocations[i]
			bFound = true
		}
	}

	/*Check if MAC was found from the list to be stored to DB*/
	if bFound == false {
		/*Not found, not storing*/
		log.Printf("%s Ruuvitag is not listed in config.yml to be stored to database.\n", MAC)
		return
	}

	tags["Device"] = MAC
	fields := measurement

	pt := influxdb2.NewPoint("HWmeasurements", tags, fields, time.Now())
	writeAPI.WritePoint(pt)

	// Force all unwritten data to be sent
	writeAPI.Flush()

}
