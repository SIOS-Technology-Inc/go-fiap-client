package cmd

import (
	// "encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	// "github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/cockroachdb/errors"
	// "github.com/spf13/cobra"
)

var (
	originalFetchLatest    = fetchLatest
	originalFetchOldest    = fetchLatest
	originalFetchDateRange = fetchDateRange
	originalCreateFile     = createFile
	originalMarshalJSON    = marshalJSON
	originalArgs           = os.Args

	mockFailLatest, mockFailOldest, mockFailDateRange        bool
	mockFailCreateFile, mockFailWriteFile, mockFailCloseFile bool

	mockResultPointSet  map[string](model.ProcessedPointSet)
	mockResultPoint     map[string]([]model.Value)
	mockResultFiapError *model.Error

	actualOut        = &strings.Builder{}
	actualErrOut     = &strings.Builder{}
	actualArguments  = mockFuncArguments{}
	mockFileInstance = &mockFile{}
)

type mockFuncArguments struct {
	connectionURL string
	fromDate      *time.Time
	untilDate     *time.Time
	ids           []string
}

type mockFile struct {
	fileName string
	opened   bool
	closed   bool
	builder  strings.Builder
}

func resetActualValues() {
	actualOut.Reset()
	actualErrOut.Reset()
	actualArguments.connectionURL = ""
	actualArguments.fromDate = nil
	actualArguments.untilDate = nil
	actualArguments.ids = nil
	mockFileInstance.fileName = ""
	mockFileInstance.opened = false
	mockFileInstance.closed = false
	mockFileInstance.builder.Reset()
}

func mockFetchLatest(connectionURL string, fromDate *time.Time, untilDate *time.Time, ids ...string) (map[string](model.ProcessedPointSet), map[string]([]model.Value), *model.Error, error) {
	if mockFailLatest {
		return nil, nil, nil, errors.New("test FetchLatest error")
	} else {
		actualArguments.connectionURL = connectionURL
		actualArguments.fromDate = fromDate
		actualArguments.untilDate = untilDate
		actualArguments.ids = ids
		return mockResultPointSet, mockResultPoint, mockResultFiapError, nil
	}
}

func mockFetchOldest(connectionURL string, fromDate *time.Time, untilDate *time.Time, ids ...string) (map[string](model.ProcessedPointSet), map[string]([]model.Value), *model.Error, error) {
	if mockFailOldest {
		return nil, nil, nil, errors.New("test FetchOldest error")
	} else {
		actualArguments.connectionURL = connectionURL
		actualArguments.fromDate = fromDate
		actualArguments.untilDate = untilDate
		actualArguments.ids = ids
		return mockResultPointSet, mockResultPoint, mockResultFiapError, nil
	}
}

func mockFetchDateRange(connectionURL string, fromDate *time.Time, untilDate *time.Time, ids ...string) (map[string](model.ProcessedPointSet), map[string]([]model.Value), *model.Error, error) {
	if mockFailDateRange {
		return nil, nil, nil, errors.New("test FetchDateRange error")
	} else {
		actualArguments.connectionURL = connectionURL
		actualArguments.fromDate = fromDate
		actualArguments.untilDate = untilDate
		actualArguments.ids = ids
		return mockResultPointSet, mockResultPoint, mockResultFiapError, nil
	}
}

func mockCreateFile(name string) (io.WriteCloser, error) {
	if mockFailCreateFile {
		return nil, errors.New("test file create error")
	} else {
		mockFileInstance.fileName = name
		mockFileInstance.builder.Reset()
		mockFileInstance.opened = true
		mockFileInstance.closed = false
		return mockFileInstance, nil
	}
}

func (f *mockFile) Write(p []byte) (int, error) {
	if mockFailWriteFile {
		return 0, errors.New("test file write error")
	} else {
		return f.builder.Write(p)
	}
}

func (f *mockFile) Close() error {
	if mockFailCloseFile {
		return errors.New("test file close error")
	} else {
		f.closed = true
		return nil
	}
}

func marshalJSONAlwayseFailed(v any) ([]byte, error) {
	return nil, errors.New("test json marshal error")
}

