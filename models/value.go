// value
package models

import (
	"strconv"
	"time"
)

type EvePraisal struct {
	Created    int64            `json:"created"`
	ID         int64            `json:"id"`
	Items      []EvePraisalItem `json:"items"`
	Kind       string           `json:"kind"`
	MarketID   int64            `json:"market_id"`
	MarketName string           `json:"market_name"`
	Totals     EvePraisalTotals `json:"totals"`
}

type EvePraisalItem struct {
	GroupID  int64            `json:"groupID"`
	Market   bool             `json:"market"`
	Name     string           `json:"name"`
	Prices   EvePraisalPrices `json:"prices"`
	Quantity int64            `json:"quantity"`
	TypeID   int64            `json:"typeID"`
	TypeName string           `json:"typeName"`
	Volume   float64          `json:"volume"`
}

type EvePraisalPrices struct {
	All  EvePraisalPrice `json:"all"`
	Buy  EvePraisalPrice `json:"buy"`
	Sell EvePraisalPrice `json:"sell"`
}

type EvePraisalPrice struct {
	Average float64 `json:"avg"`
	Maximum float64 `json:"max"`
	Minimum float64 `json:"min"`
	Price   float64 `json:"price"`
}

type EvePraisalTotals struct {
	Buy    float64 `json:"buy"`
	Sell   float64 `json:"sell"`
	Volume float64 `json:"volume"`
}

func (e EvePraisal) GetTotalBuyValue() float64 {
	return e.Totals.Buy
}

func (e EvePraisal) GetTotalSellValue() float64 {
	return e.Totals.Sell
}

func (e EvePraisal) GetTotalVolume() float64 {
	return e.Totals.Volume
}

type ZKillboard struct {
	KillID        string               `json:"killID"`
	SolarSystemID string               `json:"solarSystemID"`
	KillTime      time.Time            `json:"killTime"`
	MoonID        string               `json:"moonID"`
	Victim        ZKillboardVictim     `json:"victim"`
	Attackers     []ZKillboardAttacker `json:"attackers"`
	Items         []ZKillboardItem     `json:"items"`
	Info          ZKillboardInfo       `json:"zkb"`
}

type ZKillboardVictim struct {
	ShipTypeID      string `json:"shipTypeID"`
	DamageTaken     string `json:"damageTaken"`
	FactionName     string `json:"factionName"`
	FactionID       string `json:"factionID"`
	AllianceName    string `json:"allianceName"`
	AllianceID      string `json:"allianceID"`
	CorporationName string `json:"corporationName"`
	CorporationID   string `json:"corporationID"`
	CharacterName   string `json:"characterName"`
	CharacterID     string `json:"characterID"`
	Victim          string `json:"victim"`
}

type ZKillboardAttacker struct {
	CharacterID     string `json:"characterID"`
	CharacterName   string `json:"characterName"`
	CorporationID   string `json:"corporationID"`
	CorporationName string `json:"corporationName"`
	AllianceID      string `json:"allianceID"`
	AllianceName    string `json:"allianceName"`
	FactionID       string `json:"factionID"`
	FactionName     string `json:"factionName"`
	SecurityStatus  string `json:"securityStatus"`
	DamageDone      string `json:"damageDone"`
	FinalBlow       string `json:"finalBlow"`
	WeaponTypeID    string `json:"weaponTypeId"`
	ShipTypeID      string `json:"shipTypeID"`
}

type ZKillboardItem struct {
	TypeID            string `json:"typeID"`
	Flag              string `json:"flag"`
	QuantityDropped   string `json:"qtyDropped"`
	QuantityDestroyed string `json:"qtyDestroyed"`
	Singleton         string `json:"singleton"`
}

type ZKillboardInfo struct {
	TotalValue string `json:"totalValue"`
	Points     string `json:"points"`
	Source     string `json:"source"`
}

func (z ZKillboard) GetTotalValue() (float64, error) {
	return strconv.ParseFloat(z.Info.TotalValue, 64)
}
