package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/torlenor/AbyleEDA/quantities"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAclient"
	"github.com/torlenor/AbyleEDA/AEDAevents"
)

// This is for go-logger
var log = logging.MustGetLogger("AEDAlogger")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`)

func checkError(err error) {
	if err != nil {
		log.Error("Error: ", err)
	}
}

func prepLogging() {
	backend2 := logging.NewLogBackend(os.Stdout, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)
}

func getSensorValue(sensorfile string) float64 {
	out, err := ioutil.ReadFile(sensorfile)
	if err != nil {
		log.Fatal(err)
	}

	outstr := strings.Replace(string(out), "\n", "", -1)

	val, _ := strconv.ParseFloat(outstr, 64)
	val = val / 1000.0

	return val
}

// Command line: server address and port
var srvAddrPtr = flag.String("srvaddr", "127.0.0.1", "server address")
var srvPortPtr = flag.Int("port", 10001, "server port")

func getServerAddress() *net.UDPAddr {
	// Define the server address and port
	var srvPort = strconv.Itoa(*srvPortPtr)
	ServerAddr, err := net.ResolveUDPAddr("udp", *srvAddrPtr+":"+srvPort)
	checkError(err)

	return ServerAddr
}

func main() {
	prepLogging()

	sensorIDPtr := flag.Int("sensorid", 1001, "sensor id")
	sensorFilePtr := flag.String("sensorfile", "/sys/class/hwmon/hwmon0/temp2_input", "sensor file in sys interface")

	flag.Parse()

	client, err := AEDAclient.ConnectUDPClient(getServerAddress())
	checkError(err)
	defer AEDAclient.DisconnectUDPClient(client)

	// Send JSON stuff
	for {
		val := getSensorValue(*sensorFilePtr)
		// valstr := strconv.FormatFloat(val, 'f', -1, 64)

		var t quantities.Temperature
		t.FromFloat(val)

		event := AEDAevents.EventMessage{Id: int32(*sensorIDPtr),
			Type:       "sensor",
			Event:      AEDAevents.EventValueUpdate,
			Quantities: []quantities.Quantity{&t}}
		// Content: []AEDAevents.EventContent{
		// 	AEDAevents.EventContent{Quantity: "temperature", Value: valstr, Unit: "degC"},
		// 	AEDAevents.EventContent{Quantity: "temperature", Value: valstr, Unit: "degC"}}}

		msg, err := json.Marshal(event)
		checkError(err)

		AEDAclient.SendMessageToServer(client, msg)

		time.Sleep(time.Second * 1)
	}
}
