package app

import (
	_ "golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"wb-test/pkg/models"
)

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
			var win bool

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
				win = true
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
				if win {
					if err := app.DB.UpdateTask(taskID, id); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}

			w.WriteHeader(http.StatusOK)
		}
	} else if userRole == "loader" {
		http.Error(w, "You aren't a customer", http.StatusBadRequest)
		return
	}
}
