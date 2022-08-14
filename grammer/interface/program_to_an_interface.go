package _interface

import "fmt"

type Country struct {
	Name string
}

type City struct {
	Name string
}

type Stringable interface {
	ToString() string
}

func (c Country) ToString() string {
	return "Country = " + c.Name
}

func (c City) ToString() string {
	return "City = " + c.Name
}

func PrintStr(p Stringable) {
	fmt.Println(p.ToString())
}

func ProgramToInterface()  {
	d1 := Country{"USA"}
	d2 := City{"Los Angeles"}

	// var _ Stringable = (*Country)nil // when annotate country's ToString, we can use this way to identify if it's a interface
	PrintStr(d1)
	PrintStr(d2)
}