func TestFetchCommandRun(t *testing.T) {
	tokyoTz := time.FixedZone("Asia/Tokyo", 9*60*60)
	newYorkTz := time.FixedZone("America/New_York", -4*60*60)

	fetchLatest = mockFetchLatest
	fetchOldest = mockFetchOldest
	fetchDateRange = mockFetchDateRange
	createFile = mockCreateFile

	t.Run("Normal", func(t *testing.T) {
		mockResultPointSet = map[string](model.ProcessedPointSet){}
		mockResultPoint = map[string]([]model.Value){
			"test_id": []model.Value{
				{Time: time.Date(2004, 4, 30, 12, 15, 3, 0, tokyoTz), Value: "100"},
				{Time: time.Date(2004, 5, 2, 9, 0, 15, 0, time.UTC), Value: "200"},
				{Time: time.Date(2004, 12, 1, 0, 0, 0, 0, newYorkTz), Value: "300"},
			},
		}
		mockResultFiapError = nil

		t.Run("WithoutFileOutput", func(t *testing.T) {
			mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = true, true, true

			expectedOut := `{"points":{"test_id":[{"time":"2004-04-30T12:15:03+09:00","value":"100"},{"time":"2004-05-02T09:00:15Z","value":"200"},{"time":"2004-12-01T00:00:00-04:00","value":"300"}]}}
`
			expectedErrOut := ""

			t.Run("FetchLatest", func(t *testing.T) {
				mockFailLatest, mockFailOldest, mockFailDateRange = false, true, true

				t.Run("LeastFlags", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "http://test.url", "test_id"}

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
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate != nil {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate != nil {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
				})
				t.Run("ExplicitSelectFlag", func(t *testing.T) {
					t.Run("Short", func(t *testing.T) {
						os.Args = []string{"go-fiap-client", "fetch", "-s", "max", "http://test.url", "test_id"}

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
						if actualArguments.connectionURL != "http://test.url" {
							t.Error("assertion error of url")
						}
						if actualArguments.fromDate != nil {
							t.Error("assertion error of from date")
						}
						if actualArguments.untilDate != nil {
							t.Error("assertion error of until date")
						}
						if actualArguments.ids == nil {
							t.Error("assertion error of id")
						} else {
							if len(actualArguments.ids) != 1 {
								t.Error("assertion error of id")
							} else if actualArguments.ids[0] != "test_id" {
								t.Error("assertion error of id")
							}
						}
					})
					t.Run("Long", func(t *testing.T) {
						os.Args = []string{"go-fiap-client", "fetch", "--select", "max", "http://test.url", "test_id"}

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
						if actualArguments.connectionURL != "http://test.url" {
							t.Error("assertion error of url")
						}
						if actualArguments.fromDate != nil {
							t.Error("assertion error of from date")
						}
						if actualArguments.untilDate != nil {
							t.Error("assertion error of until date")
						}
						if actualArguments.ids == nil {
							t.Error("assertion error of id")
						} else {
							if len(actualArguments.ids) != 1 {
								t.Error("assertion error of id")
							} else if actualArguments.ids[0] != "test_id" {
								t.Error("assertion error of id")
							}
						}
					})
				})
				t.Run("Debug", func(t *testing.T) {
					t.Run("Short", func(t *testing.T) {
						os.Args = []string{"go-fiap-client", "fetch", "-d", "http://test.url", "test_id"}
						expectedDebugPrint := `url: http://test.url
id: test_id
debug: true
output: 
select: maximum
from: <nil>
until: <nil>
`

						resetActualValues()
						if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
							t.Error("failed to run command")
						}
						if actualOut.String() != expectedDebugPrint+expectedOut {
							t.Error("assertion error of stdout")
						}
						if actualErrOut.String() != expectedErrOut {
							t.Error("assertion error of stderr")
						}
						if actualArguments.connectionURL != "http://test.url" {
							t.Error("assertion error of url")
						}
						if actualArguments.fromDate != nil {
							t.Error("assertion error of from date")
						}
						if actualArguments.untilDate != nil {
							t.Error("assertion error of until date")
						}
						if actualArguments.ids == nil {
							t.Error("assertion error of id")
						} else {
							if len(actualArguments.ids) != 1 {
								t.Error("assertion error of id")
							} else if actualArguments.ids[0] != "test_id" {
								t.Error("assertion error of id")
							}
						}
					})
					t.Run("Long", func(t *testing.T) {
						os.Args = []string{"go-fiap-client", "fetch", "--debug", "http://test.url", "test_id"}
						expectedDebugPrint := `url: http://test.url
id: test_id
debug: true
output: 
select: maximum
from: <nil>
until: <nil>
`

						resetActualValues()
						if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
							t.Error("failed to run command")
						}
						if actualOut.String() != expectedDebugPrint+expectedOut {
							t.Error("assertion error of stdout")
						}
						if actualErrOut.String() != expectedErrOut {
							t.Error("assertion error of stderr")
						}
						if actualArguments.connectionURL != "http://test.url" {
							t.Error("assertion error of url")
						}
						if actualArguments.fromDate != nil {
							t.Error("assertion error of from date")
						}
						if actualArguments.untilDate != nil {
							t.Error("assertion error of until date")
						}
						if actualArguments.ids == nil {
							t.Error("assertion error of id")
						} else {
							if len(actualArguments.ids) != 1 {
								t.Error("assertion error of id")
							} else if actualArguments.ids[0] != "test_id" {
								t.Error("assertion error of id")
							}
						}
					})
				})
				t.Run("WithFrom", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "--from", "2012-01-01T00:00:00+09:00", "http://test.url", "test_id"}
					expectedFrom := time.Date(2012, 1, 1, 0, 0, 0, 0, tokyoTz)

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
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate == nil {
						t.Error("assertion error of from date")
					} else if !actualArguments.fromDate.Equal(expectedFrom) {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate != nil {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
				})
				t.Run("WithUntil", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "--until", "2012-12-31T23:59:59+09:00", "http://test.url", "test_id"}
					expectedUntil := time.Date(2012, 12, 31, 23, 59, 59, 0, tokyoTz)

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
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate != nil {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate == nil {
						t.Error("assertion error of until date")
					} else if !actualArguments.untilDate.Equal(expectedUntil) {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
				})
			})
			t.Run("FetchOldest", func(t *testing.T) {
				mockFailLatest, mockFailOldest, mockFailDateRange = true, false, true

				t.Run("Short", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "-s", "min", "http://test.url", "test_id"}

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
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate != nil {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate != nil {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
				})
				t.Run("Long", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "--select", "min", "http://test.url", "test_id"}

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
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate != nil {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate != nil {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
				})

			})
			t.Run("FetchDateRange", func(t *testing.T) {
				mockFailLatest, mockFailOldest, mockFailDateRange = true, true, false

				t.Run("Short", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "-s", "none", "http://test.url", "test_id"}

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
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate != nil {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate != nil {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
				})
				t.Run("Long", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "--select", "none", "http://test.url", "test_id"}

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
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate != nil {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate != nil {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
				})

			})
		})
		t.Run("WithFileOutput", func(t *testing.T) {
			mockFailLatest, mockFailOldest, mockFailDateRange = false, true, true
			mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = false, false, false

			expectedOut := ""
			expectedFileOut := `{"points":{"test_id":[{"time":"2004-04-30T12:15:03+09:00","value":"100"},{"time":"2004-05-02T09:00:15Z","value":"200"},{"time":"2004-12-01T00:00:00-04:00","value":"300"}]}}`
			expectedErrOut := ""

			t.Run("LeastFlags", func(t *testing.T) {
				t.Run("Short", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "http://test.url", "test_id"}

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
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate != nil {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate != nil {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
					if !mockFileInstance.opened {
						t.Error("file not opened")
					}
					if mockFileInstance.fileName != "./test/file.ext" {
						t.Error("assertion error of opened file name")
					}
					if mockFileInstance.builder.String() != expectedFileOut {
						t.Error("assertion error of file output")
					}
					if !mockFileInstance.closed {
						t.Error("file not closed")
					}
				})
				t.Run("Long", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "--output", "./test/file.ext", "http://test.url", "test_id"}

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
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate != nil {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate != nil {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
					if !mockFileInstance.opened {
						t.Error("file not opened")
					}
					if mockFileInstance.fileName != "./test/file.ext" {
						t.Error("assertion error of opened file name")
					}
					if mockFileInstance.builder.String() != expectedFileOut {
						t.Error("assertion error of file output")
					}
					if !mockFileInstance.closed {
						t.Error("file not closed")
					}
				})
			})
			t.Run("FullFlags", func(t *testing.T) {
				t.Run("Short", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "-o", "/abs/test/file.ext", "-s", "max", "-d", "--from", "2012-12-31T23:00:00+09:00", "--until", "2012-12-31T23:59:59+09:00", "http://test.url", "test_id"}
					expectedFrom := time.Date(2012, 12, 31, 23, 0, 0, 0, tokyoTz)
					expectedUntil := time.Date(2012, 12, 31, 23, 59, 59, 0, tokyoTz)
					expectedDebugPrint := `url: http://test.url
id: test_id
debug: true
output: /abs/test/file.ext
select: maximum
from: 2012-12-31 23:00:00 +0900 +0900
until: 2012-12-31 23:59:59 +0900 +0900
`

					resetActualValues()
					if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
						t.Error("failed to run command")
					}
					if actualOut.String() != expectedDebugPrint+expectedOut {
						t.Error("assertion error of stdout")
					}
					if actualErrOut.String() != expectedErrOut {
						t.Error("assertion error of stderr")
					}
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate == nil {
						t.Error("assertion error of from date")
					} else if !actualArguments.fromDate.Equal(expectedFrom) {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate == nil {
						t.Error("assertion error of until date")
					} else if !actualArguments.untilDate.Equal(expectedUntil) {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
					if !mockFileInstance.opened {
						t.Error("file not opened")
					}
					if mockFileInstance.fileName != "/abs/test/file.ext" {
						t.Error("assertion error of opened file name")
					}
					if mockFileInstance.builder.String() != expectedFileOut {
						t.Error("assertion error of file output")
					}
					if !mockFileInstance.closed {
						t.Error("file not closed")
					}
				})
				t.Run("Long", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "--until", "2012-12-31T23:59:59+09:00", "--debug", "--output", "/abs/test/file.ext", "--from", "2012-12-31T23:00:00+09:00", "--select", "max", "http://test.url", "test_id"}
					expectedFrom := time.Date(2012, 12, 31, 23, 0, 0, 0, tokyoTz)
					expectedUntil := time.Date(2012, 12, 31, 23, 59, 59, 0, tokyoTz)
					expectedDebugPrint := `url: http://test.url
id: test_id
debug: true
output: /abs/test/file.ext
select: maximum
from: 2012-12-31 23:00:00 +0900 +0900
until: 2012-12-31 23:59:59 +0900 +0900
`

					resetActualValues()
					if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
						t.Error("failed to run command")
					}
					if actualOut.String() != expectedDebugPrint+expectedOut {
						t.Error("assertion error of stdout")
					}
					if actualErrOut.String() != expectedErrOut {
						t.Error("assertion error of stderr")
					}
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate == nil {
						t.Error("assertion error of from date")
					} else if !actualArguments.fromDate.Equal(expectedFrom) {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate == nil {
						t.Error("assertion error of until date")
					} else if !actualArguments.untilDate.Equal(expectedUntil) {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
					if !mockFileInstance.opened {
						t.Error("file not opened")
					}
					if mockFileInstance.fileName != "/abs/test/file.ext" {
						t.Error("assertion error of opened file name")
					}
					if mockFileInstance.builder.String() != expectedFileOut {
						t.Error("assertion error of file output")
					}
					if !mockFileInstance.closed {
						t.Error("file not closed")
					}
				})
				t.Run("EqualSyntax", func(t *testing.T) {
					os.Args = []string{"go-fiap-client", "fetch", "--until=2012-12-31T23:59:59+09:00", "-d", "-o=/abs/test/spaced file.ext", "--from=2012-12-31T23:00:00+09:00", "--select=max", "http://test.url", "test_id"}
					expectedFrom := time.Date(2012, 12, 31, 23, 0, 0, 0, tokyoTz)
					expectedUntil := time.Date(2012, 12, 31, 23, 59, 59, 0, tokyoTz)
					expectedDebugPrint := `url: http://test.url
id: test_id
debug: true
output: /abs/test/spaced file.ext
select: maximum
from: 2012-12-31 23:00:00 +0900 +0900
until: 2012-12-31 23:59:59 +0900 +0900
`

					resetActualValues()
					if err := newRootCmd(actualOut, actualErrOut).Execute(); err != nil {
						t.Error("failed to run command")
					}
					if actualOut.String() != expectedDebugPrint+expectedOut {
						t.Error("assertion error of stdout")
					}
					if actualErrOut.String() != expectedErrOut {
						t.Error("assertion error of stderr")
					}
					if actualArguments.connectionURL != "http://test.url" {
						t.Error("assertion error of url")
					}
					if actualArguments.fromDate == nil {
						t.Error("assertion error of from date")
					} else if !actualArguments.fromDate.Equal(expectedFrom) {
						t.Error("assertion error of from date")
					}
					if actualArguments.untilDate == nil {
						t.Error("assertion error of until date")
					} else if !actualArguments.untilDate.Equal(expectedUntil) {
						t.Error("assertion error of until date")
					}
					if actualArguments.ids == nil {
						t.Error("assertion error of id")
					} else {
						if len(actualArguments.ids) != 1 {
							t.Error("assertion error of id")
						} else if actualArguments.ids[0] != "test_id" {
							t.Error("assertion error of id")
						}
					}
					if !mockFileInstance.opened {
						t.Error("file not opened")
					}
					if mockFileInstance.fileName != "/abs/test/spaced file.ext" {
						t.Error("assertion error of opened file name")
					}
					if mockFileInstance.builder.String() != expectedFileOut {
						t.Error("assertion error of file output")
					}
					if !mockFileInstance.closed {
						t.Error("file not closed")
					}
				})
			})
		})
	})
	t.Run("ArgumentError", func(t *testing.T) {
		mockFailLatest, mockFailOldest, mockFailDateRange = false, false, false
		mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = false, false, false
		mockResultPointSet = map[string](model.ProcessedPointSet){}
		mockResultPoint = map[string]([]model.Value){
			"test_id": {
				{Time: time.Date(2004, 4, 30, 12, 15, 3, 0, tokyoTz), Value: "100"},
				{Time: time.Date(2004, 5, 2, 9, 0, 15, 0, time.UTC), Value: "200"},
				{Time: time.Date(2004, 12, 1, 0, 0, 0, 0, newYorkTz), Value: "300"},
			},
		}
		mockResultFiapError = nil

		expectedOut := `Usage:
  go-fiap-client fetch [flags] URL (POINT_ID | POINTSET_ID)

Flags:
  -d, --debug           set output log level to debug
      --from string     filter query from datetime string=<Datetime in RFC 3339 format>
  -h, --help            help for fetch
  -o, --output string   specify output file path. string=<filepath>
  -s, --select string   fiap select option. string=<max|min|none> (default "max")
      --until string    filter query until datetime string=<Datetime in RFC 3339 format>

`

		t.Run("InvalidSelect", func(t *testing.T) {
			expectedErrOut := `Error: select type allows only max, min, or none
`
			expectedError := "select type allows only max, min, or none"

			t.Run("Short", func(t *testing.T) {
				os.Args = []string{"go-fiap-client", "fetch", "-s", "aaaaa", "http://test.url", "test_id"}

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected select argument error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
			})
			t.Run("Long", func(t *testing.T) {
				os.Args = []string{"go-fiap-client", "fetch", "--select", "aaaaa", "http://test.url", "test_id"}

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected select argument error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
			})
		})
		t.Run("InvalidFrom", func(t *testing.T) {
			expectedError := "from allows only datetime in RFC3339 format"

			t.Run("Format", func(t *testing.T) {
				os.Args = []string{"go-fiap-client", "fetch", "--from", "2012/01/01 00:00:00 +0900", "http://test.url", "test_id"}
				expectedErrOut := `Error: from allows only datetime in RFC3339 format: parsing time "2012/01/01 00:00:00 +0900" as "2006-01-02T15:04:05Z07:00": cannot parse "/01/01 00:00:00 +0900" as "-"
`

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected from argument error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
			})
			t.Run("Date", func(t *testing.T) {
				os.Args = []string{"go-fiap-client", "fetch", "--from", "2012-02-30T23:59:59+09:00", "http://test.url", "test_id"}
				expectedErrOut := `Error: from allows only datetime in RFC3339 format: parsing time "2012-02-30T23:59:59+09:00": day out of range
`

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected from argument error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
			})
		})
		t.Run("InvalidUntil", func(t *testing.T) {
			expectedError := "until allows only datetime in RFC3339 format"

			t.Run("Format", func(t *testing.T) {
				os.Args = []string{"go-fiap-client", "fetch", "--until", "Dec 31, 2012 11:59:59 PM JST", "http://test.url", "test_id"}
				expectedErrOut := `Error: until allows only datetime in RFC3339 format: parsing time "Dec 31, 2012 11:59:59 PM JST" as "2006-01-02T15:04:05Z07:00": cannot parse "Dec 31, 2012 11:59:59 PM JST" as "2006"
`

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected until argument error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
			})
			t.Run("Date", func(t *testing.T) {
				os.Args = []string{"go-fiap-client", "fetch", "--until", "2012-02-29T24:00:00+09:00", "http://test.url", "test_id"}
				expectedErrOut := `Error: until allows only datetime in RFC3339 format: parsing time "2012-02-29T24:00:00+09:00": hour out of range
`

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected until argument error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
			})
		})
		t.Run("FewArguments", func(t *testing.T) {
			os.Args = []string{"go-fiap-client", "fetch", "http://test.url"}
			expectedErrOut := `Error: too few arguments
`
			expectedError := "too few arguments"

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
				t.Error("expected to fail command but succeed")
			} else if !strings.Contains(err.Error(), expectedError) {
				t.Error("expected too few arguments error but not")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
		})
		t.Run("ManyArguments", func(t *testing.T) {
			os.Args = []string{"go-fiap-client", "fetch", "http://test.url", "test_id", "extra"}
			expectedErrOut := `Error: too many arguments
`
			expectedError := "too many arguments"

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
				t.Error("expected to fail command but succeed")
			} else if !strings.Contains(err.Error(), expectedError) {
				t.Error("expected too many arguments error but not")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
		})
		t.Run("Multiple", func(t *testing.T) {
			expectedSelectError := "select type allows only max, min, or none"
			expectedFromError := "from allows only datetime in RFC3339 format"
			expectedUntilError := "until allows only datetime in RFC3339 format"
			expectedFewError := "too few arguments"
			expectedManyError := "too many arguments"

			t.Run("Short", func(t *testing.T) {
				os.Args = []string{"go-fiap-client", "fetch", "-s", "aaaaa", "--from", "bbbbb", "--until", "ccccc", "http://test.url", "test_id", "extra"}
				expectedErrOut := `Error: select type allows only max, min, or none
from allows only datetime in RFC3339 format: parsing time "bbbbb" as "2006-01-02T15:04:05Z07:00": cannot parse "bbbbb" as "2006"
until allows only datetime in RFC3339 format: parsing time "ccccc" as "2006-01-02T15:04:05Z07:00": cannot parse "ccccc" as "2006"
too many arguments
`

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedSelectError) {
					t.Error("expected select argument error but not")
				} else if !strings.Contains(err.Error(), expectedFromError) {
					t.Error("expected from argument error but not")
				} else if !strings.Contains(err.Error(), expectedUntilError) {
					t.Error("expected until argument error but not")
				} else if !strings.Contains(err.Error(), expectedManyError) {
					t.Error("expected too many arguments error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
			})
			t.Run("Long", func(t *testing.T) {
				os.Args = []string{"go-fiap-client", "fetch", "--select", "aaaaa", "--from", "bbbbb", "--until", "ccccc", "http://test.url"}
				expectedErrOut := `Error: select type allows only max, min, or none
from allows only datetime in RFC3339 format: parsing time "bbbbb" as "2006-01-02T15:04:05Z07:00": cannot parse "bbbbb" as "2006"
until allows only datetime in RFC3339 format: parsing time "ccccc" as "2006-01-02T15:04:05Z07:00": cannot parse "ccccc" as "2006"
too few arguments
`

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedSelectError) {
					t.Error("expected select argument error but not")
				} else if !strings.Contains(err.Error(), expectedFromError) {
					t.Error("expected from argument error but not")
				} else if !strings.Contains(err.Error(), expectedUntilError) {
					t.Error("expected until argument error but not")
				} else if !strings.Contains(err.Error(), expectedFewError) {
					t.Error("expected too few arguments error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
			})
		})
	})
	t.Run("RuntimeError", func(t *testing.T) {
		t.Run("FileOpen", func(t *testing.T) {
			mockFailLatest, mockFailOldest, mockFailDateRange = false, false, false
			mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = true, false, false
			mockResultPointSet = map[string](model.ProcessedPointSet){}
			mockResultPoint = map[string]([]model.Value){
				"test_id": []model.Value{
					{Time: time.Date(2004, 4, 30, 12, 15, 3, 0, tokyoTz), Value: "100"},
					{Time: time.Date(2004, 5, 2, 9, 0, 15, 0, time.UTC), Value: "200"},
					{Time: time.Date(2004, 12, 1, 0, 0, 0, 0, newYorkTz), Value: "300"},
				},
			}
			mockResultFiapError = nil

			os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "http://test.url", "test_id"}
			expectedOut := ""
			expectedErrOut := `Error: cannnot open file './test/file.ext': test file create error
`
			expectedError := "cannnot open file './test/file.ext'"

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
				t.Error("expected to fail command but succeed")
			} else if !strings.Contains(err.Error(), expectedError) {
				t.Error("expected file open error but not")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
		})
		t.Run("FetchMethod", func(t *testing.T) {
			mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = false, false, false
			mockResultPointSet = map[string](model.ProcessedPointSet){}
			mockResultPoint = map[string]([]model.Value){
				"test_id": []model.Value{
					{Time: time.Date(2004, 4, 30, 12, 15, 3, 0, tokyoTz), Value: "100"},
					{Time: time.Date(2004, 5, 2, 9, 0, 15, 0, time.UTC), Value: "200"},
					{Time: time.Date(2004, 12, 1, 0, 0, 0, 0, newYorkTz), Value: "300"},
				},
			}
			mockResultFiapError = nil

			expectedOut := ""
			expectedFileOut := ""

			t.Run("FetchLatest", func(t *testing.T) {
				mockFailLatest, mockFailOldest, mockFailDateRange = true, false, false

				os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "http://test.url", "test_id"}
				expectedErrOut := `Error: failed to fetch from http://test.url: test FetchLatest error
`
				expectedError := "failed to fetch from http://test.url"

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected fetch error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
				if !mockFileInstance.opened {
					t.Error("file not opened")
				}
				if mockFileInstance.fileName != "./test/file.ext" {
					t.Error("assertion error of opened file name")
				}
				if mockFileInstance.builder.String() != expectedFileOut {
					t.Error("assertion error of file output")
				}
				if !mockFileInstance.closed {
					t.Error("file not closed")
				}
			})
			t.Run("FetchOldest", func(t *testing.T) {
				mockFailLatest, mockFailOldest, mockFailDateRange = false, true, false

				os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "-s", "min", "http://test.url", "test_id"}
				expectedErrOut := `Error: failed to fetch from http://test.url: test FetchOldest error
`
				expectedError := "failed to fetch from http://test.url"

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected fetch error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
				if !mockFileInstance.opened {
					t.Error("file not opened")
				}
				if mockFileInstance.fileName != "./test/file.ext" {
					t.Error("assertion error of opened file name")
				}
				if mockFileInstance.builder.String() != expectedFileOut {
					t.Error("assertion error of file output")
				}
				if !mockFileInstance.closed {
					t.Error("file not closed")
				}
			})
			t.Run("FetchDateRange", func(t *testing.T) {
				mockFailLatest, mockFailOldest, mockFailDateRange = false, false, true

				os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "-s", "none", "http://test.url", "test_id"}
				expectedErrOut := `Error: failed to fetch from http://test.url: test FetchDateRange error
`
				expectedError := "failed to fetch from http://test.url"

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected fetch error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
				if !mockFileInstance.opened {
					t.Error("file not opened")
				}
				if mockFileInstance.fileName != "./test/file.ext" {
					t.Error("assertion error of opened file name")
				}
				if mockFileInstance.builder.String() != expectedFileOut {
					t.Error("assertion error of file output")
				}
				if !mockFileInstance.closed {
					t.Error("file not closed")
				}
			})
		})
		t.Run("FiapError", func(t *testing.T) {
			mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = false, false, false
			mockResultPointSet = map[string]model.ProcessedPointSet{
				"test_id": {
					PointSetID: []string{"test_id_1", "test_id_2", "test_id_3"},
					PointID:    []string{"test_id_4", "test_id_5"},
				},
			}
			mockResultPoint = map[string]([]model.Value){}
			mockResultFiapError = &model.Error{Type: "test_type", Value: "test_value"}

			expectedOut := ""
			expectedFileOut := `{"point_sets":{"test_id":{"point_set_id":["test_id_1","test_id_2","test_id_3"],"point_id":["test_id_4","test_id_5"]}}}`
			expectedErrOut := `Error: fiap error: type test_type, value test_value
`
			expectedError := "fiap error: type test_type, value test_value"

			t.Run("FetchLatest", func(t *testing.T) {
				mockFailLatest, mockFailOldest, mockFailDateRange = false, true, true

				os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "http://test.url", "test_id"}

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected fiap error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
				if actualArguments.connectionURL != "http://test.url" {
					t.Error("assertion error of url")
				}
				if actualArguments.fromDate != nil {
					t.Error("assertion error of from date")
				}
				if actualArguments.untilDate != nil {
					t.Error("assertion error of until date")
				}
				if actualArguments.ids == nil {
					t.Error("assertion error of id")
				} else {
					if len(actualArguments.ids) != 1 {
						t.Error("assertion error of id")
					} else if actualArguments.ids[0] != "test_id" {
						t.Error("assertion error of id")
					}
				}
				if !mockFileInstance.opened {
					t.Error("file not opened")
				}
				if mockFileInstance.fileName != "./test/file.ext" {
					t.Error("assertion error of opened file name")
				}
				if mockFileInstance.builder.String() != expectedFileOut {
					t.Error("assertion error of file output")
				}
				if !mockFileInstance.closed {
					t.Error("file not closed")
				}
			})
			t.Run("FetchOldest", func(t *testing.T) {
				mockFailLatest, mockFailOldest, mockFailDateRange = true, false, true

				os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "-s", "min", "http://test.url", "test_id"}

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected fiap error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
				if actualArguments.connectionURL != "http://test.url" {
					t.Error("assertion error of url")
				}
				if actualArguments.fromDate != nil {
					t.Error("assertion error of from date")
				}
				if actualArguments.untilDate != nil {
					t.Error("assertion error of until date")
				}
				if actualArguments.ids == nil {
					t.Error("assertion error of id")
				} else {
					if len(actualArguments.ids) != 1 {
						t.Error("assertion error of id")
					} else if actualArguments.ids[0] != "test_id" {
						t.Error("assertion error of id")
					}
				}
				if !mockFileInstance.opened {
					t.Error("file not opened")
				}
				if mockFileInstance.fileName != "./test/file.ext" {
					t.Error("assertion error of opened file name")
				}
				if mockFileInstance.builder.String() != expectedFileOut {
					t.Error("assertion error of file output")
				}
				if !mockFileInstance.closed {
					t.Error("file not closed")
				}
			})
			t.Run("FetchDateRange", func(t *testing.T) {
				mockFailLatest, mockFailOldest, mockFailDateRange = true, true, false

				os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "-s", "none", "http://test.url", "test_id"}

				resetActualValues()
				if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
					t.Error("expected to fail command but succeed")
				} else if !strings.Contains(err.Error(), expectedError) {
					t.Error("expected fiap error but not")
				}
				if actualOut.String() != expectedOut {
					t.Error("assertion error of stdout")
				}
				if actualErrOut.String() != expectedErrOut {
					t.Error("assertion error of stderr")
				}
				if actualArguments.connectionURL != "http://test.url" {
					t.Error("assertion error of url")
				}
				if actualArguments.fromDate != nil {
					t.Error("assertion error of from date")
				}
				if actualArguments.untilDate != nil {
					t.Error("assertion error of until date")
				}
				if actualArguments.ids == nil {
					t.Error("assertion error of id")
				} else {
					if len(actualArguments.ids) != 1 {
						t.Error("assertion error of id")
					} else if actualArguments.ids[0] != "test_id" {
						t.Error("assertion error of id")
					}
				}
				if !mockFileInstance.opened {
					t.Error("file not opened")
				}
				if mockFileInstance.fileName != "./test/file.ext" {
					t.Error("assertion error of opened file name")
				}
				if mockFileInstance.builder.String() != expectedFileOut {
					t.Error("assertion error of file output")
				}
				if !mockFileInstance.closed {
					t.Error("file not closed")
				}
			})
		})
		t.Run("FileWrite", func(t *testing.T) {
			mockFailLatest, mockFailOldest, mockFailDateRange = false, true, true
			mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = false, true, false
			mockResultPointSet = map[string]model.ProcessedPointSet{
				"test_id": {
					PointSetID: []string{"test_id_1", "test_id_2", "test_id_3"},
					PointID:    []string{"test_id_4", "test_id_5"},
				},
			}
			mockResultPoint = map[string]([]model.Value){}
			mockResultFiapError = nil

			os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "http://test.url", "test_id"}
			expectedOut := ""
			expectedFileOut := ""
			expectedErrOut := `Error: failed to write file './test/file.ext': test file write error
`
			expectedError := "failed to write file './test/file.ext'"

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
				t.Error("expected to fail command but succeed")
			} else if !strings.Contains(err.Error(), expectedError) {
				t.Error("expected file write error but not")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
			if actualArguments.connectionURL != "http://test.url" {
				t.Error("assertion error of url")
			}
			if actualArguments.fromDate != nil {
				t.Error("assertion error of from date")
			}
			if actualArguments.untilDate != nil {
				t.Error("assertion error of until date")
			}
			if actualArguments.ids == nil {
				t.Error("assertion error of id")
			} else {
				if len(actualArguments.ids) != 1 {
					t.Error("assertion error of id")
				} else if actualArguments.ids[0] != "test_id" {
					t.Error("assertion error of id")
				}
			}
			if !mockFileInstance.opened {
				t.Error("file not opened")
			}
			if mockFileInstance.fileName != "./test/file.ext" {
				t.Error("assertion error of opened file name")
			}
			if mockFileInstance.builder.String() != expectedFileOut {
				t.Error("assertion error of file output")
			}
			if !mockFileInstance.closed {
				t.Error("file not closed")
			}
		})
		t.Run("FileClose", func(t *testing.T) {
			mockFailLatest, mockFailOldest, mockFailDateRange = false, true, true
			mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = false, false, true
			mockResultPointSet = map[string]model.ProcessedPointSet{
				"test_id": {
					PointSetID: []string{"test_id_1", "test_id_2", "test_id_3"},
					PointID:    []string{"test_id_4", "test_id_5"},
				},
			}
			mockResultPoint = map[string]([]model.Value){}
			mockResultFiapError = nil

			os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "http://test.url", "test_id"}
			expectedOut := ""
			expectedFileOut := `{"point_sets":{"test_id":{"point_set_id":["test_id_1","test_id_2","test_id_3"],"point_id":["test_id_4","test_id_5"]}}}`
			expectedErrOut := `Error: failed to close file './test/file.ext': test file close error
`
			expectedError := "failed to close file './test/file.ext'"

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
				t.Error("expected to fail command but succeed")
			} else if !strings.Contains(err.Error(), expectedError) {
				t.Error("expected file close error but not")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
			if actualArguments.connectionURL != "http://test.url" {
				t.Error("assertion error of url")
			}
			if actualArguments.fromDate != nil {
				t.Error("assertion error of from date")
			}
			if actualArguments.untilDate != nil {
				t.Error("assertion error of until date")
			}
			if actualArguments.ids == nil {
				t.Error("assertion error of id")
			} else {
				if len(actualArguments.ids) != 1 {
					t.Error("assertion error of id")
				} else if actualArguments.ids[0] != "test_id" {
					t.Error("assertion error of id")
				}
			}
			if !mockFileInstance.opened {
				t.Error("file not opened")
			}
			if mockFileInstance.fileName != "./test/file.ext" {
				t.Error("assertion error of opened file name")
			}
			if mockFileInstance.builder.String() != expectedFileOut {
				t.Error("assertion error of file output")
			}
			if mockFileInstance.closed {
				t.Error("expected to fail close file but succeed")
			}
		})
		t.Run("FormatToJson", func(t *testing.T) {
			marshalJSON = marshalJSONAlwayseFailed

			mockFailLatest, mockFailOldest, mockFailDateRange = false, true, true
			mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = false, false, false
			mockResultPointSet = map[string]model.ProcessedPointSet{
				"test_id": {
					PointSetID: []string{"test_id_1", "test_id_2", "test_id_3"},
					PointID:    []string{"test_id_4", "test_id_5"},
				},
			}
			mockResultPoint = map[string]([]model.Value){}
			mockResultFiapError = nil

			os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "http://test.url", "test_id"}
			expectedOut := ""
			expectedFileOut := ""
			expectedErrOut := `Error: failed to format output to json: test json marshal error
`
			expectedError := "failed to format output to json"

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
				t.Error("expected to fail command but succeed")
			} else if !strings.Contains(err.Error(), expectedError) {
				t.Error("expected json format error but not")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
			if actualArguments.connectionURL != "http://test.url" {
				t.Error("assertion error of url")
			}
			if actualArguments.fromDate != nil {
				t.Error("assertion error of from date")
			}
			if actualArguments.untilDate != nil {
				t.Error("assertion error of until date")
			}
			if actualArguments.ids == nil {
				t.Error("assertion error of id")
			} else {
				if len(actualArguments.ids) != 1 {
					t.Error("assertion error of id")
				} else if actualArguments.ids[0] != "test_id" {
					t.Error("assertion error of id")
				}
			}
			if !mockFileInstance.opened {
				t.Error("file not opened")
			}
			if mockFileInstance.fileName != "./test/file.ext" {
				t.Error("assertion error of opened file name")
			}
			if mockFileInstance.builder.String() != expectedFileOut {
				t.Error("assertion error of file output")
			}
			if !mockFileInstance.closed {
				t.Error("file not closed")
			}

			marshalJSON = originalMarshalJSON
		})
		t.Run("Multiple", func(t *testing.T) {
			mockFailLatest, mockFailOldest, mockFailDateRange = false, true, true
			mockFailCreateFile, mockFailWriteFile, mockFailCloseFile = false, true, true
			mockResultPointSet = map[string]model.ProcessedPointSet{
				"test_id": {
					PointSetID: []string{"test_id_1", "test_id_2", "test_id_3"},
					PointID:    []string{"test_id_4", "test_id_5"},
				},
			}
			mockResultPoint = map[string]([]model.Value){}
			mockResultFiapError = &model.Error{Type: "test_type", Value: "test_value"}

			os.Args = []string{"go-fiap-client", "fetch", "-o", "./test/file.ext", "http://test.url", "test_id"}
			expectedOut := ""
			expectedFileOut := ""
			expectedErrOut := `Error: fiap error: type test_type, value test_value
failed to write file './test/file.ext': test file write error
failed to close file './test/file.ext': test file close error
`
			expectedFiapError := "fiap error: type test_type, value test_value"
			expectedFileWriteError := "failed to write file './test/file.ext'"
			expectedFileCloseError := "failed to close file './test/file.ext'"

			resetActualValues()
			if err := newRootCmd(actualOut, actualErrOut).Execute(); err == nil {
				t.Error("expected to fail command but succeed")
			} else if !strings.Contains(err.Error(), expectedFiapError) {
				t.Error("expected fiap error but not")
			} else if !strings.Contains(err.Error(), expectedFileWriteError) {
				t.Error("expected file write error but not")
			} else if !strings.Contains(err.Error(), expectedFileCloseError) {
				t.Error("expected file close error but not")
			}
			if actualOut.String() != expectedOut {
				t.Error("assertion error of stdout")
			}
			if actualErrOut.String() != expectedErrOut {
				t.Error("assertion error of stderr")
			}
			if actualArguments.connectionURL != "http://test.url" {
				t.Error("assertion error of url")
			}
			if actualArguments.fromDate != nil {
				t.Error("assertion error of from date")
			}
			if actualArguments.untilDate != nil {
				t.Error("assertion error of until date")
			}
			if actualArguments.ids == nil {
				t.Error("assertion error of id")
			} else {
				if len(actualArguments.ids) != 1 {
					t.Error("assertion error of id")
				} else if actualArguments.ids[0] != "test_id" {
					t.Error("assertion error of id")
				}
			}
			if !mockFileInstance.opened {
				t.Error("file not opened")
			}
			if mockFileInstance.fileName != "./test/file.ext" {
				t.Error("assertion error of opened file name")
			}
			if mockFileInstance.builder.String() != expectedFileOut {
				t.Error("assertion error of file output")
			}
			if mockFileInstance.closed {
				t.Error("expected to fail close file but succeed")
			}
		})
	})

	fetchLatest = originalFetchLatest
	fetchOldest = originalFetchOldest
	fetchDateRange = originalFetchDateRange
	createFile = originalCreateFile
	marshalJSON = originalMarshalJSON
	os.Args = originalArgs
}
