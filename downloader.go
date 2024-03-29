package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type void struct{}

var placeHolder void

var sentImages = make(map[string]void) // set

// Picks random photo from fetched images and returns.
func GetRandomPhoto(config UnsplashConfig) Photo {
	photosJson, err := FetchRandomPhotos(config)
	if err != nil {
		return Photo{"", "", ""}
	}

	var objArr []map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(photosJson))

	for decoder.More() {
		err := decoder.Decode(&objArr)
		if err != nil {
			log.Println("Error decoding json: ", err)
		}
	}

	index := getRandomIndex(&objArr)
	if index == -1 {
		return Photo{"", "", ""}
	}

	id := fmt.Sprintf("%s", objArr[index]["id"])

	var description string
	if objArr[index]["description"] == nil {
		description = ""
	} else {
		description = fmt.Sprintf("%s", objArr[index]["description"])
	}

	imageUrl := fmt.Sprintf("%s", objArr[index]["urls"].(map[string]interface{})["regular"])

	sentImages[id] = placeHolder
	saveId(id)

	return Photo{id, description, imageUrl}
}

// From what I've experienced, Unsplash's random API is not that random.
// So, this function tries to get random index until it finds an image that hasn't been sent before.
func getRandomIndex(objArr *[]map[string]interface{}) int {

	if len(*objArr) == 0 {
		return -1
	}

	rand.Seed(time.Now().UnixNano())
	randIdx := rand.Intn(len(*objArr))

	attempt := 0
	for attempt < len(*objArr) {
		attempt += 1

		id := fmt.Sprintf("%s", (*objArr)[randIdx]["id"])

		if _, exists := sentImages[id]; exists {
			// try with another random index
			randIdx = rand.Intn(len(*objArr))
		} else {
			break
		}
	}

	return randIdx
}

// Fetches image data by given configuration and returns as json.
func FetchRandomPhotos(config UnsplashConfig) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s%s?count=%d&query=%s",
		config.ApiURL,
		config.RandomPhotosPath,
		config.ImageCount,
		strings.Join(config.Queries, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating new request: ", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", config.AccessKey))
	req.Header.Add("Accept-Version", "v1")
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Host returned error: ", err)
		return "", err
	}
	
	defer resp.Body.Close()

	//log.Println("response: ", resp)
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response: ", err)
	}
	//log.Printf("response Body: %s", bytes)
	return string(bytes), nil
}

func saveId(imageId string) {
	file, err := os.OpenFile("sentImages.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(fmt.Sprintf("%s\n", imageId)); err != nil {
		log.Println(err)
	}
}

func SetPreviouslySentImages() {
	file, err := os.Open("sentImages.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sentImages[scanner.Text()] = placeHolder
	}

	log.Printf("Set %d previously sent images.\n", len(sentImages))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
