package ydk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	logger "github.com/ywteam/ydk-go/logger"
)

type YdkSecOpsScannerVendor struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Url         string `json:"url"`
}
type YdkSecOpsScannerCapability struct {
	Formats []string `json:"formats"`
	Targets []string `json:"targets"`
}
type YdkSecOpsScannerMetadata struct {
	Tags         []string                   `json:"tags"`
	Features     []string                   `json:"features"`
	Vendor       YdkSecOpsScannerVendor     `json:"vendor"`
	Capabilities YdkSecOpsScannerCapability `json:"capabilities"`
}
type YdkSecOpsScannerPackages map[string]string
type YdkSecOpsScannerCli struct {
	Command string   `json:"cmd"`
	Args    []string `json:"args"`
	Version []string `json:"version"`
}
type YdkSecOpsScanner struct {
	Id       string                   `json:"id"`
	Env      map[string]string        `json:"env"`
	Metadata YdkSecOpsScannerMetadata `json:"metadata"`
	Packages YdkSecOpsScannerPackages `json:"packages"`
	Cli      YdkSecOpsScannerCli      `json:"cli"`
}
func (s *YdkSecOpsScanner) Entrypoint(args []string) (scanner *ScannerResult, err error) {
	logger.Log(logger.Info, fmt.Sprintf("Running scanner %s", s.Id))
	command := s.Cli.Command
	defaultArgs := append(s.Cli.Args, args...)
	for key, value := range s.Env {
		os.Setenv(key, value)
	}
	process := exec.Command(command, defaultArgs...)
	defer process.Wait()
	result := &ScannerResult{}
	result.Output = bytes.Buffer{}
	process.Stdout = io.MultiWriter(&result.Output)
	process.Stderr = io.MultiWriter(&result.Output)
	chann := make(chan int)
	defer close(chann)
	go func() {
		if err := process.Run(); err != nil {
			result.Err = err
		}
		result.ExitCode = process.ProcessState.ExitCode()
		result.Output = *bytes.NewBuffer([]byte(strings.Trim(result.Output.String(), "\n")))
		// result.Output = *bytes.NewBuffer([]byte(strings.ReplaceAll(result.Output.String(), "\n", "\\n")))
		chann <- result.ExitCode
	}()
	<-chann
	return result, result.Err
}
func (s *YdkSecOpsScanner) Version() (string, error) {
	command := s.Cli.Command
	args := s.Cli.Version
	process := exec.Command(command, args...)
	output, err := process.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

type ScannerResult struct {
	Scanner  YdkSecOpsScanner `json:"scanner"`
	Output   bytes.Buffer     `json:"output"`
	Err      error            `json:"error"`
	ExitCode int              `json:"exitCode"`
}

var scanners []YdkSecOpsScanner

func Scanners() []YdkSecOpsScanner {
	return scanners
}
func init() {
	data, err := os.ReadFile("/workspace/projects/ydk/src/go/assets/scanners.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &scanners); err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Loaded %d scanners\n", len(scanners))
	logger.Log(logger.Info, fmt.Sprintf("Loaded %d scanners", len(scanners)))
}
