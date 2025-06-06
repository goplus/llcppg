package header_test

import (
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/header"
)

func TestTrieContains(t *testing.T) {
	testCases := []struct {
		name     string
		search   string
		inserted []string
		want     bool
	}{
		{
			name:   "empty string",
			search: "abc",
			want:   false,
		},
		{
			name:     "input empty string",
			search:   "",
			inserted: []string{""},
			want:     false,
		},
		{
			name:     "one string",
			search:   "/a",
			inserted: []string{"/a"},
			want:     true,
		},
		{
			name:     "two string",
			search:   "/a",
			inserted: []string{"/a", "/b"},
			want:     true,
		},
		{
			name:     "multiple string case 1",
			search:   "/c",
			inserted: []string{"/a", "/b", "/d"},
			want:     false,
		},

		{
			name:     "multiple string case 2",
			search:   "",
			inserted: []string{"/a", "/b", "/d"},
			want:     false,
		},

		{
			name:     "multiple string case 3",
			search:   "/c",
			inserted: []string{"/a/c", "/b/c", "/c/d"},
			want:     true,
		},

		{
			name:     "multiple string case 4",
			search:   "/c/d",
			inserted: []string{"/a/c/d", "/b/c/d", "/c/d/a"},
			want:     true,
		},
		{
			name:     "substring string case 1",
			search:   "/a/b",
			inserted: []string{"/a"},
			want:     false,
		},
		{
			name:     "substring string case 2",
			search:   "/a",
			inserted: []string{"/a/b"},
			want:     true,
		},

		{
			name:     "substring string case 3",
			search:   "/a/b",
			inserted: []string{"/a/b", "/a/b/c"},
			want:     true,
		},

		{
			name:     "absolute path case 1",
			search:   "a",
			inserted: []string{"/a/b"},
			want:     false,
		},
		{
			name:     "absolute path case 2",
			search:   "/a",
			inserted: []string{"a/b", "a/b/c"},
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trie := header.NewTrie()

			for _, i := range tc.inserted {
				trie.Insert(i)
			}
			if got := trie.Contains(tc.search); got != tc.want {
				t.Fatalf("unexpected result: want %v got %v", tc.want, got)
			}
		})
	}
}

func TestTrieSearch(t *testing.T) {
	testCases := []struct {
		name     string
		search   string
		inserted []string
		want     bool
	}{
		{
			name:     "Empty string insertion and search",
			search:   "",
			inserted: []string{""},
			want:     false,
		},
		{
			name:     "Single directory exact match",
			search:   "/usr/local/bin/",
			inserted: []string{"/usr/local/bin/"},
			want:     true,
		},
		{
			name:     "Single directory partial match",
			search:   "/usr/local/bin/python",
			inserted: []string{"/usr/local/bin/"},
			want:     false,
		},
		{
			name:     "Multiple directories exact match",
			search:   "/usr/local/lib/",
			inserted: []string{"/usr/local/bin/", "/usr/local/lib/", "/usr/include/"},
			want:     true,
		},
		{
			name:     "Multiple directories partial match",
			search:   "/usr/local/lib/python",
			inserted: []string{"/usr/local/bin/", "/usr/local/lib/", "/usr/include/"},
			want:     false,
		},
		{
			name:     "Mixed path separators",
			search:   "/usr/local/bin/",
			inserted: []string{"/usr/local/bin/"},
			want:     true,
		},
		{
			name:     "Non-existent path",
			search:   "/non/existent/path",
			inserted: []string{"/usr/local/bin/", "/usr/local/lib/"},
			want:     false,
		},
		{
			name:     "Empty search string",
			search:   "",
			inserted: []string{"/usr/local/bin/"},
			want:     false,
		},
		{
			name:     "Subdirectory search",
			search:   "/usr/local/bin/",
			inserted: []string{"/usr/local/bin/"},
			want:     true,
		},
		{
			name:     "Deep directory structure",
			search:   "/a/b/c/d/e/f/g",
			inserted: []string{"/a/b/c/d/e/f/g"},
			want:     true,
		},
		{
			name:     "Long path with special characters",
			search:   "/home/user/!@#$%^&*()",
			inserted: []string{"/home/user/!@#$%^&*()"},
			want:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trie := header.NewTrie()
			for _, word := range tc.inserted {
				trie.Insert(word)
			}
			if got := trie.Search(tc.search); got != tc.want {
				t.Fatalf("Search(%q) = %v, want %v", tc.search, got, tc.want)
			}
		})
	}
}

