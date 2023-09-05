package telegram

import (
	"bytes"
	"encoding/json"
	"health-checker/helpers"
	"log"
	"net/http"
	"os"
	"strconv"
)

type TelegramMessage struct {
	ChatID              int    `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}

type TelegramSendResponse struct {
	Ok bool `json:"ok"`
}

func Send(message string, DisableNotification bool) error {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))

	if err != nil {
		log.Println("ERROR" + err.Error())
		return err
	}

	enableTelegramNotification, err := strconv.ParseBool(os.Getenv("ENABLE_TELEGRAM_NOTIFICATION"))
	if err != nil || !enableTelegramNotification {
		log.Println("[TELEGRAM] " + message)
		log.Println("env 'ENABLE_TELEGRAM_NOTIFICATION' is not found or false. Skipping telegram notification")
		return nil
	}

	botMessage := TelegramMessage{
		ChatID:              telegramChatID,
		Text:                message,
		DisableNotification: DisableNotification,
	}
	botMessageBytes, err := json.Marshal(botMessage)
	if err != nil {
		log.Println("ERROR" + err.Error())
		return err
	}
	telegramAPIURL := "https://api.telegram.org/bot" + telegramBotToken + "/sendMessage"

	resp, err := http.Post(telegramAPIURL, "application/json", bytes.NewBuffer(botMessageBytes))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	responseBody := TelegramSendResponse{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	if !responseBody.Ok {
		log.Fatalln("Response of status is not ok.")
		log.Println(telegramAPIURL)
		helpers.PrintByte(botMessageBytes)
		helpers.PrintReader(resp.Body)
	}

	return nil
}
