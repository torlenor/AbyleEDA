package quantities

import (
	"encoding/json"
	"strconv"
	"time"
)

// Ping quantity
type Ping struct {
	HostName     string
	IPAddr       string
	DidAnswer    bool
	ResponseTime time.Duration
}

// Type returns the type of the quantity
func (p *Ping) Type() string {
	return "ping"
}

func (p *Ping) MarshalJSON() (b []byte, e error) {
	return json.Marshal(map[string]string{
		"type":         p.Type(),
		"hostname":     p.HostName,
		"ipaddr":       p.IPAddr,
		"didanswer":    strconv.FormatBool(p.DidAnswer),
		"responsetime": p.ResponseTime.String(),
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
		if key == "didanswer" {
			didAnswer, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			p.DidAnswer = didAnswer
		}
		if key == "responsetime" {
			responseTime, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			p.ResponseTime = responseTime
		}
	}

	return nil
}

// String returns the value as a string
func (p *Ping) String() string {
	return p.ResponseTime.String()
}

// FromString converts a provided numerical value as a string into
// a float and stores it inside t
func (p *Ping) FromString(valstr string) {
	f, err := time.ParseDuration(valstr)
	if err == nil {
		p.ResponseTime = f
	}
}
