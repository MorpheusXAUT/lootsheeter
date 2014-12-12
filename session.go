// session
package main

import (
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
	store   *sessions.FilesystemStore
	players map[string]*models.Player
}

func NewSession() *Session {
	session := &Session{
		store:   sessions.NewFilesystemStore("web/sessions", []byte(securecookie.GenerateRandomKey(128))),
		players: make(map[string]*models.Player),
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

	characterName := session.Values["character_name"].(string)

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

	_, ok := s.players[characterName]
	if !ok {
		delete(s.players, characterName)
	}
}

func (s *Session) GetPlayerFromRequest(r *http.Request) *models.Player {
	session, _ := s.store.Get(r, "player")

	characterName, ok := session.Values["character_name"].(string)
	if !ok {
		return nil
	}

	player, ok := s.players[characterName]
	if !ok {
		p, err := database.LoadPlayerFromName(characterName)
		if err != nil {
			logger.Errorf("Failed to load player from database in session: [%v]", err)
			return nil
		}

		player = p
		s.players[characterName] = p
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
	} else {
		if strings.EqualFold(session.Values["character_name"].(string), player.Name) {
			return true
		} else {
			return false
		}
	}
}

func (s *Session) SetIdentity(w http.ResponseWriter, r *http.Request, a models.CharacterAffiliation, sh models.CorporationSheet) {
	session, _ := s.store.Get(r, "player")

	session.Values["character_id"] = a.GetCharacterId()
	session.Values["character_name"] = a.GetCharacterName()
	session.Values["corporation_id"] = a.GetCorporationId()
	session.Values["corporation_name"] = a.GetCorporationName()
	session.Values["alliance_id"] = a.GetAllianceId()
	session.Values["alliance_name"] = a.GetAllianceName()
	session.Values["faction_id"] = a.GetFactionId()
	session.Values["faction_name"] = a.GetFactionName()

	corp, err := database.LoadCorporationFromName(a.GetCorporationName())
	if err != nil {
		c, err := database.SaveCorporation(&models.Corporation{
			Id:     -1,
			Name:   a.GetCorporationName(),
			CorpId: a.GetCorporationId(),
			Ticker: sh.Ticker})
		if err != nil {
			logger.Errorf("Failed to save new corporation in session: [%v]", err)
			return
		}

		corp = c
	}

	_, err = database.LoadPlayerFromName(a.GetCharacterName())
	if err != nil {
		_, err = database.SavePlayer(&models.Player{
			Id:         -1,
			Name:       a.GetCharacterName(),
			PlayerId:   a.GetCharacterId(),
			Corp:       corp,
			AccessMask: models.AccessMaskNone,
		})
		if err != nil {
			logger.Errorf("Failed to save new player in session: [%v]", err)
			return
		}
	}

	session.Save(r, w)
}

func (s *Session) GetCorporationName(r *http.Request) string {
	session, _ := s.store.Get(r, "player")

	return session.Values["corporation_name"].(string)
}

func (s *Session) SetSSOState(w http.ResponseWriter, r *http.Request, state string) {
	session, _ := s.store.Get(r, "login")

	session.Values["sso_state"] = state

	session.Save(r, w)
}

func (s *Session) GetSSOState(r *http.Request) string {
	session, _ := s.store.Get(r, "login")
	if session.IsNew {
		return ""
	}

	return session.Values["sso_state"].(string)
}
