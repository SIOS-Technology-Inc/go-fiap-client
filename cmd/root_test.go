package cmd

import (
	"os"
	"testing"
)

func TestRootCommandRun(t *testing.T) {
	mockFailLatest, mockFailOldest, mockFailDateRange = true, true, true
	mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = true, true, true

	t.Run("help", func(t *testing.T) {
		expectedOut := `IEEE1888 (a.k.a. UGCCNet or FIAP) library for golang

Usage:
  go-fiap-client [flags]
  go-fiap-client [command]

Available Commands:
  fetch       Run FIAP fetch method once

Flags:
  -h, --help      help for go-fiap-client
  -v, --version   print version of go-fiap-client

Use "go-fiap-client [command] --help" for more information about a command.
`
		expectedErrOut := ""

		t.Run("LeastFlags", func(t *testing.T) {
			os.Args = []string{"go-fiap-client"}

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
				t.Error("failed to run command")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
		})
		t.Run("Short", func(t *testing.T) {
			os.Args = []string{"go-fiap-client", "-h"}

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
				t.Error("failed to run command")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
		})
		t.Run("Long", func(t *testing.T) {
			os.Args = []string{"go-fiap-client", "--help"}

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
				t.Error("failed to run command")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
		})
	})
	t.Run("version", func(t *testing.T) {
		expectedOut := libVersion + "\n"
		expectedErrOut := ""

		t.Run("Short", func(t *testing.T) {
			os.Args = []string{"go-fiap-client", "-v"}

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
				t.Error("failed to run command")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
		})
		t.Run("Long", func(t *testing.T) {
			os.Args = []string{"go-fiap-client", "--version"}

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
				t.Error("failed to run command")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
		})
	})
}
