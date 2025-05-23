package unmarshal_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
	llcppg "github.com/goplus/llcppg/config"
)

func TestUnmarshalPkg(t *testing.T) {
	files := `
{
  "File": {
    "_Type": "File",
    "decls": [
      {
        "_Type": "TypeDecl",
        "Loc": {
          "_Type": "Location",
          "File": "/opt/homebrew/include/lua/lua.h"
        },
        "Doc": null,
        "Parent": null,
        "Name": {
          "_Type": "Ident",
          "Name": "lua_State"
        },
        "Type": {
          "_Type": "RecordType",
          "Tag": 0,
          "Fields": {
            "_Type": "FieldList",
            "List": null
          },
          "Methods": []
        }
      },
      {
        "_Type": "TypedefDecl",
        "Loc": {
          "_Type": "Location",
          "File": "/opt/homebrew/include/lua/lua.h"
        },
        "Doc": null,
        "Parent": null,
        "Name": {
          "_Type": "Ident",
          "Name": "lua_Number"
        },
        "Type": {
          "_Type": "BuiltinType",
          "Kind": 8,
          "Flags": 16
        }
      }
    ],
    "includes": [],
    "macros": []
  },
  "FileMap": {
    "/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/sys/_types/_rsize_t.h": {
      "FileType":3
    },
    "/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/sys/_types/_seek_set.h": {
      "FileType":3
    },
    "/opt/homebrew/include/lua/lua.h": {
      "FileType":3
    },
    "/opt/homebrew/include/lua/luaconf.h": {
      "FileType":1
    }
  }
}
`

	expected := &llcppg.Pkg{
		File: &ast.File{
			Decls: []ast.Decl{
				&ast.TypeDecl{
					Object: ast.Object{
						Loc: &ast.Location{
							File: "/opt/homebrew/include/lua/lua.h",
						},
						Name: &ast.Ident{Name: "lua_State"},
					},
					Type: &ast.RecordType{
						Tag:     0,
						Fields:  &ast.FieldList{},
						Methods: []*ast.FuncDecl{},
					},
				},
				&ast.TypedefDecl{
					Object: ast.Object{
						Loc: &ast.Location{
							File: "/opt/homebrew/include/lua/lua.h",
						},
						Name: &ast.Ident{Name: "lua_Number"},
					},
					Type: &ast.BuiltinType{
						Kind:  8,
						Flags: 16,
					},
				},
			},
			Includes: []*ast.Include{},
			Macros:   []*ast.Macro{},
		},
		FileMap: map[string]*llcppg.FileInfo{
			"/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/sys/_types/_rsize_t.h": {
				FileType: llcppg.Third,
			},
			"/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/sys/_types/_seek_set.h": {
				FileType: llcppg.Third,
			},
			"/opt/homebrew/include/lua/lua.h": {
				FileType: llcppg.Third,
			},
			"/opt/homebrew/include/lua/luaconf.h": {
				FileType: llcppg.Inter,
			},
		},
	}

	fileSet, err := unmarshal.Pkg([]byte(files))

	if err != nil {
		t.Fatalf("UnmarshalNode failed: %v", err)
	}

	resultJSON, err := json.MarshalIndent(fileSet, "", " ")
	if err != nil {
		t.Fatalf("Failed to marshal result to JSON: %v", err)
	}

	expectedJSON, err := json.MarshalIndent(expected, "", " ")

	if err != nil {
		t.Fatalf("Failed to marshal expected result to JSON: %v", err)
	}

	if string(resultJSON) != string(expectedJSON) {
		t.Errorf("JSON mismatch.\nExpected: %s\nGot: %s", string(expectedJSON), string(resultJSON))
	}
}

func TestUnmarshalPkgErrors(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{
			name:        "Invalid JSON",
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in Pkg into unmarshal.pkgTemp",
		},
		{
			name:        "not *ast.File",
			input:       `{"File": {"_Type": "File","decls": [{"_Type": "Token", "Token": 1, "Lit": "test"}]}}`,
			expectedErr: "unmarshal error in Pkg when converting File of unmarshal.pkgTemp",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := unmarshal.Pkg([]byte(tc.input))
			if tc.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error containing %q, but got nil", tc.expectedErr)
				} else if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Errorf("Expected error containing %q, but got: %v", tc.expectedErr, err)
				}
			}
		})
	}
}
