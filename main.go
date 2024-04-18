package main

import (
	"os"

	"github.com/SIOS-Technology-Inc/go-fiap-client/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
