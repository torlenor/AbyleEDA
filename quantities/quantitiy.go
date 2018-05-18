package quantities

type Quantity interface {
	String() string
	FromString(valstr string)
}

type NumericalQuantity interface {
	Quantity
	FromFloat(val float64)
}
