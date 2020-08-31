package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type UnsplashConfig struct {
	AccessKey        string
	ApiURL           string
	RandomPhotosPath string
	ImageCount       int
	Queries          []string
}

type TelegramConfig struct {
	BotToken                 string
	ApiURL                   string
	SendingIntervalInMinutes int
	ChatIds                  []int64
	DisableNotification      bool
}

type Configuration struct {
	Unsplash UnsplashConfig
	Telegram TelegramConfig
}

type Photo struct {
	id          string
	description string
	url         string
}

func main() {
	config := readConfig()

	interval := time.Duration(config.Telegram.SendingIntervalInMinutes)
	c := time.Tick(interval * time.Minute)

	for range c {
		photo := GetRandomPhoto(config.Unsplash)
		sendPhoto(photo, config.Telegram)
	}
}

func readConfig() Configuration {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("error:", err)
	}
	log.Println("Successfully read the configuration.")

	return configuration
}
