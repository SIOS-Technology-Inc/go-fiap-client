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

type cmdFetchFunc func(string, *time.Time, *time.Time, ...string) (map[string](model.ProcessedPointSet), map[string]([]model.Value), *model.Error, error)

var (
	fetchLatest    cmdFetchFunc = fiap.FetchLatest
	fetchOldest    cmdFetchFunc = fiap.FetchOldest
	fetchDateRange cmdFetchFunc = fiap.FetchDateRange

	createFile func(string) (io.WriteCloser, error) = func(name string) (io.WriteCloser, error) {
		return os.Create(name)
	}
)

func newFetchCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	var (
		debug        bool
		outputString string
		selectString string
		fromString   string
		untilString  string

		output     io.WriteCloser
		selectType model.SelectType = model.SelectTypeMaximum
		fromDate   *time.Time
		untilDate  *time.Time
	)

	cmd := &cobra.Command{
		Use:   "fetch [flags] URL (POINT_ID | POINTSET_ID)",
		Short: "Run FIAP fetch method once",
		RunE: func(cmd *cobra.Command, args []string) error {
			argumentErrors := make([]error, 0, 4)

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
			runtimeErrors := make([]error, 0, 3)

			connectionURL := args[0]
			id := args[1]
			tools.DEBUG = debug
			if outputString != "" {
				if f, err := createFile(outputString); err == nil {
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

			if jsonResult, fErr, err := executeFetch(connectionURL, id, fromDate, untilDate, selectType); err == nil {
				if fErr != nil {
					runtimeErrors = append(runtimeErrors, fErr)
				}
				if output != nil {
					if _, err := output.Write(jsonResult); err != nil {
						runtimeErrors = append(runtimeErrors, errors.Wrapf(err, "failed to write file '%s'", outputString))
					}
				} else {
					cmd.Println(string(jsonResult))
				}
			} else {
				if fErr != nil {
					runtimeErrors = append(runtimeErrors, fErr)
				}
				runtimeErrors = append(runtimeErrors, err)
			}

			if output != nil {
				if err := output.Close(); err != nil {
					runtimeErrors = append(runtimeErrors, errors.Wrapf(err, "failed to close file '%s'", outputString))
				}
			}
			if len(runtimeErrors) > 0 {
				return errors.Join(runtimeErrors...)
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

func executeFetch(connectionURL string, id string, fromDate, untilDate *time.Time, selectType model.SelectType) ([]byte, error, error) {
	var result struct {
		PointSets map[string](model.ProcessedPointSet) `json:"point_sets,omitempty"`
		Points    map[string]([]model.Value)           `json:"points,omitempty"`
	}
	var fiapError error = nil

	switch selectType {
	case model.SelectTypeMaximum:
		if pointSets, points, fiapErr, err := fetchLatest(connectionURL, fromDate, untilDate, id); err == nil {
			result.PointSets = pointSets
			result.Points = points
			if fiapErr != nil {
				fiapError = errors.Newf("fiap error: type %s, value %s", fiapErr.Type, fiapErr.Value)
			}
		} else {
			return nil, nil, errors.Wrapf(err, "failed to fetch from %s", connectionURL)
		}
	case model.SelectTypeMinimum:
		if pointSets, points, fiapErr, err := fetchOldest(connectionURL, fromDate, untilDate, id); err == nil {
			result.PointSets = pointSets
			result.Points = points
			if fiapErr != nil {
				fiapError = errors.Newf("fiap error: type %s, value %s", fiapErr.Type, fiapErr.Value)
			}
		} else {
			return nil, nil, errors.Wrapf(err, "failed to fetch from %s", connectionURL)
		}
	case model.SelectTypeNone:
		if pointSets, points, fiapErr, err := fetchDateRange(connectionURL, fromDate, untilDate, id); err == nil {
			result.PointSets = pointSets
			result.Points = points
			if fiapErr != nil {
				fiapError = errors.Newf("fiap error: type %s, value %s", fiapErr.Type, fiapErr.Value)
			}
		} else {
			return nil, nil, errors.Wrapf(err, "failed to fetch from %s", connectionURL)
		}
	}

	if b, err := json.Marshal(&result); err == nil {
		return b, fiapError, nil
	} else {
		return nil, fiapError, errors.Wrap(err, "failed to format output to json")
	}
}
