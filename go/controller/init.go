// Package controller holds all handlers and handler functions
// as well as necessary infrastructure for session management
// and security
package controller

import (
	"encoding/gob"
	"github.com/geobe/https-proxy/go/model"
	scc "github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"os"
)

var sessionStore = makeStore()

// keys for the session store
const S_PROXY = "Proxy-App-Session"

// helper function to create a gorilla session store with
// a strong set of keys
func makeStore() sessions.Store {
	// store sessions in temp directory to allow sessions stores larger than 4 kB
	// IE restricts cookie stores to 4 kB
	store := sessions.NewFilesystemStore("",
		scc.GenerateRandomKey(32),
		scc.GenerateRandomKey(32))
	registerTypes()
	// set session store of unlimited length
	store.MaxLength(0)
	log.Printf("storing sessions to %s\n", os.TempDir())
	return store
}

// accessor for the gorilla session store
func SessionStore() sessions.Store {
	return sessionStore
}

// register application types for serialization/deserialization
// necessary for session store
func registerTypes() {
	gob.Register(model.User{})
}
