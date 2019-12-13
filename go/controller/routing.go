package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func err(writer http.ResponseWriter, request *http.Request) {
	//error := request.Header.Get("Status")
	fmt.Fprintf(writer, "Error with path %s",
		request.URL.Path[1:])
}

func SetRouting() *mux.Router {
	mux := mux.NewRouter()

	//requestLogging := alice.New(RequestLogger)
	//csrfChecking := alice.New(nosurf.NewPure)
	//anyChecking := alice.New(RequestLogger, nosurf.NewPure, SessionChecker, AuthTarget)

	// error
	mux.HandleFunc("/err", err)
	// index
	/*	mux.Handle("/", csrfChecking.ThenFunc(HandleIndex))
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
	*/
	return mux
}
