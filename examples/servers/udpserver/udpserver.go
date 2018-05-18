package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAcrypt"
	"github.com/torlenor/AbyleEDA/AEDAevents"
	"github.com/torlenor/AbyleEDA/AEDAserver"
	"github.com/torlenor/AbyleEDA/quantities"
)

// This is for go-logger
var log = logging.MustGetLogger("AEDAlogger")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`,
)
var starttime = time.Now()

// Simple OK/NOTOK for the client
var rcvOK = []byte("0")
var rcvFAIL = []byte("1")

func checkError(err error) {
	if err != nil {
		log.Error(err.Error())
		os.Exit(0)
	}
}

func initInterruptHandling(srv *AEDAserver.UDPServer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("")
		writeStatsToStdout(srv)
		writeStatsToLog(srv)
		os.Exit(1)
	}()
}

func writeStatsToStdout(srv *AEDAserver.UDPServer) {
	uptime := time.Now().Sub(starttime)
	fmt.Println(srv.Stats.Pktsrecvcnt, "pkts sent,", srv.Stats.Pktsrecvcnt, "pkts received,", srv.Stats.Pktserrcnt, "pkts with errors")
	fmt.Println("Uptime:", uptime.String())
}

func writeStatsToLog(srv *AEDAserver.UDPServer) {
	uptime := time.Now().Sub(starttime)
	log.Info(srv.Stats.Pktsrecvcnt, "pkts sent,", srv.Stats.Pktsrecvcnt, "pkts received,", srv.Stats.Pktserrcnt, "pkts with errors")
	log.Info("Uptime:", uptime.String())
}

func initStatsWrite(srv *AEDAserver.UDPServer) {
	go func() {
		for {
			time.Sleep(time.Second * 120)
			writeStatsToLog(srv)
		}
	}()
}

type config struct {
	debugMode bool
	port      int
	ccfg      AEDAcrypt.CryptCfg
}

var cfg config

func parseCmdLine() {
	numbPtr := flag.Int("port", 10001, "server port to listen on")
	boolPtr := flag.Bool("debug", false, "debug output")

	flag.Parse()

	cfg.debugMode = *boolPtr
	cfg.port = *numbPtr

	nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")
	cfg.ccfg = AEDAcrypt.CryptCfg{Key: []byte("AES256Key-32Characters1234567890"),
		Nonce: nonce}
}

func prepLogging() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

type page struct {
	Title string
	Body  []byte
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p := page{Title: "Hello world"}
	t, _ := template.ParseFiles("sensor.html")
	t.Execute(w, p)
}

func webShowSensors(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "<h1>Sensors</h1>")

	for clientID, sensors := range sensorsMap {
		fmt.Fprintf(w, "<h2>Client ID: %d</h2>", clientID)
		for sensorID, sensor := range sensors {
			fmt.Fprintf(w, "\tSensor ID: %d<br>", sensorID)
			fmt.Fprintf(w, "\tValue = %.2f<br>", sensor.Value)
		}
	}
}

func startWebServer() {
	// Setup simple web server
	log.Info("Starting web server on port 10080")
	http.HandleFunc("/", webShowSensors)  // set router
	go http.ListenAndServe(":10080", nil) // set listen port
}

// Sensor struct is used to store sensor data
type Sensor struct {
	ID          int32
	SensorType  string
	Quantity    string
	Value       float64
	Unit        string
	LastUpdated time.Time
}

var sensorsMap = map[int32]map[int32]Sensor{}

func updateSensorValue(event AEDAevents.EventMessage) {
	if sensorsMap[event.ClientID] == nil {
		sensorsMap[event.ClientID] = map[int32]Sensor{}
	}

	for _, content := range event.Quantities {
		if _, ok := sensorsMap[event.ClientID][event.EventID]; ok {
			log.Infof("Received sensor update clientID: %d, sensorID: %d", event.ClientID, event.EventID)

			switch v := content.(type) {
			case *quantities.Temperature:
				log.Info("Content (numeric):", v.Degrees(), "°C")
				sensorsMap[event.ClientID][event.EventID] = Sensor{ID: event.EventID,
					SensorType: "temperature",
					Quantity:   v.Type(),
					Value:      v.Degrees(),
					Unit:       "°C"}
			default:
				log.Info("Content:", content.String())
			}
		} else {
			log.Infof("Registering new sensor clientID: %d, sensorID: %d", event.ClientID, event.EventID)

			switch v := content.(type) {
			case *quantities.Temperature:
				log.Info("Content (numeric):", v.Degrees(), "°C")
				sensorsMap[event.ClientID][event.EventID] = Sensor{ID: event.EventID,
					SensorType: "temperature",
					Quantity:   v.Type(),
					Value:      v.Degrees(),
					Unit:       "°C"}
			default:
				log.Info("Content:", content.String())
			}
		}
	}
}

func doA(i int) {
	log.Infof("[doA]: I'm param is: %d\n", i)
}

func main() {
	// Prepare logging with go-logging
	prepLogging()

	// Command line flags parsing
	parseCmdLine()

	// Create an AEDA UDP server
	srv, err := AEDAserver.CreateUDPServer(cfg.port, cfg.ccfg)
	checkError(err)

	initStatsWrite(srv)

	// Prepare interupt handling
	initInterruptHandling(srv)

	// Start web server
	startWebServer()

	if cfg.debugMode {
		log.Debug("Debug mode on")
		srv.DebugMode = true
	}

	log.Info("AbyleEDA server prepared on", srv.Addr)
	log.Info("Starting to listen...")
	go srv.Start() // start the server in an own thread

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Server has to be set in AEDAevents to make it possible
	// to send messages to clients in event system
	AEDAevents.SetAEDAserver(srv)

	// Register custom event callbacks
	AEDAevents.AddCustomEvent(1001, 1, updateSensorValue)
	AEDAevents.AddCustomEvent(1002, 1, updateSensorValue)

	for {
		select {
		// Fetch messages from AEDAserver
		case clientMsg := <-srv.ResQueue:
			var event AEDAevents.EventMessage
			if err := json.Unmarshal(clientMsg.Msg, &event); err != nil {
				log.Error("Error in decoding JSON:", err)
				continue
			}

			if cfg.debugMode {
				log.Debug(clientMsg.Addr, "sent:")
				log.Debug(string(clientMsg.Msg))
				log.Debug("------------------------------------------------------------")
			}

			AEDAevents.EventInterpreter(event)

		// or quit if os.Interrupt
		case <-c:
			break
		}
	}

	srv.Close()
}
