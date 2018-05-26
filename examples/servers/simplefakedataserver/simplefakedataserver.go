package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/op/go-logging"

	"github.com/torlenor/AbyleEDA/AEDAcrypt"
	"github.com/torlenor/AbyleEDA/AEDAevents"
	"github.com/torlenor/AbyleEDA/AEDAserver"
	"github.com/torlenor/AbyleEDA/eventmessage"
	"github.com/torlenor/AbyleEDA/quantities"
)

// This is for go-logger
var log = logging.MustGetLogger("AEDAlogger")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`,
)

type config struct {
	srvPort int
	ccfg    AEDAcrypt.CryptCfg
}

var cfg config

func parseCmdLine() {
	srvPortPtr := flag.Int("port", 20001, "server port to listen on")
	flag.Parse()

	cfg.srvPort = *srvPortPtr

	nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")
	cfg.ccfg = AEDAcrypt.CryptCfg{Key: []byte("AES256Key-32Characters1234567890"),
		Nonce: nonce}
}

func prepLogging() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

func setupEvents() {
	// Register custom event callbacks with the events system
	AEDAevents.AddCustomEvent(101, 1, printEvent)
}

func printEvent(event eventmessage.EventMessage) {
	log.Info("ClientID =", event.ClientID)
	log.Info("EventID =", event.EventID)
	log.Info("Timestamp =", event.Timestamp, "(", time.Unix(0, event.Timestamp), ")")
	cnt := 0
	for _, content := range event.Quantities {
		cnt++

		switch v := content.(type) {
		case *quantities.Floating:
			log.Info("Content (", v.Type(), ")", cnt, ":", v.Float())
		default:
			log.Info("Content (", v.Type(), ")", cnt, ":", content.String())
		}
	}
}

func main() {
	// Prepare logging with go-logging
	prepLogging()

	// Command line flags parsing
	parseCmdLine()

	// Create a channel to handle SigINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Create an AEDA UDP server
	srv, err := AEDAserver.CreateUDPServer(cfg.srvPort, cfg.ccfg)
	if err != nil {
		log.Error(err.Error())
		os.Exit(0)
	}
	defer srv.Close()

	log.Info("AbyleEDA server prepared on", srv.Addr)
	log.Info("Starting to listen...")
	go srv.Start() // start the server in an own thread

	// Server has to be set in AEDAevents to make it possible
	// to send messages to clients in event system
	AEDAevents.SetAEDAserver(srv)

	// Setup custom event callbacks
	setupEvents()

	for {
		select {
		// Fetch messages from AEDAserver
		case clientMsg := <-srv.ResQueue:
			AEDAevents.EventInterpreter(clientMsg.Event)
		// or quit if os.Interrupt
		case <-c:
			fmt.Printf("Server quit\n")
			os.Exit(0)
		}
	}
}
