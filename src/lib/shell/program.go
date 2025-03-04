package shell

import "errors"

type IShellProgram interface {
	Name() string
	Command(args ...string) *ShellCommand
	Exec(args ...string) *ShellCommand
	Run(args ...string) *ShellCommandResult
	Install() error
	Uninstall() error
	Version() (string, error)
}

type ShellProgram struct {}
func (c *ShellProgram) Name() string {
	return ""
}
func (c *ShellProgram) Command(args ...string) *ShellCommand {
	return NewCommand(c.Name(), args...)
}
func (c *ShellProgram) Exec(args ...string) *ShellCommand {
	return c.Command(args...).Run(nil)
}
func (c *ShellProgram) Run(args ...string) *ShellCommandResult {
	return c.Exec(args...).Result
}
func (c *ShellProgram) Install() error {
	return errors.New("unsupported shell")
}
func (c *ShellProgram) Uninstall() error {
	return errors.New("unsupported shell")
}
func (c *ShellProgram) Version() (string, error) {
	return "", errors.New("unsupported shell")
}
