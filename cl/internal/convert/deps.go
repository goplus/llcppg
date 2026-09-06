package convert

import (
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/goplus/gogen"
	llcppg "github.com/xgo-dev/llcppg/config"
)

type PkgDepLoader struct {
	root     string
	pkg      *gogen.Package
	pkgCache map[string]*PkgInfo // pkgPath -> *PkgInfo
	pkgs     map[string]string   // pkgPath -> pkgDir
	regCache map[string]struct{} // pkgPath
}

func NewPkgDepLoader(root string, pkg *gogen.Package, deps []string) (*PkgDepLoader, error) {
	ret := &PkgDepLoader{
		root:     root,
		pkg:      pkg,
		pkgCache: make(map[string]*PkgInfo),
		regCache: make(map[string]struct{}),
	}
	if err := ret.loadPkgDirs(deps); err != nil {
		return nil, err
	}
	return ret, nil
}

// for current package & dependent packages
type PkgInfo struct {
	PkgBase
	Deps []*PkgInfo
}

type PkgBase struct {
	PkgPath string            // package path, e.g. github.com/goplus/llgo/cjson
	Deps    []string          // dependent packages
	Pubs    map[string]string // llcppg.pub
}

func NewPkgInfo(pkgPath string, deps []string, pubs map[string]string) *PkgInfo {
	if pubs == nil {
		pubs = make(map[string]string)
	}
	return &PkgInfo{
		PkgBase: PkgBase{PkgPath: pkgPath, Deps: deps, Pubs: pubs},
	}
}

// LoadDeps loads direct dependencies of the current package and recursively loads their
// dependencies, to get the complete dependency.
func (pm *PkgDepLoader) LoadDeps(p *PkgInfo) ([]*PkgInfo, error) {
	deps, err := pm.Imports(p.PkgBase.Deps)
	if err != nil {
		return nil, err
	}
	return deps, nil
}

func (pm *PkgDepLoader) Imports(pkgPaths []string) (pkgs []*PkgInfo, err error) {
	pkgs = make([]*PkgInfo, len(pkgPaths))
	for i, pkgPath := range pkgPaths {
		pkgs[i], err = pm.Import(pkgPath)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (pm *PkgDepLoader) Import(pkgPath string) (*PkgInfo, error) {
	// standard C library paths
	pkgPath, isStd := IsDepStd(pkgPath)
	pkgPath, _ = splitPkgPath(pkgPath)

	if pkg, exist := pm.pkgCache[pkgPath]; exist {
		return pkg, nil
	}

	pkgDir := pm.pkgs[pkgPath]
	if pkgDir == "" {
		return nil, fmt.Errorf("%w: go list cache has no dir for package %q", llcppg.ErrConfig, pkgPath)
	}
	pkgDir, err := filepath.Abs(pkgDir)
	if err != nil {
		return nil, err
	}

	pubs, err := llcppg.ReadPubFile(filepath.Join(pkgDir, llcppg.LLCPPG_PUB))
	if err != nil {
		return nil, err
	}

	var conf llcppg.Config
	var deps []string
	if !isStd {
		conf, err = llcppg.GetConfFromFile(filepath.Join(pkgDir, llcppg.LLCPPG_CFG))
		if err != nil {
			return nil, err
		}
		deps = conf.Deps
	}

	newPkg := NewPkgInfo(pkgPath, deps, pubs)
	pm.pkgCache[pkgPath] = newPkg

	if len(conf.Deps) > 0 {
		deps, err := pm.LoadDeps(newPkg)
		newPkg.Deps = deps
		if err != nil {
			return nil, fmt.Errorf("failed to get deps for package %s: %w", pkgPath, err)
		}
	}
	return newPkg, nil
}

func (pm *PkgDepLoader) loadPkgDirs(deps []string) error {
	args := []string{"list", "-deps", "-f={{.ImportPath}}={{.Dir}}"}

	// Warm the go-list cache with explicit dependency patterns.
	//
	// We intentionally avoid relying on `go list ... all` when deps are provided:
	// `all` only includes packages reachable from packages in the current main module.
	// At this stage, we may have finished `go get`, but generated code has not been
	// written yet, so those dependencies are not necessarily reachable by imports.
	// In that case, `all` can miss entries such as github.com/goplus/lib/c.
	//
	// So we normalize configured deps first (for example, c -> github.com/goplus/lib/c
	// and pkg@version -> pkg), de-duplicate them, and pass them as explicit patterns
	// to get stable ImportPath -> Dir mappings for later dependency loading.
	patterns := listPatterns(deps)
	args = append(args, patterns...)

	data, err := runGoCommand(pm.root, args...)
	if err != nil {
		return errors.New("go " + strings.Join(args, " ") + ": " + err.Error() + ": " + strings.TrimSpace(string(data)))
	}
	pm.pkgs = make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] != "" {
			pm.pkgs[parts[0]] = parts[1]
		}
	}
	return nil
}

func listPatterns(deps []string) []string {
	seen := make(map[string]struct{})
	patterns := make([]string, 0, len(deps))
	for _, dep := range deps {
		dep, _ = IsDepStd(dep)
		dep, _ = splitPkgPath(dep)
		if _, ok := seen[dep]; ok {
			continue
		}
		seen[dep] = struct{}{}
		patterns = append(patterns, dep)
	}
	if len(patterns) == 0 {
		patterns = append(patterns, "all")
	}
	return patterns
}

func runGoCommand(root string, args ...string) ([]byte, error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = root
	return cmd.CombinedOutput()
}

func splitPkgPath(pkgPath string) (string, string) {
	parts := strings.Split(pkgPath, "@")
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], "@" + parts[1]
}

func (pm *PkgDepLoader) InitDeps(p *PkgInfo) error {
	deps, err := pm.LoadDeps(p)
	p.Deps = deps
	if err != nil {
		return err
	}
	pm.RegisterDeps(p)
	return nil
}

// RegisterDeps registers types from dependent packages into the current conversion project's scope
func (pm *PkgDepLoader) RegisterDeps(p *PkgInfo) {
	for _, dep := range p.Deps {
		pm.RegisterDep(dep)
	}
}

func (pm *PkgDepLoader) RegisterDep(dep *PkgInfo) {
	if _, ok := pm.regCache[dep.PkgPath]; ok {
		return
	}
	pm.regCache[dep.PkgPath] = struct{}{}
	genPkg := pm.pkg
	scope := genPkg.Types.Scope()
	depPkg := genPkg.Import(dep.PkgPath)
	pm.RegisterDeps(dep)
	for cName, pubGoName := range dep.Pubs {
		if pubGoName == "" {
			pubGoName = cName
		}
		if obj := depPkg.TryRef(pubGoName); obj != nil {
			var preObj types.Object
			if pubGoName == cName {
				preObj = obj
			} else {
				preObj = gogen.NewSubst(token.NoPos, genPkg.Types, cName, obj)
			}
			if old := scope.Insert(preObj); old != nil {
				log.Printf("conflicted name `%v` in %v, previous definition is %v\n", pubGoName, dep.PkgPath, old)
			}
		}
	}
}

func IsDepStd(pkgPath string) (string, bool) {
	if pkgPath == "c" || strings.HasPrefix(pkgPath, "c/") {
		return path.Join("github.com/goplus/lib/", pkgPath), true
	}
	return pkgPath, false
}
