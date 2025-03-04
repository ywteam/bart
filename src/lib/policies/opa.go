// represents a POC to use OPA Policy
// 1 - Use docker to for all
// 2 - Start OPA server
// 3 - Write to policies one to allow and another to deny
// 4 - Write a go code to use OPA client to check the policy
// 5 - Write a test to check the policy
// 6 - Write a k6 script to ckeck the policies

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type OpaService struct {
	Url string
}
func (s *OpaService) Bundle() *CommandResult {
	return Bash(`
		docker run --rm -v $(pwd):/data openpolicyagent/opa:latest build -t rego /data/policy.rego
	`)
}
func (s *OpaService) ZeroTrustPolicy() string {
	return `
package tests.zerotrust

import rego.v1

# Only owner can update the pet's information
# Ownership information is provided as part of OPA's input
default allow := false

allow if {
	input.method == "PUT"
	some petid
	input.path = ["pets", petid]
	input.user == input.owner
}

	`
}
func (s *OpaService) PublicPolicy() string {
	return `
package tests.public

import rego.v1

# Only owner can update the pet's information
# Ownership information is provided as part of OPA's input
default allow := false

allow if {
	input.method == "PUT"
	some petid
	input.path = ["pets", petid]
	input.user == input.owner
}

	`
}
func (s *OpaService) CreatePolicies() *CommandResult {
	return Bash(`
		curl -X PUT -H "Content-Type: text/plain" --data-binary @- http://localhost:8181/v1/policies/tests/authz/zero-trust <<EOF
		%s
		EOF
		curl -X PUT -H "Content-Type: text/plain" --data-binary @- http://localhost:8181/v1/policies/tests/authz/public <<EOF
		%s
		EOF
	`, s.ZeroTrustPolicy(), s.PublicPolicy())
}

func (s *OpaService) Start() *CommandResult {
	return Bash(`
		# docker run --name opa -d -p 8181:8181 openpolicyagent/opa run --server --addr :8181

		echo "Starting OPA server"
	`).Watch(func(line string) {
		fmt.Println(line)
	})
}
func (s *OpaService) Stop() *CommandResult {
	return Bash(`
		docker stop opa
		docker rm opa
		docker image rm openpolicyagent/opa:latest
	`)
}

type CommandResult struct {
	Stdout []string
	Error  error
	State  *os.ProcessState `json:"-"`
	Chan   chan string      `json:"-"`
}

func (c *CommandResult) String() string {
	return strings.Join(c.Stdout, "\n")
}
func (c *CommandResult) HasError() bool {
	return c.Error != nil
}
func (c *CommandResult) Watch(watcher func(string)) *CommandResult {
	for line := range c.Chan {
		watcher(line)
	}
	return c
}
func (c *CommandResult) Wait() *CommandResult{
	<-c.Chan
	return c
}
func Bash(script string, cleanup ...string) *CommandResult {
	fmt.Printf("Running bash script: %s\n", script)	
	return Command("bash", "-c", fmt.Sprintf(`
		set -e -o pipefail
		function cleanup() {
			%s
			return 0
		}
		trap cleanup EXIT
		function task(){
			%s
			return $?
		}
		task
		[[ $? -ne 0 ]] && echo "command ${BASH_COMMAND} failed with code $?" && exit 255
		exit 0
	`, cleanup, script))
}
func Command(command string, args ...string) *CommandResult {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return &CommandResult{Error: err}
	}
	if err := cmd.Start(); err != nil {
		return &CommandResult{Error: err}
	}
	result := &CommandResult{
		Chan:   make(chan string),
		Stdout: make([]string, 0),
	}
	go func() {
		defer close(result.Chan)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			text := scanner.Text()
			result.Stdout = append(result.Stdout, text)
			result.Chan <- text
		}
		if err := scanner.Err(); err != nil {
			result.Error = err
		}
		if err := cmd.Wait(); err != nil {
			result.Error = err
		}
		// close(result.Chan)
	}()
	return result
}

func main() {
	opa := &OpaService{}
	fmt.Println(opa.Start())
	fmt.Println(opa.Stop())
}
