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

func CheckError(err error) {
	if err != nil {
		log.Error("Error: ", err)
	}
}

func prepLogging() {
	backend2 := logging.NewLogBackend(os.Stdout, "", 0)
	// backend3, _ := logging.NewSyslogBackend("AbyleEDA")

	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	logging.SetBackend(backend2Formatter)
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
	var srvPort string = strconv.Itoa(*srvPortPtr)
	ServerAddr, err := net.ResolveUDPAddr("udp", *srvAddrPtr+":"+srvPort)
	CheckError(err)

	client, err := AEDAclient.ConnectUDPClient(ServerAddr)
	CheckError(err)
	defer AEDAclient.DisconnectUDPClient(client)

	// Send JSON stuff
	for {
		var eventEvent int32
		if mrand.Float64() > 0.5 {
			eventEvent = 0
		} else {
			eventEvent = 1
		}
		event := AEDAevents.EventMessage{Id: 1001,
			Value:    float64(mrand.Intn(250)) + mrand.Float64(),
			Type:     "sensor",
			Event:    eventEvent,
			Quantity: "temperature",
			Unit:     "degC"}
		jsonMapA, _ := json.Marshal(event)

		msg := jsonMapA

		AEDAclient.SendMessageToServer(client, msg)

		time.Sleep(time.Second * 1)
	}
}
