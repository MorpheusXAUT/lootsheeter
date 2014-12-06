// fleet
package models

import (
	"fmt"
	"time"
)

type Fleet struct {
	Id                int64
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
	ReportId          int64
}

func NewFleet(id int64, name string, system string, systemNick string, profit float64, losses float64, sites int, start time.Time, end time.Time, payout float64, complete bool, report int64) *Fleet {
	fleet := &Fleet{
		Id:                id,
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
		ReportId:          report,
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

func (fleet *Fleet) FleetCommander() *FleetMember {
	for _, member := range fleet.Members {
		if member.Role == FleetRoleFleetCommander {
			return member
		}
	}

	return nil
}

func (fleet *Fleet) AddMember(member *FleetMember) error {
	if fleet.HasMember(member.Player.Name) {
		return fmt.Errorf("Member %q already exists in fleet, cannot add twice")
	}

	fleet.Members[member.Name] = member

	return nil
}

func (fleet *Fleet) RemoveMember(player string) error {
	if !fleet.HasMember(player) {
		return fmt.Errorf("Member %q does not exists in fleet, cannot remove")
	}

	delete(fleet.Members, player)

	return nil
}

func (fleet *Fleet) SetMemberSiteModifier(player string, modifier int) error {
	if !fleet.HasMember(player) {
		return fmt.Errorf("Member %q does not exists in fleet, cannot set site modifier")
	}

	fleet.Members[player].SiteModifier = modifier

	return nil
}

func (fleet *Fleet) GetMemberSiteModifier(player string) (int, error) {
	if !fleet.HasMember(player) {
		return 0, fmt.Errorf("Member %q does not exists in fleet, cannot get site modifier")
	}

	return fleet.Members[player].SiteModifier, nil
}

func (fleet *Fleet) GetMemberSitesFinished(player string) (int, error) {
	if !fleet.HasMember(player) {
		return 0, fmt.Errorf("Member %q does not exists in fleet, cannot get sites finished")
	}

	modifier, err := fleet.GetMemberSiteModifier(player)
	if err != nil {
		return 0, err
	}

	return (fleet.SitesFinished - modifier), nil
}

func (fleet *Fleet) GetMemberPaymentModifier(player string) (float64, error) {
	if !fleet.HasMember(player) {
		return 0, fmt.Errorf("Member %q does not exists in fleet, cannot get payment modifier")
	}

	return fleet.Members[player].PaymentModifier, nil
}

func (fleet *Fleet) CalculatePayouts() {
	var totalPoints float64
	var corpPayment float64
	var payout float64

	totalPoints = 0
	corpPayment = (fleet.Profit - fleet.Losses) * 0.28
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
