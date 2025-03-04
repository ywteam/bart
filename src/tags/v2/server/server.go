package ydk

import (
	"gofr.dev/pkg/gofr"
	logger "github.com/ywteam/ydk-go/logger"
)

type YdkServerConfig struct {
	Port int
}
type YdkServer struct {
	App *gofr.App
	Channel chan bool
}
func (s *YdkServer) Listen() (err error) {
	s.App.Run()
	return
}
func init() {
	logger.Log(logger.Info, "YdkServer initialized")
}
func NewServer(config YdkServerConfig) *YdkServer {
	app := gofr.New()
	if config.Port == 0 {
		config.Port = 3000
	}
	app.GET("/", func(c *gofr.Context) (interface{}, error) {
		return "Hello World!", nil
	})
	return &YdkServer{
		App: app,
	}
}
