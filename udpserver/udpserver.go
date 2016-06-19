package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAevents"
	"github.com/torlenor/AbyleEDA/AEDAserver"
	"os"
	"os/signal"
	"time"
)

// This is for go-logger
var log = logging.MustGetLogger("AEDAlogger")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`,
)
var starttime = time.Now()

// Simple OK/NOTOK for the client
var rcvOK []byte = []byte("0")
var rcvFAIL []byte = []byte("1")

func CheckError(err error) {
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

type Config struct {
	debugMode bool
	port      int
}

var cfg Config

func parseCmdLine() {
	numbPtr := flag.Int("port", 10001, "server port to listen on")
	boolPtr := flag.Bool("debug", false, "debug output")

	flag.Parse()

	cfg.debugMode = *boolPtr
	cfg.port = *numbPtr
}

type SrvStats struct {
	pktssentcnt int
	pktsrecvcnt int
	pktserrcnt  int
}

var stats SrvStats

func prepLogging() {
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	// backend2, err := logging.NewSyslogBackend("AbyleEDA")
	// CheckError(err)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)
}

func main() {
	// Prepare logging with go-logging
	prepLogging()

	// Command line flags parsing
	parseCmdLine()

	// Create an AEDA UDP server
	srv, err := AEDAserver.CreateUDPServer(cfg.port)
	CheckError(err)

	initStatsWrite(srv)

	// Prepare interupt handling
	initInterruptHandling(srv)

	if cfg.debugMode {
		log.Debug("Debug mode on")
		srv.DebugMode = true
	}

	log.Info("AbyleEDA server prepared on", srv.Addr)
	log.Info("Starting to listen...")
	go AEDAserver.Start(srv) // start the server in an own thread

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

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

			if event.Id != 0 {
				AEDAevents.EventInterpreter(event)
			}

		// or quit if os.Interrupt
		case <-c:
			break
		}
	}

	srv.Conn.Close()
}
