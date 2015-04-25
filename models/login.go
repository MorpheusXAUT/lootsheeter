// login
package models

import (
	"encoding/xml"
)

type SSOToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Expiry       int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SSOVerification struct {
	CharacterID        int64     `json:"CharacterID"`
	CharacterName      string    `json:"CharacterName"`
	ExpiresOn          string `json:"ExpiresOn"`
	Scopes             string    `json:"Scopes"`
	TokenType          string    `json:"TokenType"`
	CharacterOwnerHash string    `json:"CharacterOwnerHash"`
}

type CharacterAffiliation struct {
	XMLName xml.Name                  `xml:"eveapi"`
	Rows    []CharacterAffiliationRow `xml:"result>rowset>row"`
}

type CharacterAffiliationRow struct {
	CharacterID     int64  `xml:"characterID,attr"`
	CharacterName   string `xml:"characterName,attr"`
	CorporationID   int64  `xml:"corporationID,attr"`
	CorporationName string `xml:"corporationName,attr"`
	AllianceID      int64  `xml:"allianceID,attr"`
	AllianceName    string `xml:"allianceName,attr"`
	FactionID       int64  `xml:"factionID,attr"`
	FactionName     string `xml:"factionName,attr"`
}

type CorporationSheet struct {
	XMLName         xml.Name `xml:"eveapi"`
	CorporationID   int64    `xml:"result>corporationID"`
	CorporationName string   `xml:"result>corporationName"`
	Ticker          string   `xml:"result>ticker"`
	CEOId           int64    `xml:"result>ceoID"`
	CEOName         string   `xml:"result>ceoName"`
	StationID       int64    `xml:"result>stationID"`
	StationName     string   `xml:"result>stationName"`
	Description     string   `xml:"result>description"`
	Homepage        string   `xml:"result>url"`
	AllianceID      int64    `xml:"result>allianceID"`
	AllianceName    string   `xml:"result>allianceName"`
	FactionID       int64    `xml:"result>factionID"`
	TaxRate         float64  `xml:"result>taxRate"`
	MemberCount     int64    `xml:"result>memberCount"`
	Shares          int64    `xml:"result>shares"`
}

func (a CharacterAffiliation) GetCharacterID() int64 {
	if len(a.Rows) == 0 {
		return -1
	}

	return a.Rows[0].CharacterID
}

func (a CharacterAffiliation) GetCharacterName() string {
	if len(a.Rows) == 0 {
		return ""
	}

	return a.Rows[0].CharacterName
}

func (a CharacterAffiliation) GetCorporationID() int64 {
	if len(a.Rows) == 0 {
		return -1
	}

	return a.Rows[0].CorporationID
}

func (a CharacterAffiliation) GetCorporationName() string {
	if len(a.Rows) == 0 {
		return ""
	}

	return a.Rows[0].CorporationName
}

func (a CharacterAffiliation) GetAllianceID() int64 {
	if len(a.Rows) == 0 {
		return -1
	}

	return a.Rows[0].AllianceID
}

func (a CharacterAffiliation) GetAllianceName() string {
	if len(a.Rows) == 0 {
		return ""
	}

	return a.Rows[0].AllianceName
}

func (a CharacterAffiliation) GetFactionID() int64 {
	if len(a.Rows) == 0 {
		return -1
	}

	return a.Rows[0].FactionID
}

func (a CharacterAffiliation) GetFactionName() string {
	if len(a.Rows) == 0 {
		return ""
	}

	return a.Rows[0].FactionName
}
