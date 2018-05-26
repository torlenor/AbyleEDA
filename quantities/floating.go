package quantities

import (
	"encoding/json"
	"strconv"
)

// Floating point quantity
type Floating struct {
	Val float64
}

func (t *Floating) MarshalJSON() (b []byte, e error) {
	return json.Marshal(map[string]string{
		"type":  "floating",
		"value": t.String(),
	})
}

func (t *Floating) UnmarshalJSON(b []byte) error {
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

// Type returns the type of the quantity
func (t *Floating) Type() string {
	return "floating"
}

// String returns the value as a string
func (t *Floating) String() string {
	return strconv.FormatFloat(t.Val, 'f', -1, 64)
}

// FromString converts a provided numerical value as a string into
// a float and stores it inside t
func (t *Floating) FromString(valstr string) {
	f, err := strconv.ParseFloat(valstr, 64)
	if err == nil {
		t.Val = f
	}
}

// FromFloat stores the provided value inside t
func (t *Floating) FromFloat(val float64) {
	t.Val = val
}

// Float returns the value as a float64
func (t *Floating) Float() float64 {
	return t.Val
}
