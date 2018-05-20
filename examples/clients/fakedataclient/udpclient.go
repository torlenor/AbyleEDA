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
	"github.com/torlenor/AbyleEDA/eventmessage"
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

func receiveFromServer(client *AEDAclient.UDPClient) {
	for {
		select {
		// Fetch messages from AEDAclient
		case clientMsg := <-client.ResQueue:
			var event eventmessage.EventMessage
			if err := json.Unmarshal(clientMsg.Msg, &event); err != nil {
				log.Error("Error in decoding JSON:", err)
				continue
			}

			AEDAevents.EventInterpreter(event)
		}
	}
}

func main() {
	// Prepare logging with go-logging
	prepLogging()

	// Command line flags parsing
	srvAddrPtr := flag.String("srvaddr", "127.0.0.1", "server address")
	srvPortPtr := flag.Int("port", 10001, "server port")
	clientIDPtr := flag.Int("clientid", 2001, "client id")

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
		var t1 quantities.Temperature
		t1.FromFloat(float64(mrand.Intn(250)) + mrand.Float64())
		var t2 quantities.Temperature
		t2.FromFloat(float64(mrand.Intn(250)) + mrand.Float64())

		event := eventmessage.EventMessage{ClientID: int32(*clientIDPtr),
			EventID:    1,
			Quantities: []quantities.Quantity{&t1, &t2}}

		AEDAclient.SendMessageToServer(client, event)

		time.Sleep(time.Second * 1)
	}

}
