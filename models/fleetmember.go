// fleetmember
package models

import (
	"fmt"
	"strings"
)

type FleetMember struct {
	ID      int64
	FleetID int64
	*Player
	Role            FleetRole
	Ship            string
	SiteModifier    int
	PaymentModifier float64
	Payout          float64
	PayoutComplete  bool
	ReportID        int64
}

func NewFleetMember(id int64, fleetID int64, player *Player, role FleetRole, ship string, site int, payment float64, payout float64, complete bool, report int64) *FleetMember {
	member := &FleetMember{
		ID:              id,
		FleetID:         fleetID,
		Player:          player,
		Role:            role,
		Ship:            ship,
		SiteModifier:    site,
		PaymentModifier: payment,
		Payout:          payout,
		PayoutComplete:  complete,
		ReportID:        report,
	}

	return member
}

func (member *FleetMember) HasRole(role string) bool {
	return strings.EqualFold(role, fmt.Sprintf("%s", member.Role))
}
