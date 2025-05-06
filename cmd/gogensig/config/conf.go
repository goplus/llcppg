package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	llcppg "github.com/goplus/llcppg/config"
)

// llcppg.cfg
func GetCppgCfgFromPath(filePath string) (*llcppg.Config, error) {
	bytes, err := ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	conf := llcppg.NewDefault()
	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// llcppg.pub
func GetPubFromPath(filePath string) (map[string]string, error) {
	return ReadPubFile(filePath)
}

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
	cmd := exec.Command("llcppsigfetch", append(args, "-ClangResourceDir="+ClangResourceDir())...)
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

func ReadPubFile(pubfile string) (ret map[string]string, err error) {
	b, err := os.ReadFile(pubfile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return
	}

	text := string(b)
	lines := strings.Split(text, "\n")
	ret = make(map[string]string, len(lines))
	for i, line := range lines {
		flds := strings.Fields(line)
		goName := ""
		switch len(flds) {
		case 1:
		case 2:
			goName = flds[1]
		case 0:
			continue
		default:
			err = fmt.Errorf("%s:%d: too many fields", pubfile, i+1)
			return
		}
		ret[flds[0]] = goName
	}
	return
}

func WritePubFile(file string, public map[string]string) (err error) {
	if len(public) == 0 {
		return
	}
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()
	ret := make([]string, 0, len(public))
	for name, goName := range public {
		if goName == "" {
			ret = append(ret, name)
		} else {
			ret = append(ret, name+" "+goName)
		}
	}
	sort.Strings(ret)
	_, err = f.WriteString(strings.Join(ret, "\n"))
	return
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

// temp to avoid call exec.Command in llcppsigfetch
func ClangResourceDir() string {
	res, err := exec.Command("clang", "-print-resource-dir").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(res))
}
