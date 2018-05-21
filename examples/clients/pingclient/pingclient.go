package main

import (
	"flag"
	mrand "math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAclient"
	"github.com/torlenor/AbyleEDA/eventmessage"
	"github.com/torlenor/AbyleEDA/quantities"
)

// This is for go-logger
var log = logging.MustGetLogger("AEDAlogger")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`)

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

func main() {
	// Prepare logging with go-logging
	prepLogging()

	// Command line flags parsing
	srvAddrPtr := flag.String("srvaddr", "127.0.0.1", "server address")
	srvPortPtr := flag.Int("port", 10002, "server port")
	clientIDPtr := flag.Int("clientid", 3001, "client id")

	flag.Parse()

	mrand.Seed(time.Now().UnixNano())

	// Define the server address and port
	var srvPort = strconv.Itoa(*srvPortPtr)
	ServerAddr, err := net.ResolveUDPAddr("udp", *srvAddrPtr+":"+srvPort)
	checkError(err)

	client, err := AEDAclient.ConnectUDPClient(ServerAddr)
	checkError(err)
	defer AEDAclient.DisconnectUDPClient(client)

	urls := []string{
		"www.orf.at",
		"www.kernel.org",
		"www.heise.de",
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			os.Exit(0)
		case syscall.SIGINT:
			os.Exit(0)
		}
	}()

	// pingResponses := make(chan *ping.Statistics)

	for {
		event := eventmessage.EventMessage{ClientID: int32(*clientIDPtr),
			EventID:    1, // host group nummer
			Quantities: []quantities.Quantity{},
		}

		// var wg sync.WaitGroup

		checkError(err)

		// for _, url := range urls {
		// 	wg.Add(1)
		// 	go func(url string) {
		// 		defer wg.Done()
		// 		fmt.Printf(url + "\n")
		// 		stats, err := pingHost(url)
		// 		fmt.Printf("done")
		// 		if err != nil {
		// 			log.Error(err)
		// 		} else {
		// 			pingResponses <- stats
		// 		}
		// 	}(url)
		// }

		// go func() {
		// 	for stat := range pingResponses {
		// 		var pingResults quantities.Ping
		// 		pingResults.DidAnswer = true
		// 		pingResults.ResponseTime = stat.AvgRtt
		// 		pingResults.HostName = stat.Addr
		// 		pingResults.IPAddr = stat.IPAddr.String()
		// 		event.Quantities = append(event.Quantities, &pingResults)
		// 	}
		// }()

		// wg.Wait()

		for _, url := range urls {
			// wg.Add(1)
			stats, err := pingHost(url)
			if err != nil {
				log.Error(err)
			} else {
				var pingResults quantities.Ping
				pingResults.DidAnswer = true
				pingResults.ResponseTime = stats.AvgRtt
				pingResults.HostName = stats.Addr
				pingResults.IPAddr = stats.IPAddr.String()
				event.Quantities = append(event.Quantities, &pingResults)
			}
		}

		AEDAclient.SendMessageToServer(client, event)

		time.Sleep(time.Second * 1)
	}

}
