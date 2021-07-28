# Code Search

This is a fork of the [Google Code Search](github.com/google/codesearch),
that supports indexing dot files.

Google Code Search is a tool for indexing and then performing
regular expression searches over large bodies of source code.
It is a set of command-line programs written in Go.

For background and an overview of the commands,
see http://swtch.com/~rsc/regexp/regexp4.html.

To install my fork:

```sh
	go get github.com/mojotx/codesearch/cmd/...@index_dot_files
```

Use "go get -u" to update an existing installation.

Michael Jarvis
July 2021