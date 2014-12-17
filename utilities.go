// utilities
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/morpheusxaut/lootsheeter/models"
)

func GetEvepraisalValue(raw string) (float64, error) {
	url := strings.ToLower(raw)

	if strings.HasPrefix(url, "http://evepraisal.com/e/") || strings.HasPrefix(url, "http://evepraisal.com/estimate/") {
		if !strings.HasSuffix(url, ".json") {
			url += ".json"
		}

		resp, err := http.Get(url)
		if err != nil {
			return 0, err
		}

		defer resp.Body.Close()

		jsonContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		var evePraisal models.EvePraisal
		err = json.Unmarshal(jsonContent, &evePraisal)
		if err != nil {
			return 0, err
		}

		return evePraisal.GetTotalBuyValue(), nil
	}

	return 0, fmt.Errorf("Invalid evepraisal link, cannot parse")
}

func GetZKillboardValue(raw string) (float64, error) {
	url := strings.TrimRight(strings.ToLower(raw), "/")

	if strings.HasPrefix(url, "https://zkillboard.com/kill/") {
		killID := url[strings.LastIndex(url, "/")+1 : len(url)]

		resp, err := http.Get(fmt.Sprintf("https://zkillboard.com/api/killID/%s", killID))
		if err != nil {
			return 0, err
		}

		defer resp.Body.Close()

		jsonContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		var zKillboard models.ZKillboard
		err = json.Unmarshal(jsonContent, &zKillboard)
		if err != nil {
			return 0, err
		}

		return zKillboard.GetTotalValue()
	}

	return 0, fmt.Errorf("Invalid zKillboard link, cannot parse")
}

func GetPasteValue(raw string) (float64, error) {
	data := url.Values{}
	data.Set("raw_paste", raw)

	req, err := http.NewRequest("POST", "http://evepraisal.com/estimate", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	reg := regexp.MustCompile("Result #([0-9]+)")
	resultID := reg.FindStringSubmatch(string(body))

	value, err := GetEvepraisalValue(fmt.Sprintf("http://evepraisal.com/e/%s.json", resultID[1]))
	if err != nil {
		return 0, err
	}

	return value, nil
}

func GenerateRandomString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

func FetchCharacterAffiliation(v models.SSOVerification) (models.CharacterAffiliation, error) {
	assocReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.eveonline.com/eve/CharacterAffiliation.xml.aspx?ids=%d", v.CharacterID), nil)

	client := &http.Client{}
	assocResp, err := client.Do(assocReq)
	if err != nil {
		return models.CharacterAffiliation{}, err
	}
	defer assocResp.Body.Close()

	assocBody, err := ioutil.ReadAll(assocResp.Body)
	if err != nil {
		return models.CharacterAffiliation{}, err
	}

	var a models.CharacterAffiliation

	err = xml.Unmarshal(assocBody, &a)
	if err != nil {
		return models.CharacterAffiliation{}, err
	}

	return a, nil
}

func FetchCorporationSheet(a models.CharacterAffiliation) (models.CorporationSheet, error) {
	sheetReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.eveonline.com/corp/CorporationSheet.xml.aspx?corporationID=%d", a.GetCorporationID()), nil)

	client := &http.Client{}
	sheetResp, err := client.Do(sheetReq)
	if err != nil {
		return models.CorporationSheet{}, err
	}
	defer sheetResp.Body.Close()

	sheetBody, err := ioutil.ReadAll(sheetResp.Body)
	if err != nil {
		return models.CorporationSheet{}, err
	}

	var s models.CorporationSheet

	err = xml.Unmarshal(sheetBody, &s)
	if err != nil {
		return models.CorporationSheet{}, err
	}

	return s, nil
}

func SendJSONResponse(w http.ResponseWriter, response map[string]interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("Failed to encode response to JSON: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))

	w.WriteHeader(http.StatusOK)

	w.Write(jsonResponse)
}

func ParseFleetCompositionRows(fleetID int64, rows []string) ([]*models.FleetMember, error) {
	var members []*models.FleetMember

	for _, row := range rows {
		splitRow := strings.Split(row, "\t")
		if len(splitRow) != 7 {
			return members, fmt.Errorf("Invalid fleet composition row: %q", row)
		}

		name := splitRow[0]
		ship := splitRow[2]
		fleetBoss := false

		if strings.Contains(splitRow[4], "(Boss)") {
			fleetBoss = true
		}

		player, err := database.LoadPlayerFromName(name)
		if err != nil {
			return members, err
		}

		role, err := ParseFleetRole(ship, fleetBoss)
		if err != nil {
			return members, err
		}

		member := models.NewFleetMember(-1, fleetID, player, role, ship, 0, 1, 0, false, -1)

		members = append(members, member)
	}

	return members, nil
}

func ParseFleetRole(ship string, fleetBoss bool) (models.FleetRole, error) {
	if fleetBoss {
		return models.FleetRoleFleetCommander, nil
	}

	role, err := database.QueryShipRole(ship)
	if err != nil {
		if strings.EqualFold(err.Error(), "sql: no rows in result set") {
			return models.FleetRoleNone, nil
		} else {
			return models.FleetRoleUnknown, err
		}
	}

	return role, nil
}
