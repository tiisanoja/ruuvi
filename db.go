package main

import (
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
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

	if sensorData.ValidData.Temp == true && sensorData.ValidData.Humidity == true {
		measurement["AbsoluteHumidity"] = sensorData.AbsHumidity
		measurement["Dewpoint"] = sensorData.DewPoint
		measurement["SSI"] = sensorData.SSI
	}

	insert(measurement, sensorData.MAC)
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

// Insert points to database
// Uses: Measurement table
func insert(measurement map[string]interface{}, MAC string) {
	// Create client and set batch size to 2
	c := influxdb2.NewClientWithOptions(ConnectionString, DBToken, influxdb2.DefaultOptions().SetBatchSize(2))
	defer c.Close()

	// user blocking write client for writes to desired bucket
	writeAPI := c.WriteAPI(DBOrg, DBBucket)

	// Get errors channel
	errorsCh := writeAPI.Errors()
	// Create go proc for reading and logging errors
	go func() {
		for err := range errorsCh {
			log.Printf("DB Error: %s\n", err.Error())
		}
	}()

	//
	tags := map[string]string{"Address": Address}

	bFound := false
	// Create a point and add to batch
	for i := 0; i != len(aSensors); i++ {
		if MAC == aSensors[i] {
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
	fields := measurement

	pt := influxdb2.NewPoint("measurements", tags, fields, time.Now())
	writeAPI.WritePoint(pt)

	// Force all unwritten data to be sent
	writeAPI.Flush()
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
	for i := 0; i != len(aSensors); i++ {
		if MAC == aSensors[i] {
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
