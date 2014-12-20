// corporation
package models

type Corporation struct {
	ID             int64
	CorporationID  int64
	Name           string
	Ticker         string
	CorporationCut float64
	APIID          int64
	APICode        string
}

func NewCorporation(id int64, corpID int64, name string, ticker string, cut float64, apiID int64, code string) *Corporation {
	corp := &Corporation{
		ID:             id,
		CorporationID:  corpID,
		Name:           name,
		Ticker:         ticker,
		CorporationCut: cut,
		APIID:          apiID,
		APICode:        code,
	}

	return corp
}
