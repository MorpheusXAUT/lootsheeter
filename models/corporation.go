// corporation
package models

type Corporation struct {
	Id     int64
	CorpId int64
	Name   string
	Ticker string
}

func NewCorporation(id int64, corpId int64, name string, ticker string) *Corporation {
	corp := &Corporation{
		Id:     id,
		CorpId: corpId,
		Name:   name,
		Ticker: ticker,
	}

	return corp
}
