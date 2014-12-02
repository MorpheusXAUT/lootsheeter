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
		return "Unknown"
	case FleetRoleNone:
		return "None"
	case FleetRoleScout:
		return "Scout"
	case FleetRoleSalvage:
		return "Salvage"
	case FleetRoleLogistics:
		return "Logistics"
	case FleetRoleDPS:
		return "DPS"
	case FleetRoleFleetCommander:
		return "Fleetcommander"
	default:
		return "Invalid"
	}
}

func (role FleetRole) LabelType() string {
	switch role {
	case FleetRoleUnknown:
		return ""
	case FleetRoleNone:
		return ""
	case FleetRoleScout:
		return "label-default"
	case FleetRoleSalvage:
		return "label-info"
	case FleetRoleLogistics:
		return "label-success"
	case FleetRoleDPS:
		return "label-primary"
	case FleetRoleFleetCommander:
		return "label-warning"
	default:
		return "label-danger"
	}
}
