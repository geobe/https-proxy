package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"net/http"
	"os"
)

func err(writer http.ResponseWriter, request *http.Request) {
	//error := request.Header.Get("Status")
	fmt.Fprintf(writer, "Error with path %s",
		request.URL.Path[1:])
}

func SetRouting() *mux.Router {
	mux := mux.NewRouter()
	// finde Working directory = GOPATH
	docbase, _ := os.Getwd()
	docbase += "/"
	// FileServer ist ein Handler, der dieses Verzeichnis bedient
	// Funktionsvariablen für alice anlegen
	files := http.StripPrefix("/pages/", http.FileServer(http.Dir(docbase + pagedir)))
	resources := http.FileServer(http.Dir(docbase + resourcedir))
	pages := http.FileServer(http.Dir(docbase + pagedir))

	requestLogging := alice.New(RequestLogger)
	csrfChecking := alice.New(nosurf.NewPure)
	resultsChecking := alice.New(RequestLogger, nosurf.NewPure, SessionChecker, AuthProjectOffice)
	enroleChecking := alice.New(RequestLogger, nosurf.NewPure, SessionChecker, AuthEnrol)
	anyChecking := alice.New(RequestLogger, nosurf.NewPure, SessionChecker, AuthAny)

	// Zugriff auf das Verzeichnis via Präfic /pages/
	mux.PathPrefix("/pages/").Handler(requestLogging.Then(files))
	// Zugriff auf die Resourcen-Verzeichnisse mit regular expression
	mux.PathPrefix("/{dir:(?:css|fonts|js|images)}/").Handler(requestLogging.Then(resources))
	// Zugriff auf *.htm[l] Dateien im /pages Verzeichnis
	mux.Handle("/{dir:\\w+\\.html?}", requestLogging.Then(pages))
	// error
	mux.HandleFunc("/err", err)
	// index
	mux.Handle("/", csrfChecking.ThenFunc(HandleIndex))
	mux.Handle("/index", csrfChecking.ThenFunc(HandleIndex))
	// login
	mux.Handle("/login", csrfChecking.ThenFunc(HandleLogin))
	// logout
	mux.HandleFunc("/logout", HandleLogout)
	// work
	mux.Handle("/work", anyChecking.ThenFunc(HandleWork))
	// find
	mux.Handle("/find/applicant", enroleChecking.ThenFunc(FindApplicant))
	// show results edit form
	mux.Handle("/results/show", resultsChecking.ThenFunc(ShowResults))
	// submit results edit form
	mux.Handle("/results/submit", resultsChecking.ThenFunc(SubmitResults))
	// show enrol form
	mux.Handle("/enrol/show", enroleChecking.ThenFunc(ShowEnrol))
	// process enrol form
	mux.Handle("/enrol/submit", enroleChecking.ThenFunc(SubmitEnrol))
	// process enrol delete
	mux.Handle("/enrol/delete", enroleChecking.ThenFunc(SubmitApplicantDelete))
	// show cancellation form
	mux.Handle("/cancellation/show", enroleChecking.ThenFunc(ShowCancellation))
	// process cancellation form
	mux.Handle("/cancellation/submit", enroleChecking.ThenFunc(SubmitCancelation))
	// process edit form
	mux.Handle("/edit/submit", enroleChecking.ThenFunc(SubmitApplicantEdit))
	// register
	mux.Handle("/register", csrfChecking.ThenFunc(ShowRegistration))
	// register
	mux.Handle("/register/submit", csrfChecking.ThenFunc(SubmitRegistration))
	return mux
}
