// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || freebsd || openbsd || netbsd
// +build darwin freebsd openbsd netbsd

package index

import (
	"os"
	"syscall"

	"github.com/rs/zerolog/log"
)

// missing from package syscall on freebsd, openbsd
const (
	_PROT_READ  = 1
	_MAP_SHARED = 1
)

func mmapFile(f *os.File) mmapData {
	st, err := f.Stat()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to stat file")
	}
	size := st.Size()
	if int64(int(size+4095)) != size+4095 {
		log.Fatal().Str("file", f.Name()).Msg("file size too large for mmap")
	}
	n := int(size)
	if n == 0 {
		return mmapData{f, nil}
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, (n+4095)&^4095, _PROT_READ, _MAP_SHARED)
	if err != nil {
		log.Fatal().Err(err).Str("file", f.Name()).Msg("failed to mmap file")
	}
	return mmapData{f, data[:n]}
}
