package basic

import (
	"github.com/goplus/llcppg/cmd/gogensig/convert"
	"github.com/goplus/llcppg/cmd/gogensig/processor"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
	"github.com/goplus/llcppg/cmd/gogensig/visitor"
)

// For a default full convert processing,for main logic
type Config struct {
	convert.AstConvertConfig
	// PkgPreprocessor is called after package initialization but before creating FileSet
	PkgPreprocessor func(*convert.Package)
}

func ConvertProcesser(cfg *Config) (*processor.DocFileSetProcessor, *convert.Package, error) {
	astConvert, err := convert.NewAstConvert(&convert.AstConvertConfig{
		PkgName:   cfg.PkgName,
		SymbFile:  cfg.SymbFile,
		CfgFile:   cfg.CfgFile,
		OutputDir: cfg.OutputDir,
		PubFile:   cfg.PubFile,
	})
	if err != nil {
		return nil, nil, err
	}
	if cfg.PkgPreprocessor != nil {
		cfg.PkgPreprocessor(astConvert.Pkg)
	}
	docVisitors := []visitor.DocVisitor{astConvert}
	visitorManager := processor.NewDocVisitorManager(docVisitors)

	return processor.NewDocFileSetProcessor(&processor.ProcesserConfig{
		Exec: func(file *unmarshal.FileEntry) error {
			visitorManager.Visit(file.Doc, file.Path, file.IncPath, file.IsSys)
			return nil
		},
		DepIncs: astConvert.Pkg.AllDepIncs(),
		Done: func() {
			astConvert.WriteLinkFile()
			astConvert.WritePubFile()
		},
	}), astConvert.Pkg, nil
}
