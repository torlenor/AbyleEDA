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
        
        val1 := float64(mrand.Intn(250)) + mrand.Float64();
        val1str := strconv.FormatFloat(val1, 'f', -1, 64)
        val2 := float64(mrand.Intn(250)) + mrand.Float64();
        val2str := strconv.FormatFloat(val2, 'f', -1, 64)
            
        event := AEDAevents.EventMessage{Id: 1002,
			Type:     "sensor",
		 	Event:    eventEvent,
            Content: []AEDAevents.EventContent {
                            AEDAevents.EventContent{ Quantity: "temperature", Value: val1str, Unit: "degC" },
                            AEDAevents.EventContent{ Quantity: "temperature", Value: val2str, Unit: "degC" } } }
        
        msgng, _ := json.Marshal(event)

        AEDAclient.SendMessageToServer(client, msgng)

		time.Sleep(time.Second * 1)
	}
}
