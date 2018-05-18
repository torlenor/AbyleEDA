package quantities

import (
	"encoding/json"
	"strconv"
)

// Temperature quantity
type Temperature struct {
	Val float64
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

// Type returns the type of the quantity
func (t *Temperature) Type() string {
	return "temperature"
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

// Float returns the value as a float64
func (t *Temperature) Float() float64 {
	return t.Val
}

// Degrees returns the temperature value in degrees centigrate
func (t *Temperature) Degrees() float64 {
	return t.Val
}

// Fahrenheit returns the temperature value in fahrenheit
func (t *Temperature) Fahrenheit() float64 {
	return t.Val*1.8 + 32
}

// Kelvin returns the temperature value in kelvin
func (t *Temperature) Kelvin() float64 {
	return t.Val + 273.15
}
