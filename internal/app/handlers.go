package app

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"time"
	"wb-test/pkg/models"
)

func (app *Application) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//get username, password and role
		username := r.FormValue("username")
		password := r.FormValue("password")
		role := r.FormValue("role")

		//hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		//create the user
		if role == "loader" {
			loader := models.Loader{
				ID:        0,
				Username:  username,
				Password:  string(hash),
				MaxWeight: 0,
				Drunk:     false,
				Fatigue:   0,
				Salary:    0,
			}

			id, err := app.DB.InsertLoader(loader)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(id)
		} else {
			customer := models.Customer{
				ID:           0,
				Username:     username,
				Password:     string(hash),
				StartCapital: 0,
				Tasks:        models.Task{},
			}

			id, err := app.DB.InsertCustomer(customer)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(id)
		}
		//respond
		w.WriteHeader(http.StatusCreated)
	} else {
		//find html template
		ts, err := template.ParseFiles("./web/html/register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//exec html page
		err = ts.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (app *Application) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//get username and password
		username := r.FormValue("username")
		password := r.FormValue("password")

		//look up for requested user
		//get from DB
		if userLoader, err := app.DB.GetLoader(username); err == nil {
			//compare pass with pass in table
			err := bcrypt.CompareHashAndPassword([]byte(userLoader.Password), []byte(password))
			if err != nil {
				http.Error(w, "Invalid password", http.StatusBadRequest)
			}

			//generate jwt token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": userLoader.ID,
				"exp": time.Now().Add(time.Hour).Unix(),
			})

			tokenString, err := token.SignedString([]byte(app.Secret))
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to create token", http.StatusInternalServerError)
			}

			//send token back
			//cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "Auth",
				Value:    tokenString,
				Path:     "",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})

			w.WriteHeader(http.StatusOK)

		} else if _, err := app.DB.GetCustomer(username); err == nil {
			//compare pass with pass in table

			//generate jwt token

			//send token back
		} else {
			log.Println(err)
			http.Error(w, "Cant find any user with this username", http.StatusBadRequest)
		}
	} else {
		//find html template
		ts, err := template.ParseFiles("./web/html/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//exec html page
		err = ts.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
