// corporation
package models

type Corporation struct {
	ID             int64
	CorporationID  int64
	Name           string
	Ticker         string
	CorporationCut float64
}

func NewCorporation(id int64, corpID int64, name string, ticker string, cut float64) *Corporation {
	corp := &Corporation{
		ID:             id,
		CorporationID:  corpID,
		Name:           name,
		Ticker:         ticker,
		CorporationCut: cut,
	}

	return corp
}
