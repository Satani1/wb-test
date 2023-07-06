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

	//private
	//loader and customer profile page
	rMux.Handle("/me", app.RequireAuth(http.HandlerFunc(app.ProfilePage)))

	//generate OR view tasks for loader OR customer
	rMux.HandleFunc("/tasks", app.GenerateTasks).Methods("POST")
	rMux.Handle("/tasks", app.RequireAuth(http.HandlerFunc(app.GenerateTasks)))

	// start the game
	rMux.Handle("/start", app.RequireAuth(http.HandlerFunc(app.Start)))

	return rMux
}
