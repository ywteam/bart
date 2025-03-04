package tools

import (
	"errors"

	"yellowteam/lib/shell"
)

type TrivyProgram struct {
	shell.ShellProgram
}

func (c *TrivyProgram) Name() string {
	return "trivy"
}

//	func (c *TrivyProgram) Command(args ...string) *shell.ShellCommand {
//		return shell.NewCommand(c.Name(), args...)
//	}
//
//	func (c *TrivyProgram) Exec(args ...string) *shell.ShellCommand {
//		return c.Command(args...).Run(nil)
//	}
//
//	func (c *TrivyProgram) Run(args ...string) *shell.ShellCommandResult {
//		return c.Exec(args...).Result
//	}
func (c *TrivyProgram) Install() error {
	if shell.Features["trivy"] {
		return nil
	}
	if shell.Features["bash"] {
		cmd := shell.Bash(`
			! uinstaller:install trivy && exit 1
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
		shell.Features["trivy"] = true
	} else if shell.Features["powershell"] {
		cmd := shell.Powershell(`
			PowershelUInstaller -packageNames @("trivy")
			if (-not (Get-Command trivy -ErrorAction SilentlyContinue)) {
				exit 1
			}
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
		shell.Features["trivy"] = true
	}
	return errors.New("unsupported shell")
}
func (c *TrivyProgram) Uninstall() error {
	if shell.Features["bash"] {
		cmd := shell.Bash(`
			! uinstaller:remove trivy && exit 1
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
	} else if shell.Features["powershell"] {
		cmd := shell.Powershell(`
			PowershelUInstaller-Remove -packageNames @("trivy")
			if (Get-Command trivy -ErrorAction SilentlyContinue) {
				exit 1
			}
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
	}
	return errors.New("unsupported shell")
}
func (c *TrivyProgram) Version() (string, error) {
	cmd := shell.Run("trivy", "--version")
	if cmd.Error != nil {
		return "", cmd.Error
	}
	return cmd.String(), nil
}
func init() {
	shell.Register(&TrivyProgram{})
}
