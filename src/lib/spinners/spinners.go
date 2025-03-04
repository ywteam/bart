package spinners

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	logger "yellowteam/lib/logger"
)

type Spinner struct {
	Interval float64  `json:"interval"`
	Frames   []string `json:"frames"`
}

var spinners map[string]Spinner = make(map[string]Spinner)

type SpinOperation func() bool

func (spinner *Spinner) Start(message string, operation SpinOperation) {
	startAt := time.Now()
	chann := make(chan bool)
	defer close(chann)
	go func() {
		for {
			for _, frame := range spinner.Frames {
				elapsed := time.Since(startAt)
				fmt.Printf("\r%s %s %.3f 🩳", frame, message, elapsed.Seconds())
				time.Sleep(time.Duration(spinner.Interval) * time.Millisecond)
				select {
				case <-chann:
					fmt.Printf("\r\033[K")
					return
				default:
				}
			}
		}
	}()
	result := operation()
	chann <- result
	elapsed := fmt.Sprintf("%.3fs", time.Since(startAt).Seconds())
	if result {
		logger.Success(fmt.Sprintf("%s in %s", message, elapsed))
	} else {
		logger.Error(fmt.Sprintf("%s in %s", message, elapsed))
	}
}
func New(name string, interval float64, frames []string) *Spinner {
	if _, ok := spinners[name]; ok {
		spinner := spinners[name]
		return &spinner
	} else {
		spinner := Spinner{
			Interval: interval,
			Frames:   frames,
		}
		spinners[name] = spinner
		return &spinner
	}
}
func Spin(name string) *Spinner {
	if spinner, ok := spinners[name]; ok {
		return &spinner
	}
	return nil
}
func Spinners() []string {
	keys := make([]string, 0, len(spinners))
	for key := range spinners {
		keys = append(keys, key)
	}
	return keys
}
func Download() {
	url := "https://assets.yellowteam.cloud/spinners.json"
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	cacheFile := path.Join(cacheDir, "spinners.json")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		resp, err := http.Get(url)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		defer resp.Body.Close()
		out, err := os.Create(cacheFile)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
	file, err := os.Open(cacheFile)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var jsonData = make(map[string]interface{})
	err = decoder.Decode(&jsonData)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for key, value := range jsonData {
		interval := value.(map[string]interface{})["interval"].(float64)
		framesInterface := value.(map[string]interface{})["frames"].([]interface{})
		frames := make([]string, len(framesInterface))
		for i, frame := range framesInterface {
			frames[i] = frame.(string)
		}
		New(key, interval, frames)
		// spinners[key] = Spinner{
		// 	Interval: interval,
		// 	Frames:   frames,
		// }
	}
}

func init() {
	New("default", 100, []string{"-", "\\", "|", "/"}).Start(
		"Dowloading spinners",
		func() bool {
			Download()
			return true
		},
	)

	// test all spinners
	// for key, _ := range spinners {
	// 	Spin(key).Start(
	// 		fmt.Sprintf("Testing spinner %s", key),
	// 		func() bool {
	// 			time.Sleep(5 * time.Second)
	// 			// return random boolean
	// 			return len(key)%2 == 0
	// 			// return true
	// 		},
	// 	)
	// }

	// logger.Notice("%d spinners loaded %s", len(spinners), Spinners())
	// Spin("hearts").Start(
	// 	"Loading",
	// 	func() bool {
	// 		// sleep for 5 seconds
	// 		time.Sleep(5 * time.Second)
	// 		return true
	// 	},
	// )
	// New("spinner-2", 100, []string{"🩳", "🩲", "🩱", "🩰"}).Start(
	// 	"Loading",
	// 	func() bool {
	// 		// sleep for 5 seconds
	// 		time.Sleep(5 * time.Second)
	// 		return true
	// 	},
	// )
}
