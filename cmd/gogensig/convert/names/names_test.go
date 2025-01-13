package names

import (
	"testing"
)

func Test_getSuffixUndercores(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"_",
			args{"_PyCfgBasicblock_"},
			"_",
		},
		{
			"__",
			args{"_PyCfgBasicblock__"},
			"__",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSuffixUndercores(tt.args.name); got != tt.want {
				t.Errorf("getSuffixUndercores() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"_name",
			args{"_PyCfgBasicblock"},
			"X_PyCfgBasicblock",
		},
		{
			"_name_",
			args{"_PyCfgBasicblock_"},
			"X_PyCfgBasicblock_",
		},
		{
			"_name__",
			args{"_PyCfgBasicblock__"},
			"X_PyCfgBasicblock__",
		},
		{
			"__name__",
			args{"__PyCfgBasicblock__"},
			"X__PyCfgBasicblock__",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PubName(tt.args.name); got != tt.want {
				t.Errorf("PubName() = %v, want %v", got, tt.want)
			}
		})
	}
}
