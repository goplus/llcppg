package header

import (
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Segmenter func(s string) iter.Seq[string]

type TrieNode struct {
	isLeaf    bool                 // Indicates if this node represents the end of a word
	linkCount int                  // Number of children nodes
	children  map[string]*TrieNode // Map of child nodes by segment
}

// Creates a new TrieNode with empty children map
func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[string]*TrieNode)}
}

type Trie struct {
	root      *TrieNode // Root node of the trie
	segmenter Segmenter // Function to split strings into segments
}
type Options func(*Trie) // Function type for configuring Trie options

func skipEmpty(s []string) []string {
	for len(s) > 0 && s[0] == "" {
		s = s[1:]
	}
	return s
}

func splitPathAbsSafe(path string) (paths []string) {
	originalPath := filepath.Clean(path)

	sep := string(os.PathSeparator)

	// keep absolute path info
	if filepath.IsAbs(originalPath) {
		i := strings.Index(originalPath[1:], sep)
		if i > 0 {
			// bound edge: if i is greater than zero, which means there's second separator
			// for example, /usr/, i: 3, with first separator what we just skipped, i: 4
			paths = append(paths, originalPath[0:i+1])
			paths = append(paths, skipEmpty(strings.Split(originalPath[i+1:], sep))...)
		} else {
			// start with / but no other / is found, like /usr
			paths = append(paths, originalPath)
		}
	}

	if len(paths) == 0 {
		paths = skipEmpty(strings.Split(originalPath, sep))
	}

	return
}

// Returns an option to configure path segmenter
// Splits strings by OS path separator and yields each segment
func WithPathSegmenter() Options {
	return func(t *Trie) {
		t.segmenter = func(s string) iter.Seq[string] {
			return func(yield func(string) bool) {
				for _, path := range splitPathAbsSafe(s) {
					if path != "" && !yield(path) {
						return
					}
				}
			}
		}
	}
}

// Returns an option to configure reverse path segmenter
// Splits and reverses strings by OS path separator
func WithReversePathSegmenter() Options {
	return func(t *Trie) {
		t.segmenter = func(s string) iter.Seq[string] {
			return func(yield func(string) bool) {
				paths := splitPathAbsSafe(s)

				slices.Reverse(paths)

				for _, path := range paths {
					if path != "" && !yield(path) {
						return
					}
				}
			}
		}
	}
}

// Creates a new Trie with default path segmenter
// Applies all provided options to configure the Trie
func NewTrie(opts ...Options) *Trie {
	t := &Trie{root: NewTrieNode()}

	WithPathSegmenter()(t)

	for _, o := range opts {
		o(t)
	}

	return t
}

// Inserts a string into the trie
// Creates nodes for each segment in the string
func (t *Trie) Insert(s string) {
	if s == "" {
		return
	}
	node := t.root

	for segment := range t.segmenter(s) {
		child, ok := node.children[segment]
		if !ok {
			child = NewTrieNode()
			node.children[segment] = child
			node.linkCount++
		}
		node = child
	}
	node.isLeaf = true
}

// Searches for a prefix in the trie
// Returns the node at the end of the prefix or nil if not found
func (t *Trie) searchPrefix(s string) *TrieNode {
	if s == "" {
		return nil
	}
	node := t.root

	for segment := range t.segmenter(s) {
		child, ok := node.children[segment]
		if !ok {
			return nil
		}
		node = child
	}

	return node
}

// Finds the longest common prefix of the given string
// Returns the longest prefix that exists in the trie
//
// Implement Source: https://leetcode.com/problems/longest-common-prefix/solutions/127449/longest-common-prefix
func (t *Trie) LongestPrefix(s string) string {
	var prefix []string

	node := t.root

	for segment := range t.segmenter(s) {
		child := node.children[segment]

		isLongestPrefix := child != nil && node.linkCount == 1 && !node.isLeaf

		if !isLongestPrefix {
			break
		}

		prefix = append(prefix, segment)
		node = child
	}

	return filepath.Join(prefix...)
}

// Checks if the trie contains the given string as a prefix
func (t *Trie) Contains(s string) bool {
	if s == "" {
		return false
	}
	node := t.root

	for segment := range t.segmenter(s) {
		child, ok := node.children[segment]
		if !ok {
			if node == t.root {
				node = nil
			}
			break
		}
		node = child
	}

	return node != nil
}

// Checks if the trie contains the exact string
// Returns true if the string exists in the trie
func (t *Trie) Search(s string) bool {
	node := t.searchPrefix(s)
	return node != nil && node.isLeaf
}
