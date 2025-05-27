package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ReadSigfetchFile(sigfetchFile string) ([]byte, error) {
	_, file := filepath.Split(sigfetchFile)
	var data []byte
	var err error
	if file == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(sigfetchFile)
	}
	return data, err
}

type SigfetchExtractConfig struct {
	File   string
	IsTemp bool
	IsCpp  bool
	Dir    string
}

func SigfetchExtract(cfg *SigfetchExtractConfig) ([]byte, error) {
	args := []string{"--extract", cfg.File}

	if cfg.IsTemp {
		args = append(args, "-temp=true")
	}

	if cfg.IsCpp {
		args = append(args, "-cpp=true")
	} else {
		args = append(args, "-cpp=false")
	}

	return executeSigfetch(args, cfg.Dir, cfg.IsCpp)
}

func SigfetchConfig(configFile string, dir string, isCpp bool) ([]byte, error) {
	args := []string{configFile}
	return executeSigfetch(args, dir, isCpp)
}

func executeSigfetch(args []string, dir string, isCpp bool) ([]byte, error) {
	cmd := exec.Command("llcppsigfetch", args...)
	if dir != "" {
		cmd.Dir = dir
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error running llcppsigfetch: %v\nStderr: %s\nArgs: %s", err, stderr.String(), strings.Join(args, " "))
	}

	return out.Bytes(), nil
}

func ReadFile(filePath string) ([]byte, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	return io.ReadAll(jsonFile)
}

func RunCommand(dir, cmdName string, args ...string) error {
	execCmd := exec.Command(cmdName, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Dir = dir
	return execCmd.Run()
}

func CreateTmpJSONFile(filename string, data any) (string, error) {
	filePath := filepath.Join(os.TempDir(), filename)
	err := CreateJSONFile(filePath, data)
	return filePath, err
}

func CreateJSONFile(filepath string, data any) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
