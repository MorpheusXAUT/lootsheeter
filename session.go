// session
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/morpheusxaut/lootsheeter/models"
)

var (
	session *Session
)

type Session struct {
	store *sessions.FilesystemStore
}

func NewSession() *Session {
	session := &Session{
		store: sessions.NewFilesystemStore("web/sessions", []byte(securecookie.GenerateRandomKey(128))),
	}

	return session
}

func InitialiseSessions() {
	CleanSessions()

	session = NewSession()
}

func CleanSessions() {
	sessions, err := filepath.Glob("web/sessions/session_*")
	if err != nil {
		logger.Errorf("Failed to find old sessions: [%v]", err)
		return
	}

	for _, s := range sessions {
		err = os.Remove(s)
		if err != nil {
			logger.Errorf("Failed to delete session: [%v]", err)
		}
	}
}

func (s *Session) DestroySession(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "player")

	playerIDInterface, ok := session.Values["playerID"]
	if ok {
		playerID := playerIDInterface.(int64)
		database.RemovePlayerFromCache(playerID)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "player",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "login",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func (s *Session) GetPlayerFromRequest(r *http.Request) *models.Player {
	session, _ := s.store.Get(r, "player")

	playerIDInterface, ok := session.Values["playerID"]
	if !ok {
		return nil
	}

	playerID := playerIDInterface.(int64)

	player, err := database.LoadPlayer(playerID)
	if err != nil {
		logger.Errorf("Failed to load player from database in session: [%v]", err)
		return nil
	}

	return player
}

func (s *Session) IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	session, _ := s.store.Get(r, "player")
	if session.IsNew {
		return false
	}

	player := s.GetPlayerFromRequest(r)
	if player == nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "player",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		http.SetCookie(w, &http.Cookie{
			Name:   "login",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		return false
	}

	characterName, ok := session.Values["characterName"]
	if !ok {
		return false
	}

	if strings.EqualFold(characterName.(string), player.Name) {
		return true
	}

	return false
}

func (s *Session) SetIdentity(w http.ResponseWriter, r *http.Request, a models.CharacterAffiliation, sh models.CorporationSheet) error {
	session, _ := s.store.Get(r, "player")

	session.Values["characterID"] = a.GetCharacterID()
	session.Values["characterName"] = a.GetCharacterName()
	session.Values["corporationID"] = a.GetCorporationID()
	session.Values["corporationName"] = a.GetCorporationName()
	session.Values["allianceID"] = a.GetAllianceID()
	session.Values["allianceName"] = a.GetAllianceName()
	session.Values["factionID"] = a.GetFactionID()
	session.Values["factionName"] = a.GetFactionName()

	corp, err := database.LoadCorporationFromName(a.GetCorporationName())
	if err != nil {
		if len(a.GetCorporationName()) > 0 && a.GetCorporationID() > 0 {
			c, err := database.SaveCorporation(&models.Corporation{
				ID:            -1,
				Name:          a.GetCorporationName(),
				CorporationID: a.GetCorporationID(),
				Ticker:        sh.Ticker})
			if err != nil {
				return fmt.Errorf("Failed to save new corporation in session: [%v]", err)
			}

			corp = c
		} else {
			return fmt.Errorf("Failed to save new corporation in session: name was empty or ID was < 0")
		}
	}

	player, err := database.LoadPlayerFromName(a.GetCharacterName())
	if err != nil {
		if len(a.GetCharacterName()) > 0 && a.GetCharacterID() > 0 {
			_, err = database.SavePlayer(&models.Player{
				ID:         -1,
				Name:       a.GetCharacterName(),
				PlayerID:   a.GetCharacterID(),
				Corp:       corp,
				AccessMask: models.AccessMaskMember,
			})
			if err != nil {
				return fmt.Errorf("Failed to save new player in session: [%v]", err)
			}
		} else {
			return fmt.Errorf("Failed to save new player in session: name was empty or ID was < 0")
		}
	}

	session.Values["playerID"] = player.ID
	session.Values["corpID"] = corp.ID

	return session.Save(r, w)
}

func (s *Session) GetCorporationName(r *http.Request) string {
	session, _ := s.store.Get(r, "player")

	corporationName, ok := session.Values["corporationName"]
	if !ok {
		return ""
	}

	return corporationName.(string)
}

func (s *Session) GetCorpID(r *http.Request) int64 {
	session, _ := s.store.Get(r, "player")

	corpID, ok := session.Values["corpID"]
	if !ok {
		return -1
	}

	return corpID.(int64)
}

func (s *Session) GetPlayerID(r *http.Request) int64 {
	session, _ := s.store.Get(r, "player")

	playerID, ok := session.Values["playerID"]
	if !ok {
		return -1
	}

	return playerID.(int64)
}

func (s *Session) SetLoginRedirect(w http.ResponseWriter, r *http.Request, redirect string) {
	session, _ := s.store.Get(r, "login")

	session.Values["redirect"] = redirect

	session.Save(r, w)
}

func (s *Session) GetLoginRedirect(r *http.Request) string {
	session, _ := s.store.Get(r, "login")
	if session.IsNew {
		return "/"
	}

	redirectInterface, ok := session.Values["redirect"]
	if !ok {
		return "/"
	}

	return redirectInterface.(string)
}

func (s *Session) SetSSOState(w http.ResponseWriter, r *http.Request, state string) {
	session, _ := s.store.Get(r, "login")

	session.Values["ssoState"] = state

	session.Save(r, w)
}

func (s *Session) GetSSOState(r *http.Request) string {
	session, _ := s.store.Get(r, "login")
	if session.IsNew {
		return ""
	}

	stateInterface, ok := session.Values["ssoState"]
	if !ok {
		return ""
	}

	return stateInterface.(string)
}
