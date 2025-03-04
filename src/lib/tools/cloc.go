package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	reports "yellowteam/lib/report"
	"yellowteam/lib/shell"
)

type ClocProgram struct {
	shell.ShellProgram
}

func (c ClocProgram) Name() string {
	return "cloc"
}

//	func (c *ClocProgram) Command(args ...string) *shell.ShellCommand {
//		return shell.NewCommand(c.Name(), args...)
//	}
//
//	func (c *ClocProgram) Exec(args ...string) *shell.ShellCommand {
//		return c.Command(args...).Run(nil)
//	}
//
//	func (c *ClocProgram) Run(args ...string) *shell.ShellCommandResult {
//		return c.Exec(args...).Result
//	}
func (c *ClocProgram) Install() error {
	if shell.Features["cloc"] {
		return nil
	}
	if shell.Features["bash"] {
		cmd := shell.Bash(`
			! uinstaller:install cloc && exit 1
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
		shell.Features["cloc"] = true
	} else if shell.Features["powershell"] {
		cmd := shell.Powershell(`
			PowershelUInstaller -packageNames @("cloc")
			if (-not (Get-Command cloc -ErrorAction SilentlyContinue)) {
				exit 1
			}
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
		shell.Features["cloc"] = true
	}
	return errors.New("unsupported shell")
}
func (c *ClocProgram) Uninstall() error {
	if shell.Features["bash"] {
		cmd := shell.Bash(`
			! uinstaller:remove cloc && exit 1
		`)
		cmd.Run(nil)
		if cmd.Result.Error != nil {
			return cmd.Result.Error
		}
	} else if shell.Features["powershell"] {
		cmd := shell.Powershell(`
			PowershelUInstaller-Remove -packageNames @("cloc")
			if (Get-Command cloc -ErrorAction SilentlyContinue) {
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
func (c *ClocProgram) Version() (string, error) {
	cmd := shell.Run("cloc", "--version")
	if cmd.Error != nil {
		return "", cmd.Error
	}
	return cmd.String(), nil
}
func (c *ClocProgram) CountLinesOfCode(target string) *ClocResult {
	if fileInfo, err := os.Stat(target); err != nil {
		return &ClocResult{Error: err.Error()}
	} else if !fileInfo.IsDir() {
		return &ClocResult{Error: "target is not a directory"}
	}
	resultFile, err := os.CreateTemp("", "rap-*.json")
	if err != nil {
		return &ClocResult{Error: err.Error()}
	}
	cmd := shell.Bash(`
		if ! command -v cloc &> /dev/null; then
			! uinstaller:install cloc && exit 1
			! command -v cloc &> /dev/null && exit 2
		fi
		target="` + target + `"
		resultFile="` + resultFile.Name() + `"
		if [ -f "$target/.gitignore" ]; then
			cloc --json --quiet --exclude-list-file="$target/.gitignore" "$target" > "$resultFile"
		else
			cloc --json --quiet --exclude-dir=node_modules,vendor,tests,examples,docs,build,bin,assets,assets "$target" > "$resultFile"
		fi
		[[ $? -ne 0 ]] && rm -f "$resultFile" && exit 3
		echo "$resultFile"
	`)
	cmd.Run(nil)
	if cmd.Result.Error != nil {
		return &ClocResult{Error: cmd.Result.Error.Error()}
	}
	result := NewClocResult(cmd)
	if result.Error != "" {
		return result
	}
	os.Remove(resultFile.Name())
	return result
}

type Language struct {
	NFiles  int `json:"nFiles"`
	Blank   int `json:"blank"`
	Comment int `json:"comment"`
	Code    int `json:"code"`
}
type ClocResult struct {
	Error   string `json:"error,omitempty"`
	Metrics struct {
		ClocUrl        string  `json:"-"`
		Version        string  `json:"-"`
		ElapsedSeconds float64 `json:"elapsed_seconds"`
		NFiles         int     `json:"-"`
		NLines         int     `json:"-"`
		FilesPerSecond float64 `json:"files_per_second"`
		LinesPerSecond float64 `json:"lines_per_second"`
	} `json:"header"`
	Sum struct {
		Blank   int `json:"blank"`
		Comment int `json:"comment"`
		Code    int `json:"code"`
		NFiles  int `json:"nFiles"`
	} `json:"totals"`
	Languages map[string]Language `json:"laguages,inline,omitempty"`
}

func (c *ClocResult) Json(indent bool) string {
	if indent {
		raw, _ := json.MarshalIndent(c, "", "  ")
		return string(raw)
	}
	raw, _ := json.Marshal(c)
	return string(raw)
}
func (c *ClocResult) Sarif(report *reports.Report) {
	run := report.NewRun("cloc", "")
	for lang, value := range c.Languages {
		run.PropertyBag.AddInteger(fmt.Sprintf("cloc.%s.nFiles", lang), value.NFiles)
		run.PropertyBag.AddInteger(fmt.Sprintf("cloc.%s.blank", lang), value.Blank)
		run.PropertyBag.AddInteger(fmt.Sprintf("cloc.%s.comment", lang), value.Comment)
		run.PropertyBag.AddInteger(fmt.Sprintf("cloc.%s.code", lang), value.Code)
	}
}
func NewClocResult(arg interface{}) *ClocResult {
	switch v := arg.(type) {
	case map[string]json.RawMessage:
		var result = ClocResult{Languages: make(map[string]Language)}
		for key, value := range v {
			if key == "header" {
				if err := json.Unmarshal(value, &result.Metrics); err != nil {
					result.Error = err.Error()
					return &result
				}
			} else if key == "SUM" {
				if err := json.Unmarshal(value, &result.Sum); err != nil {
					result.Error = err.Error()
					return &result
				}
			} else {
				var lang Language
				if err := json.Unmarshal(value, &lang); err != nil {
					result.Error = err.Error()
					return &result
				}
				result.Languages[key] = lang
			}
		}
		return &result
	case string:
		// if is file path
		if fileInfo, err := os.Stat(v); err == nil {
			if fileInfo.IsDir() {
				return &ClocResult{Error: "target is a directory"}
			}
			file, err := os.Open(v)
			if err != nil {
				return &ClocResult{Error: err.Error()}
			}
			defer file.Close()
			var raw map[string]json.RawMessage
			if err := json.NewDecoder(file).Decode(&raw); err != nil {
				return &ClocResult{Error: err.Error()}
			}
			return NewClocResult(raw)
		}
		// if is json
		if json.Valid([]byte(v)) {
			var raw map[string]json.RawMessage
			if err := json.Unmarshal([]byte(v), &raw); err != nil {
				return &ClocResult{Error: err.Error()}
			}
			return NewClocResult(raw)
		}
	case *shell.ShellCommand:
		return NewClocResult(v.Result.String())
	case *os.File:
		file, err := os.Open(v.Name())
		if err != nil {
			return &ClocResult{Error: err.Error()}
		}
		defer file.Close()
		var raw map[string]json.RawMessage
		if err := json.NewDecoder(file).Decode(&raw); err != nil {
			return &ClocResult{Error: err.Error()}
		}
		return NewClocResult(raw)
	}
	return &ClocResult{Error: "unsupported argument type"}
}

func init() {
	shell.Register(&ClocProgram{})
}
