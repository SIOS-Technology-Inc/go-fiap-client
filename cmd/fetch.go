package cmd

import (
	"io"
	"os"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

func newFetchCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	var (
		debug bool
		outputString string
		selectString string
		fromString string
		untilString string
		neString string

		output *os.File
	)

	cmd := &cobra.Command{
		Use:   "fetch [flags] URL (POINT_ID | POINTSET_ID)",
		Short: "Run FIAP fetch method once",
		RunE: func(cmd *cobra.Command, args []string) error {
			argumentErrors := make([]error, 0, 5)
			keys := make([]model.UserInputKey, 1, 1)

			switch selectString {
				case "max":
					keys[0].MinMaxIndicator = model.SelectTypeMaximum
				case "min":
					keys[0].MinMaxIndicator = model.SelectTypeMinimum
				case "none":
					keys[0].MinMaxIndicator = model.SelectTypeNone
				default:
					argumentErrors = append(argumentErrors, errors.New("select type allows only max, min, or none"))
			}
			if fromString != "" {
				if dt, err := time.Parse(time.RFC3339, fromString); err == nil {
					keys[0].Gteq = &dt
				} else {
					argumentErrors = append(argumentErrors, errors.Wrap(err, "from allows only datetime in RFC3339 format"))
				}
			}
			if untilString != "" {
				if dt, err := time.Parse(time.RFC3339, untilString); err == nil {
					keys[0].Lteq = &dt
				} else {
					argumentErrors = append(argumentErrors, errors.Wrap(err, "until allows only datetime in RFC3339 format"))
				}
			}
			if neString != "" {
				if dt, err := time.Parse(time.RFC3339, neString); err == nil {
					keys[0].Neq = &dt
				} else {
					argumentErrors = append(argumentErrors, errors.Wrap(err, "ne allows only datetime in RFC3339 format"))
				}
			}
			if len(args) < 2 {
				argumentErrors = append(argumentErrors, errors.New("too few arguments"))
			} else if len(args) > 2 {
				argumentErrors = append(argumentErrors, errors.New("too few arguments"))
			}

			if len(argumentErrors) > 0 {
				return errors.Join(argumentErrors...)
			}

			connectionURL := args[0]
			keys[0].ID = args[1]
			tools.DEBUG = debug
			if outputString != "" {
				if f, err := os.Create(outputString); err == nil {
					output = f
				} else {
					return errors.Wrapf(err, "cannnot open file '%s'", outputString)
				}
			}

			cmd.Println("url:", connectionURL)
			cmd.Println("id:", keys[0].ID)
			cmd.Println("debug", tools.DEBUG)
			cmd.Println("output", outputString)
			cmd.Println("select", keys[0].MinMaxIndicator)
			cmd.Println("from", keys[0].Gteq)
			cmd.Println("until", keys[0].Lteq)
			cmd.Println("ne", keys[0].Neq)
			if output != nil {
				if err := output.Close(); err != nil {
					return errors.Wrapf(err, "failed to close file '%s'", outputString)
				}
			}
			return nil
		},
	}

	cmd.SetOut(out)
	cmd.SetErr(errOut)

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "set output log level to debug")
	cmd.Flags().StringVarP(&outputString, "output", "o", "", "specify output file path. string=<filepath>")
	cmd.Flags().StringVarP(&selectString, "select", "s", "max", "fiap select option. string=<max|min|none>")
	cmd.Flags().StringVar(&fromString, "from", "", "filter query from datetime string=<Datetime in RFC 3339 format>")
	cmd.Flags().StringVar(&untilString, "until", "", "filter query until datetime string=<Datetime in RFC 3339 format>")
	cmd.Flags().StringVar(&neString, "ne", "", "filter query not equal datetime string=<Datetime in RFC 3339 format>")

	return cmd
}
