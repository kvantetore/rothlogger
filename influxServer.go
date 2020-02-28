
package main

import (
	"time"
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/kvantetore/rothTouchline"
)

type InfluxSettings struct {
	serverURL string
	dbName string
	username string
	password string
	measurementName string
}

 //StoreSensorData saves the current state of the sensors to an influxdb measurement
func StoreSensorData(settings InfluxSettings, sensors []roth.Sensor) error {
	//create influx client
	cli, err := client.NewHTTPClient(client.HTTPConfig {
		Addr: settings.serverURL,
		Username: settings.username,
		Password: settings.password,
	})
	if err != nil {
		return fmt.Errorf("Failed to create HTTP Client, %v", err)
	}
	defer cli.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  settings.dbName,
		Precision: "s",
	})
	if err != nil {
		return fmt.Errorf("Error creating batch points %v", err)
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
			"mode": fmt.Sprintf("%v", sensor.Mode),
		}
		fields := map[string]interface{}{
			"roomTemperature": sensor.RoomTemperature,
			"targetTemperature": sensor.TargetTemperature,
			"valve_value": sensor.GetValveValue(),
		}

		pt, err := client.NewPoint(settings.measurementName, tags, fields, currentTime)
		if err != nil {
			return fmt.Errorf("Error creating new point, %v", err)
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	if err := cli.Write(bp); err != nil {
		return fmt.Errorf("error writing points, %v", err)
	}

	return nil;
}
