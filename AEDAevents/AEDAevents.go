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

type EventContent struct {
    Quantity string
    Value string
    Unit string
}

type EventMessage struct {
    Id       int32
    Type     string
    Event    int32
    Content  []EventContent
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
        printEvent(event);

		// M[event.Id] = Sensor{Id: event.Id,
		// 	SensorType: "temperature",
		// 	Quantity:   "temperature",
		// 	Value:      event.Value,
		// 	Unit:       "event.Unit"}
	} else {
		log.Info("Registering new sensor:")
        printEvent(event);
        
		// M[event.Id] = Sensor{Id: event.Id,
		// 	SensorType: "temperature",
		// 	Quantity:   "temperature",
		// 	Value:      event.Value,
		// 	Unit:       "event.Unit"}
	}
}

func eventTrigger(event EventMessage) {
	log.Info("Received trigger event:")
    printEvent(event);
}

func EventInterpreter(event EventMessage) {
	switch event.Event {
	case EventValueUpdate:
		eventValueUpdate(event)
	case EventTrigger:
		eventTrigger(event)
	default:
		log.Info("Received unknown event:")
        printEvent(event);
	}
}

func printEvent(event EventMessage) {
    log.Info("Id =", event.Id)
    log.Info("Type =", event.Type)
    log.Info("Event =", event.Event)
    cnt := 0
    for _, content := range event.Content {
        cnt++
        log.Info("Content", cnt, ": ", content.Quantity, content.Value, content.Unit)
    }
}
