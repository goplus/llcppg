package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

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
