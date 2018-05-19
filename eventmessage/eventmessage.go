package eventmessage

import (
	"encoding/json"
	"errors"

	"github.com/torlenor/AbyleEDA/quantities"
)

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
