package main

import (
	"encoding/json"
	"flag"
	mrand "math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAclient"
	"github.com/torlenor/AbyleEDA/AEDAevents"
	"github.com/torlenor/AbyleEDA/quantities"
)

// This is for go-logger
var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`,
)

var rcvOK = []byte("0")
var rcvFAIL = []byte("1")

func checkError(err error) {
	if err != nil {
		log.Error("Error: ", err)
	}
}

func prepLogging() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

func main() {
	// Prepare logging with go-logging
	prepLogging()

	// Command line flags parsing
	srvAddrPtr := flag.String("srvaddr", "127.0.0.1", "server address")
	srvPortPtr := flag.Int("port", 10001, "server port")

	flag.Parse()

	mrand.Seed(time.Now().UnixNano())

	// Define the server address and port
	var srvPort = strconv.Itoa(*srvPortPtr)
	ServerAddr, err := net.ResolveUDPAddr("udp", *srvAddrPtr+":"+srvPort)
	checkError(err)

	client, err := AEDAclient.ConnectUDPClient(ServerAddr)
	checkError(err)
	defer AEDAclient.DisconnectUDPClient(client)

	// Send JSON stuff
	for {
		var eventEvent int32
		if mrand.Float64() > 0.5 {
			eventEvent = 0
		} else {
			eventEvent = 1
		}

		var t1 quantities.Temperature
		t1.FromFloat(float64(mrand.Intn(250)) + mrand.Float64())
		var t2 quantities.Temperature
		t2.FromFloat(float64(mrand.Intn(250)) + mrand.Float64())

		event := AEDAevents.EventMessage{Id: 1002,
			Type:       "sensor",
			Event:      eventEvent,
			Quantities: []quantities.Quantity{&t1, &t2}}

		msgng, _ := json.Marshal(event)

		AEDAclient.SendMessageToServer(client, msgng)

		time.Sleep(time.Second * 1)
	}
}
