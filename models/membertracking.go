// membertracking
package models

import (
	"encoding/xml"
)

type MemberTracking struct {
	XMLName     xml.Name            `xml:"eveapi"`
	Rows        []MemberTrackingRow `xml:"result>rowset>row"`
	CachedUntil string              `xml:"cachedUntil"`
}

type MemberTrackingRow struct {
	CharacterID   int64  `xml:"characterID,attr"`
	Name          string `xml:"name,attr"`
	StartDateTime string `xml:"startDateTime,attr"`
	BaseID        int64  `xml:"baseID,attr"`
	Base          string `xml:"base,attr"`
	Title         string `xml:"title,attr"`
}
