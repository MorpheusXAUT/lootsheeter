// reportpayout
package models

type ReportPayout struct {
	ID             int64
	ReportID       int64
	Player         *Player
	Payout         float64
	PayoutComplete bool
}

func NewReportPayout(id int64, report int64, player *Player, total float64, complete bool) *ReportPayout {
	payout := &ReportPayout{
		ID:             id,
		ReportID:       report,
		Player:         player,
		Payout:         total,
		PayoutComplete: complete,
	}

	return payout
}
