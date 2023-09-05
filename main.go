package main

import (
	"health-checker/env"
	cronService "health-checker/services/cron"
	"health-checker/services/telegram"

	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	appVersion := "0.1.2"
	log.Println("🚀 Starting health-checker v" + appVersion)

	godotenv.Load(".env")

	db, err := ConnectDB()

	if err != nil {
		panic(err.Error())
	}
	err = telegram.Send("🚀 health-checker v"+appVersion+" has started", false)
	if err != nil {
		log.Fatal(err)
	}

	env := &env.Env{Db: db}

	cronService.Start(env)

	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	http.Handle("/", router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Application served at :8080")
	log.Fatal(srv.ListenAndServe())
}
