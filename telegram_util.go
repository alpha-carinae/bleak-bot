package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func sendPhoto(photo Photo, config TelegramConfig) {

	if len(photo.url) == 0 {
		log.Println("Photo url is empty.")
		return
	}

	url := fmt.Sprintf("%s%s/sendPhoto", config.ApiURL, config.BotToken)

	for _, chatId := range config.ChatIds {
		str := fmt.Sprintf("{\"chat_id\":%d,\"photo\":\"%s\",\"caption\":\"%s\",\"disable_notification\":%t}",
			chatId, photo.url, photo.description, config.DisableNotification)
		resp, err := http.Post(url, "application/json", bytes.NewBufferString(str))
		if err != nil {
			log.Println("Error sending photo to chat: ", chatId, err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			log.Println("Photo sent.")
		} else {
			log.Println("Could not send photo: ", resp.StatusCode)

			log.Println("response: ", resp)
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading response: ", err)
				continue
			}
			log.Printf("response Body: %s", bytes)
			resp.Body.Close()
		}
	}
}

func getUpdates(config TelegramConfig) {
	url := fmt.Sprintf("%s%s/getUpdates", config.ApiURL, config.BotToken)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error getting updates: ", err)
		return
	}
	defer resp.Body.Close()

	log.Println("response: ", resp)
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response: ", err)
		return
	}
	log.Printf("response Body: %s", bytes)
}
