package srv

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"sync"
)

var _instance *WebSrv
var once sync.Once

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func Inst() *WebSrv {
	once.Do(func() {
		_instance = newSrv()
	})
	return _instance
}

func newSrv() *WebSrv {
	var ws = &WebSrv{}
	return ws
}

type WebSrv struct {
}

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
func (s *WebSrv) Start() error {
	http.HandleFunc("/login", login)

	go func() {
		http.ListenAndServe(":80", nil)
	}()
	return nil
}

func (s *WebSrv) ShutDown() {

}