func TestTrieLongestPrefix(t *testing.T) {
	tests := []struct {
		name     string
		inserted []string
		input    string
		want     string
	}{
		{
			name:     "Empty trie",
			inserted: []string{},
			input:    "/usr/local/bin",
			want:     "",
		},
		{
			name:     "Single directory exact match",
			inserted: []string{"/usr/local/bin/"},
			input:    "/usr/local/bin",
			want:     "/usr/local/bin",
		},
		{
			name:     "Single directory partial match",
			inserted: []string{"/usr/local/bin/"},
			input:    "/usr/local/bin/python",
			want:     "/usr/local/bin",
		},
		{
			name:     "Multiple directories with common prefix",
			inserted: []string{"/usr/local/bin/", "/usr/local/lib/", "/usr/include/"},
			input:    "/usr/local/bin/python",
			want:     "/usr",
		},
		{
			name:     "No common prefix",
			inserted: []string{"/home/user/", "/var/log/", "/tmp/"},
			input:    "/etc/passwd",
			want:     "",
		},
		{
			name:     "Reverse path match",
			inserted: []string{"bin", "lib", "include"},
			input:    "include/lib/bin",
			want:     "",
		},
		{
			name:     "Longer input than stored",
			inserted: []string{"/short/"},
			input:    "/shorter/path",
			want:     "",
		},
		{
			name:     "Empty input",
			inserted: []string{"/test/"},
			input:    "",
			want:     "",
		},
		{
			name:     "No match",
			inserted: []string{"/apple/", "/banana/"},
			input:    "/cherry/",
			want:     "",
		},
		{
			name:     "Partial reverse match",
			inserted: []string{"bin", "lib", "include"},
			input:    "lib/bin",
			want:     "",
		},
		{
			name: "normal case 1",
			inserted: []string{
				"/opt/homebrew/Cellar/cjson/1.7.18/include/cJSON.h",
				"/opt/homebrew/Cellar/cjson/1.7.18/include/zlib/zlib.h",
			},
			input: "/opt/homebrew/Cellar/cjson/1.7.18/include/cJSON/cJSON.h",
			want:  "/opt/homebrew/Cellar/cjson/1.7.18/include",
		},
		{
			name:     "absolute path case 1",
			inserted: []string{"/usr", "usr", "/usr/include"},
			input:    "/usr",
			want:     "",
		},
		{
			name:     "absolute path case 2",
			inserted: []string{"usr/share", "/usr", "usr/include"},
			input:    "usr/include/share",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := header.NewTrie()
			for _, word := range tt.inserted {
				trie.Insert(word)
			}
			result := trie.LongestPrefix(tt.input)
			if result != tt.want {
				t.Errorf("LongestPrefix(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

func TestTrieReverse(t *testing.T) {
	testCases := []struct {
		name     string
		search   string
		inserted []string
		want     bool
	}{
		{
			name:   "empty string",
			search: "abc",
			want:   false,
		},
		{
			name:     "input empty string",
			search:   "",
			inserted: []string{""},
			want:     false,
		},
		{
			name:     "one string",
			search:   "/a",
			inserted: []string{"/a"},
			want:     true,
		},
		{
			name:     "two string",
			search:   "/a",
			inserted: []string{"/a", "/b"},
			want:     true,
		},
		{
			name:     "multiple string case 1",
			search:   "/c",
			inserted: []string{"/a", "/b", "/d"},
			want:     false,
		},

		{
			name:     "multiple string case 2",
			search:   "",
			inserted: []string{"/a", "/b", "/d"},
			want:     false,
		},

		{
			name:     "multiple string case 3",
			search:   "/c",
			inserted: []string{"/a/c", "/b/c", "/c/d"},
			want:     false,
		},

		{
			name:     "multiple string case 4",
			search:   "/c/d",
			inserted: []string{"/a/c/d", "/b/c/d", "/c/d/a"},
			want:     false,
		},

		{
			name:     "multiple string case 5",
			search:   "c",
			inserted: []string{"/a/c", "/b/c", "/c/d"},
			want:     true,
		},
		{
			name:     "substring string case 1",
			search:   "/a/b",
			inserted: []string{"/a"},
			want:     false,
		},
		{
			name:     "substring string case 2",
			search:   "b",
			inserted: []string{"/a/b"},
			want:     true,
		},

		{
			name:     "substring string case 3",
			search:   "/a/b",
			inserted: []string{"/a/b", "/a/b/c"},
			want:     true,
		},

		{
			name:   "normal case 1",
			search: "libxslt/variables.h",
			inserted: []string{
				"/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/libxslt/imports.h",
				"/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/libxslt/xsltexports.h",
				"/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/libxslt/variables.h",
			},
			want: true,
		},

		{
			name:   "normal case 2",
			search: "libxslt/c14n.h",
			inserted: []string{
				"/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/libxslt/imports.h",
				"/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/libxslt/xsltexports.h",
				"/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include/libxslt/variables.h",
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trie := header.NewTrie(header.WithReversePathSegmenter())

			for _, i := range tc.inserted {
				trie.Insert(i)
			}
			if got := trie.Contains(tc.search); got != tc.want {
				t.Fatalf("unexpected result: want %v got %v", tc.want, got)
			}
		})
	}
}
