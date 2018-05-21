package main

import (
	"flag"
	"fmt"
	mrand "math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
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
	pingCountPtr := flag.Int("pingcount", 3, "number of pings per iteration")
	hostsPtr := flag.String("hosts", "", "hosts separated by , (comma)")
	waitBetweenPingsPtr := flag.Float64("waitInterval", 1, "Wait between ping iterations")

	flag.Parse()

	if len(*hostsPtr) == 0 {
		fmt.Printf("Please specifiy hosts to ping with -hosts flag\n")
		os.Exit(1)
	}

	mrand.Seed(time.Now().UnixNano())

	// Define the server address and port
	var srvPort = strconv.Itoa(*srvPortPtr)
	ServerAddr, err := net.ResolveUDPAddr("udp", *srvAddrPtr+":"+srvPort)
	checkError(err)

	client, err := AEDAclient.ConnectUDPClient(ServerAddr)
	checkError(err)
	defer AEDAclient.DisconnectUDPClient(client)

	hostUrls := strings.Replace(*hostsPtr, " ", "", -1)
	urls := strings.Split(hostUrls, ",")

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

	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup

	for {
		event := eventmessage.EventMessage{ClientID: int32(*clientIDPtr),
			EventID:    1, // host group nummer
			Quantities: []quantities.Quantity{},
		}

		checkError(err)

		wg.Add(len(urls))
		for _, url := range urls {
			go func(url string) {
				defer wg.Done()
				stats, err := pingHost(url, *pingCountPtr)
				if err != nil {
					log.Error(err)
					var pingResults quantities.Ping
					pingResults.HostName = url
					mutex.Lock()
					event.Quantities = append(event.Quantities, &pingResults)
					mutex.Unlock()
				} else {
					var pingResults quantities.Ping
					pingResults.HostName = stats.Addr
					pingResults.IPAddr = stats.IPAddr.String()
					pingResults.PacketsRecv = int64(stats.PacketsRecv)
					pingResults.PacketsSent = int64(stats.PacketsSent)
					pingResults.MinRtt = stats.MinRtt
					pingResults.MaxRtt = stats.MaxRtt
					pingResults.AvgRtt = stats.AvgRtt
					pingResults.StdDevRtt = stats.StdDevRtt
					mutex.Lock()
					event.Quantities = append(event.Quantities, &pingResults)
					mutex.Unlock()
				}
			}(url)
		}

		wg.Wait()

		AEDAclient.SendMessageToServer(client, event)

		time.Sleep(time.Millisecond * time.Duration(*waitBetweenPingsPtr*1000))
	}

}
