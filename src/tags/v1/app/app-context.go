package ydk
import (
	// "log"
	"os"
	// "fmt"
)

type YdkAppContextProcess struct {
	Args []string
	Env []string
	Stdin *os.File
	Stdout *os.File
	Stderr *os.File
	Dir string
	Path string
	Signal os.Signal
	ExitCode int
	Id int
	Name string
	Owner string
}