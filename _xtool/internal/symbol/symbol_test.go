package symbol

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFindLibFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "symbol_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name         string
		libName      string
		mode         Mode
		shouldCreate bool
		expectError  bool
	}{
		{
			name:         "find dynamic",
			libName:      "test",
			mode:         ModeDynamic,
			shouldCreate: true,
			expectError:  false,
		},
		{
			name:         "find static",
			libName:      "test",
			mode:         ModeStatic,
			shouldCreate: true,
			expectError:  false,
		},
		{
			name:         "library not found - dynamic",
			libName:      "nonexistent",
			mode:         ModeDynamic,
			shouldCreate: false,
			expectError:  true,
		},
		{
			name:         "library not found - static",
			libName:      "nonexistent",
			mode:         ModeStatic,
			shouldCreate: false,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedFileName string

			// Dynamically create test library file if needed
			if tt.shouldCreate {
				var extension string
				if tt.mode == ModeStatic {
					extension = ".a"
				} else {
					if runtime.GOOS == "linux" {
						extension = ".so"
					} else {
						extension = ".dylib"
					}
				}

				expectedFileName = "lib" + tt.libName + extension
				filePath := filepath.Join(tempDir, expectedFileName)

				file, err := os.Create(filePath)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", filePath, err)
				}
				file.Close()
				defer os.Remove(filePath)
			}

			result, err := FindLibFile(tempDir, tt.libName, tt.mode)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if result != expectedFileName {
					t.Errorf("Expected %s, got %s", expectedFileName, result)
				}
			}
		})
	}
}
