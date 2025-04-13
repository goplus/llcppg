package dbg

import (
	"testing"
)

func TestSetDebugSymbolNotFound(t *testing.T) {
	tests := []struct {
		name string
	}{{
		"TestSetDebugSymbolNotFound",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugSymbolNotFound()
			if !GetDebugSymbolNotFound() {
				t.Errorf("GetDebugSymbolNotFound() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"TestSetDebugError",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugError()
			if !GetDebugError() {
				t.Errorf("GetDebugError() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugLog",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugLog()
			if !GetDebugLog() {
				t.Errorf("GetDebugLog() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugAll(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugAll",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugAll()
			if flags != DbgFlagAll {
				t.Errorf("flags = %v, want %v", flags, DbgFlagAll)
			}
		})
	}
}

func TestSetDebugSetCurFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugSetCurFile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugSetCurFile()
			if !GetDebugSetCurFile() {
				t.Errorf("GetDebugSetCurFile() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugNew",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugNew()
			if !GetDebugNew() {
				t.Errorf("GetDebugNew() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugWrite(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugWrite",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugWrite()
			if !GetDebugWrite() {
				t.Errorf("GetDebugWrite() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugUnmarshalling(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugUnmarshalling",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugUnmarshalling()
			if !GetDebugUnmarshalling() {
				t.Errorf("GetDebugUnmarshalling() = %v, want %v", false, true)
			}
		})
	}
}
