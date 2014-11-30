// fleetrole
package models

type FleetRole int

const (
	FleetRoleUnknown FleetRole = 1 << iota
	FleetRoleNone
	FleetRoleScout
	FleetRoleSalvage
	FleetRoleLogistics
	FleetRoleDPS
	FleetRoleFleetCommander
)

func (role FleetRole) String() string {
	switch role {
	case FleetRoleUnknown:
		return "unknown"
	case FleetRoleNone:
		return "none"
	case FleetRoleScout:
		return "scout"
	case FleetRoleSalvage:
		return "salvage"
	case FleetRoleLogistics:
		return "logistics"
	case FleetRoleDPS:
		return "dps"
	case FleetRoleFleetCommander:
		return "fleetcommander"
	default:
		return "invalid"
	}
}
