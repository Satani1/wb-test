package app

import (
	"html/template"
	"net/http"
	"wb-test/pkg/models"
)

func (app *Application) ProfilePage(w http.ResponseWriter, r *http.Request) {
	userName := r.Header.Get("name")
	userRole := r.Header.Get("role")

	switch userRole {
	case "loader":
		userLoader, err := app.DB.GetLoaderByName(userName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		//find template
		ts, err := template.ParseFiles("./web/html/profileLoader.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//execute html template with user data
		err = ts.Execute(w, userLoader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case "customer":

		userCustomer, err := app.DB.GetCustomer(userName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		loaders, err := app.DB.GetLoaders()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var cUser struct {
			Customer *models.Customer
			Loaders  []models.Loader
		}

		cUser.Customer, cUser.Loaders = userCustomer, loaders

		w.Header().Set("Content-Type", "text/html")
		//find template
		ts, err := template.ParseFiles("./web/html/profileCustomer.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//execute html template with user data
		err = ts.Execute(w, cUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "No role", http.StatusBadRequest)
		return
	}
}
