// store
package main

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var store = sessions.NewFilesystemStore("web/sessions", []byte(securecookie.GenerateRandomKey(64)))
