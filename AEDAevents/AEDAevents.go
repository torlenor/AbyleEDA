package AEDAevents

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/torlenor/AbyleEDA/quantities"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAserver"
)

var log = logging.MustGetLogger("AEDAlogger")

// EventMessage contains the data of an event
type EventMessage struct {
	ID         int32                 `json:"id"`
	Quantities []quantities.Quantity `json:"quantities"`
}

// UnmarshalJSON is part of the json interface for EventMessage
func (ce *EventMessage) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var idval int32
	err = json.Unmarshal(*objMap["id"], &idval)
	if err != nil {
		return err
	}
	ce.ID = idval

	var rawMessagesForEventMessage []*json.RawMessage
	err = json.Unmarshal(*objMap["quantities"], &rawMessagesForEventMessage)
	if err != nil {
		return err
	}

	ce.Quantities = make([]quantities.Quantity, len(rawMessagesForEventMessage))

	var m map[string]string
	for index, rawMessage := range rawMessagesForEventMessage {
		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}

		if m["type"] == "temperature" {
			var t quantities.Temperature
			err := json.Unmarshal(*rawMessage, &t)
			if err != nil {
				return err
			}
			ce.Quantities[index] = &t
		} else {
			return errors.New("unsupported type found: " + m["type"])
		}
	}

	return nil
}

// Sensor struct is used to store sensor data
type Sensor struct {
	ID          int32
	SensorType  string
	Quantity    string
	Value       float64
	Unit        string
	LastUpdated time.Time
}

// M temporarily stores Sensor data until we have a more sophisticated
// system in place
var M map[int32]Sensor

func init() {
	M = make(map[int32]Sensor)
}

var myWriter AEDAserver.ServerWriter

// SetAEDAserver defines the serverWriter to use for the event system in
// case of sending a message back to the clients or to another server
func SetAEDAserver(serverWriter AEDAserver.ServerWriter) {
	myWriter = serverWriter
}

func eventValueUpdate(event EventMessage) {
	if _, ok := M[event.ID]; ok {
		log.Info("Received sensor update:")
		printEvent(event)

		// M[event.ID] = Sensor{ID: event.ID,
		// 	SensorType: "temperature",
		// 	Quantity:   "temperature",
		// 	Value:      event.Value,
		// 	Unit:       "event.Unit"}
	} else {
		log.Info("Registering new sensor:")
		printEvent(event)

		// M[event.ID] = Sensor{ID: event.ID,
		// 	SensorType: "temperature",
		// 	Quantity:   "temperature",
		// 	Value:      event.Value,
		// 	Unit:       "event.Unit"}
	}
}

func eventTrigger(event EventMessage) {
	log.Info("Received trigger event:")
	printEvent(event)
}

// EventInterpreter should be called when a new message comes in and
// it will be the entry point to the event handling process
func EventInterpreter(event EventMessage) {
	eventValueUpdate(event)
	eventTrigger(event)
}

func printEvent(event EventMessage) {
	log.Info("ID =", event.ID)
	cnt := 0
	for _, content := range event.Quantities {
		cnt++
		log.Info("Content", cnt, ": ", content.String())
	}
}
