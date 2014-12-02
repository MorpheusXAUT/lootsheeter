// fleetmember
package models

import (
	"fmt"
	"strings"
)

type FleetMember struct {
	Id      int64
	FleetId int64
	*Player
	Role            FleetRole
	SiteModifier    int
	PaymentModifier float64
}

func NewFleetMember(id int64, fleetId int64, player *Player, role FleetRole, site int, payment float64) *FleetMember {
	member := &FleetMember{
		Id:              id,
		FleetId:         fleetId,
		Player:          player,
		Role:            role,
		SiteModifier:    site,
		PaymentModifier: payment,
	}

	return member
}

func (member *FleetMember) HasRole(role string) bool {
	return strings.EqualFold(role, fmt.Sprintf("%s", member.Role))
}

func (member *FleetMember) TickSiteModifier() {
	member.SiteModifier -= 1
}
