package shell

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
	"yellowteam/lib/config"
	"yellowteam/lib/logger"
)

type ShellCommandResult struct {
	Stdout []string
	Error  error
}

func (c *ShellCommandResult) String() string {
	return strings.Join(c.Stdout, "\n")
}
func (c *ShellCommandResult) CheckError() *ShellCommandResult {
	if c.Error != nil {
		panic(c.Error)
	}
	return c
}

type CommandMetric struct {
	StartAt int64
	EndAt   int64
}

func (c *CommandMetric) Duration() int64 {
	return c.EndAt - c.StartAt
}

type ShellCommand struct {
	Cmd  string    `json:"cmd"`
	Args []string  `json:"args"`
	Exec *exec.Cmd `json:"-"`
	// Process *os.Process
	// State   *os.ProcessState `json:"-"`
	// Chan    chan string         `json:"-"`
	Result  *ShellCommandResult `json:"result"`
	Metrics *CommandMetric      `json:"metrics"`
}

func (c *ShellCommand) IsRunning() bool {
	return c.Exec.ProcessState == nil || !c.Exec.ProcessState.Exited()
}
func (c *ShellCommand) Write(message string) (err error) {
	stdin := c.Exec.Stdin.(io.WriteCloser)
	_, err = stdin.Write([]byte(message))
	return err
}
func (c *ShellCommand) RunAsync(watcher func(string)) *ShellCommand {
	go c.Run(watcher)
	return c
}
func (c *ShellCommand) Run(watcher func(string)) *ShellCommand {
	if watcher == nil {
		watcher = func(text string) {
			if text != "" {
				logger.Trace("%s stdout: %s", c.Cmd, strings.Replace(text, "\n", "\\n", -1))
			}
		}
	}
	if c.Exec == nil {
		c.Exec = exec.Command(c.Cmd, c.Args...)
	}
	c.Metrics.StartAt = time.Now().Unix()
	stdout, err := c.Exec.StdoutPipe()
	if err != nil {
		c.Result.Error = err
		return c
	}
	if err := c.Exec.Start(); err != nil {
		c.Result.Error = err
		return c
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		text := scanner.Text()
		c.Result.Stdout = append(c.Result.Stdout, text)
		watcher(text)
	}
	if err := scanner.Err(); err != nil {
		c.Result.Error = err
	}
	c.Result.Error = c.Exec.Wait()
	c.Metrics.EndAt = time.Now().Unix()
	if c.Exec.ProcessState != nil {
		if c.Exec.ProcessState.Success() {
			c.Result.Error = nil
		} else if c.Exec.ProcessState.Exited() {
			c.Result.Error = errors.New("process exited")
		} else {
			c.Result.Error = errors.New("process killed")
		}
	}
	return c
	// done := make(chan string)
	// go func() {
	// 	defer close(done)
	// 	scanner := bufio.NewScanner(stdout)
	// 	for scanner.Scan() {
	// 		text := scanner.Text()
	// 		c.Result.Stdout = append(c.Result.Stdout, text)
	// 		done <- text
	// 	}
	// 	if err := scanner.Err(); err != nil {
	// 		c.Result.Error = err
	// 	}
	// 	c.Result.Error = c.Exec.Wait()
	// 	c.Metrics.EndAt = time.Now().Unix()
	// }()
	// for text := range done {
	// 	watcher(text)
	// }
	// c.Metrics.EndAt = time.Now().Unix()
	// return c
}
func NewCommand(cmd string, args ...string) *ShellCommand {
	return &ShellCommand{
		Cmd:  cmd,
		Args: args,
		Exec: exec.Command(cmd, args...),
		Result: &ShellCommandResult{
			Stdout: make([]string, 0),
		},
		Metrics: &CommandMetric{
			StartAt: time.Now().Unix(),
		},
	}
}
func Bash(scripts ...string) *ShellCommand {
	const source = "/tmp/opsh.sh"
	os.Remove(source)
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		sourceFile, err := os.Create("/tmp/opsh.sh")
		if err != nil {
			panic(err)
		}
		defer sourceFile.Close()
		_, err = sourceFile.WriteString(strings.Join([]string{
			"#!/usr/bin/env bash",
			"# Created by ydk-go",
			"# Create At " + time.Now().String(),
			"set -e -o pipefail",
			BashInit(),
			// "rap:sdk --opsh-set:logger.level=info --opsh-set:logger.template=\"%{icon} %{message}\"",
		}, "\n"))
		if err != nil {
			panic(err)
		}
	}
	scripts = append(scripts, "exit $?")
	scripts = append([]string{"source " + source}, scripts...)
	// logger.Trace("Running bash scripts: \n%s", strings.Join(scripts, "\n"))
	return NewCommand("bash", "-c", strings.Join(scripts, "\n"))
	// return NewCommand("bash", "-c", strings.Join(scripts, "\n"))
}
func Powershell(scripts ...string) *ShellCommand {
	const source = "/tmp/opsh.ps1"
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		sourceFile, err := os.Create("/tmp/opsh.ps1")
		if err != nil {
			panic(err)
		}
		defer sourceFile.Close()
		_, err = sourceFile.WriteString(strings.Join([]string{
			"#!/usr/bin/env pwsh",
			"# Created by ydk-go",
			"# Create At " + time.Now().String(),
			"Set-StrictMode -Version Latest",
			PowershellInit(),
		}, "\n"))
		if err != nil {
			panic(err)
		}
	}
	scripts = append(scripts, "exit $?")
	scripts = append([]string{"source " + source}, scripts...)
	return NewCommand("pwsh", "-c", strings.Join(scripts, "\n"))
}
func Script(scripts ...string) *ShellCommand {
	logger.Trace("Running scripts: \n%s", strings.Join(scripts, "\n"))
	if Features["bash"] {
		return Bash(scripts...)
	}
	if Features["pwsh"] {
		return Powershell(scripts...)
	}
	return NewCommand("echo", "unsupported shell")
}
func Exec(cmd string, args ...string) *ShellCommand {
	return Script(
		strings.Join(append([]string{cmd}, QuoteArgs(args...)...), " "),
	).Run(nil)
	// return Script(cmd + " " + strings.Join(QuoteArgs(args), " ")).Run(nil)
}
func Run(cmd string, args ...string) *ShellCommandResult {
	return Exec(cmd, args...).Result
	// return Exec(cmd + " " + strings.Join(QuoteArgs(args), " ")).Result
}
func Witch(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
func QuoteArgs(args ...string) []string {
	if args == nil {
		return make([]string, 0)
	}
	quoted := make([]string, len(args))
	for i, arg := range args {
		if arg == "" {
			quoted[i] = "''"
		} else if strings.Contains(arg, " ") {
			quoted[i] = "\"" + arg + "\""
		} else {
			quoted[i] = arg
		}
	}
	return quoted
}

func Register(program IShellProgram) {
	name := program.Name()
	Features[name] = Witch(name)
	Programs[name] = program
}

var Programs = map[string]IShellProgram{}
var Features = map[string]bool{
	"bash":   Witch("bash"),
	"pwsh":   Witch("pwsh"),
	"docker": Witch("docker"),
}

func init() {
	config.WithEnv("SHELL", false, "bash")
	config.WithEnv("SHELL_AUTOINTALL", false, "true")

}
