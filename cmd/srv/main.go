package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"wb-test/internal/app"
	"wb-test/pkg/db"
)

//const dsn string = "root:YaPoc290302@/wb-test?parseTime=true"

func main() {
	//read config variables
	cfg := LoadEnvVariables()

	//open DB
	dbURL := fmt.Sprintf("%s:%s@/%s?parseTime=true", cfg.MyUser, cfg.MyPassword, cfg.MyDB)

	appDB, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer appDB.Close()

	//server setup and start
	App := &app.Application{
		Addr:   cfg.ServerAddr,
		DB:     appDB,
		Secret: cfg.Secret,
	}
	srv := http.Server{
		Addr:    App.Addr,
		Handler: App.Routes(),
	}

	//running http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}

	log.Println("Shutdown the server.")
}
