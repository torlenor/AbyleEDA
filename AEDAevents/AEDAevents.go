package AEDAevents

import (
	"errors"
	"strconv"
	"sync"

	"github.com/torlenor/AbyleEDA/eventmessage"
	"github.com/torlenor/AbyleEDA/quantities"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAserver"
)

var log = logging.MustGetLogger("AEDAlogger")

var customEventLookupMap = map[int32]map[int32]func(eventmessage.EventMessage){}
var customEventLookupMapMutex = &sync.Mutex{}

// AddCustomEvent lets you add a new callback for a certain clientID and eventID
func AddCustomEvent(clientID int32, eventID int32, f func(eventmessage.EventMessage)) {
	customEventLookupMapMutex.Lock()
	if customEventLookupMap[clientID] == nil {
		customEventLookupMap[clientID] = map[int32]func(eventmessage.EventMessage){}
	}

	customEventLookupMap[clientID][eventID] = f
	customEventLookupMapMutex.Unlock()
}

// RemoveCustomEvent removes a custom event callback for a certain clientID and eventID
func RemoveCustomEvent(clientID int32, eventID int32) {
	customEventLookupMapMutex.Lock()
	if _, found := customEventLookupMap[clientID][eventID]; found {
		delete(customEventLookupMap[clientID], eventID)
	}
	customEventLookupMapMutex.Unlock()
}

func executeCustomEvent(event eventmessage.EventMessage) error {
	customEventLookupMapMutex.Lock()
	if cb, found := customEventLookupMap[event.ClientID][event.EventID]; found {
		cb(event)
	} else {
		log.Warning("clientID: " + strconv.Itoa(int(event.ClientID)) + ", eventID: " + strconv.Itoa(int(event.EventID)) + " does not exist in custom event lookup table")
		return errors.New("clientID: " + strconv.Itoa(int(event.ClientID)) + ", eventID: " + strconv.Itoa(int(event.EventID)) + " does not exist in custom event lookup table")
	}
	customEventLookupMapMutex.Unlock()
	return nil
}

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
