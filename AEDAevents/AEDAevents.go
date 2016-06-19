package AEDAevents

import (
	"github.com/op/go-logging"
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

func eventValueUpdate(event EventMessage) {
	log.Info("Received valueupdate event:")
	log.Info("Id =", event.Id)
	log.Info("Type =", event.Type)
	log.Info("Event =", event.Event)
	log.Info("Quantity =", event.Quantity)
	log.Info("Value =", event.Value, event.Unit)
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
