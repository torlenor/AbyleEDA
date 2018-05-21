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
	"github.com/torlenor/AbyleEDA/AEDAserver"
	"github.com/torlenor/AbyleEDA/eventmessage"
	"github.com/torlenor/AbyleEDA/quantities"
)

// This is for go-logger
var log = logging.MustGetLogger("AEDAlogger")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`,
)

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

type config struct {
	debugMode bool
	port      int
	ccfg      AEDAcrypt.CryptCfg
}

var cfg config

func parseCmdLine() {
	numbPtr := flag.Int("port", 10002, "server port to listen on")
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

func printEvent(event eventmessage.EventMessage) {
	log.Info("ClientID =", event.ClientID)
	log.Info("EventID =", event.EventID)
	log.Info("Timestamp =", event.Timestamp, "(", time.Unix(0, event.Timestamp), ")")
	cnt := 0
	for _, content := range event.Quantities {
		cnt++

		switch v := content.(type) {
		case *quantities.Ping:
			log.Info("Content (", v.Type(), ")", cnt, ": Hostname =", v.HostName,
				"IPAddr =", v.IPAddr,
				"ResponseTime = ", v.ResponseTime)
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

	for {
		select {
		// Fetch messages from AEDAserver
		case clientMsg := <-srv.ResQueue:
			printEvent(clientMsg.Event)
		// or quit if os.Interrupt
		case <-c:
			break
		}
	}

	srv.Close()
}
