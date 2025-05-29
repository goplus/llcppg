package ld_test

import (
	"reflect"
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/ld"
)

func TestLdOutput(t *testing.T) {
	res := ld.ParseOutput(
		`GNU ld (GNU Binutils for Ubuntu) 2.42
	  Supported emulations:
	   aarch64linux
	   aarch64elf
	   aarch64elf32
	   aarch64elf32b
	   aarch64elfb
	   armelf
	   armelfb
	   aarch64linuxb
	   aarch64linux32
	   aarch64linux32b
	   armelfb_linux_eabi
	   armelf_linux_eabi
	using internal linker script:
	==================================================
	/* Script for -z combreloc */
	/* Copyright (C) 2014-2024 Free Software Foundation, Inc.
	   Copying and distribution of this script, with or without modification,
	   are permitted in any medium without royalty provided the copyright
	   notice and this notice are preserved.  */
	OUTPUT_FORMAT("elf64-littleaarch64", "elf64-bigaarch64",
				  "elf64-littleaarch64")
	OUTPUT_ARCH(aarch64)
	ENTRY(_start)
	SEARCH_DIR("=/usr/local/lib/aarch64-linux-gnu"); SEARCH_DIR("=/lib/aarch64-linux-gnu"); SEARCH_DIR("=/usr/lib/aarch64-linux-gnu"); SEARCH_DIR("=/usr/local/lib"); SEARCH_DIR("=/lib"); SEARCH_DIR("=/usr/lib"); SEARCH_DIR("=/usr/aarch64-linux-gnu/lib");
	SECTIONS
	{
	  /* Read-only sections, merged into text segment: */
	  PROVIDE (__executable_start = SEGMENT_START("text-segment", 0x400000)); . = SEGMENT_START("text-segment", 0x400000) + SIZEOF_HEADERS;
	  .interp         : { *(.interp) }
	  .note.gnu.build-id  : { *(.note.gnu.build-id) }
	  .hash           : { *(.hash) }
	  .gnu.hash       : { *(.gnu.hash) }
	  .dynsym         : { *(.dynsym) }
	  .dynstr         : { *(.dynstr) }
	`)
	expect := []string{
		"/usr/local/lib/aarch64-linux-gnu",
		"/lib/aarch64-linux-gnu",
		"/usr/lib/aarch64-linux-gnu",
		"/usr/local/lib",
		"/lib",
		"/usr/lib",
		"/usr/aarch64-linux-gnu/lib",
	}
	if !reflect.DeepEqual(res, expect) {
		t.Fatalf("expect %v, but got %v", expect, res)
	}
}
