// Yellow Open Source Scanners API
package ydk

import (
	"bytes"
	"encoding/json"
	flags "flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	// "reflect"
	console "github.com/ywteam/ydk-go/console"
	logger "github.com/ywteam/ydk-go/logger"
)

// Declare argv variable
// var argv = os.Args[1:]

type Input struct {
	Argv   string `json:"argv"`
	Env    string `json:"env"`
	Config string `json:"config"`
	Query  string `json:"query"`
	Prompt string `json:"prompt"`
	JSON   string `json:"json"`
}

func (input *Input) HttpBodyValue(flag string, request *http.Request) (string, error) {
	// Parse JSON body
	model := make(map[string]interface{})
	if err := json.NewDecoder(request.Body).Decode(&model); err != nil {
		logger.Default().Error("Error parsing JSON body: " + err.Error())
		if jsonValue, err := input.JsonValue(flag, model); err != nil {
			logger.Default().Error("Error parsing JSON value: " + err.Error())
		} else {
			return jsonValue, nil
		}
	}
	return "", fmt.Errorf("flag %s not found", flag)
}
func (input *Input) JsonValue(flag string, json map[string]interface{}) (string, error) {
	logger.Default().Debug("Parsing json input" + flag)
	if input.JSON != "" {
		if value, ok := json[flag]; ok {
			return value.(string), nil
		}
	}
	return "", fmt.Errorf("json %s not found", flag)
}
func (input *Input) QueryStringValue(flag string, query string) (string, error) {
	logger.Default().Debug("Parsing query input" + flag)
	if input.Query != "" {
		querstring := strings.Split(query, "&")
		for _, q := range querstring {
			raw := strings.Split(q, "=")
			if len(raw) == 2 {
				if raw[0] == flag {
					return raw[1], nil
				}
			}
		}
	}
	return "", fmt.Errorf("query %s not found", flag)
}
func (input *Input) PromptValue(flag string) (string, error) {
	logger.Default().Debug("Parsing prompt input" + flag)
	if input.Prompt != "" {
		if value, err := console.Prompt(input.Prompt); err != nil {
			if value != "" {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("prompt %s not found", flag)
}
func (input *Input) EnvValue(flag string) (string, error) {
	logger.Default().Debug("Parsing <env> input: " + flag)
	if input.Env != "" {
		if value, error := os.LookupEnv(input.Env); error {
			return value, nil
		}
	}
	return "", fmt.Errorf("env %s not found", flag)
}
func (input *Input) ArgvValue(flag string, description string) (string, error) {
	logger.Default().Debug("Parsing <argv> input: " + flag)
	if input.Argv != "" {
		flags.
		flag = strings.Replace(flag, "--", "", 1)
		argv := flags.String(flag, "", description)
		flags.Parse()
		if *argv != "" {
			return *argv, nil
		}
	}
	return "", fmt.Errorf("flag %s not found", flag)
}

// Parameter struct represents each parameter under "params"
type Parameter struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     string `json:"default"`
	Input       Input  `json:"input"`
}
func (param *Parameter) Raw(flag string) (value string, origin string, err error) {
	// logger.Default().Debug("Parsing value to: " + flag)
	if raw, err := param.Input.ArgvValue(flag, param.Description); err == nil {
		return raw, "cli", nil
	}
	if raw, err := param.Input.EnvValue(flag); err == nil {
		return raw, "env", nil
	}
	// if raw, err := param.Input.PromptValue(flag); err == nil {
	// 	return raw, "input", nil
	// }
	// if raw, err := param.Input.QueryStringValue(flag, "?test=value"); err == nil {
	// 	return raw, "querstring", nil
	// }
	// if raw, err := param.Input.HttpBodyValue(flag, &http.Request{}); err == nil {
	// 	return raw, "body", nil
	// }
	if param.Default != "" {
		return param.Default, "default", nil
	}
	if param.Required {
		return "", "required", fmt.Errorf("required %s not found", flag)
	}
	return "", "default", nil
}

// Params struct for the "params" object
// type ParamKey string
// type Params struct {
// 	Path Parameter `json:"path"`
// }

// Vendor struct for vendor details
type Vendor struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

// Capabilities struct lists what the scanner can do
type Capabilities struct {
	Formats  []string `json:"formats"`
	Scopes   []string `json:"scopes"`
	Targets  []string `json:"targets"`
	Features []string `json:"features"`
}

// Packages struct for package management systems
type Packages struct {
	Apt    []string `json:"apt"`
	Brew   []string `json:"brew"`
	Docker []string `json:"docker"`
	NPM    []string `json:"npm"`
	Shell  []string `json:"shell"`
	Choco  []string `json:"choco"`
}

// CLI struct represents the command line interface details
type CLI struct {
	Cmd     string   `json:"cmd"`
	Args    []string `json:"args"`
	Version []string `json:"version"`
}

// // ScannerParam defines individual CLI parameters
//
//	type ScannerParam struct {
//		Description string `json:"description"`
//		Type        string `json:"type"`
//		Required    bool   `json:"required"`
//		Default     string `json:"default"` // Use pointer to allow null values
//		Input       Input  `json:"input"`
//	}

// AapiEndpoint struct for the AapiEndpoint details of scanners
type AapiEndpoint struct {
	Help string               `json:"help"`
	Args map[string]Parameter `json:"args"`
}

func (ae *AapiEndpoint) HandleArgs() []string {
	logger.Default().Debug("Handling args from endpoint: " + ae.Help)
	params := make([]Parameter, 0, len(ae.Args))
	args := []string{}
	for flag := range ae.Args {
		param := ae.Args[flag]
		raw, origin, err := param.Raw(flag)
		if err != nil {
			logger.Default().Error("Error handing input: " + err.Error())
		} else {
			logger.Default().Debug(fmt.Sprintf("Raw: %s, Origin: %s", raw, origin))
		}
		logger.Default().Debug(fmt.Sprintf("%s = %s", flag, raw))
		params = append(params, param)
		// break
	}
	logger.Default().Debug(fmt.Sprintf("Args: %d", len(params)))
	return args
}

type ScannerResult struct {
	Output   bytes.Buffer
	Err      error
	ExitCode int
}

//	func (w *ScannerResult) Write(p []byte) (n int, err error) {
//		w.Chann <- *bytes.NewBuffer(p)
//		return len(p), nil
//	}
//
// Scanner struct for each scanner
type Scanner struct {
	ID           string                  `json:"id"`
	Vendor       Vendor                  `json:"vendor"`
	Hooks        []string                `json:"hooks"`
	Capabilities Capabilities            `json:"capabilities"`
	Packages     Packages                `json:"packages"`
	CLI          CLI                     `json:"cli"`
	Params       map[string]Parameter    `json:"params"`
	API          map[string]AapiEndpoint `json:"api"`
}

func (s *Scanner) Cli(args []string) (scanner *ScannerResult, err error) {
	logger.Default().Debug("Ossa Scanner CLI")
	command := s.CLI.Cmd
	defaultArgs := append(s.CLI.Args, args...)
	process := exec.Command(command, defaultArgs...)
	defer process.Wait()
	// var out bytes.Buffer
	// chann := make(chan bytes.Buffer)
	// defer close(chann)
	os.Setenv("LC_ALL", "C")
	result := &ScannerResult{}
	result.Output = bytes.Buffer{}
	process.Stdout = io.MultiWriter(&result.Output)
	process.Stderr = io.MultiWriter(&result.Output)
	chann := make(chan int)
	defer close(chann)
	go func() {
		if err := process.Run(); err != nil {
			// logger.Default().Error("Error running command. " + err.Error())
			result.Err = err
		}
		// else {
		// 	logger.Default().Info("Command ran successfully")
		// }
		result.ExitCode = process.ProcessState.ExitCode()
		// replace \n with <br> from result.Output
		result.Output = *bytes.NewBuffer([]byte(strings.ReplaceAll(result.Output.String(), "\n", "\\n")))
		chann <- result.ExitCode
	}()
	<-chann
	return result, result.Err
}
func (s *Scanner) Version() (scanner *ScannerResult, err error) {
	return s.Cli(s.CLI.Version)
}
func (s *Scanner) Endpoints() {
	logger.Default().Debug("Listing endpoints")
	// endpoints := make([]string, 0, len(s.API))
	headers := []string{"Endpoint", "Help", "* | Type", "Default", "Argv", "Env", "Query | Form | Path | JSON"}
	var data [][]string
	printArg := func(arg Parameter) []string {
		required := ""
		if arg.Required {
			required = "*"
		} else {
			required = "?"
		}
		return []string{
			fmt.Sprintf("%s | %s", required, arg.Type),
			arg.Default,
			arg.Input.Argv,
			arg.Input.Env,
			arg.Input.Query,
		}
	}
	for endpoint := range s.API {
		// endpointRef := s.API[endpoint]
		argsLen := len(s.API[endpoint].Args)
		logger.Default().Debug(fmt.Sprintf("Endpoint: %s, Args: %d", endpoint, argsLen))
		args := make([]string, 0, argsLen)
		for arg := range s.API[endpoint].Args {
			args = append(args, arg)
		}
		firstArg := []string{
			endpoint,
			s.API[endpoint].Help,
		}
		if argsLen > 0 {
			// firstArg = append(firstArg, args[0])
			firstArg = append(firstArg, printArg(s.API[endpoint].Args[args[0]])...)
		} else {
			firstArg = append(firstArg, "", "", "", "")
		}
		data = append(data, firstArg)
		if len(args) > 1 {
			for arg := range args {
				if arg == 0 {
					continue
				}
				data = append(data, append([]string{
					" \" ",
					" \" ",
					// args[arg],
				}, printArg(s.API[endpoint].Args[args[arg]])...))
			}
			// , printArg(s.API[endpoint].Args[args[arg]])
		}

	}
	console.PrintTable(headers, data)
}

type ScannerApiAdapter interface {
	HandleRequest() (err error)
}
type ScannerApiCliAdapter struct {
	Cmd *exec.Cmd
}

func (s *ScannerApiCliAdapter) HandleArgs() (err error) {
	return nil
}

type ScannerApiHttpAdapter struct {
	Req *http.Request
}

func (s *ScannerApiHttpAdapter) HandleArgs() (err error) {
	return nil
}

type ScannerApiArgs struct {
}
var schema OssaSchema
func privateHandleArgs(arguments map[string]Parameter) []string {
	returnValue := make([]string, 0, len(arguments))
	for arg := range arguments {
		parameter := arguments[arg]
		raw, origin, err := parameter.Raw(arg)
		if err != nil {
			logger.Default().Error("Error handing input: " + err.Error())
		} else {
			logger.Default().Debug(fmt.Sprintf("Raw: %s, Origin: %s", raw, origin))
		}
		logger.Default().Debug(fmt.Sprintf("%s = %s", arg, raw))
		returnValue = append(returnValue, arg)
	}
	return returnValue
}

// args []string) (scanner *ScannerResult, err error) {
func (s *Scanner) Api(path string, adapters []ScannerApiAdapter) (scanner *ScannerResult, err error) {
	if _, ok := s.API[path]; !ok {
		return nil, fmt.Errorf("endpoint %s does not exist", path)
	}
	endpoint := s.API[path]
	logger.Default().Debug(fmt.Sprintf("Endpoint: %s, Args: %d", path, len(endpoint.Args)))
	args := privateHandleArgs(schema.Params)
	// append(
	// 	privateHandleArgs(schema.Params),
	// 	privateHandleArgs(endpoint.Args)...,
	// )
	// args := endpoint.HandleArgs()
	for _, arg := range args {
		logger.Default().Debug(fmt.Sprintf("Arg: %s", arg))
	}
	// for _, adapter := range adapters {
	// 	t := reflect.TypeOf(adapter)
	// 	if t.Kind() == reflect.Ptr {
	//         t = t.Elem()
	//     }
	// 	logger.Default().Debug("Handling request with adapter " + t.Name())
	// 	err := adapter.HandleRequest()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// inputs := make([]string, 0, len(endpoint.Args))
	// for inputName := range endpoint.Args {
	// 	input := endpoint.Args[inputName]

	// 	// args = append(args, arg)
	// }

	return nil, nil
}

// OssaSchema struct to wrap the entire JSON structure
type OssaSchema struct {
	Params   map[string]Parameter `json:"params"`
	Scanners []Scanner            `json:"scanners"`
}



func FirstScanner() {
	logger.Default().Debug("First scanner")
	scanner := schema.Scanners[0]
	if version, err := scanner.Version(); err != nil {
		logger.Default().Error("Error getting version. " + err.Error())
	} else {
		logger.Default().Info(fmt.Sprintf("Version: %s", version.Output.String()))
	}
	scanner.Endpoints()
	logger.Default().Info("First scanner")
	scanner.Api("sast/sca/count-lines-of-code/v2", []ScannerApiAdapter{})
}

func Scanners() {
	logger.Default().Debug("Listing scanners")
	// schema.Scanners[0].Version()
	headers := []string{"ID", "Description", "Scopes", "Description", "Endpoints"}
	var data [][]string
	for _, scanner := range schema.Scanners {
		// get api keys
		keys := make([]string, 0, len(scanner.API))
		for k := range scanner.API {
			keys = append(keys, k)
		}
		data = append(data, []string{
			scanner.ID[:8],
			scanner.Vendor.Description,
			strings.Join(scanner.Capabilities.Scopes, ", "),
			scanner.Vendor.Description,
			func() string {
				if len(keys) > 0 {
					return keys[0]
				}
				return ""
			}(),
			// strings.Join(scanner.Capabilities.Features, ", "),
			// func() string {
			// 	keys := make([]string, 0, len(scanner.API))
			// 	for k := range scanner.API {
			// 		keys = append(keys, k)
			// 	}
			// 	return strings.Join(keys, " | ")
			// }(),
		})
		// for each key
		for key := range keys {
			if key == 0 {
				continue
			}
			apiKey := keys[key]
			data = append(data, []string{
				"-> ",
				"-> ",
				"-> ",
				"-> ",
				apiKey,
			})
		}
	}
	console.PrintTable(headers, data)
	logger.Default().Info("Scanners listed")
}

func init() {
	data, err := os.ReadFile("/workspace/sdk/go/src/ossa/scanners.json") // Assume JSON data is in data.json
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, &schema); err != nil {
		log.Fatal(err)
	}
}
