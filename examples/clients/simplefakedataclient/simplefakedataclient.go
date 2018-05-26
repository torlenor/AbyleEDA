package main

import (
	"encoding/hex"
	"flag"
	mrand "math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAclient"
	"github.com/torlenor/AbyleEDA/AEDAcrypt"
	"github.com/torlenor/AbyleEDA/eventmessage"
	"github.com/torlenor/AbyleEDA/quantities"
)

// This is for go-logger
var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`,
)

func checkError(err error) {
	if err != nil {
		log.Error(err.Error())
		os.Exit(0)
	}
}

type config struct {
	srvaddr  string
	srvport  int
	clientid int
	ccfg     AEDAcrypt.CryptCfg
}

var cfg config

func parseCmdLine() {
	srvAddrPtr := flag.String("srvaddr", "127.0.0.1", "server address")
	srvPortPtr := flag.Int("port", 20001, "server port")
	clientIDPtr := flag.Int("clientid", 101, "client id")
	flag.Parse()

	cfg.srvaddr = *srvAddrPtr
	cfg.srvport = *srvPortPtr
	cfg.clientid = *clientIDPtr

	nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")
	cfg.ccfg = AEDAcrypt.CryptCfg{Key: []byte("AES256Key-32Characters1234567890"),
		Nonce: nonce}
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
	parseCmdLine()

	// We want to send always "new" random data
	mrand.Seed(time.Now().UnixNano())

	// Define the server address and port
	serverAddr, err := net.ResolveUDPAddr("udp", cfg.srvaddr+":"+strconv.Itoa(cfg.srvport))
	checkError(err)

	client, err := AEDAclient.ConnectUDPClient(serverAddr, cfg.ccfg)
	checkError(err)
	defer AEDAclient.DisconnectUDPClient(client)

	// Send two quantity entries every time
	var t1 quantities.Floating
	var t2 quantities.Floating
	for {
		t1.FromFloat(float64(mrand.Intn(250)) + mrand.Float64())
		t2.FromFloat(float64(mrand.Intn(250)) + mrand.Float64())

		event := eventmessage.EventMessage{ClientID: int32(cfg.clientid),
			EventID:    1,
			Quantities: []quantities.Quantity{&t1, &t2}}

		AEDAclient.SendMessageToServer(client, event)

		time.Sleep(time.Second * 1)
	}
}
