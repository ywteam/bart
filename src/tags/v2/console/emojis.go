package ydk

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var emojis map[string]string

// RandomEmoji returns a random emoji key-value pair from the emojis map.
// If the emojis map is empty, it returns the default key-value pair "0" and "0".
func RandomEmoji() (key string, value string) {
	keys := make([]string, 0, len(emojis))
	for k := range emojis {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return "0", "0"
	}
	randomIndex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(keys))
	return keys[randomIndex], emojis[keys[randomIndex]]
}

// ListEmojis returns a map of emojis, where the key is the emoji name and the value is the corresponding emoji code.
func ListEmojis() map[string]string {
	return emojis
}

// InterpolateEmojis replaces emoji codes in the given text with their corresponding emoji characters.
// It iterates over the emojis map and replaces each occurrence of a code with its corresponding emoji character.
// The function returns the updated text with emojis.
func InterpolateEmojis(text string) string {
	// match all emoji name :emoji_name: in the text using regex
	emojiRegex := regexp.MustCompile(`:[a-zA-Z0-9_]+:`)
	matches := emojiRegex.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		// remove the colons from the emoji name
		emojiName := strings.Trim(match[0], ":")
		if emojiName == "random" {
			emojiName, _ = RandomEmoji()
		}
		// get the emoji character from the emojis map
		emoji, ok := emojis[emojiName]
		if ok {
			// replace the emoji name with the emoji character in the text
			text = strings.ReplaceAll(text, match[0], emoji)
		}
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
