package app

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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
			//loader := models.Loader{
			//	ID:        0,
			//	Username:  username,
			//	Password:  string(hash),
			//	MaxWeight: 0,
			//	Drunk:     false,
			//	Fatigue:   0,
			//	Salary:    0,
			//}
			var loader models.Loader
			if err := gofakeit.Struct(&loader); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			loader.Username, loader.Password = username, string(hash)

			id, err := app.DB.InsertLoader(loader)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(id)
		} else {
			customer := models.Customer{
				Username:     username,
				Password:     string(hash),
				StartCapital: gofakeit.Number(10000, 30000),
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

func (app *Application) GenerateTasks(w http.ResponseWriter, r *http.Request) {
	userRole := r.Header.Get("role")

	if userRole != "customer" && userRole != "loader" {
		numberOfTasks := rand.Intn(10)

		for i := 0; i < numberOfTasks; i++ {
			//generate random task
			var task models.Task
			if err := gofakeit.Struct(&task); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//insert into DB
			if err := app.DB.InsertTask(task); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte("Tasks are created")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		user := r.Header.Get("user")
		userID, err := strconv.Atoi(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userRole := r.Header.Get("role")
		var tasksAvailable struct {
			Tasks []models.Task
		}

		switch userRole {
		case "loader":
			tasks, err := app.DB.GetTaskCompleted(userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tasksAvailable.Tasks = tasks
		case "customer":
			tasks, err := app.DB.GetTaskAvailable()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tasksAvailable.Tasks = tasks
		}

		//find template
		ts, err := template.ParseFiles("./web/html/tasks.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//execute html template with user data
		err = ts.Execute(w, tasksAvailable)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (app *Application) Start(w http.ResponseWriter, r *http.Request) {
	//check user role
	userRole := r.Header.Get("role")
	userName := r.Header.Get("name")
	log.Println(userRole)

	if userRole == "customer" {
		if r.Method == "GET" {
			//find template
			ts, err := template.ParseFiles("./web/html/start.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//execute html template
			if err := ts.Execute(w, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			//Start the game
			loadersIDstr := r.FormValue("loaders")
			taskIDstr := r.FormValue("task")
			log.Println(loadersIDstr)
			log.Println(taskIDstr)

			//loaders IDs convert to []int
			loadersStr := strings.Split(loadersIDstr, ",")
			loaders := make([]int, len(loadersStr))
			for i := 0; i < len(loadersStr); i++ {
				//if there are spaces
				loadersStr[i] = strings.Trim(loadersStr[i], " ")

				//convert to int
				var err error
				loaders[i], err = strconv.Atoi(loadersStr[i])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			//task ID convert to int
			taskID, err := strconv.Atoi(taskIDstr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Println(loaders, reflect.TypeOf(loaders))
			log.Println(taskID, reflect.TypeOf(taskID))

			//game

			//loaders map
			userLoaders := make(map[int]models.Loader)

			//fill loaders map
			for _, id := range loaders {
				loader, err := app.DB.GetLoaderByID(id)
				if err != nil {
					log.Println("Loader", id, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				userLoaders[loader.ID] = *loader
			}

			//calculate winnable
			//get customer data from DB
			userCustomer, err := app.DB.GetCustomer(userName)
			if err != nil {
				log.Println("Customer", userName, "||", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			//get task data from DB
			task, err := app.DB.GetTask(taskID)
			log.Println("TASK", task)
			if err != nil {
				log.Println("task", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var LoadersPrice int
			var CarryWeight float64

			for i := 0; i < len(loaders); i++ {
				//calculate price
				userLoader := userLoaders[loaders[i]]
				LoadersPrice += userLoader.Salary

				//calculate carry weight by loaders
				if userLoader.Drunk {
					log.Println(float64(userLoader.MaxWeight), float64(userLoader.Fatigue)/100, float64(userLoader.Fatigue+50)/100)
					CarryWeight += float64(userLoader.MaxWeight) * ((100 - float64(userLoader.Fatigue)) / 100) * (float64(userLoader.Fatigue+50) / 100)
				} else {
					CarryWeight += float64(userLoader.MaxWeight) * ((100 - float64(userLoader.Fatigue)) / 100)
				}
			}

			//lose the game if customer cant afford the loaders price
			log.Println(LoadersPrice, userCustomer.StartCapital)
			if LoadersPrice > userCustomer.StartCapital {
				if _, err := w.Write([]byte("You lose!")); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

			} else if float64(task.Weight) > CarryWeight {
				log.Println(CarryWeight, task.Weight)
				//lose the game if loaders cant carry weight
				if _, err := w.Write([]byte("You lose!")); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

			} else {
				//win the game
				if _, err := w.Write([]byte("You win!")); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			//customer change capital
			if err := app.DB.UpdateCustomer(userCustomer.ID, userCustomer.StartCapital-LoadersPrice); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//change loaders fatigue
			//update loaders win done tasks
			for id, loader := range userLoaders {
				if err := app.DB.UpdateLoader(loader.Fatigue+20, id); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if err := app.DB.UpdateTask(taskID, id); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
		}
	} else if userRole == "loader" {
		http.Error(w, "You aren't a customer", http.StatusBadRequest)
		return
	}
}
