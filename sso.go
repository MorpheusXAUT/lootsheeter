// sso
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/morpheusxaut/lootsheeter/models"
)

func FetchSSOToken(authorizationCode string) (models.SSOToken, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.SSOClientID, config.SSOClientSecret)))

	verifyData := url.Values{}
	verifyData.Set("grant_type", "authorization_code")
	verifyData.Set("code", authorizationCode)

	verifyReq, err := http.NewRequest("POST", "https://login.eveonline.com/oauth/token", bytes.NewBufferString(verifyData.Encode()))
	verifyReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	verifyReq.Header.Add("Content-Length", strconv.Itoa(len(verifyData.Encode())))
	verifyReq.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))

	client := &http.Client{}
	verifyResp, err := client.Do(verifyReq)
	if err != nil {
		return models.SSOToken{}, err
	}
	defer verifyResp.Body.Close()

	verifyBody, err := ioutil.ReadAll(verifyResp.Body)
	if err != nil {
		return models.SSOToken{}, err
	}

	var t models.SSOToken

	err = json.Unmarshal(verifyBody, &t)
	if err != nil {
		return models.SSOToken{}, err
	}

	return t, nil
}

func FetchSSOVerification(t models.SSOToken) (models.SSOVerification, error) {
	charReq, err := http.NewRequest("GET", "https://login.eveonline.com/oauth/verify", nil)
	charReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))

	client := &http.Client{}
	charResp, err := client.Do(charReq)
	if err != nil {
		return models.SSOVerification{}, err
	}
	defer charResp.Body.Close()

	charBody, err := ioutil.ReadAll(charResp.Body)
	if err != nil {
		return models.SSOVerification{}, err
	}

	var v models.SSOVerification

	err = json.Unmarshal(charBody, &v)
	if err != nil {
		return models.SSOVerification{}, err
	}

	return v, nil
}
