// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/google/codesearch/regexp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var usageMessage = `usage: cgrep [-c] [-h] [-i] [-l] [-n] [-v] regexp [file...]

Cgrep behaves like grep, searching for regexp, an RE2 (nearly PCRE) regular expression.

The -c, -h, -i, -l, -n, and -v flags are as in grep, although note that as per Go's
flag parsing convention, they cannot be combined: the option pair -i -n
cannot be abbreviated to -in.
`

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339Nano,
		NoColor:    false,
	}).With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func usage() {
	fmt.Fprint(os.Stderr, usageMessage)
	os.Exit(2)
}

var (
	iflag      = flag.Bool("i", false, "case-insensitive match")
	cpuProfile = flag.String("cpuprofile", "", "write cpu profile to this file")
)

func main() {
	var g regexp.Grep
	g.AddFlags()
	g.AddVFlag()
	g.Stdout = os.Stdout
	g.Stderr = os.Stderr
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
	}

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create CPU profile file")
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Error().Err(err).Msg("failed to close CPU profile file")
			}
		}()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal().Err(err).Msg("failed to start CPU profile")
		}
		defer pprof.StopCPUProfile()
	}

	pat := "(?m)" + args[0]
	if *iflag {
		pat = "(?i)" + pat
	}
	re, err := regexp.Compile(pat)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to compile regular expression")
	}
	g.Regexp = re
	if len(args) == 1 {
		g.Reader(os.Stdin, "<standard input>")
	} else {
		for _, arg := range args[1:] {
			g.File(arg)
		}
	}
	if !g.Match {
		os.Exit(1)
	}
}
