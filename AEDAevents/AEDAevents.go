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

const ( // iota is reset to 0
	EventValueUpdate = iota // = 0
	EventTrigger     = iota // = 1
)

type EventMessage struct {
	Id         int32
	Type       string
	Event      int32
	Quantities []quantities.Quantity
}

func (ce *EventMessage) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var eventId int32
	err = json.Unmarshal(*objMap["Id"], &eventId)
	if err != nil {
		return err
	}
	ce.Id = eventId

	var typevar string
	err = json.Unmarshal(*objMap["Type"], &typevar)
	if err != nil {
		return err
	}
	ce.Type = typevar

	var event int32
	err = json.Unmarshal(*objMap["Event"], &event)
	if err != nil {
		return err
	}
	ce.Event = event

	var rawMessagesForEventMessage []*json.RawMessage
	err = json.Unmarshal(*objMap["Quantities"], &rawMessagesForEventMessage)
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
			return errors.New("Unsupported type found!")
		}
	}

	return nil
}

type Sensor struct {
	Id          int32
	SensorType  string
	Quantity    string
	Value       float64
	Unit        string
	LastUpdated time.Time
}

var M map[int32]Sensor

func init() {
	M = make(map[int32]Sensor)
}

var myWriter AEDAserver.ServerWriter

func SetAEDAserver(serverWriter AEDAserver.ServerWriter) {
	myWriter = serverWriter
}

func eventValueUpdate(event EventMessage) {
	if _, ok := M[event.Id]; ok {
		log.Info("Received sensor update:")
		printEvent(event)

		// M[event.Id] = Sensor{Id: event.Id,
		// 	SensorType: "temperature",
		// 	Quantity:   "temperature",
		// 	Value:      event.Value,
		// 	Unit:       "event.Unit"}
	} else {
		log.Info("Registering new sensor:")
		printEvent(event)

		// M[event.Id] = Sensor{Id: event.Id,
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

func EventInterpreter(event EventMessage) {
	switch event.Event {
	case EventValueUpdate:
		eventValueUpdate(event)
	case EventTrigger:
		eventTrigger(event)
	default:
		log.Info("Received unknown event:")
		printEvent(event)
	}
}

func printEvent(event EventMessage) {
	log.Info("Id =", event.Id)
	log.Info("Type =", event.Type)
	log.Info("Event =", event.Event)
	cnt := 0
	for _, content := range event.Quantities {
		cnt++
		log.Info("Content", cnt, ": ", content.String())
	}
}
