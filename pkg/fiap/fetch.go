package fiap

import (
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)



func Fetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	pointSets = make(map[string](model.ProcessedPointSet))
  points = make(map[string](model.ProcessedPoint))
	
	// cursorの初期化
	cursor := ""
	// 初回のFetchOnceを実行
	fetchOnceOption := &model.FetchOnceOption{AcceptableSize: option.AcceptableSize, Cursor: &cursor}
	fetchOncePointSets, fetchOncePoints, cursor, err := FetchOnce(connectionURL, keys, fetchOnceOption)
	if err != nil {
		return nil, nil, err
	}
	// pointSetsとpointsにデータを追加
	for key, value := range fetchOncePointSets {
		pointSets[key] = value
	}
	for key, value := range fetchOncePoints {
		points[key] = value
	}

	// cursorが空でない限り、繰り返し処理を行う
	for cursor != "" {
		// FetchOnceを実行
		fetchOnceOption := &model.FetchOnceOption{AcceptableSize: option.AcceptableSize,	Cursor: &cursor}
		fetchOncePointSets, fetchOncePoints, cursor, err := FetchOnce(connectionURL, keys, fetchOnceOption)
		if err != nil {
			return nil, nil, err
		}
		// pointSetsとpointsにデータを追加
		for key, value := range fetchOncePointSets {
			pointSets[key] = value
		}
		for key, value := range fetchOncePoints {
			points[key] = value
		}
		if cursor == "" {
			break
		}
	}
	return pointSets, points, nil
}


func FetchOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), cursor string, err error) {
	_, body, err := fiapFetch(connectionURL, keys, option)
	if err != nil {
		return nil, nil, "", err	
	}

	pointSets, points, cursor, err = QueryRSToProcessedDatas(body)
	if err != nil {
		return nil, nil, "", err
	} else if cursor == "" {
		return pointSets, points, "", nil
	} else {
		return pointSets, points, cursor, nil
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

func FetchDateRange(connectionURL string, fromDate time.Time, untilDate time.Time, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	// ...
	return
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKey, option *model.FetchOption, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	// ...
	return
}
