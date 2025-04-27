package dbg

import (
	"testing"
)

func TestSetDebugSymbol(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugSymbol",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugSymbol()
			if !GetDebugSymbol() {
				t.Errorf("GetDebugSymbol() = got %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugParseIsMethod(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugParseIsMethod",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugParseIsMethod()
			if !GetDebugParseIsMethod() {
				t.Errorf("GetDebugParseIsMethod() = got %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugEditSymMap(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugEditSymMap",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugEditSymMap()
			if !GetDebugEditSymMap() {
				t.Errorf("GetDebugEditSymMap() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugVisitTop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"TestSetDebugVisitTop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugVisitTop()
			if !GetDebugVisitTop() {
				t.Errorf("GetDebugVisitTop() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugCollectFuncInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugCollectFuncInfo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugCollectFuncInfo()
			if !GetDebugCollectFuncInfo() {
				t.Errorf("GetDebugCollectFuncInfo() got = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugNewSymbol(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugNewSymbol",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugNewSymbol()
			if !GetDebugNewSymbol() {
				t.Errorf("GetDebugNewSymbol() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugFileType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugFileType",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugFileType()
			if !GetDebugFileType() {
				t.Errorf("GetDebugFileType() = %v, want %v", false, true)
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
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugAll()
			if dbgFlags != DbgFlagAll {
				t.Errorf("dbgFlags = %v, want %v", dbgFlags, DbgFlagAll)
			}
		})
	}
}