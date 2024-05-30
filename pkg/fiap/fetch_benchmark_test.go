package fiap

import (
	"testing"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/testutil"
)

var blackhole interface{}
var benchMarkConnectionURL = "http://fiap-benchmark-server.eastus.cloudapp.azure.com:8080/mockFIAPServiceSoap"

func BenchmarkWithOutFetchOnce(b *testing.B) {
	for _ , id := range []string{
		"http://go-fiap-client/perf/test-no-perf/1/point-10MB",
		"http://go-fiap-client/perf/test-no-perf/1/point-100MB",
		"http://go-fiap-client/perf/test-no-perf/1/point-1GB",
		"http://go-fiap-client/perf/test-no-perf/100/point-10MB",
		"http://go-fiap-client/perf/test-no-perf/10/point-100MB",
		"http://go-fiap-client/perf/test-no-perf/1/pointSet-10MB",
		"http://go-fiap-client/perf/test-no-perf/1/pointSet-100MB",
		"http://go-fiap-client/perf/test-no-perf/1/pointSet-1GB",
		"http://go-fiap-client/perf/test-no-perf/100/pointSet-10MB",
		"http://go-fiap-client/perf/test-no-perf/10/pointSet-100MB",
	}	{
		f := FetchClient{ConnectionURL: benchMarkConnectionURL}

		// Fetchのテスト
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				pointSets, points, fiapErr, err := f.Fetch([]model.UserInputKey{
					{ID: id},
				}, &model.FetchOption{})
				if err != nil {
					b.Fatal(err)
				}

				blackhole = []interface{}{pointSets, points, fiapErr}
				_ = blackhole
			}
		})

		// FetchByIdsWithKeyのテスト
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				pointSets, points, fiapErr, err := f.FetchByIdsWithKey(model.UserInputKeyNoID{}, id)
				if err != nil {
					b.Fatal(err)
				}

				blackhole = []interface{}{pointSets, points, fiapErr}
				_ = blackhole
			}
		})

		// FetchLatestのテスト
		b.Run(id, func(b *testing.B) {
			fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
			toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
		
			for i := 0; i < b.N; i++ {
				pointSets, points, fiapErr, err := f.FetchLatest(
					fromDate,
					toDate,
					id)
				if err != nil {
					b.Fatal(err)
				}
		
				blackhole = []interface{}{pointSets, points, fiapErr}
				_ = blackhole
			}
		})

		// FetchOldestのテスト
		b.Run(id, func(b *testing.B) {
			fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
			toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
		
			for i := 0; i < b.N; i++ {
				pointSets, points, fiapErr, err := f.FetchOldest(
					fromDate,
					toDate,
					id)
				if err != nil {
					b.Fatal(err)
				}
		
				blackhole = []interface{}{pointSets, points, fiapErr}
				_ = blackhole
			}
		})

		// FetchDateRangeのテスト
		b.Run(id, func(b *testing.B) {
			fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
			toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
		
			for i := 0; i < b.N; i++ {
				pointSets, points, fiapErr, err := f.FetchDateRange(
					fromDate,
					toDate,
					id)
				if err != nil {
					b.Fatal(err)
				}
		
				blackhole = []interface{}{pointSets, points, fiapErr}
				_ = blackhole
			}
		})
	}
}

func BenchmarkFetchOnce(b *testing.B){
	for _ , id := range []string{
		"http://go-fiap-client/perf/test-perf/1/point-10MB",
		"http://go-fiap-client/perf/test-perf/1/point-100MB",
		"http://go-fiap-client/perf/test-perf/1/point-1GB",
	}	{
		f := FetchClient{ConnectionURL: benchMarkConnectionURL}

		// FetchOnceのテスト
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				pointSets, points, cursor, fiapErr, err := f.FetchOnce([]model.UserInputKey{
					{ID: id},
				}, &model.FetchOnceOption{})
				if err != nil {
					b.Fatal(err)
				}

				blackhole = []interface{}{pointSets, points, cursor, fiapErr}
				_ = blackhole
			}
		})
	}
}