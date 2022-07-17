package main

import (
    "log"
    "time"

    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
    // Specifies name of bucket where data is stored
    MyBucket = "weather"
    DBERROR  = "DB Error: "
)

// Insert points to database
// Uses: Measurement table
func Insert(measurement map[string]interface{}, MAC string) {
    // Create client and set batch size to 2
    c := influxdb2.NewClientWithOptions(ConnectionString, DBToken, influxdb2.DefaultOptions().SetBatchSize(2))
    defer c.Close()

    // user blocking write client for writes to desired bucket
    writeAPI := c.WriteAPI(DBOrg, MyBucket)

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
        log.Printf("%s Ruuvitag is not stored to DB.\n", MAC)
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
//Uses Hardware table
func InsertHW(measurement map[string]interface{}, MAC string) {
    // Create client and set batch size to 2
    c := influxdb2.NewClientWithOptions(ConnectionString, "my-token", influxdb2.DefaultOptions().SetBatchSize(2))
    defer c.Close()

    // user blocking write client for writes to desired bucket
    writeAPI := c.WriteAPI("my-org", MyBucket)

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
        log.Printf("%s Ruuvitag is not stored to DB.\n", MAC)
        return
    }

    tags["Device"] = MAC
    fields := measurement

    pt := influxdb2.NewPoint("HWmeasurements", tags, fields, time.Now())
    writeAPI.WritePoint(pt)

    // Force all unwritten data to be sent
    writeAPI.Flush()

}
