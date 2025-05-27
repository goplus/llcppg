package config_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/goplus/llcppg/config"
)

func TestReadPubFile(t *testing.T) {
	pub, err := config.ReadPubFile(path.Join("./testdata", config.LLCPPG_PUB))
	if err != nil {
		t.Fatal(err)
	}
	if len(pub) != 3 {
		t.Fatalf("expect 3 entries, got %d", len(pub))
	}
	if pub["file"] != "FILE" || pub["err"] != "Err" || pub["stdio"] != "" {
		t.Fatalf("expect file, err, stdio, got %v", pub)
	}
}

func TestReadPubFileError(t *testing.T) {
	pub, err := config.ReadPubFile("./testdata/llcppg.txt")
	if !(pub != nil && len(pub) == 0 && err == nil) {
		t.Fatalf("expect empty map for llcppg.txt")
	}
	temp, err := os.CreateTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(temp.Name())
	content := `a b c`
	_, err = temp.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.ReadPubFile(temp.Name())
	if err == nil {
		t.Fatalf("expect error, got nil")
	}
}

func TestWritePubFile(t *testing.T) {
	pub := map[string]string{
		"file":  "FILE",
		"err":   "Err",
		"stdio": "",
	}
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	pubFile := filepath.Join(tempDir, config.LLCPPG_PUB)
	err = config.WritePubFile(pubFile, pub)
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(pubFile)
	if err != nil {
		t.Fatal(err)
	}
	expect :=
		`err Err
file FILE
stdio`
	if string(content) != expect {
		t.Fatalf("expect %s, got %s", expect, string(content))
	}

	notExistFilePath := filepath.Join(tempDir, "not_exit_dir", "not_exist_file.pub")
	err = config.WritePubFile(notExistFilePath, pub)
	if err == nil {
		t.Fatalf("expect error, got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expect os.IsNotExist error, got %v", err)
	}

	notExistFile := filepath.Join(tempDir, "not_exist_file.pub")
	err = config.WritePubFile(notExistFile, make(map[string]string, 0))
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	_, err = os.Stat(notExistFile)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("expect file %s, got error %v", notExistFile, err)
	}
}
