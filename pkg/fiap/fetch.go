package fiap

import (
	"time"
	"fmt"

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

	pointSets, points, cursor, err = ProcessQueryRS(body)
	if err != nil {
		return nil, nil, "", err
	} else if cursor == "" {
		return pointSets, points, "", nil
	} else {
		return pointSets, points, cursor, nil
	}
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKey, option *model.FetchOption, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	// ...
	return
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

func ProcessQueryRS(data *model.QueryRS) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), cursor string, err error){
	// エラー処理
	if data == nil {
		return nil, nil, "", fmt.Errorf("QueryRS is nil")
	}	

	// BodyにPointSetが返っていれば、それを処理する
	if data.Transport.Body.PointSet != nil {
		// 内部処理のために、IDを格納する配列を作成
		var pointSetSecondLayerPointSetIds []string ;
		var pointSetSecondLayerPointIds []string;
		// pointSetsを初期化
		pointSets = make(map[string](model.ProcessedPointSet))
		// PointSetの数だけ処理を繰り返す
		for _, pointSetFirstLayer := range data.Transport.Body.PointSet {
			// 初期化
			pointSetSecondLayerPointSetIds = nil
			pointSetSecondLayerPointIds = nil
			// PointSet直下のPointSetをループ処理
			for _, pointSetSecondLayerPoinSet := range pointSetFirstLayer.PointSet {
				pointSetSecondLayerPointSetIds = append(pointSetSecondLayerPointSetIds, string(pointSetSecondLayerPoinSet.Id))
			}
			// PointSet直下のPointをループ処理
			for _, pointSetSecondLayerPoint := range pointSetFirstLayer.Point {
				pointSetSecondLayerPointIds = append(pointSetSecondLayerPointIds, string(pointSetSecondLayerPoint.Id))
			}
			// キーが重複していれば、データを結合してpointSetsに格納
			if existingPointSet, ok := pointSets[string(pointSetFirstLayer.Id)]; ok {
				existingPointSet.PointSetID = append(existingPointSet.PointSetID, pointSetSecondLayerPointSetIds...)
				existingPointSet.PointID = append(existingPointSet.PointID, pointSetSecondLayerPointIds...)
				pointSets[string(pointSetFirstLayer.Id)] = existingPointSet
			// キーが重複していなければ、pointSetsにデータを格納
			} else {
				pointSets[string(pointSetFirstLayer.Id)] = model.ProcessedPointSet{
					PointSetID: pointSetSecondLayerPointSetIds,
					PointID: pointSetSecondLayerPointIds,
				}
			}
		}
	} else {
		pointSets = nil
	}

	// BodyにPointが返っていれば、それを処理する
	if data.Transport.Body.Point != nil {
		// 内部処理のために、IDを格納する配列を作成
		var pointValues []model.ProcessedValue
		// pointsを初期化
		points = make(map[string](model.ProcessedPoint))
		// Pointの数だけ処理を繰り返す
		for _, point := range data.Transport.Body.Point {
			// 初期化
			pointValues = nil
			// Point直下のValueをループ処理
			for _, value := range point.Value {
				// ValueのTimeとValueを格納
				pointValues = append(pointValues, model.ProcessedValue{
					Time: *value.Time,
					Value: value.Value,
				})
			}
			// キーが重複していれば、データを結合してpointsに格納
			if existingPoint, ok := points[string(point.Id)]; ok {
				existingPoint.Values = append(existingPoint.Values, pointValues...)
				points[string(point.Id)] = existingPoint
			// キーが重複していなければ、pointsにデータを格納
			} else {
				points[string(point.Id)] = model.ProcessedPoint{
					Values: pointValues,
				}
			}
		}
	} else {
		points = nil
	}

	// QueryクラスにCursorがあれば、それを処理する
	if data.Transport.Header.Query.Cursor != nil {
		cursor = string(*data.Transport.Header.Query.Cursor)
	} else {
		cursor = ""
	}

	return pointSets, points, cursor, nil
}