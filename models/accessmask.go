// accessmask
package models

import (
	"strings"
)

type AccessMask int

const (
	AccessMaskUnknown AccessMask = 1 << iota
	AccessMaskNone
	AccessMaskMember
	AccessMaskJuniorFleetCommander
	AccessMaskSeniorFleetCommander
	AccessMaskOfficer
	AccessMaskPayoutOfficer
	AccessMaskDirector
	AccessMaskCEO
	AccessMaskAdmin
)

func (mask AccessMask) String() string {
	str := ""

	if mask&AccessMaskNone == AccessMaskNone {
		str += "None|"
	}
	if mask&AccessMaskMember == AccessMaskMember {
		str += "Member|"
	}
	if mask&AccessMaskJuniorFleetCommander == AccessMaskJuniorFleetCommander {
		str += "Junior Fleetcommander|"
	}
	if mask&AccessMaskSeniorFleetCommander == AccessMaskSeniorFleetCommander {
		str += "Senior Fleetcommander|"
	}
	if mask&AccessMaskOfficer == AccessMaskOfficer {
		str += "Officer|"
	}
	if mask&AccessMaskPayoutOfficer == AccessMaskPayoutOfficer {
		str += "Payout Officer|"
	}
	if mask&AccessMaskDirector == AccessMaskDirector {
		str += "Director|"
	}
	if mask&AccessMaskCEO == AccessMaskCEO {
		str += "CEO|"
	}
	if mask&AccessMaskAdmin == AccessMaskAdmin {
		str += "Admin|"
	}

	str = strings.TrimRight(str, "|")

	return str
}
