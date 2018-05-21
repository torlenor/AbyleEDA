package quantities

import (
	"encoding/json"
	"strconv"
	"time"
)

// Ping quantity
type Ping struct {
	HostName    string
	IPAddr      string
	PacketsRecv int64
	PacketsSent int64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	StdDevRtt   time.Duration
}

// Type returns the type of the quantity
func (p *Ping) Type() string {
	return "ping"
}

func (p *Ping) MarshalJSON() (b []byte, e error) {
	return json.Marshal(map[string]string{
		"type":        p.Type(),
		"hostname":    p.HostName,
		"ipaddr":      p.IPAddr,
		"packetsrecv": strconv.FormatInt(p.PacketsRecv, 10),
		"packetssent": strconv.FormatInt(p.PacketsSent, 10),
		"minrtt":      p.MinRtt.String(),
		"maxrtt":      p.MaxRtt.String(),
		"avgrtt":      p.AvgRtt.String(),
		"stddevrtt":   p.StdDevRtt.String(),
	})
}

func (p *Ping) UnmarshalJSON(b []byte) error {
	var m map[string]string
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	for key, value := range m {
		if key == "hostname" {
			p.HostName = value
		}
		if key == "ipaddr" {
			p.IPAddr = value
		}
		if key == "ipaddr" {
			p.IPAddr = value
		}

		if key == "packetsrecv" {
			packetsRecv, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			p.PacketsRecv = packetsRecv
		}

		if key == "packetssent" {
			packetsSent, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			p.PacketsSent = packetsSent
		}

		if key == "minrtt" {
			minRtt, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			p.MinRtt = minRtt
		}
		if key == "maxrtt" {
			maxRtt, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			p.MaxRtt = maxRtt
		}
		if key == "avgrtt" {
			avgRtt, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			p.AvgRtt = avgRtt
		}
		if key == "stddevrtt" {
			stdDevRtt, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			p.StdDevRtt = stdDevRtt
		}
	}

	return nil
}

// String returns the value as a string
func (p *Ping) String() string {
	return p.AvgRtt.String()
}

// FromString converts a provided numerical value as a string into
// a float and stores it inside t
func (p *Ping) FromString(valstr string) {
	f, err := time.ParseDuration(valstr)
	if err == nil {
		p.AvgRtt = f
	}
}
