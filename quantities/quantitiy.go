package quantities

type Quantity interface {
	Type() string
	String() string
	FromString(valstr string)
}

type NumericalQuantity interface {
	Quantity
	Float() float64
	FromFloat(val float64)
}
