package fiap

import (
	"testing"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/testutil"

	"github.com/globusdigital/soap"
	"context"
)

var benchMarkConnectionURL = "http://fiap-benchmark-server.eastus.cloudapp.azure.com:8080/mockFIAPServiceSoap"


func BenchmarkFiapFetch(b *testing.B) {
	for _ , id := range []string{
		"http://go-fiap-client/perf/test-no-perf/1/point-1KB",
		"http://go-fiap-client/perf/test-no-perf/1/point-10KB",
		"http://go-fiap-client/perf/test-no-perf/1/point-100KB",
		"http://go-fiap-client/perf/test-no-perf/1/point-1MB",
		"http://go-fiap-client/perf/test-perf/1/point-10MB",
		"http://go-fiap-client/perf/test-perf/1/point-100MB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-1KB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-10KB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-100KB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-1MB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-10MB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-100MB",
	}	{
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := fiapFetch(benchMarkConnectionURL, []model.UserInputKey{
					{ID: id},
				},&model.FetchOnceOption{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkSoapCall(b *testing.B) {
	for _ , id := range []string{
		"http://go-fiap-client/perf/test-no-perf/1/point-1KB",
		"http://go-fiap-client/perf/test-no-perf/1/point-10KB",
		"http://go-fiap-client/perf/test-no-perf/1/point-100KB",
		"http://go-fiap-client/perf/test-no-perf/1/point-1MB",
		"http://go-fiap-client/perf/test-perf/1/point-10MB",
		"http://go-fiap-client/perf/test-perf/1/point-100MB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-1KB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-10KB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-100KB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-1MB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-10MB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-100MB",
	}	{
		client := soap.NewClient(benchMarkConnectionURL, nil)
		queryRQ := newQueryRQ(&model.FetchOnceOption{}, []model.UserInputKey{{ID: id}})
		resBody := &model.QueryRS{}

		b.ResetTimer()
		
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {

				_, err := client.Call(context.Background(), "http://soap.fiap.org/query", queryRQ, resBody)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkWithCursor(b *testing.B) {
	for _ , id := range []string{
		"http://go-fiap-client/perf/test-no-perf/10/point-10MB",
		"http://go-fiap-client/perf/test-no-perf/10/pointSet-10MB",
	}	{
		f := FetchClient{ConnectionURL: benchMarkConnectionURL}

		// Fetchのテスト
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.Fetch([]model.UserInputKey{
					{ID: id},
				}, &model.FetchOption{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// FetchByIdsWithKeyのテスト
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.FetchByIdsWithKey(model.UserInputKeyNoID{}, id)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// FetchLatestのテスト
		b.Run(id, func(b *testing.B) {
			fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
			toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
		
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.FetchLatest(
					fromDate,
					toDate,
					id)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// FetchOldestのテスト
		b.Run(id, func(b *testing.B) {
			fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
			toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
		
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.FetchOldest(
					fromDate,
					toDate,
					id)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// FetchDateRangeのテスト
		b.Run(id, func(b *testing.B) {
			fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
			toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
		
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.FetchDateRange(
					fromDate,
					toDate,
					id)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkFetchers(b *testing.B){
	for _ , id := range []string{
		"http://go-fiap-client/perf/test-no-perf/1/point-1KB",
		"http://go-fiap-client/perf/test-no-perf/1/point-10KB",
		"http://go-fiap-client/perf/test-no-perf/1/point-100KB",
		"http://go-fiap-client/perf/test-no-perf/1/point-1MB",
		"http://go-fiap-client/perf/test-perf/1/point-10MB",
		"http://go-fiap-client/perf/test-perf/1/point-100MB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-1KB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-10KB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-100KB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-1MB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-10MB",
		"http://go-fiap-client/perf/test-perf/1/pointSet-100MB",
	}	{
		f := FetchClient{ConnectionURL: benchMarkConnectionURL}

		// FetchOnceのテスト
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _, _, err := f.FetchOnce([]model.UserInputKey{
					{ID: id},
				}, &model.FetchOnceOption{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})

			// Fetchのテスト
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.Fetch([]model.UserInputKey{
					{ID: id},
				}, &model.FetchOption{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// FetchByIdsWithKeyのテスト
		b.Run(id, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.FetchByIdsWithKey(model.UserInputKeyNoID{}, id)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// FetchLatestのテスト
		b.Run(id, func(b *testing.B) {
			fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
			toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
		
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.FetchLatest(
					fromDate,
					toDate,
					id)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// FetchOldestのテスト
		b.Run(id, func(b *testing.B) {
			fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
			toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
		
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.FetchOldest(
					fromDate,
					toDate,
					id)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// FetchDateRangeのテスト
		b.Run(id, func(b *testing.B) {
			fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
			toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
		
			for i := 0; i < b.N; i++ {
				_, _, _, err := f.FetchDateRange(
					fromDate,
					toDate,
					id)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}