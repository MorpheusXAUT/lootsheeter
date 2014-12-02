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

func (role FleetRole) PaymentRate() float64 {
	switch role {
	case FleetRoleUnknown:
		return 0
	case FleetRoleNone:
		return 0
	case FleetRoleScout:
		return 1
	case FleetRoleSalvage:
		return 1
	case FleetRoleLogistics:
		return 1
	case FleetRoleDPS:
		return 1
	case FleetRoleFleetCommander:
		return 1.25
	default:
		return 0
	}
}
