package ydk

import (
	"encoding/json"
	"fmt"
	logger "github.com/ywteam/ydk-go/logger"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

type YdkSpinnerData struct {
	interval float64  `json:"interval"`
	frames   []string `json:"frames"`
}

// var spinners = make(map[string]interface{})
var spinners map[string]YdkSpinnerData

func ListSpinners() map[string]YdkSpinnerData {
	return spinners
}
func RandomSpinner() (key string, value YdkSpinnerData) {
	for key, value := range spinners {
		return key, value
	}
	return "0", YdkSpinnerData{}
}
func SpinnerFrames(key string) []string {
	return spinners[key].frames
}
func SpinnerInterval(key string) float64 {
	return spinners[key].interval
}

type SpinOperation func(message *string) bool

func Start(key string, message string, operation SpinOperation) {
	initialMessage := message
	if key == "random" {
		key, _ = RandomSpinner()
	}
	_, ok := spinners[key]
	if !ok {
		logger.Default().Warning(fmt.Sprintf("Spinner with key %s not found", key))
		operation(&message)
		return
	}
	spinner := spinners[key]
	startAt := time.Now()
	// crate new channel
	chann := make(chan bool)
	defer close(chann)
	go func() {
		for {
			for _, frame := range spinner.frames {
				elapsed := time.Since(startAt)
				fmt.Printf("\r%s", frame)
				fmt.Printf(" %s %.3f ", message, elapsed.Seconds())
				// logger.Default().Info(message)
				time.Sleep(time.Duration(spinner.interval) * time.Millisecond)
				select {
				case <-chann:
					Stop()
					return
				default:
				}
			}
		}

	}()
	result := operation(&message)
	chann <- result
	elapsed := fmt.Sprintf("%.3fs", time.Since(startAt).Seconds())
	if result {
		logger.Default().Success(initialMessage + ". Operation completed successfully. " + elapsed)
	} else {
		logger.Default().Error(initialMessage + ". Operation failed. " + elapsed)
	}
}
func Stop() {
	fmt.Printf("\r\033[K")
}
func init() {
	spinners = make(map[string]YdkSpinnerData)
	url := "https://assets.yellowteam.cloud/spinners.json"
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Printf("Error getting cache dir: %s", err)
		return
	}
	cacheFile := path.Join(cacheDir, "spinners.json")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error getting spinners: %s", err)
			return
		}
		defer resp.Body.Close()
		out, err := os.Create(cacheFile)
		if err != nil {
			log.Printf("Error creating cache file: %s", err)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Printf("Error copying spinners to cache: %s", err)
			return

		}
	}
	// log.Printf("Opening cache file: %s", cacheFile)
	file, err := os.Open(cacheFile)
	if err != nil {
		log.Printf("Error opening cache file: %s", err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var jsonData = make(map[string]interface{})
	err = decoder.Decode(&jsonData)
	if err != nil {
		log.Printf("Error decoding emojis: %s", err)
		return
	}
	for key, value := range jsonData {
		interval := value.(map[string]interface{})["interval"].(float64)
		framesInterface := value.(map[string]interface{})["frames"].([]interface{})
		frames := make([]string, len(framesInterface))
		for i, frame := range framesInterface {
			frames[i] = frame.(string)
		}
		spinners[key] = YdkSpinnerData{
			interval: interval,
			frames:   frames,
		}
	}
}
