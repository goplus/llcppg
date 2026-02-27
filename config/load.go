package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// llcppg.cfg
func GetConfFromStdin() (conf Config, err error) {
	return ConfigFromReader(os.Stdin)
}

func GetConfFromFile(cfgFile string) (conf Config, err error) {
	fileReader, err := os.Open(cfgFile)
	if err != nil {
		return
	}
	defer fileReader.Close()

	return ConfigFromReader(fileReader)
}

func ConfigFromReader(reader io.Reader) (Config, error) {
	var config Config

	if err := json.NewDecoder(reader).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// llcppg.sym.json
func GetSymTableFromFile(symFile string) (symTable *SymTable, err error) {
	bytes, err := os.ReadFile(symFile)
	if err != nil {
		return nil, err
	}
	var syms []SymbolInfo
	err = json.Unmarshal(bytes, &syms)
	if err != nil {
		return nil, err
	}
	return NewSymTable(syms), nil
}

type SymTable struct {
	t map[string]SymbolInfo // mangle -> SymbolInfo
}

func NewSymTable(syms []SymbolInfo) *SymTable {
	symTable := &SymTable{
		t: make(map[string]SymbolInfo),
	}
	for _, sym := range syms {
		symTable.t[sym.Mangle] = sym
	}
	return symTable
}

func (t *SymTable) LookupSymbol(mangle string) (*SymbolInfo, error) {
	symbol, ok := t.t[mangle]
	if ok {
		return &symbol, nil
	}
	return nil, fmt.Errorf("symbol %s not found", mangle)
}

// llcppg.pub
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
	_, err = f.Write([]byte(strings.Join(ret, "\n")))
	return
}
