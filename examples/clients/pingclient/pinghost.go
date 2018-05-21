package main

import ping "github.com/SewanDevs/go-ping"

func pingHost(host string, pingCount int) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, err
	}
	pinger.SetPrivileged(true)

	pinger.Count = pingCount
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats

	return stats, nil
}
