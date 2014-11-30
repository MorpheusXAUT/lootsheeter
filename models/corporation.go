// corporation
package models

type Corporation struct {
	Id     int64
	Name   string
	Ticker string
}

func NewCorporation(name string) Corporation {
	corp := Corporation{
		Id:     -1,
		Name:   name,
		Ticker: "",
	}

	return corp
}
