// fleetmemberpayout
package models

type FleetMemberPayout struct {
	FleetID        int64
	PlayerID       int64
	Payout         float64
	PayoutComplete bool
}

func NewFleetMemberPayout(fleet int64, player int64, payout float64, complete bool) *FleetMemberPayout {
	pay := &FleetMemberPayout{
		FleetID:        fleet,
		PlayerID:       player,
		Payout:         payout,
		PayoutComplete: complete,
	}

	return pay
}
