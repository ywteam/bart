package ydk

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var emojis map[string]string

// RandomEmoji returns a random emoji key-value pair from the emojis map.
// If the emojis map is empty, it returns the default key-value pair "0" and "0".
func RandomEmoji() (key string, value string) {
	for key, value := range emojis {
		return key, value
	}
	return "0", "0"
}

// ListEmojis returns a map of emojis, where the key is the emoji name and the value is the corresponding emoji code.
func ListEmojis() map[string]string {
	return emojis
}

// InterpolateEmojis replaces emoji codes in the given text with their corresponding emoji characters.
// It iterates over the emojis map and replaces each occurrence of a code with its corresponding emoji character.
// The function returns the updated text with emojis.
func InterpolateEmojis(text string) string {
	for key, value := range emojis {
		text = strings.ReplaceAll(text, ":"+key+":", value)
	}
	return text
}

// init initializes the emojis map by fetching the emojis JSON file from a remote URL and decoding it.
// It also caches the emojis JSON file locally for future use.
func init() {
	emojis = make(map[string]string)
	url := "https://assets.yellowteam.cloud/emojis.json"
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Printf("Error getting cache dir: %s", err)
		return
	}
	cacheFile := path.Join(cacheDir, "emojis.json")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error getting emojis: %s", err)
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
			log.Printf("Error copying emojis to cache: %s", err)
			return

		}
	}
	file, err := os.Open(cacheFile)	
	if err != nil {
		log.Printf("Error opening cache file: %s", err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&emojis)
	if err != nil {
		log.Printf("Error decoding emojis: %s", err)
		return
	}
}
