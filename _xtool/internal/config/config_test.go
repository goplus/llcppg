package config_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/config"
	llconfig "github.com/goplus/llcppg/config"
)

func TestPkgHfileInfo(t *testing.T) {
	cases := []struct {
		conf *llconfig.Config
		want *config.PkgHfilesInfo
	}{
		{
			conf: &llconfig.Config{
				CFlags:  "-I./testdata/hfile -I ./testdata/thirdhfile",
				Include: []string{"temp1.h", "temp2.h"},
			},
			want: &config.PkgHfilesInfo{
				Inters: []string{"testdata/hfile/temp1.h", "testdata/hfile/temp2.h"},
				Impls:  []string{"testdata/hfile/tempimpl.h"},
			},
		},
		{
			conf: &llconfig.Config{
				CFlags:  "-I./testdata/hfile -I ./testdata/thirdhfile",
				Include: []string{"temp1.h", "temp2.h"},
				Mix:     true,
			},
			want: &config.PkgHfilesInfo{
				Inters: []string{"testdata/hfile/temp1.h", "testdata/hfile/temp2.h"},
				Impls:  []string{},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			info := config.PkgHfileInfo(tc.conf, []string{})
			if !reflect.DeepEqual(info.Inters, tc.want.Inters) {
				t.Fatalf("inter expected %v, but got %v", tc.want.Inters, info.Inters)
			}
			if !reflect.DeepEqual(info.Impls, tc.want.Impls) {
				t.Fatalf("impl expected %v, but got %v", tc.want.Impls, info.Impls)
			}

			thirdhfile, err := filepath.Abs("./testdata/thirdhfile/third.h")
			if err != nil {
				t.Fatalf("failed to get abs path: %w", err)
			}
			tfileFound := false
			stdioFound := false
			for _, tfile := range info.Thirds {
				absTfile, err := filepath.Abs(tfile)
				if err != nil {
					t.Fatalf("failed to get abs path: %w", err)
				}
				if absTfile == thirdhfile {
					tfileFound = true
				}
				if strings.HasSuffix(absTfile, "stdio.h") {
					stdioFound = true
				}
			}
			if !tfileFound || !stdioFound {
				t.Fatalf("third hfile or std hfile not found")
			}
		})
	}
}
