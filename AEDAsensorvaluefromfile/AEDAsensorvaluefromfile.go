package main

import (
	"encoding/json"
	"flag"
	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAclient"
	"github.com/torlenor/AbyleEDA/AEDAevents"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// This is for go-logger
var log = logging.MustGetLogger("AEDAlogger")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`)

func CheckError(err error) {
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
	var srvPort string = strconv.Itoa(*srvPortPtr)
	ServerAddr, err := net.ResolveUDPAddr("udp", *srvAddrPtr+":"+srvPort)
	CheckError(err)

	return ServerAddr
}

func main() {
	prepLogging()

	sensorFilePtr := flag.String("sensorfile", "/sys/class/hwmon/hwmon0/temp2_input", "senosr file in sys interface")

	flag.Parse()

	client, err := AEDAclient.ConnectUDPClient(getServerAddress())
	CheckError(err)
	defer AEDAclient.DisconnectUDPClient(client)

	// Send JSON stuff
	for {
		event := AEDAevents.EventMessage{Id: 1001,
			Value:    getSensorValue(*sensorFilePtr),
			Type:     "sensor",
			Event:    AEDAevents.EventValueUpdate,
			Quantity: "temperature",
			Unit:     "degC"}

		msg, err := json.Marshal(event)
		CheckError(err)

		AEDAclient.SendMessageToServer(client, msg)

		time.Sleep(time.Second * 1)
	}
}
