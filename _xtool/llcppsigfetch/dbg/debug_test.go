package dbg

import "testing"

func TestSetDebugParse(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugParse",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugParse()
			if !GetDebugParse() {
				t.Errorf("GetDebugParse() = %v, want %v", false, true)
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

func TestSetDebugVisitTop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugVisitTop",
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

func TestSetDebugProcess(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugProcess",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugProcess()
			if !GetDebugProcess() {
				t.Errorf("GetDebugProcess() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugGetCurFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugGetCurFile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugGetCurFile()
			if !GetDebugGetCurFile() {
				t.Errorf("GetDebugGetCurFile() = %v, want %v", false, true)
			}
		})
	}
}

func TestSetDebugMacro(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestSetDebugMacro",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDebugMacro()
			if !GetDebugMacro() {
				t.Errorf("GetDebugMacro() = %v, want %v", false, true)
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
