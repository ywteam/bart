package ydk

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"math/rand"
	config "github.com/ywteam/ydk-go/config"
	console "github.com/ywteam/ydk-go/console"
	emojis "github.com/ywteam/ydk-go/emojis"
	logger "github.com/ywteam/ydk-go/logger"
	spinners "github.com/ywteam/ydk-go/spinners"
	shell "github.com/ywteam/ydk-go/ydk-shell"
	// ossa "github.com/ywteam/ydk-go/ossa"
)

func TestYdkShell() {
	fmt.Println("Running command with args " + os.Args[1])
	if len(os.Args) < 2 {
		fmt.Println("Usage: ", os.Args[0], "host")
		os.Exit(1)
	}
	command := os.Args[1]
	args := os.Args[2:]
	cmd := shell.YDkCommand{Command: command, Args: args}
	cmd.Run()
}
func TestConfig() {
	ydkConfig := config.YDkConfig{}
	path, _ := os.Getwd()
	dotenvs := []string{
		filepath.Join(path, "..", ".env"),
		filepath.Join(path, "..", ".env.example"),
	}
	ydkConfig.DotEnv(dotenvs...)
	fmt.Println(ydkConfig.Host)
}
func TestConsole() {
	terminal := console.YdkConsole{}
	terminal.Log("This is a log message")
	terminal.Print("This is a print message")
	terminal.Error("This is an error message")
	// console.Fatal("This is a fatal message")
	// console.Exit(1)
	// console.Clear()
	terminal.Title("This is a title")
	fmt.Println(terminal.Colorize("YDK Go SDK", console.Yellow, console.Reverse))
	// for key, value := range emojis {
	// 	fmt.Println(key, value)
	// }
}
func TestLogger() *logger.YdkLogger {
	logger := logger.Default()
	// os.Setenv("YDK_ENV", "production")
	// logger.State.Env = "production"
	logger.SetContext("GO")
	// logger.Log(ydkLogger.Trace, "This is a trace message")
	logger.Trace("This is a trace message")
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Notice("This is a notice message")
	logger.Output("This is an output message")
	logger.Success("This is a success message")
	logger.Warning("This is a warning message")
	logger.Alert("This is an alert message")
	logger.Error("This is an error message")
	logger.Critical("This is a critical message")
	logger.Emergency("This is an emergency message")
	logger.Panic("This is a panic message")
	logger.Fatal("This is a fatal message")
	// logger.Log(ydkLogger.Debug, "This is a debug message")
	// logger.Log(ydkLogger.Info, "This is an info message")
	// logger.Log(ydkLogger.Notice, "This is a notice message")
	// logger.Log(ydkLogger.Warning, "This is a warning message")
	// logger.Log(ydkLogger.Alert, "This is an alert message")
	// logger.Log(ydkLogger.Error, "This is an error message")
	// logger.Log(ydkLogger.Critical, "This is a critical message")
	// logger.Log(ydkLogger.Emergency, "This is an emergency message")
	// logger.Log(ydkLogger.Panic, "This is a panic message")
	// logger.Log(ydkLogger.Fatal, "This is a fatal message")
	// logger.Log(ydkLogger.Debug, "This is a debug message")

	return logger

}
func TestEmojis() {
	key, value := emojis.RandomEmoji()
	fmt.Println(key, value)
	// list := emojis.ListEmojis()
	// for key, value := range list {
	// 	fmt.Println(key, value)
	// }
	text := "Hello, :shorts:!"
	fmt.Println(emojis.InterpolateEmojis(text))

}
func TestSpinners(){
	spinnersList := spinners.ListSpinners()
	for key := range spinnersList {
		text := fmt.Sprintf("Operation in progress with spinner %s {{.Log.Level}}", key)		
		rand.Seed(time.Now().UnixNano())
		randomTime := 1 + time.Duration(rand.Intn(3)) - 1
		spinners.Start(key, text, func(message *string) bool {
			// logger.Default().Info("Running operation")
			time.Sleep(randomTime * time.Second)
			*message = text + ". Still working, please wait..." // Fix: Assign the string value to the pointer variable using the dereference operator (*message)
			time.Sleep(randomTime * time.Second)
			// logger.Default().Success("Operation completed")
			spinners.Stop()
			// return random boolean
			return randomTime%2 == 0
		})
	}
	// spinners.Start("random", "Time travel in progress", func(message *string) bool {
	// 	logger.Default().Info("Running time travel operation")
	// 	time.Sleep(5 * time.Second)
	// 	*message = "Time travel in progress, please wait..." // Fix: Assign the string value to the pointer variable using the dereference operator (*message)
	// 	time.Sleep(30 * time.Hour)
	// 	logger.Default().Success("Time travel operation completed")
	// 	return true
	// })	
}
func TestOssa() {
	// ossa.Scanners()
	// ossa.FirstScanner()
	return
}
func main() {
	// TestYdkShell()
	// TestConfig()
	// TestConsole()
	// TestLogger()
	// TestEmojis()
	TestSpinners()
	// TestOssa()
}
