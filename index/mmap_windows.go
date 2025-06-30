// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/rs/zerolog/log"
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
	if size == 0 {
		return mmapData{f, nil}
	}
	h, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READONLY, uint32(size>>32), uint32(size), nil)
	if err != nil {
		log.Fatal().Err(err).Str("file", f.Name()).Msg("failed to create file mapping")
	}

	addr, err := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		log.Fatal().Err(err).Str("file", f.Name()).Msg("failed to map view of file")
	}
	data := (*[1 << 30]byte)(unsafe.Pointer(addr))
	return mmapData{f, data[:size]}
}
