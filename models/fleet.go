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
	Commander      FleetMember
	Members        []FleetMember
	System         string
	SystemNickname string
	StartTime      time.Time
	EndTime        time.Time
	Profit         float64
	Losses         float64
	SitesFinished  int

	fleetMembersMutex sync.RWMutex
}

func NewFleet(id int64, name string, system string, systemNick string) Fleet {
	fleet := Fleet{
		Id:             id,
		Name:           name,
		Members:        make([]FleetMember, 0),
		System:         system,
		SystemNickname: systemNick,
		StartTime:      time.Now(),
		Profit:         0.0,
		Losses:         0.0,
		SitesFinished:  0,
	}

	return fleet
}

func (fleet Fleet) FinishFleet() {
	fleet.EndTime = time.Now()
}

func (fleet Fleet) AddProfit(profit float64) {
	fleet.Profit += profit
}

func (fleet Fleet) AddLoss(loss float64) {
	fleet.Losses += loss
}

func (fleet Fleet) TickSitesFinished() {
	fleet.SitesFinished += 1
}

func (fleet Fleet) AddCommander(player Player) {
	fleet.Commander = NewFleetMember(-1, fleet.Id, player, FleetRoleFleetCommander)
}

func (fleet Fleet) HasMember(player string) bool {
	fleet.fleetMembersMutex.RLock()
	defer fleet.fleetMembersMutex.RUnlock()

	for _, member := range fleet.Members {
		if strings.EqualFold(member.Name, player) {
			return true
		}
	}

	return false
}

func (fleet Fleet) AddMember(player Player, role FleetRole) error {
	if fleet.HasMember(player.Name) {
		return fmt.Errorf("Member %q already exists in fleet, cannot add twice")
	}

	fleet.fleetMembersMutex.Lock()
	defer fleet.fleetMembersMutex.Unlock()

	fleet.Members = append(fleet.Members, NewFleetMember(-1, fleet.Id, player, role))

	return nil
}

func (fleet Fleet) RemoveMember(player string) error {
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

func (fleet Fleet) TickMemberSiteModifier(player string) error {
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

func (fleet Fleet) GetMemberSiteModifier(player string) (int, error) {
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

func (fleet Fleet) GetMemberSitesFinished(player string) (int, error) {
	if !fleet.HasMember(player) {
		return 0, fmt.Errorf("Member %q does not exists in fleet, cannot get sites finished")
	}

	modifier, err := fleet.GetMemberSiteModifier(player)
	if err != nil {
		return 0, err
	}

	return (fleet.SitesFinished - modifier), nil
}
