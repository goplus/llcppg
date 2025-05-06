package gen

import (
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestNewObjFile(t *testing.T) {
	type args struct {
		oFile string
		hFile string
	}
	tests := []struct {
		name string
		args args
		want *ObjFile
	}{
		{
			"cJSON",
			args{
				"cJSON.o",
				"cJSON.h",
			},
			&ObjFile{
				OFile: "cJSON.o",
				HFile: "cJSON.h",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewObjFile(tt.args.oFile, tt.args.hFile); !got.IsEqual(tt.want) {
				t.Errorf("NewObjFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewObjFileString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want *ObjFile
	}{
		{
			"cjson_ok",
			args{
				"cJSON.o:cJSON.h",
			},
			&ObjFile{
				OFile: "cJSON.o",
				HFile: "cJSON.h",
			},
		}, {
			"cjson_fail",
			args{
				"cJSON.o cJSON.h",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewObjFileString(tt.args.str)
			if got != nil && tt.want != nil && !got.IsEqual(tt.want) {
				t.Errorf("NewObjFileString() = %v, want %v", got, tt.want)
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("NewObjFileString() = %v, want nil", got)
				}
			}
		})
	}
}

func TestObjFile_String(t *testing.T) {
	tests := []struct {
		name string
		o    *ObjFile
		want string
	}{
		{
			"cjson",
			&ObjFile{
				OFile: "cJSON.o",
				HFile: "cJSON.h",
			},
			"{OFile:cJSON.o, HFile:cJSON.h, Deps:[]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.String(); got != tt.want {
				t.Errorf("ObjFile.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCflagEntry_String(t *testing.T) {
	tests := []struct {
		name string
		c    *CflagEntry
		want string
	}{
		{
			"cjson",
			&CflagEntry{
				Include: "/user/local/opt/cjson/include",
				ObjFiles: []*ObjFile{
					{
						OFile: "cJSON.o",
						HFile: "cJSON.h",
					},
					{
						OFile: "cJSON_Utils.o",
						HFile: "cJSON_Utils.h",
					},
				},
			},
			"{Include:/user/local/opt/cjson/include, ObjFiles:[{OFile:cJSON.o, HFile:cJSON.h, Deps:[]} {OFile:cJSON_Utils.o, HFile:cJSON_Utils.h, Deps:[]}]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("CflagEntry.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjFile_IsEqual(t *testing.T) {
	type args struct {
		o *ObjFile
	}
	tests := []struct {
		name string
		p    *ObjFile
		args args
		want bool
	}{
		{
			"hfile_not_equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{},
			},
			args{
				&ObjFile{
					OFile: "a.o",
					HFile: "b.h",
					Deps:  []string{},
				},
			},
			false,
		},
		{
			"ofile_not_equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{},
			},
			args{
				&ObjFile{
					OFile: "b.o",
					HFile: "a.h",
					Deps:  []string{},
				},
			},
			false,
		},
		{
			"dep_len_not_equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{},
			},
			args{
				&ObjFile{
					OFile: "a.o",
					HFile: "a.h",
					Deps:  []string{"c.h"},
				},
			},
			false,
		},
		{
			"dep_not_equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{"b.h"},
			},
			args{
				&ObjFile{
					OFile: "a.o",
					HFile: "a.h",
					Deps:  []string{"c.h"},
				},
			},
			false,
		},
		{
			"equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{"b.h"},
			},
			args{
				&ObjFile{
					OFile: "a.o",
					HFile: "a.h",
					Deps:  []string{"b.h"},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsEqual(tt.args.o); got != tt.want {
				t.Errorf("ObjFile.IsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIncludeList(t *testing.T) {
	tests := []struct {
		name string
		want *IncludeList
	}{
		{
			"new",
			NewIncludeList(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIncludeList()
			if got == nil || tt.want == nil {
				t.Errorf("NewIncludeList() returns nil, want not nil")
			}
		})
	}
}

func TestIncludeList_AddCflagEntry(t *testing.T) {
	cjsonExpandCflags := ExpandName("libcjson", "", cfgCflagsKey)
	cjsonCflagsList := strings.Fields(cjsonExpandCflags)
	lenCjsonCflagsList := len(cjsonCflagsList)
	trimCjsonCflagList := make([]string, 2)
	if lenCjsonCflagsList > 0 {
		trimCjsonCflagList[0] = strings.TrimPrefix(strings.TrimSpace(cjsonCflagsList[0]), "-I")
	}
	if lenCjsonCflagsList > 1 {
		trimCjsonCflagList[1] = strings.TrimPrefix(strings.TrimSpace(cjsonCflagsList[1]), "-I")
	}
	_, inc := newCflags("cfg_test_data/same_rel")
	inc0 := filepath.Join(inc, "libcjson/include")
	inc1 := filepath.Join(inc, "stdcjson/include")
	_, depsInc := newCflags("cfg_test_data/deps")
	case3Inc := filepath.Join(depsInc, "case3")
	type fields struct {
		include    []string
		absPathMap map[string]struct{}
		relPathMap map[string]struct{}
	}
	type args struct {
		index int
		entry *CflagEntry
	}
	tests := []struct {
		name   string
		fields fields
		args   []args
		want   []string
	}{
		{
			"cjson",
			fields{
				make([]string, 0),
				make(map[string]struct{}),
				make(map[string]struct{}),
			},
			[]args{
				{
					0,
					&CflagEntry{
						Include: trimCjsonCflagList[0],
						ObjFiles: []*ObjFile{
							{OFile: "cJSON_Utils.o", HFile: "cjson/cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
							{OFile: "cJSON.o", HFile: "cjson/cJSON.h", Deps: []string{}},
						},
					},
				},
				{
					1,
					&CflagEntry{
						Include: trimCjsonCflagList[1],
						ObjFiles: []*ObjFile{
							{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cJSON.h"}},
							{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
						},
					},
				},
			},
			[]string{
				"cjson/cJSON_Utils.h",
				"cjson/cJSON.h",
			},
		},
		{
			"same_rel",
			fields{
				make([]string, 0),
				make(map[string]struct{}),
				make(map[string]struct{}),
			},
			[]args{
				{
					0,
					&CflagEntry{
						Include: inc0,
						ObjFiles: []*ObjFile{
							{OFile: "cjson.o", HFile: "cjson.h", Deps: []string{}},
						},
					},
				},
				{
					1,
					&CflagEntry{
						Include: inc1,
						ObjFiles: []*ObjFile{
							{OFile: "cjson.o", HFile: "cjson.h", Deps: []string{}},
						},
					},
				},
			},
			[]string{
				"cjson.h",
				"1:cjson.h",
			},
		},
		{
			"nil",
			fields{
				make([]string, 0),
				make(map[string]struct{}),
				make(map[string]struct{}),
			},
			[]args{
				{
					0,
					nil,
				},
			},
			[]string{},
		},
		{
			"empty",
			fields{
				make([]string, 0),
				make(map[string]struct{}),
				make(map[string]struct{}),
			},
			[]args{
				{
					0,
					&CflagEntry{
						Include: "",
						ObjFiles: []*ObjFile{
							{OFile: "cjson.o", HFile: "cjson.h", Deps: []string{}},
						},
					},
				},
			},
			[]string{},
		},
		{
			"deps/case3",
			fields{
				make([]string, 0),
				make(map[string]struct{}),
				make(map[string]struct{}),
			},
			[]args{
				{
					0,
					&CflagEntry{
						Include:  case3Inc,
						ObjFiles: []*ObjFile{},
						InvalidObjFiles: []*ObjFile{
							{OFile: "a.h", HFile: "a.h", Deps: []string{}},
							{OFile: "b.h", HFile: "b.h", Deps: []string{}},
							{OFile: "c.h", HFile: "c.h", Deps: []string{}},
						},
					},
				},
			},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &IncludeList{
				include:    tt.fields.include,
				absPathMap: tt.fields.absPathMap,
				relPathMap: tt.fields.relPathMap,
			}
			for _, arg := range tt.args {
				p.AddCflagEntry(arg.index, arg.entry)
			}
			sort.Strings(tt.want)
			sort.Strings(p.include)
			if !reflect.DeepEqual(p.include, tt.want) {
				t.Errorf("AddCflagEntry got %v, want %v", p.include, tt.want)
			}
		})
	}
}

func TestIncludeList_AddIncludeForObjFile(t *testing.T) {
	type fields struct {
		include    []string
		absPathMap map[string]struct{}
		relPathMap map[string]struct{}
	}
	type args struct {
		objFile *ObjFile
		entryID int
	}
	tests := []struct {
		name   string
		fields fields
		args   []args
		want   []string
	}{
		{
			"cjson",
			fields{
				make([]string, 0),
				make(map[string]struct{}),
				make(map[string]struct{}),
			},
			[]args{
				{
					&ObjFile{OFile: "cJSON_Utils.o", HFile: "cjson/cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
					0,
				},
				{
					&ObjFile{OFile: "cJSON.o", HFile: "cjson/cJSON.h", Deps: []string{}},
					0,
				},
			},
			[]string{
				"cjson/cJSON_Utils.h",
				"cjson/cJSON.h",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &IncludeList{
				include:    tt.fields.include,
				absPathMap: tt.fields.absPathMap,
				relPathMap: tt.fields.relPathMap,
			}
			for _, arg := range tt.args {
				p.AddIncludeForObjFile(arg.objFile, arg.entryID)
			}
			sort.Strings(tt.want)
			sort.Strings(p.include)
			if !reflect.DeepEqual(p.include, tt.want) {
				t.Errorf("AddIncludeForObjFile got %v, want %v", p.include, tt.want)
			}
		})
	}
}

func TestCflagEntry_IsEmpty(t *testing.T) {
	type fields struct {
		Include  string
		ObjFiles []*ObjFile
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"not empty",
			fields{"/usr/local/include", []*ObjFile{{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{}}}},
			false,
		},
		{
			"Include empty",
			fields{"", []*ObjFile{{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{}}}},
			true,
		},
		{
			"ObjFiles empty",
			fields{"/usr/local/include", []*ObjFile{}},
			true,
		},
		{
			"Include & ObjFiles empty",
			fields{"", []*ObjFile{}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CflagEntry{
				Include:  tt.fields.Include,
				ObjFiles: tt.fields.ObjFiles,
			}
			if got := c.IsEmpty(); got != tt.want {
				t.Errorf("CflagEntry.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
