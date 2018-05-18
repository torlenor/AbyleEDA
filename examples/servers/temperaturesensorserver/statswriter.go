package main

import (
	"fmt"
	"time"

	"github.com/torlenor/AbyleEDA/AEDAserver"
)

var starttime = time.Now()

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
