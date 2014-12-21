// reportpayout
package models

type ReportPayout struct {
	Player         *Player
	Payouts        []*FleetMemberPayout
	TotalPayout    float64
	PayoutComplete bool
}

func NewReportPayout(player *Player, complete bool) *ReportPayout {
	payout := &ReportPayout{
		Player:         player,
		Payouts:        make([]*FleetMemberPayout, 0),
		TotalPayout:    0,
		PayoutComplete: complete,
	}

	return payout
}

func (pay *ReportPayout) AddPayout(payout *FleetMemberPayout) {
	pay.Payouts = append(pay.Payouts, payout)
	pay.TotalPayout += payout.Payout
	if pay.PayoutComplete && !payout.PayoutComplete {
		pay.PayoutComplete = false
	}
}

func (pay *ReportPayout) AllPayoutsComplete() bool {
	if pay.PayoutComplete {
		return true
	}

	pay.PayoutComplete = true

	for _, payout := range pay.Payouts {
		if !payout.PayoutComplete {
			pay.PayoutComplete = false
		}
	}

	return pay.PayoutComplete
}
