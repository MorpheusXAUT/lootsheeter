// templates
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/morpheusxaut/lootsheeter/models"
)

var (
	templates = template.Must(template.New("").Funcs(TemplateFunctions(nil)).ParseGlob("web/template/*"))
)

func TemplateFunctions(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"FormatFloat":                 func(f float64) string { return FormatFloat(f) },
		"IsPositiveFloat":             func(f float64) bool { return IsPositiveFloat(f) },
		"FloatEquals":                 func(f1 float64, f2 float64) bool { return FloatEquals(f1, f2) },
		"HasFleetRole":                func(fleet *models.Fleet, roleInt int) bool { return HasFleetRole(r, fleet, roleInt) },
		"MemberHasFleetRole":          func(member *models.FleetMember, roleInt int) bool { return MemberHasFleetRole(member, roleInt) },
		"IsFleetCommander":            func(fleet *models.Fleet) bool { return IsFleetCommander(r, fleet) },
		"IsReportCreator":             func(report *models.Report) bool { return IsReportCreator(r, report) },
		"IsPlayerName":                func(name string) bool { return IsPlayerName(r, name) },
		"HasAccessMask":               func(access int) bool { return HasAccessMask(r, access) },
		"GetFleetRolePaymentModifier": func(role *models.FleetRole) float64 { return GetFleetRolePaymentModifier(role) },
	}
}

func FormatFloat(f float64) string {
	fString := humanize.Ftoa(f)

	var formattedFloat string

	if strings.Contains(fString, ".") {
		fInt, err := strconv.ParseInt(fString[:strings.Index(fString, ".")], 10, 64)
		if err != nil {
			return fString
		}

		digitsAfterPoint := len(fString) - strings.Index(fString, ".")
		if digitsAfterPoint > 3 {
			digitsAfterPoint = 3
		}

		formattedFloat = fmt.Sprintf("%s%s", humanize.Comma(fInt), fString[strings.Index(fString, "."):strings.Index(fString, ".")+digitsAfterPoint])
	} else {
		fInt, err := strconv.ParseInt(fString, 10, 64)
		if err != nil {
			return fString
		}

		formattedFloat = fmt.Sprintf("%s.00", humanize.Comma(fInt))
	}

	return formattedFloat
}

func IsPositiveFloat(f float64) bool {
	return f > 0
}

func FloatEquals(f1 float64, f2 float64) bool {
	return f1 == f2
}

func HasFleetRole(r *http.Request, fleet *models.Fleet, roleInt int) bool {
	player := session.GetPlayerFromRequest(r)
	if player == nil {
		return false
	}

	role, err := fleet.GetMemberRole(player.Name)
	if err != nil {
		return false
	}

	return role == models.FleetRole(roleInt)
}

func MemberHasFleetRole(member *models.FleetMember, roleInt int) bool {
	return member.Role == models.FleetRole(roleInt)
}

func IsFleetCommander(r *http.Request, fleet *models.Fleet) bool {
	player := session.GetPlayerFromRequest(r)
	if player == nil {
		return false
	}

	if strings.EqualFold(fleet.FleetCommander().Name, player.Name) {
		return true
	}

	return false
}

func IsReportCreator(r *http.Request, report *models.Report) bool {
	player := session.GetPlayerFromRequest(r)
	if player == nil {
		return false
	}

	if strings.EqualFold(report.Creator.Name, player.Name) {
		return true
	}

	return false
}

func IsPlayerName(r *http.Request, name string) bool {
	player := session.GetPlayerFromRequest(r)
	if player == nil {
		return false
	}

	if strings.EqualFold(name, player.Name) {
		return true
	}

	return false
}

func HasAccessMask(r *http.Request, access int) bool {
	player := session.GetPlayerFromRequest(r)
	if player == nil {
		return false
	}

	return player.AccessMask == models.AccessMask(access)
}

func HasHigherAccessMask(r *http.Request, access int) bool {
	player := session.GetPlayerFromRequest(r)
	if player == nil {
		return false
	}

	return player.AccessMask >= models.AccessMask(access)
}

func GetFleetRolePaymentModifier(role *models.FleetRole) float64 {
	return role.PaymentRate()
}
