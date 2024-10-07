package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var libVersion = "1.0.0"

func newRootCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	version := false

	cmd := &cobra.Command{
		Use:   "go-fiap-client",
		Short: "IEEE1888 (a.k.a. UGCCNet or FIAP) library for golang",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if version {
				cmd.Println(libVersion)
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

/*
Execute parses command line arguments and executes commands.

Executeは、コマンドライン引数を解析し、コマンドを実行します。

main.goでこの関数を呼び出し、コマンドラインの処理を行います。
*/
func Execute() error {
	return newRootCmd(os.Stdout, os.Stderr).Execute()
}
