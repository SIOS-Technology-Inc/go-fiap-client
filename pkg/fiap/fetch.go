package fiap

import (
	"log"
	"net/http"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)



func Fetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), err error) {
	// ...
	return
}

func FetchRaw(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (raw string, err error) {
	// ...
	return
}

// func FetchOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), cursor string, err error) {
// 	res, err := fetchStructOnce(connectionURL, keys, option)
	
// 	if err != nil {
// 		return nil, nil, "", err	
// 	}
// 	if cursor == "" {
// 		points, pointSet, err := RawQueryRSToProcessedDatas(res)

// 		return pointSets, nil, "", nil
// 	}

// }

func FetchRawOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (raw *http.Response, err error) {
	// クエリを実行
	raw, _ , err = fiapFetch(connectionURL, keys, option)

	// エラーがあればログを出力して終了
	if err != nil {
		log.Fatalf("couldn't get point data: %v", err)
		return nil, err
	// エラーがなければ結果を返す
	} else {
		return raw, nil
	}
}

func FetchLatest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	// ...
	return
}

func FetchOldest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	// ...
	return
}

func FetchDateRange(connectionURL string, fromDate time.Time, untilDate time.Time, option *model.FetchOption) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), err error) {
	// ...
	return
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKey, option *model.FetchOption, ids ...string) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), err error) {
	// ...
	return
}
