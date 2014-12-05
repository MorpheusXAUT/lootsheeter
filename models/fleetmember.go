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
	Payout          float64
	PayoutComplete  bool
	ReportId        int64
}

func NewFleetMember(id int64, fleetId int64, player *Player, role FleetRole, site int, payment float64, payout float64, complete bool, report int64) *FleetMember {
	member := &FleetMember{
		Id:              id,
		FleetId:         fleetId,
		Player:          player,
		Role:            role,
		SiteModifier:    site,
		PaymentModifier: payment,
		Payout:          payout,
		PayoutComplete:  complete,
		ReportId:        report,
	}

	return member
}

func (member *FleetMember) HasRole(role string) bool {
	return strings.EqualFold(role, fmt.Sprintf("%s", member.Role))
}
