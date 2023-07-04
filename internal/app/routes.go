package app

import (
	"github.com/gorilla/mux"
)

func (app *Application) Routes() *mux.Router {
	rMux := mux.NewRouter()

	//public

	rMux.HandleFunc("/login", app.SignIn)
	rMux.HandleFunc("/register", app.SignUp)
	//rMux.HandleFunc("/tasks",)

	//loader
	// me
	// tasks

	//customer
	// me
	// tasks
	// start

	return rMux
}
