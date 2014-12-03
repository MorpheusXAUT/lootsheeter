// fleet
package models

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Fleet struct {
	Id             int64
	Name           string
	Members        []*FleetMember
	System         string
	SystemNickname string
	StartTime      time.Time
	EndTime        time.Time
	Profit         float64
	Losses         float64
	SitesFinished  int
	PayoutComplete bool

	fleetMembersMutex sync.RWMutex
}

func NewFleet(id int64, name string, system string, systemNick string, profit float64, losses float64, sites int, start time.Time, end time.Time, complete bool) *Fleet {
	fleet := &Fleet{
		Id:             id,
		Name:           name,
		Members:        make([]*FleetMember, 0),
		System:         system,
		SystemNickname: systemNick,
		Profit:         profit,
		Losses:         losses,
		SitesFinished:  sites,
		StartTime:      start,
		EndTime:        end,
		PayoutComplete: complete,
	}

	return fleet
}

func (fleet *Fleet) IsFleetFinished() bool {
	return !fleet.EndTime.IsZero()
}

func (fleet *Fleet) FinishFleet() {
	fleet.EndTime = time.Now()
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

func (fleet *Fleet) TickSitesFinished() {
	fleet.SitesFinished += 1
}

func (fleet *Fleet) HasMember(player string) bool {
	fleet.fleetMembersMutex.RLock()
	defer fleet.fleetMembersMutex.RUnlock()

	for _, member := range fleet.Members {
		if strings.EqualFold(member.Name, player) {
			return true
		}
	}

	return false
}

func (fleet *Fleet) AddMember(member *FleetMember) error {
	if fleet.HasMember(member.Player.Name) {
		return fmt.Errorf("Member %q already exists in fleet, cannot add twice")
	}

	fleet.fleetMembersMutex.Lock()
	defer fleet.fleetMembersMutex.Unlock()

	fleet.Members = append(fleet.Members, member)

	return nil
}

func (fleet *Fleet) RemoveMember(player string) error {
	if !fleet.HasMember(player) {
		return fmt.Errorf("Member %q does not exists in fleet, cannot remove")
	}

	fleet.fleetMembersMutex.Lock()
	defer fleet.fleetMembersMutex.Unlock()

	var index int

	for idx, member := range fleet.Members {
		if strings.EqualFold(member.Name, player) {
			index = idx
			break
		}
	}

	fleet.Members[index], fleet.Members = fleet.Members[len(fleet.Members)-1], fleet.Members[:len(fleet.Members)-1]

	return nil
}

func (fleet *Fleet) TickMemberSiteModifier(player string) error {
	if !fleet.HasMember(player) {
		return fmt.Errorf("Member %q does not exists in fleet, cannot tick modifier")
	}

	fleet.fleetMembersMutex.RLock()
	defer fleet.fleetMembersMutex.RUnlock()

	for _, member := range fleet.Members {
		if strings.EqualFold(member.Name, player) {
			member.TickSiteModifier()
		}
	}

	return nil
}

func (fleet *Fleet) GetMemberSiteModifier(player string) (int, error) {
	if !fleet.HasMember(player) {
		return 0, fmt.Errorf("Member %q does not exists in fleet, cannot get modifier")
	}

	var modifier int

	for _, member := range fleet.Members {
		if strings.EqualFold(member.Name, player) {
			modifier = member.SiteModifier
		}
	}

	return modifier, nil
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

	var modifier float64

	for _, member := range fleet.Members {
		if strings.EqualFold(member.Name, player) {
			modifier = member.PaymentModifier
		}
	}

	return modifier, nil
}

func (fleet *Fleet) CalculatePayments() map[string]float64 {
	payments := make(map[string]float64)
	var totalPoints float64
	var corpPayment float64
	var payout float64

	totalPoints = 0
	corpPayment = (fleet.Profit - fleet.Losses) * 0.28
	payout = fleet.Profit - corpPayment - fleet.Losses

	payments["CORPORATION"] = corpPayment

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

		payments[member.Name] = isk
	}

	return payments
}
