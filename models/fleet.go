// fleet
package models

import (
	"fmt"
	"time"
)

type Fleet struct {
	ID                int64
	Corporation       *Corporation
	Name              string
	Members           map[string]*FleetMember
	System            string
	SystemNickname    string
	StartTime         time.Time
	EndTime           time.Time
	Profit            float64
	Losses            float64
	SitesFinished     int
	CorporationPayout float64
	PayoutComplete    bool
	Notes             string
	ReportID          int64
}

func NewFleet(id int64, corp *Corporation, name string, system string, systemNick string, profit float64, losses float64, sites int, start time.Time, end time.Time, payout float64, complete bool, notes string, report int64) *Fleet {
	fleet := &Fleet{
		ID:                id,
		Corporation:       corp,
		Name:              name,
		Members:           make(map[string]*FleetMember),
		System:            system,
		SystemNickname:    systemNick,
		Profit:            profit,
		Losses:            losses,
		SitesFinished:     sites,
		StartTime:         start,
		EndTime:           end,
		CorporationPayout: payout,
		PayoutComplete:    complete,
		Notes:             notes,
		ReportID:          report,
	}

	return fleet
}

func (fleet *Fleet) IsFleetFinished() bool {
	return !fleet.EndTime.IsZero()
}

func (fleet *Fleet) FinishFleet() {
	fleet.EndTime = time.Now()
	fleet.CalculatePayouts()
}

func (fleet *Fleet) AddProfit(profit float64) {
	fleet.Profit += profit
}

func (fleet *Fleet) AddLoss(loss float64) {
	fleet.Losses += loss
}

func (fleet *Fleet) GetSurplus() float64 {
	return fleet.Profit - fleet.Losses
}

func (fleet *Fleet) HasMember(player string) bool {
	_, ok := fleet.Members[player]

	return ok
}

func (fleet *Fleet) FleetCommanders() []*FleetMember {
	var fleetCommanders []*FleetMember

	for _, member := range fleet.Members {
		if member.Role == FleetRoleFleetCommander {
			fleetCommanders = append(fleetCommanders, member)
		}
	}

	return fleetCommanders
}

func (fleet *Fleet) AddMember(member *FleetMember) error {
	if fleet.HasMember(member.Player.Name) {
		return fmt.Errorf("Member %q already exists in fleet, cannot add twice", member.Player.Name)
	}

	fleet.Members[member.Player.Name] = member

	return nil
}

func (fleet *Fleet) RemoveMember(player string) error {
	if !fleet.HasMember(player) {
		return fmt.Errorf("Member %q does not exists in fleet, cannot remove", player)
	}

	delete(fleet.Members, player)

	return nil
}

func (fleet *Fleet) UpdateMember(member *FleetMember) {
	fleet.Members[member.Player.Name] = member
}

func (fleet *Fleet) GetMemberRole(player string) (FleetRole, error) {
	if !fleet.HasMember(player) {
		return FleetRoleUnknown, fmt.Errorf("Member %q does not exists in fleet, cannot get role", player)
	}

	return fleet.Members[player].Role, nil
}

func (fleet *Fleet) SetMemberSiteModifier(player string, modifier int) error {
	if !fleet.HasMember(player) {
		return fmt.Errorf("Member %q does not exists in fleet, cannot set site modifier", player)
	}

	fleet.Members[player].SiteModifier = modifier

	return nil
}

func (fleet *Fleet) GetMemberSiteModifier(player string) (int, error) {
	if !fleet.HasMember(player) {
		return 0, fmt.Errorf("Member %q does not exists in fleet, cannot get site modifier", player)
	}

	return fleet.Members[player].SiteModifier, nil
}

func (fleet *Fleet) GetMemberSitesFinished(player string) (int, error) {
	if !fleet.HasMember(player) {
		return 0, fmt.Errorf("Member %q does not exists in fleet, cannot get sites finished", player)
	}

	modifier, err := fleet.GetMemberSiteModifier(player)
	if err != nil {
		return 0, err
	}

	return (fleet.SitesFinished - modifier), nil
}

func (fleet *Fleet) GetMemberPaymentModifier(player string) (float64, error) {
	if !fleet.HasMember(player) {
		return 0, fmt.Errorf("Member %q does not exists in fleet, cannot get payment modifier", player)
	}

	return fleet.Members[player].PaymentModifier, nil
}

func (fleet *Fleet) CalculatePayouts() {
	var totalPoints float64
	var corpPayment float64
	var payout float64

	totalPoints = 0
	corpPayment = 0

	if fleet.Corporation.CorporationCut > 0 {
		corpPayment = (fleet.Profit - fleet.Losses) * (fleet.Corporation.CorporationCut / 100)
	}

	payout = fleet.Profit - corpPayment - fleet.Losses

	fleet.CorporationPayout = corpPayment

	for _, member := range fleet.Members {
		var points float64
		if member.PaymentModifier != 1 {
			points = float64((fleet.SitesFinished + member.SiteModifier)) * member.PaymentModifier
		} else {
			points = float64((fleet.SitesFinished + member.SiteModifier)) * member.Role.PaymentRate()
		}

		totalPoints += points
	}

	for _, member := range fleet.Members {
		var points float64
		var isk float64

		if member.PaymentModifier != 1 {
			points = float64((fleet.SitesFinished + member.SiteModifier)) * member.PaymentModifier
		} else {
			points = float64((fleet.SitesFinished + member.SiteModifier)) * member.Role.PaymentRate()
		}

		if totalPoints > 0 {
			isk = payout * (points / totalPoints)
		} else {
			isk = 0
		}

		member.Payout = isk
	}
}
