package AEDAevents

import (
	"testing"

	"github.com/torlenor/AbyleEDA/eventmessage"
	"github.com/torlenor/AbyleEDA/quantities"
)

var doAexecuted bool
var doBexecuted bool

func doA(event eventmessage.EventMessage) {
	doAexecuted = true
}

func doB(event eventmessage.EventMessage) {
	doBexecuted = true
}

func TestCustomEvent(t *testing.T) {
	var temp quantities.Temperature
	temp.FromFloat(12.1)
	event := eventmessage.EventMessage{ClientID: 1001,
		EventID:    1,
		Quantities: []quantities.Quantity{&temp}}

	AddCustomEvent(1001, 1, doA)

	doAexecuted = false
	doBexecuted = false
	EventInterpreter(event)
	if doAexecuted != true {
		t.Errorf("doAexecuted was not executed")
	}
	if doBexecuted != false {
		t.Errorf("doBexecuted was executed")
	}

	event = eventmessage.EventMessage{ClientID: 1001,
		EventID:    2,
		Quantities: []quantities.Quantity{&temp}}

	doAexecuted = false
	doBexecuted = false
	EventInterpreter(event)
	if doAexecuted != false {
		t.Errorf("doAexecuted was executed")
	}
	if doBexecuted != false {
		t.Errorf("doBexecuted was executed")
	}

	event = eventmessage.EventMessage{ClientID: 1002,
		EventID:    1,
		Quantities: []quantities.Quantity{&temp}}

	doAexecuted = false
	doBexecuted = false
	EventInterpreter(event)
	if doAexecuted != false {
		t.Errorf("doAexecuted was executed")
	}
	if doBexecuted != false {
		t.Errorf("doBexecuted was executed")
	}

	AddCustomEvent(1002, 2, doB)

	event = eventmessage.EventMessage{ClientID: 1002,
		EventID:    2,
		Quantities: []quantities.Quantity{&temp}}

	doAexecuted = false
	doBexecuted = false
	EventInterpreter(event)
	if doAexecuted != false {
		t.Errorf("doAexecuted was executed")
	}
	if doBexecuted != true {
		t.Errorf("doBexecuted was not executed")
	}

	RemoveCustomEvent(1002, 2)

	doAexecuted = false
	doBexecuted = false
	EventInterpreter(event)
	if doAexecuted != false {
		t.Errorf("doAexecuted was executed")
	}
	if doBexecuted != false {
		t.Errorf("doBexecuted was executed")
	}

}
