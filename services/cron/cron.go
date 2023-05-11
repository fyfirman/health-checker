package cron

import (
	"health-checker/env"
	"health-checker/services/telegram"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
)

func checkHealth(url string) (bool, error) {
	log.Println("Checking the service, url :", url)

	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer resp.Body.Close()

	return true, nil
}

type HealthCheckResult struct {
	ID        string
	Name      string
	Error     string
	Status    bool
	CheckedAt time.Time
}

func constructMessage(data []HealthCheckResult) string {
	if len(data) == 0 {
		return "👨🏻‍⚕️ Health Check Report - no service found"
	}

	date := data[0].CheckedAt.Format("02 Jan 2006, 15.04.05")
	message := "👨🏻‍⚕️ Health Check Report " + "\n"
	message += "🗓️ " + date + "\n\n"

	for _, res := range data {
		message += "• " + res.Name + " "

		if res.Status {
			message += "✅\n"
			continue
		}

		message += "❌ (" + res.Error + ")\n"
	}

	return message
}

func Start(env *env.Env) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Minutes().Do(func() {
		log.Println("Cron every 5 minutes starting...")

		res, err := GetAll(env)

		if err != nil {
			log.Println("Error : " + err.Error())
			return
		}

		var results []HealthCheckResult

		for _, service := range res {
			errMessage := ""
			res, err := checkHealth(service.URL)
			if err != nil {
				log.Println("ERROR: " + err.Error())
				errMessage = err.Error()
			}

			var payload map[string]interface{} = map[string]interface{}{
				"id":              service.ID,
				"last_status":     res,
				"last_checked_at": time.Now(),
				"updated_at":      time.Now(),
			}
			err = UpdateById(env, service.ID, payload)

			if err != nil {
				log.Println("ERROR: " + err.Error())
				errMessage = err.Error()
			}

			log.Println("Append the array")
			results = append(results, HealthCheckResult{
				ID:        service.ID,
				Name:      service.Name,
				Error:     errMessage,
				Status:    res,
				CheckedAt: time.Now(),
			})

			log.Println("Loop of " + service.Name + " is end")
		}

		message := constructMessage(results)
		log.Println(message)
		telegram.Send(message, false)
	})

	s.StartAsync()
}
