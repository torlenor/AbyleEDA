package quantities

import (
	"encoding/json"
	"strconv"
)

// Temperature Unit definition
const ( // iota is reset to 0
	Unknown    = iota // = 0
	DegreeC    = iota // = 1
	Fahrenheit = iota // = 2
	Kelvin     = iota // = 3
)

// Temperature quantity
type Temperature struct {
	Val  float64
	Unit uint
}

func (t *Temperature) MarshalJSON() (b []byte, e error) {
	return json.Marshal(map[string]string{
		"type":  "temperature",
		"value": t.String(),
	})
}

func (t *Temperature) UnmarshalJSON(b []byte) error {
	var m map[string]string
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	for key, value := range m {
		if key == "value" {
			t.FromString(value)
		}
	}

	return nil
}

// String returns the value as a string
func (t *Temperature) String() string {
	return strconv.FormatFloat(t.Val, 'f', -1, 64)
}

// FromString converts a provided numerical value as a string into
// a float and stores it inside t
func (t *Temperature) FromString(valstr string) {
	f, err := strconv.ParseFloat(valstr, 64)
	if err == nil {
		t.Val = f
	}
}

// FromFloat stores the provided value inside t
func (t *Temperature) FromFloat(val float64) {
	t.Val = val
}
