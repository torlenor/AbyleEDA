package AEDAevents

import (
	"time"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAserver"
)

var log = logging.MustGetLogger("AEDAlogger")

const ( // iota is reset to 0
	EventValueUpdate = iota // = 0
	EventTrigger     = iota // = 1
)

type EventMessage struct {
	Id       int32
	Type     string
	Event    int32
	Quantity string
	Value    float64
	Unit     string
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

var myServer *AEDAserver.UDPServer

func SetAEDAserver(srv *AEDAserver.UDPServer) {
	myServer = srv
}

func eventValueUpdate(event EventMessage) {
	if _, ok := M[event.Id]; ok {
		log.Info("Received sensor update:")
		log.Info("Id =", event.Id)
		log.Info("Type =", event.Type)
		log.Info("Event =", event.Event)
		log.Info("Quantity =", event.Quantity)
		log.Info("Value =", event.Value, event.Unit)
		M[event.Id] = Sensor{Id: event.Id,
			SensorType: "temperature",
			Quantity:   "temperature",
			Value:      event.Value,
			Unit:       "event.Unit"}
	} else {
		log.Info("Registering new sensor:")
		log.Info("Id =", event.Id)
		log.Info("Type =", event.Type)
		log.Info("Event =", event.Event)
		log.Info("Quantity =", event.Quantity)
		log.Info("Value =", event.Value, event.Unit)
		M[event.Id] = Sensor{Id: event.Id,
			SensorType: "temperature",
			Quantity:   "temperature",
			Value:      event.Value,
			Unit:       "event.Unit"}
	}
}

func eventTrigger(event EventMessage) {
	log.Info("Received trigger event:")
	log.Info("Id =", event.Id)
	log.Info("Type =", event.Type)
	log.Info("Event =", event.Event)
	log.Info("Quantity =", event.Quantity)
	log.Info("Value =", event.Value, event.Unit)
}

func EventInterpreter(event EventMessage) {
	switch event.Event {
	case EventValueUpdate:
		eventValueUpdate(event)
	case EventTrigger:
		eventTrigger(event)
	default:
		log.Info("Received unknown event:")
		log.Info("Id =", event.Id)
		log.Info("Type =", event.Type)
		log.Info("Event =", event.Event)
		log.Info("Quantity =", event.Quantity)
		log.Info("Value =", event.Value, event.Unit)
	}
}
