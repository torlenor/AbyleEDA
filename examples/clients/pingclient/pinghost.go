package main

import ping "github.com/SewanDevs/go-ping"

func pingHost(host string) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(host)
	pinger.SetPrivileged(true)
	if err != nil {
		return nil, err
	}

	pinger.Count = 3
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats

	return stats, nil
}
