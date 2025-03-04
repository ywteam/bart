package main

import (
	"fmt"

	config "github.com/ywteam/ydk-go/config"
	console "github.com/ywteam/ydk-go/console"
	logger "github.com/ywteam/ydk-go/logger"
	secops "github.com/ywteam/ydk-go/secops"
	server "github.com/ywteam/ydk-go/server"
	"gofr.dev/pkg/gofr"
)

func main() {
	logger.Log(logger.Info, console.Colorize("YDK Go :shorts:", console.Yellow, console.NoBackground)+" "+config.GetEnvironment().Emoji+" "+config.GetEnvironment().Name)
	server := server.NewServer(server.YdkServerConfig{Port: 3000})
	server.App.GET("/whoami", func(c *gofr.Context) (interface{}, error) {
		// return json whoami
		whoami := map[string]string{
			"Name":        "YDK Go",
			"Version":     "1.0.0",
			"Author":      "YDK Team",
			"Description": "YDK Go is a Go SDK for YDK",
		}
		return whoami, nil
	})
	secops.Ydk()
	go func() {
		if err := server.Listen(); err != nil {
			fmt.Printf("Error starting server: %v", err)
		}
		fmt.Println("Server started on port 3000")
	}()
	
}
