package srv

import (
	"html/template"
	"net/http"
	"sync"
)

var _instance *WebSrv
var once sync.Once

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("super-secret-key")
	//store = sessions.NewCookieStore(key)
	templates *template.Template
)

type IndexData struct {
	Title   string
	Message string
}

func Inst() *WebSrv {
	once.Do(func() {
		_instance = newSrv()
	})
	return _instance
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := IndexData{
		Title:   "Welcome to Our Site!",
		Message: "This is some data from Go!",
	}

	templates.ExecuteTemplate(w, "index.html", data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Your login handler logic
}
func newSrv() *WebSrv {
	var ws = &WebSrv{}
	templates = template.Must(template.ParseGlob("web/templates/*.html"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)

	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	return ws
}

type WebSrv struct {
}

func (s *WebSrv) Start() error {
	go func() {
		http.ListenAndServe(":80", nil)
	}()
	return nil
}

func (s *WebSrv) ShutDown() {

}
