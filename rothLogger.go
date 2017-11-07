package main

import (
	"time"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	rothManagementURL = "http://ROTH-01A6D5"
	influxServer = "http://pi:8086"
	influxDb = "test"
)

func main() {
	fmt.Printf("Hello World!\n")

	sensorCount, err := GetSensorCount(rothManagementURL)
	if err != nil {
		fmt.Printf("Error fetching sensors, %v\n", err)
		return
	}

	sensors, err := GetSensors(rothManagementURL, sensorCount)
	if err != nil {
		fmt.Printf("Error fetching sensors, %v\n", err)
		return
	}

	for i := 0; i<len(sensors); i++ {
		fmt.Printf("%v: %v\n", sensors[i].Name, sensors[i].RoomTemperature)
	}

	//create influx client
	cli, err := client.NewHTTPClient(client.HTTPConfig {
		Addr: influxServer,
	})
	if err != nil {
		fmt.Printf("Failed to create HTTP Client %v\n", err)
		return
	}
	defer cli.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDb,
		Precision: "s",
	})
	if err != nil {
		fmt.Printf("Error creating batch points %v\n", err)
		return
	}

	//create points
	currentTime := time.Now();
	for sensorIndex:=0; sensorIndex<len(sensors); sensorIndex++ {
		sensor := sensors[sensorIndex];

		// Create a point and add to batch
		tags := map[string]string{
			"sensor_name": sensor.Name,
			"sensor_id": fmt.Sprintf("%v", sensor.Id),
			"valve_state": fmt.Sprintf("%v", sensor.GetValveState()),
			"program": fmt.Sprintf("%v", sensor.Program),
		}
		fields := map[string]interface{}{
			"roomTemperature": sensor.RoomTemperature,
			"targetTemperature": sensor.TargetTemperature,
		}

		pt, err := client.NewPoint("thermostats", tags, fields, currentTime)
		if err != nil {
			fmt.Printf("Error creating new point, %v\n", err)
			return
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	if err := cli.Write(bp); err != nil {
		fmt.Printf("error writing points, %v", err)
		return
	}


}