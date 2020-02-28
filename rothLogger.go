package main

import (
	"os"
	"time"
	"fmt"
	"github.com/kvantetore/rothTouchline"
)

const (
	rothManagementURL = "http://ROTH-01A6D5"
	influxServer = "http://pi:8086"
	influxDb = "home"
	thermostatMeasurement = "thermostats"
 )

func main() {
	fmt.Printf("Starting logging thermostat data!\n")
	fmt.Printf("    ROTH Touchline management unit: %v\n", rothManagementURL)
	fmt.Printf("    Storing thermostat data in influx server %v, database %v, measurement %v\n", influxServer, influxDb, thermostatMeasurement)
	
	//setup
	fmt.Printf("Setting up...")
	sensorCount, err := roth.GetSensorCount(rothManagementURL)
	if err != nil {
		fmt.Printf("Error fetching sensors, %v\n", err)
		return
	}

	rothUrl := os.Getenv("ROTH_URL")

	influxSettings := InfluxSettings {
		serverURL: os.Getenv("INFLUX_URL"),
		dbName: os.Getenv("INFLUX_DB"),
		measurementName: os.Getenv("INFLUX_MEASUREMENT"),
		username: os.Getenv("INFLUX_USERNAME"),
		password: os.Getenv("INFLUX_PASSWORD"),
	}	
	fmt.Printf("done\n")
	
	//do measurement
	performMeasurement := func() {
		sensors, err := roth.GetSensors(rothUrl, sensorCount)
		if err != nil {
			fmt.Printf("Error fetching sensors, %v\n", err)
			return
		}

		fmt.Printf("Got sensor data %v\n", time.Now())
		for i := 0; i<len(sensors); i++ {
		 	fmt.Printf("%v: %v %v (%s)\n", sensors[i].Name, sensors[i].RoomTemperature, sensors[i].TargetTemperature, sensors[i].GetValveState())
		}
		fmt.Println("")

		err = StoreSensorData(influxSettings, sensors);
		if err != nil {
			fmt.Printf("Error storing sensor data, %v\n", err)
		}
	}

	//create timer that runs measurement 
	interval := time.Minute * 1
	fmt.Printf("Running logger every %v...\n", interval)
	ticker := time.NewTicker(interval)
	for {
		performMeasurement()
		<- ticker.C
	}

}