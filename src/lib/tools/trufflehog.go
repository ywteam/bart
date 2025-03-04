package tools

import (
	"errors"

	"yellowteam/lib/shell"
)

type TrufflehogProgram struct {
	shell.ShellProgram
}

func (c *TrufflehogProgram) Name() string {
	return "trufflehog"
}

//	func (c *TrufflehogProgram) Command(args ...string) *shell.ShellCommand {
//		return shell.NewCommand(c.Name(), args...)
//	}
//
//	func (c *TrufflehogProgram) Exec(args ...string) *shell.ShellCommand {
//		return c.Command(args...).Run(nil)
//	}
//
//	func (c *TrufflehogProgram) Run(args ...string) *shell.ShellCommandResult {
//		return c.Exec(args...).Result
//	}
func (c *TrufflehogProgram) Install() error {
	if shell.Features["trufflehog"] {
		return nil
	}
	if shell.Features["bash"] {
		cmd := shell.Bash(`
			! uinstaller:fetch https://raw.githubusercontent.com/trufflesecurity/trufflehog/main/scripts/install.sh | sh -s -- -b /usr/local/bin
			if [ ! -x "$(command -v trufflehog)" ]; then
				exit 1
			fi
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
		shell.Features["trufflehog"] = true
	} else if shell.Features["powershell"] {
		cmd := shell.Powershell(`
			PowershelUInstaller -packageNames @("trufflehog")
			if (-not (Get-Command trufflehog -ErrorAction SilentlyContinue)) {
				exit 1
			}
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
		shell.Features["trufflehog"] = true
	}
	return errors.New("unsupported shell")
}
func (c *TrufflehogProgram) Uninstall() error {
	if shell.Features["bash"] {
		cmd := shell.Bash(`
			! uinstaller:remove trufflehog && exit 1
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
	} else if shell.Features["powershell"] {
		cmd := shell.Powershell(`
			PowershelUInstaller-Remove -packageNames @("trufflehog")
			if (Get-Command trufflehog -ErrorAction SilentlyContinue) {
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
func (c *TrufflehogProgram) Version() (string, error) {
	cmd := shell.Run("trufflehog", "--version")
	if cmd.Error != nil {
		return "", cmd.Error
	}
	return cmd.String(), nil
}
func init() {
	shell.Register(&TrufflehogProgram{})
}
