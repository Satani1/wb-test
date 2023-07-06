package app

import (
	"github.com/brianvoe/gofakeit/v6"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"wb-test/pkg/models"
)

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
