package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goplus/llcppg/cmd/gogensig/config"
	llconfig "github.com/goplus/llcppg/config"
)

func TestMain(t *testing.T) {
	testFromDir(t, "testdata", false)
}

func testFromDir(t *testing.T, relDir string, gen bool) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, relDir)
	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		name := fi.Name()
		t.Run(name, func(t *testing.T) {
			testFrom(t, dir+"/"+name, name, gen)
		})
	}
}

func testFrom(t *testing.T, dir string, modulePath string, gen bool) {
	tempDir, err := os.MkdirTemp("", "gogensig-test")
	if err != nil {
		t.Fatal("MkdirTemp failed:", err)
	}
	defer os.RemoveAll(tempDir)

	astFile := path.Join(tempDir, llconfig.LLCPPG_SIGFETCH)
	cfgFile := path.Join(tempDir, llconfig.LLCPPG_CFG)
	symbFile := path.Join(tempDir, llconfig.LLCPPG_SYMB)
	confDir := path.Join(dir, "conf")
	expect := filepath.Join(dir, "gogensig.expect")

	sourceFile, err := os.Open(path.Join(confDir, llconfig.LLCPPG_SYMB))
	if err != nil {
		t.Fatal("Open failed:", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(symbFile)
	if err != nil {
		t.Fatal("Create failed:", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		t.Fatal("Copy failed:", err)
	}

	var expectContent []byte
	if !gen {
		var err error
		expectContent, err = os.ReadFile(expect)
		if err != nil {
			t.Fatal("Expect file not found")
		}
	}

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatal("Chdir failed:", err)
	}

	conf, err := config.GetCppgCfgFromPath(path.Join(confDir, llconfig.LLCPPG_CFG))
	if err != nil {
		t.Fatal("GetCppgCfgFromPath failed:", err)
	}
	conf.CFlags = " -I" + filepath.Join(dir, "hfile")
	jsonData, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		t.Fatal("MarshalIndent failed:", err)
	}
	os.WriteFile(cfgFile, jsonData, 0644)

	astData, err := config.SigfetchConfig(cfgFile, confDir, conf.Cplusplus)
	if err != nil {
		t.Fatal("SigfetchConfig failed:", err)
	}

	if err := os.WriteFile(astFile, astData, 0644); err != nil {
		t.Fatal("WriteFile failed:", err)
	}

	// mock sigfetch config file
	os.Args = []string{"gogensig", "-mod=" + modulePath, llconfig.LLCPPG_SIGFETCH}
	main()

	var res strings.Builder

	outputDir := path.Join(tempDir, conf.Name)
	outDir, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range outDir {
		if strings.HasSuffix(fi.Name(), "go.mod") || strings.HasSuffix(fi.Name(), "go.sum") || strings.HasSuffix(fi.Name(), llconfig.LLCPPG_PUB) {
			continue
		} else {
			content, err := os.ReadFile(filepath.Join(outputDir, fi.Name()))
			if err != nil {
				t.Fatal(err)
			}
			res.WriteString(fmt.Sprintf("===== %s =====\n", fi.Name()))
			res.Write(content)
			res.WriteString("\n")
		}
	}

	pub, err := os.ReadFile(filepath.Join(outputDir, llconfig.LLCPPG_PUB))
	if err == nil {
		res.WriteString("===== llcppg.pub =====\n")
		res.Write(pub)
	}

	if gen {
		if err := os.WriteFile(expect, []byte(res.String()), 0644); err != nil {
			t.Fatal(err)
		}
	} else {
		expect := string(expectContent)
		got := res.String()
		if strings.TrimSpace(expect) != strings.TrimSpace(got) {
			t.Errorf("does not match expected.\nExpected:\n%s\nGot:\n%s", expect, got)
		}
	}
}
