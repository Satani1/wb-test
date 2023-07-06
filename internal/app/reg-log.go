package app

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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
			return
		}

		//create the user
		if role == "loader" {
			var loader models.Loader
			if err := gofakeit.Struct(&loader); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			loader.Username, loader.Password = username, string(hash)

			err := app.DB.InsertLoader(loader)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			customer := models.Customer{
				Username:     username,
				Password:     string(hash),
				StartCapital: gofakeit.Number(10000, 30000),
			}

			err := app.DB.InsertCustomer(customer)
			if err != nil {
				log.Fatalln(err)
			}

		}
		//respond
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte("Successful registration")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		//find html template
		ts, err := template.ParseFiles("./web/html/register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//exec html page
		err = ts.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		if userLoader, err := app.DB.GetLoaderByName(username); err == nil {
			//compare pass with pass in table
			err := bcrypt.CompareHashAndPassword([]byte(userLoader.Password), []byte(password))
			if err != nil {
				http.Error(w, "Invalid password", http.StatusBadRequest)
				return
			}

			//generate jwt token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": userLoader.Username,
				"exp": time.Now().Add(time.Hour).Unix(),
			})

			tokenString, err := token.SignedString([]byte(app.Secret))
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to create token", http.StatusInternalServerError)
				return
			}

			//send token back
			//cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "Authorization",
				Value:    tokenString,
				Path:     "",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})

			w.WriteHeader(http.StatusOK)

		} else if userCustomer, err := app.DB.GetCustomer(username); err == nil {
			//compare pass with pass in table
			err := bcrypt.CompareHashAndPassword([]byte(userCustomer.Password), []byte(password))
			if err != nil {
				http.Error(w, "Invalid password", http.StatusBadRequest)
				return
			}

			//generate jwt token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": userCustomer.Username,
				"exp": time.Now().Add(time.Hour).Unix(),
			})

			tokenString, err := token.SignedString([]byte(app.Secret))
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to create token", http.StatusInternalServerError)
				return
			}

			//send token back
			//cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "Authorization",
				Value:    tokenString,
				Path:     "",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})

			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("Successful login")); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			log.Println(err)
			http.Error(w, "Cant find any user with this username", http.StatusBadRequest)
			return
		}
	} else {
		//find html template
		ts, err := template.ParseFiles("./web/html/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//exec html page
		err = ts.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (app *Application) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("logged in")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
