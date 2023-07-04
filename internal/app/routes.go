package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *Application) Routes() *mux.Router {
	rMux := mux.NewRouter()

	//public

	rMux.HandleFunc("/login", app.SignIn)
	rMux.HandleFunc("/register", app.SignUp)
	//rMux.HandleFunc("/tasks",)
	rMux.Handle("/test", app.RequireAuth(http.HandlerFunc(app.Test)))

	//loader
	// me
	// tasks

	//customer
	// me
	// tasks
	// start

	return rMux
}
