package shell_test

import (
	"testing"

	"yellowteam/lib/shell"

	"github.com/stretchr/testify/assert"
)

func TestDockerRun(t *testing.T) {
	d := shell.Docker{}
	opts := shell.DockerRunOptions{
		ContainerName: "test-container",
		Volumes:       [][]string{{"/host/path", "/container/path"}},
		Network:       []string{"test-network"},
		Ports:         [][]string{{"8080", "80"}},
		// Env:           []string{"ENV_VAR=value"},
		Env:         map[string]string{"ENV_VAR": "value"},
		Detach:      true,
		Interactive: true,
		Tty:         true,
		Remove:      true,
		Privileged:  true,
		Workdir:     "/workdir",
		User:        "user",
		Entrypoint:  "/entrypoint",
		Cmd:         "command",
		CmdArgs:     []string{"arg1", "arg2"},
	}
	cmd := d.Run("test-image", opts)
	assert.NotNil(t, cmd)
	assert.Contains(t, cmd.Args, "run")
	assert.Contains(t, cmd.Args, "--name")
	assert.Contains(t, cmd.Args, "test-container")
	assert.Contains(t, cmd.Args, "-v")
	assert.Contains(t, cmd.Args, "/host/path:/container/path")
	assert.Contains(t, cmd.Args, "--network")
	assert.Contains(t, cmd.Args, "test-network")
	assert.Contains(t, cmd.Args, "-p")
	assert.Contains(t, cmd.Args, "8080:80")
	assert.Contains(t, cmd.Args, "-e")
	assert.Contains(t, cmd.Args, "ENV_VAR=value")
	assert.Contains(t, cmd.Args, "-d")
	assert.Contains(t, cmd.Args, "-i")
	assert.Contains(t, cmd.Args, "-t")
	assert.Contains(t, cmd.Args, "--rm")
	assert.Contains(t, cmd.Args, "--privileged")
	assert.Contains(t, cmd.Args, "-w")
	assert.Contains(t, cmd.Args, "/workdir")
	assert.Contains(t, cmd.Args, "-u")
	assert.Contains(t, cmd.Args, "user")
	assert.Contains(t, cmd.Args, "--entrypoint")
	assert.Contains(t, cmd.Args, "/entrypoint")
	assert.Contains(t, cmd.Args, "test-image")
	assert.Contains(t, cmd.Args, "command")
	assert.Contains(t, cmd.Args, "arg1")
	assert.Contains(t, cmd.Args, "arg2")
}

func TestDockerExec(t *testing.T) {
	d := shell.Docker{}
	cmd := d.Exec("test-container", "ls", "-la")
	assert.NotNil(t, cmd)
	assert.Equal(t, []string{"exec", "test-container", "ls", "-la"}, cmd.Args)
}
