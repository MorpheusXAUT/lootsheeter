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
	CreatedBy      *Player
	Fleets         []*Fleet
}

func NewReport(id int64, payout float64, start time.Time, end time.Time, complete bool, created *Player, fleets []*Fleet) *Report {
	report := &Report{
		Id:             id,
		TotalPayout:    payout,
		StartRange:     start,
		EndRange:       end,
		PayoutComplete: complete,
		CreatedBy:      created,
		Fleets:         fleets,
	}

	return report
}
