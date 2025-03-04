package shell_test

import (
	"encoding/json"
	"testing"
	"time"
	"yellowteam/lib/shell"
)

func Test_RunAsync(t *testing.T) {
	cmd := shell.Bash(`
		for i in {1..5}; do
			echo "Hello Async $i"
			sleep 1
		done
	`)
	cmd.RunAsync(nil)
	if !cmd.IsRunning() {
		t.Errorf("expected command to be running, got %v", cmd.IsRunning())
	}
	time.Sleep(6 * time.Second) // wait for the command to finish
}
func Test_RunSync(t *testing.T) {
	cmd := shell.Bash(`
		for i in {1..5}; do
			echo "Hello sync $i"
			sleep 1
		done
	`)
	cmd.Run(nil)
	if cmd.IsRunning() {
		t.Errorf("expected command to be finished, got %v", cmd.IsRunning())
	}
}
func Test_Metrics(t *testing.T) {
	cmd := shell.Bash("sleep 2")
	cmd.Run(nil)
	if cmd.Metrics.Duration() < 2 {
		t.Errorf("expected duration > 2, got %d", cmd.Metrics.Duration())
	}
}
func TestNewCommand(t *testing.T) {
	cmd := shell.NewCommand("echo", "hello")
	if cmd.Cmd != "echo" {
		t.Errorf("expected command 'echo', got %s", cmd.Cmd)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "hello" {
		t.Errorf("expected args ['hello'], got %v", cmd.Args)
	}
}

func TestCommand_Run(t *testing.T) {
	cmd := shell.NewCommand("echo", "hello")
	cmd.Run(nil)
	time.Sleep(1 * time.Second) // wait for the command to finish

	if len(cmd.Result.Stdout) != 1 || cmd.Result.Stdout[0] != "hello" {
		t.Errorf("expected stdout ['hello'], got %v", cmd.Result.Stdout)
	}
	if cmd.Result.Error != nil {
		t.Errorf("expected no error, got %v", cmd.Result.Error)
	}
}

// func TestCommand_Write(t *testing.T) {
// 	cmd := Bash("read input && echo $input")
// 	cmd.Run(false)
// 	err := cmd.Write("hello\n")
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 		return
// 	}
// 	time.Sleep(1 * time.Second) // wait for the command to process input

// 	if len(cmd.Result.Stdout) != 1 || cmd.Result.Stdout[0] != "hello" {
// 		t.Errorf("expected stdout ['hello'], got %v", cmd.Result.Stdout)
// 	}
// }

func TestExec(t *testing.T) {
	result := shell.Exec("echo", "hello")
	if len(result.Result.Stdout) != 1 || result.Result.Stdout[0] != "hello" {
		t.Errorf("expected stdout ['hello'], got %v", result.Result.Stdout)
	}
	if result.Result.Error != nil {
		t.Errorf("expected no error, got %v", result.Result.Error)
	}
}
func TestCoproc(t *testing.T) {
	cmd := shell.Bash(`
	# ! command -v bc && echo "bc not found" && exit 1
	# echo "running coproc A"
	# coproc:run a "while sleep 1; do echo 'waiting for input'; done"
	# echo "running coproc B"
	# coproc:run b "while sleep 1; do date; done"
	# echo "Show the output of coproc A and B"
	# // paste <(coproc:get-stdout-only "$a") <(coproc:get-stdout-only "$b")
	# # coproc:get-stdout-only "$a" 
	# # echo "A " < <(coproc:get-stdout-only "$a")
	# # echo "B >" "$(coproc:get-stdout-only "$b")"
	# echo "coprocs" "$a" "$b"
	# // coproc:get-stdin-fd "$a" stdin
	# echo "stdin" "$stdin"
	# coproc:stop "$a"
	# coproc:stop "$b"
	# echo "coprocs stopped"
	exit 0
	// coproc:get-stdout-fd "$b" stdout
	// coproc:get-stderr-fd "$b" stderr
	// echo "1+1" >&$stdin
	// read -u $stdout result
	// echo $result
`)
	cmd.Run(nil)
	// time.Sleep(1 * time.Second) // wait for the command to finish

	// if len(cmd.Result.Stdout) != 1 || cmd.Result.Stdout[0] != "2" {
	// 	t.Errorf("expected stdout ['2'], got %v", cmd.Result.Stdout)
	// }
	if cmd.Result.Error != nil {
		t.Errorf("expected no error, got %v", cmd.Result.Error)
	}
	t.Logf("Stdout: %s", cmd.Result.String())
}
func TestBash(t *testing.T) {
	cmd := shell.Bash("echo hello")
	cmd.Run(nil)
	time.Sleep(1 * time.Second) // wait for the command to finish

	if len(cmd.Result.Stdout) != 1 || cmd.Result.Stdout[0] != "hello" {
		t.Errorf("expected stdout ['hello'], got %v", cmd.Result.Stdout)
	}
	if cmd.Result.Error != nil {
		t.Errorf("expected no error, got %v", cmd.Result.Error)
	}
}
func TestResult_String(t *testing.T) {
	command := shell.Bash("echo hello; sleep 1")
	command.Run(nil)
	if command.Result.String() != "hello" {
		t.Errorf("expected 'hello', got %s", command.Result.String())
	}
	if command.Metrics.Duration() == 0 {
		t.Errorf("expected duration > 0, got %d", command.Metrics.Duration())
	}
	t.Logf("Command: %v", command)
	jsonData, err := json.MarshalIndent(command, "", "  ")
	if err != nil {
		t.Errorf("failed to marshal command to JSON: %v", err)
	}
	t.Logf("Json: %s", jsonData)
}
