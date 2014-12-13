// corporation
package models

type Corporation struct {
	Id             int64
	CorporationId  int64
	Name           string
	Ticker         string
	CorporationCut float64
}

func NewCorporation(id int64, corpId int64, name string, ticker string, cut float64) *Corporation {
	corp := &Corporation{
		Id:             id,
		CorporationId:  corpId,
		Name:           name,
		Ticker:         ticker,
		CorporationCut: cut,
	}

	return corp
}
