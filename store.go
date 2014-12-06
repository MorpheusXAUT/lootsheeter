// store
package main

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/morpheusxaut/lootsheeter/models"
)

var sessionStore = sessions.NewFilesystemStore("web/sessions", []byte(securecookie.GenerateRandomKey(64)))

func GetPlayerFromRequest(r *http.Request) *models.Player {
	session, _ := sessionStore.Get(r, "player")

	playerName, ok := session.Values["username"].(string)
	if !ok {
		return nil
	}

	player, err := database.LoadPlayerFromName(playerName)
	if err != nil {
		return nil
	}

	return player
}

func IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	session, _ := sessionStore.Get(r, "player")
	if session.IsNew {
		return false
	}

	player := GetPlayerFromRequest(r)
	if player == nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "player",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		return false
	} else {
		return true
	}
}
