// fleetmember
package models

type FleetMember struct {
	Id      int64
	FleetId int64
	Player
	Role            FleetRole
	SiteModifier    int
	PaymentModifier float64
}

func NewFleetMember(fleet int64, player Player, role FleetRole) FleetMember {
	member := FleetMember{
		Id:              -1,
		FleetId:         fleet,
		Player:          player,
		Role:            role,
		SiteModifier:    0,
		PaymentModifier: 0,
	}

	return member
}

func (member FleetMember) TickSiteModifier() {
	member.SiteModifier -= 1
}
