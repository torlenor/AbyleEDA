package main

import (
	"encoding/hex"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/torlenor/AbyleEDA/AEDAclient"
	"github.com/torlenor/AbyleEDA/AEDAcrypt"
	"github.com/torlenor/AbyleEDA/eventmessage"
	"github.com/torlenor/AbyleEDA/quantities"
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

	clientIDPtr := flag.Int("clientid", 1001, "sensor id")
	sensorIDPtr := flag.Int("sensorid", 1, "sensor id")
	sensorFilePtr := flag.String("sensorfile", "/sys/class/hwmon/hwmon0/temp2_input", "sensor file in sys interface")
	updateIntervalPtr := flag.Float64("updateInterval", 1, "Update Interval for the sensor data in seconds")

	flag.Parse()

	nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")
	ccfg := AEDAcrypt.CryptCfg{Key: []byte("AES256Key-32Characters1234567890"),
		Nonce: nonce}

	client, err := AEDAclient.ConnectUDPClient(getServerAddress(), ccfg)
	checkError(err)
	defer AEDAclient.DisconnectUDPClient(client)

	// Send JSON stuff
	for {
		val := getSensorValue(*sensorFilePtr)

		var t quantities.Temperature
		t.FromFloat(val)

		event := eventmessage.EventMessage{ClientID: int32(*clientIDPtr),
			EventID:    int32(*sensorIDPtr),
			Quantities: []quantities.Quantity{&t}}

		AEDAclient.SendMessageToServer(client, event)

		time.Sleep(time.Millisecond * time.Duration(*updateIntervalPtr*1000))
	}
}
