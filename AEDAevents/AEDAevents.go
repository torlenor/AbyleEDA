package AEDAevents

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/torlenor/AbyleEDA/quantities"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAserver"
)

var log = logging.MustGetLogger("AEDAlogger")

// EventMessage contains the data of an event
type EventMessage struct {
	ClientID   int32                 `json:"clientid"`
	EventID    int32                 `json:"eventid"`
	Quantities []quantities.Quantity `json:"quantities"`
}

// UnmarshalJSON is part of the json interface for EventMessage
func (ce *EventMessage) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	if value, found := objMap["clientid"]; found {
		var clientidval int32
		err = json.Unmarshal(*value, &clientidval)
		if err != nil {
			return err
		}
		ce.ClientID = clientidval
	} else {
		return errors.New("clientid does not exist in json data")
	}

	if value, found := objMap["eventid"]; found {
		var eventidval int32
		err = json.Unmarshal(*value, &eventidval)
		if err != nil {
			return err
		}
		ce.EventID = eventidval
	} else {
		return errors.New("clientid does not exist in json data")
	}

	var rawMessagesForEventMessage []*json.RawMessage

	if value, found := objMap["quantities"]; found {
		err = json.Unmarshal(*value, &rawMessagesForEventMessage)
		if err != nil {
			return err
		}
		ce.Quantities = make([]quantities.Quantity, len(rawMessagesForEventMessage))
	} else {
		return errors.New("clientid does not exist in json data")
	}

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

var customEventLookupMap = map[int32]map[int32]func(EventMessage){}

// AddCustomEvent lets you add a new callback for a certain clientID and eventID
func AddCustomEvent(clientID int32, eventID int32, f func(EventMessage)) {
	if customEventLookupMap[clientID] == nil {
		customEventLookupMap[clientID] = map[int32]func(EventMessage){}
	}

	customEventLookupMap[clientID][eventID] = f
}

func executeCustomEvent(event EventMessage) error {
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
func EventInterpreter(event EventMessage) {
	printEvent(event)
	executeCustomEvent(event)
}

func printEvent(event EventMessage) {
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
