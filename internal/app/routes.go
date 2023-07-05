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
	rMux.Handle("/test", app.RequireAuth(http.HandlerFunc(app.Test)))
	rMux.Handle("/tasks", app.RequireAuth(http.HandlerFunc(app.GenerateTasks)))

	//private
	//loader and customer profile page
	rMux.Handle("/me", app.RequireAuth(http.HandlerFunc(app.ProfilePage)))

	// available tasks for loader

	// complete tasks for customer
	// start the game

	return rMux
}
