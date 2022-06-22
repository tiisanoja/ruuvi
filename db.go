package main


import (
	"log"
	"time"
        "github.com/influxdata/influxdb1-client/v2"
)

const (
	// MyDB specifies name of database
        MyDB = "weather"
        DBERROR = "DB Error: "
)

// Insert points to database
// Uses: Measurement table
func Insert(measurement map[string]interface{}, MAC string) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: ConnectionString,
	})
	if err != nil {
	        log.Printf(DBERROR)
		log.Println(err)
                return
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
	        log.Printf(DBERROR)
                log.Println(err)
                return
	}

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

	pt, err := client.NewPoint("measurements", tags, fields, time.Now())
	if err != nil {
	        log.Printf(DBERROR)
                log.Println(err)
                return
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
	        log.Printf(DBERROR)
                log.Println(err)
                return
	}

   	// Close client resources
	if err := c.Close(); err != nil {
	        log.Printf(DBERROR)
                log.Println(err)
                return
	}
}

// Insert points to database
//Uses Hardware table
func InsertHW(measurement map[string]interface{}, MAC string) {
        c, err := client.NewHTTPClient(client.HTTPConfig{
                Addr: ConnectionString,
        })
        if err != nil {
	        log.Printf(DBERROR)
                log.Println(err)
                return
        }
        defer c.Close()

        // Create a new point batch
        bp, err := client.NewBatchPoints(client.BatchPointsConfig{
                Database:  MyDB,
                Precision: "s",
        })
        if err != nil {
	        log.Printf(DBERROR)
                log.Println(err)
                return
        }

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

        pt, err := client.NewPoint("HWmeasurements", tags, fields, time.Now())
        if err != nil {
	        log.Printf(DBERROR)
                log.Println(err)
                return
        }
        bp.AddPoint(pt)

        // Write the batch
        if err := c.Write(bp); err != nil {
	        log.Printf(DBERROR)
                log.Println(err)
                return
        }

        // Close client resources
        if err := c.Close(); err != nil {
	        log.Printf(DBERROR)
                log.Println(err)
                return
        }
}
