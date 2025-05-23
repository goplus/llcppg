package unmarshal

import (
	"encoding/json"
	"fmt"

	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llcppg/internal/unmarshal"
)

func Pkg(data []byte) (*llcppg.Pkg, error) {
	type pkgTemp struct {
		File    json.RawMessage
		FileMap map[string]*llcppg.FileInfo
	}
	var pkgData pkgTemp

	if err := json.Unmarshal(data, &pkgData); err != nil {
		return nil, fmt.Errorf("unmarshal error in Pkg into unmarshal.pkgTemp: %w", err)
	}
	file, err := unmarshal.File(pkgData.File)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error in Pkg when converting File of unmarshal.pkgTemp: %w", err)
	}
	return &llcppg.Pkg{
		File:    file,
		FileMap: pkgData.FileMap,
	}, nil
}
