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
	AccessMaskDirector
	AccessMaskCEO
	AccessMaskAdmin
)

func (mask AccessMask) String() string {
	str := "unknown|"

	if mask&AccessMaskNone == AccessMaskNone {
		str += "none|"
	}
	if mask&AccessMaskMember == AccessMaskMember {
		str += "member|"
	}
	if mask&AccessMaskJuniorFleetCommander == AccessMaskJuniorFleetCommander {
		str += "juniorfleetcommander|"
	}
	if mask&AccessMaskSeniorFleetCommander == AccessMaskSeniorFleetCommander {
		str += "seniorfleetcommander|"
	}
	if mask&AccessMaskOfficer == AccessMaskOfficer {
		str += "officer|"
	}
	if mask&AccessMaskDirector == AccessMaskDirector {
		str += "director|"
	}
	if mask&AccessMaskCEO == AccessMaskCEO {
		str += "ceo|"
	}
	if mask&AccessMaskAdmin == AccessMaskAdmin {
		str += "admin|"
	}

	str = strings.TrimRight(str, "|")

	return str
}
