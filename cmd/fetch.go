package cmd

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

func newFetchCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	var (
		debug        bool
		outputString string
		selectString string
		fromString   string
		untilString  string

		output     *os.File
		selectType model.SelectType = model.SelectTypeMaximum
		fromDate   *time.Time
		untilDate  *time.Time
	)

	cmd := &cobra.Command{
		Use:   "fetch [flags] URL (POINT_ID | POINTSET_ID)",
		Short: "Run FIAP fetch method once",
		RunE: func(cmd *cobra.Command, args []string) error {
			argumentErrors := make([]error, 0, 5)

			switch selectString {
			case "max":
				selectType = model.SelectTypeMaximum
			case "min":
				selectType = model.SelectTypeMinimum
			case "none":
				selectType = model.SelectTypeNone
			default:
				argumentErrors = append(argumentErrors, errors.New("select type allows only max, min, or none"))
			}
			if fromString != "" {
				if dt, err := time.Parse(time.RFC3339, fromString); err == nil {
					fromDate = &dt
				} else {
					argumentErrors = append(argumentErrors, errors.Wrap(err, "from allows only datetime in RFC3339 format"))
				}
			}
			if untilString != "" {
				if dt, err := time.Parse(time.RFC3339, untilString); err == nil {
					untilDate = &dt
				} else {
					argumentErrors = append(argumentErrors, errors.Wrap(err, "until allows only datetime in RFC3339 format"))
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
			cmd.SilenceUsage = true

			connectionURL := args[0]
			id := args[1]
			tools.DEBUG = debug
			if outputString != "" {
				if f, err := os.Create(outputString); err == nil {
					output = f
				} else {
					return errors.Wrapf(err, "cannnot open file '%s'", outputString)
				}
			}

			if debug {
				cmd.Println("url:", connectionURL)
				cmd.Println("id:", id)
				cmd.Println("debug:", tools.DEBUG)
				cmd.Println("output:", outputString)
				cmd.Println("select:", selectType)
				cmd.Println("from:", fromDate)
				cmd.Println("until:", untilDate)
			}

			if jsonResult, err := executeFetch(connectionURL, id, fromDate, untilDate, selectType); err == nil {
				if output != nil {
					if _, err := output.Write(jsonResult); err != nil {
						argumentErrors = append(argumentErrors, errors.Wrapf(err, "failed to write file '%s'", outputString))
					}
				} else {
					cmd.Println(string(jsonResult))
				}
			} else {
				argumentErrors = append(argumentErrors, err)
			}

			if output != nil {
				if err := output.Close(); err != nil {
					argumentErrors = append(argumentErrors, errors.Wrapf(err, "failed to close file '%s'", outputString))
				}
			}
			if len(argumentErrors) > 0 {
				return errors.Join(argumentErrors...)
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

	return cmd
}

func executeFetch(connectionURL string, id string, fromDate, untilDate *time.Time, selectType model.SelectType) ([]byte, error) {
	var result struct {
		PointSets map[string](model.ProcessedPointSet) `json:"point_sets,omitempty"`
		Points    map[string]([]model.Value)           `json:"points,omitempty"`
		Datas     map[string](string)                  `json:"datas,omitempty"`
	}

	switch selectType {
	case model.SelectTypeMaximum:
		if datas, err := fiap.FetchLatest(connectionURL, fromDate, untilDate, id); err == nil {
			result.Datas = datas
		} else {
			return nil, errors.Wrapf(err, "failed to fetch from %s", connectionURL)
		}
	case model.SelectTypeMinimum:
		if datas, err := fiap.FetchOldest(connectionURL, fromDate, untilDate, id); err == nil {
			result.Datas = datas
		} else {
			return nil, errors.Wrapf(err, "failed to fetch from %s", connectionURL)
		}
	case model.SelectTypeNone:
		if pointSets, points, err := fiap.FetchDateRange(connectionURL, fromDate, untilDate, id); err == nil {
			result.PointSets = pointSets
			result.Points = points
		} else {
			return nil, errors.Wrapf(err, "failed to fetch from %s", connectionURL)
		}
	}

	if b, err := json.Marshal(&result); err == nil {
		return b, nil
	} else {
		return nil, errors.Wrap(err, "failed to format output to json")
	}
}
