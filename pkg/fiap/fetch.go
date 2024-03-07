package fiap

import (
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)



func Fetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	// ...
	return
}


// func FetchOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), cursor string, err error) {
// 	_, body, err := fiapFetch(connectionURL, keys, option)
	
// 	if err != nil {
// 		return nil, nil, "", err	
// 	}
// 	if cursor == "" {
// 		points, pointSets, err := RawQueryRSToProcessedDatas(body)
// 		if err != nil {
// 			return nil, nil, "", err
// 		}
// 		return pointSets, points, "", nil
// 	}
// 	return
// }


func FetchLatest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	// ...
	return
}

func FetchOldest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	// ...
	return
}

func FetchDateRange(connectionURL string, fromDate time.Time, untilDate time.Time, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	// ...
	return
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKey, option *model.FetchOption, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	// ...
	return
}
