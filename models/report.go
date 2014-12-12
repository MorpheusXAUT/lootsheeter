// report
package models

import (
	"time"
)

type Report struct {
	Id             int64
	TotalPayout    float64
	StartRange     time.Time
	EndRange       time.Time
	PayoutComplete bool
	Creator        *Player
	Fleets         []*Fleet
	Payouts        map[string]*ReportPayout
}

func NewReport(id int64, payout float64, start time.Time, end time.Time, complete bool, creator *Player, fleets []*Fleet) *Report {
	report := &Report{
		Id:             id,
		TotalPayout:    payout,
		StartRange:     start,
		EndRange:       end,
		PayoutComplete: complete,
		Creator:        creator,
		Fleets:         fleets,
		Payouts:        make(map[string]*ReportPayout),
	}

	return report
}

func (report *Report) CalculatePayouts() {
	report.TotalPayout = 0

	for _, fleet := range report.Fleets {
		fleet.CalculatePayouts()

		for _, member := range fleet.Members {
			_, ok := report.Payouts[member.Name]
			if !ok {
				report.Payouts[member.Name] = NewReportPayout(member.Player, false)
			}

			report.Payouts[member.Name].AddPayout(NewFleetMemberPayout(fleet.Id, member.Player.Id, member.Payout, member.PayoutComplete))

			report.TotalPayout += member.Payout
		}
	}
}

func (report *Report) AllPayoutsComplete() bool {
	if report.PayoutComplete {
		return true
	}

	report.PayoutComplete = true

	for _, payout := range report.Payouts {
		if report.PayoutComplete && !payout.AllPayoutsComplete() {
			report.PayoutComplete = false
		}
	}

	return report.PayoutComplete
}
