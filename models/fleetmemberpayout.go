// fleetmemberpayout
package models

type FleetMemberPayout struct {
	FleetId        int64
	PlayerId       int64
	Payout         float64
	PayoutComplete bool
}

func NewFleetMemberPayout(fleet int64, player int64, payout float64, complete bool) *FleetMemberPayout {
	pay := &FleetMemberPayout{
		FleetId:        fleet,
		PlayerId:       player,
		Payout:         payout,
		PayoutComplete: complete,
	}

	return pay
}
