package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	version := false

	cmd := &cobra.Command{
		Use:   "go-fiap-client",
		Short: "IEEE1888 (a.k.a. UGCCNet or FIAP) library for golang",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if version {
				cmd.Println("1.0.0")
				return nil
			} else {
				return cmd.Help()
			}
		},
	}

	cmd.SetOut(out)
	cmd.SetErr(errOut)

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.AddCommand(newFetchCmd(out, errOut))

	cmd.Flags().BoolVarP(&version, "version", "v", false, "print version of go-fiap-client")

	return cmd
}

func Execute() error {
	return newRootCmd(os.Stdout, os.Stderr).Execute()
}
