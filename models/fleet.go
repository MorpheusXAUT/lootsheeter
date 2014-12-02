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

	fleetMembersMutex sync.RWMutex
}

func NewFleet(id int64, name string, system string, systemNick string, profit float64, losses float64, sites int, start time.Time, end time.Time) *Fleet {
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

func (fleet *Fleet) HasProfit() bool {
	return fleet.Profit > 0
}

func (fleet *Fleet) HasLosses() bool {
	return fleet.Losses > 0
}

func (fleet *Fleet) GetProfitString() string {
	return FormatFloat(fleet.Profit)
}

func (fleet *Fleet) GetLossesString() string {
	return FormatFloat(fleet.Losses)
}

func (fleet *Fleet) HasPositiveSurplus() bool {
	return (fleet.Profit - fleet.Losses) > 0
}

func (fleet *Fleet) GetSurplus() float64 {
	return fleet.Profit - fleet.Losses
}

func (fleet *Fleet) GetSurplusString() string {
	return FormatFloat(fleet.Profit - fleet.Losses)
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
