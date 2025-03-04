package ydk

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

type YDkCommand struct {
	Command string
	Args    []string
}

func (c *YDkCommand) Run() {
	// Add a spinner to the command
	commandWithSpinner := `
	(sleep 1 && printf '.' && sleep 3 && printf '.' && sleep 5 && printf '.' && printf '\n') &
	` + c.Command
	// log.Printf("Running command: %s\n", commandWithSpinner)
	cmd := exec.Command("/bin/bash", "-c", commandWithSpinner)

	var out bytes.Buffer
	cmd.Stdout = io.MultiWriter(&out, os.Stdout)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	log.Printf("combined out:\n%s\n", out.String())
}