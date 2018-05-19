package AEDAevents

import (
	"errors"
	"strconv"

	"github.com/torlenor/AbyleEDA/eventmessage"
	"github.com/torlenor/AbyleEDA/quantities"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAserver"
)

var log = logging.MustGetLogger("AEDAlogger")

var customEventLookupMap = map[int32]map[int32]func(eventmessage.EventMessage){}

// AddCustomEvent lets you add a new callback for a certain clientID and eventID
func AddCustomEvent(clientID int32, eventID int32, f func(eventmessage.EventMessage)) {
	if customEventLookupMap[clientID] == nil {
		customEventLookupMap[clientID] = map[int32]func(eventmessage.EventMessage){}
	}

	customEventLookupMap[clientID][eventID] = f
}

func executeCustomEvent(event eventmessage.EventMessage) error {
	if cb, found := customEventLookupMap[event.ClientID][event.EventID]; found {
		cb(event)
	} else {
		log.Warning("clientID: " + strconv.Itoa(int(event.ClientID)) + ", eventID: " + strconv.Itoa(int(event.EventID)) + " does not exist in custom event lookup table")
		return errors.New("clientID: " + strconv.Itoa(int(event.ClientID)) + ", eventID: " + strconv.Itoa(int(event.EventID)) + " does not exist in custom event lookup table")
	}
	return nil
}

// https://play.golang.org/p/vEy-GPulXIN

var myWriter AEDAserver.ServerWriter

// SetAEDAserver defines the serverWriter to use for the event system in
// case of sending a message back to the clients or to another server
func SetAEDAserver(serverWriter AEDAserver.ServerWriter) {
	myWriter = serverWriter
}

// EventInterpreter should be called when a new message comes in and
// it will be the entry point to the event handling process
func EventInterpreter(event eventmessage.EventMessage) {
	printEvent(event)
	executeCustomEvent(event)
}

func printEvent(event eventmessage.EventMessage) {
	log.Info("ClientID =", event.ClientID)
	log.Info("EventID =", event.EventID)
	cnt := 0
	for _, content := range event.Quantities {
		cnt++

		switch v := content.(type) {
		case *quantities.Temperature:
			log.Info("Content (numeric)", cnt, ":", v.Degrees(), "°C")
		default:
			log.Info("Content", cnt, ":", content.String())
		}
	}
}
