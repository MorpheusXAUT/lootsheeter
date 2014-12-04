// utilities
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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

		var jsonInterface interface{}
		err = json.Unmarshal(jsonContent, &jsonInterface)
		if err != nil {
			return 0, err
		}

		jsonMap, ok := jsonInterface.(map[string]interface{})
		if ok {
			jsonTotals, ok := jsonMap["totals"].(map[string]interface{})
			if ok {
				jsonBuyValue, ok := jsonTotals["buy"].(float64)
				if ok {
					return jsonBuyValue, nil
				} else {
					return 0, fmt.Errorf("Failed to convert JSON into jsonBuyValue for evepraisal")
				}
			} else {
				return 0, fmt.Errorf("Failed to convert JSON into jsonTotals for evepraisal")
			}
		} else {
			return 0, fmt.Errorf("Failed to convert JSON into jsonMap for evepraisal")
		}
	} else {
		return 0, fmt.Errorf("Invalid evepraisal link, cannot parse")
	}

	return 0, fmt.Errorf("Failed to parse JSON response from evepraisal")
}

func GetzKillboardValue(raw string) (float64, error) {
	url := strings.TrimRight(strings.ToLower(raw), "/")

	if strings.HasPrefix(url, "https://zkillboard.com/kill/") {
		killId := url[strings.LastIndex(url, "/")+1 : len(url)]

		resp, err := http.Get(fmt.Sprintf("https://zkillboard.com/api/killID/%s", killId))
		if err != nil {
			return 0, err
		}

		defer resp.Body.Close()

		jsonContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		var jsonInterface interface{}
		err = json.Unmarshal(jsonContent, &jsonInterface)
		if err != nil {
			return 0, err
		}

		jsonArray, ok := jsonInterface.([]interface{})
		if ok {
			jsonMap, ok := jsonArray[0].(map[string]interface{})
			if ok {
				jsonZkb, ok := jsonMap["zkb"].(map[string]interface{})
				if ok {
					jsonTotalValue, ok := jsonZkb["totalValue"].(string)
					if ok {
						value, err := strconv.ParseFloat(jsonTotalValue, 64)
						if err != nil {
							return 0, err
						}

						return value, nil
					} else {
						return 0, fmt.Errorf("Failed to convert JSON into jsonTotalValue for zKillboard")
					}
				} else {
					return 0, fmt.Errorf("Failed to convert JSON into jsonZkb for zKillboard")
				}
			} else {
				return 0, fmt.Errorf("Failed to convert JSON into jsonMap for zKillboard")
			}
		} else {
			return 0, fmt.Errorf("Failed to convert JSON into jsonArray for zKillboard")
		}
	} else {
		return 0, fmt.Errorf("Invalid zKillboard link, cannot parse")
	}

	return 0, nil
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
	resultId := reg.FindStringSubmatch(string(body))

	value, err := GetEvepraisalValue(fmt.Sprintf("http://evepraisal.com/e/%s.json", resultId[1]))
	if err != nil {
		return 0, err
	}

	return value, nil
}
