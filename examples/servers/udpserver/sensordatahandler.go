package main

import (
	"time"

	"github.com/torlenor/AbyleEDA/AEDAevents"
	"github.com/torlenor/AbyleEDA/quantities"
)

// Sensor struct is used to store sensor data
type Sensor struct {
	ID          int32
	SensorType  string
	Quantity    string
	Value       float64
	Unit        string
	LastUpdated time.Time
}

var sensorsMap = map[int32]map[int32]Sensor{}

func updateSensorValue(event AEDAevents.EventMessage) {
	if sensorsMap[event.ClientID] == nil {
		sensorsMap[event.ClientID] = map[int32]Sensor{}
	}

	for _, content := range event.Quantities {
		if _, ok := sensorsMap[event.ClientID][event.EventID]; ok {
			log.Infof("Received sensor update clientID: %d, sensorID: %d", event.ClientID, event.EventID)

			switch v := content.(type) {
			case *quantities.Temperature:
				log.Info("Content (numeric):", v.Degrees(), "°C")
				sensorsMap[event.ClientID][event.EventID] = Sensor{ID: event.EventID,
					SensorType: "temperature",
					Quantity:   v.Type(),
					Value:      v.Degrees(),
					Unit:       "°C"}
			default:
				log.Info("Content:", content.String())
			}
		} else {
			log.Infof("Registering new sensor clientID: %d, sensorID: %d", event.ClientID, event.EventID)

			switch v := content.(type) {
			case *quantities.Temperature:
				log.Info("Content (numeric):", v.Degrees(), "°C")
				sensorsMap[event.ClientID][event.EventID] = Sensor{ID: event.EventID,
					SensorType: "temperature",
					Quantity:   v.Type(),
					Value:      v.Degrees(),
					Unit:       "°C"}
			default:
				log.Info("Content:", content.String())
			}
		}
	}
}
