package fiap

import (
	"log"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/cockroachdb/errors"
)

func Fetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	tools.DebugLogPrintf("Debug: Fetch start, connectionURL: %s, keys: %v, option: %#v\n", connectionURL, keys, option)
	
	pointSets = make(map[string](model.ProcessedPointSet))
	points = make(map[string](model.ProcessedPoint))
	
	// cursorの初期化
	var cursor *string
	cursor = nil
	
	// 初回のFetchOnceを実行
	fetchOnceOption := &model.FetchOnceOption{AcceptableSize: option.AcceptableSize, Cursor: cursor}
	fetchOncePointSets, fetchOncePoints, cursor, err := FetchOnce(connectionURL, keys, fetchOnceOption)
	if err != nil {
		err = errors.Wrap(err, "FetchOnce error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	// TODO: Fetchの処理が、keyが重複しているデータを上書きしてしまうため、修正が必要

	// pointSetsとpointsにデータを追加
	for key, value := range fetchOncePointSets {
		pointSets[key] = value
	}
	for key, value := range fetchOncePoints {
		points[key] = value
	}

	// cursorが空でない限り、繰り返し処理を行う
	i := 0
	for cursor != nil {
		i++
		// FetchOnceを実行
		fetchOnceOption := &model.FetchOnceOption{AcceptableSize: option.AcceptableSize,	Cursor: cursor}
		fetchOncePointSets, fetchOncePoints, cursor, err := FetchOnce(connectionURL, keys, fetchOnceOption)
		if err != nil {
			err = errors.Wrapf(err, "FetchOnce error on loop iteration %d", i)
			log.Printf("Error: %+v\n", err)
			return nil, nil, err
		}
		// pointSetsとpointsにデータを追加
		for key, value := range fetchOncePointSets {
			pointSets[key] = value
		}
		for key, value := range fetchOncePoints {
			points[key] = value
		}
		if cursor == nil {
			break
		}
	}
	tools.DebugLogPrintf("Debug: Fetch end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, nil
}


func FetchOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), cursor *string, err error) {
	tools.DebugLogPrintf("Debug: FetchOnce start, connectionURL: %s, keys: %v, option: %#v\n", connectionURL, keys, option)

	_, body, err := fiapFetch(connectionURL, keys, option)
	if err != nil {
		err = errors.Wrap(err, "fiapFetch error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, nil, err	
	}

	pointSets, points, cursor, err = processQueryRS(body)
	if err != nil {
		err = errors.Wrap(err, "processQueryRS error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, nil, err
	} else if cursor == nil {
		tools.DebugLogPrintf("Debug: FetchOnce end without cursor, pointSets: %v, points: %v\n", pointSets, points)
		return pointSets, points, nil, nil
	} else {
		tools.DebugLogPrintf("Debug: FetchOnce end with cursor, pointSets: %v, points: %v, cursor: %v\n", pointSets, points, *cursor)
		return pointSets, points, cursor, nil
	}
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKeyNoID, option *model.FetchOption, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	tools.DebugLogPrintf("Debug: FetchByIdsWithKey start, connectionURL: %s, key: %v, option: %#v, ids: %v\n", connectionURL, key, option, ids)
	// Fetchのためのキーを作成
	var keys []model.UserInputKey
	for _, id := range ids {
		keys = append(keys, model.UserInputKey{
			ID: id,
			Eq: key.Eq,
			Neq: key.Neq,
			Lt: key.Lt,
			Gt: key.Gt,
			Lteq: key.Lteq,
			Gteq: key.Gteq,
			MinMaxIndicator: key.MinMaxIndicator,
		})
	}
	// Fetchを実行
	pointSets, points, err = Fetch(connectionURL, keys, option)
	if err != nil {
		err = errors.Wrap(err, "Fetch error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchByIdsWithKey end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, err
}


func FetchLatest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	tools.DebugLogPrintf("Debug: FetchLatest start, connectionURL: %s, ids: %v\n", connectionURL, ids)
	var points map[string]model.ProcessedPoint
	_, points, err = FetchByIdsWithKey(connectionURL, model.UserInputKeyNoID{MinMaxIndicator: tools.Stringp("maximum")}, &model.FetchOption{}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		log.Printf("Error: %+v\n", err)
		return nil, err
	}
	datas = make(map[string]string)
	for id, point := range points {
		datas[id] = point.Values[0].Value
	}
	tools.DebugLogPrintf("Debug: FetchLatest end, datas: %v\n", datas)
	return datas, nil
}

func FetchOldest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	tools.DebugLogPrintf("Debug: FetchOldest start, connectionURL: %s, ids: %v\n", connectionURL, ids)
	var points map[string]model.ProcessedPoint
	_, points, err = FetchByIdsWithKey(connectionURL, model.UserInputKeyNoID{MinMaxIndicator: tools.Stringp("minimum")}, &model.FetchOption{}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		log.Printf("Error: %+v\n", err)
		return nil, err
	}
	datas = make(map[string]string)
	for id, point := range points {
		datas[id] = point.Values[0].Value
	}
	tools.DebugLogPrintf("Debug: FetchOldest end, datas: %v\n", datas)
	return datas, nil
}

func FetchDateRange(connectionURL string, fromDate time.Time, untilDate time.Time, option *model.FetchOption, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), err error) {
	tools.DebugLogPrintf("Debug: FetchDateRange start, connectionURL: %s, fromDate: %v, untilDate: %v, option: %#v, ids: %v\n", connectionURL, fromDate, untilDate, option, ids)
	pointSets, points, err = FetchByIdsWithKey(connectionURL, model.UserInputKeyNoID{Gteq: &fromDate, Lteq: &untilDate}, &model.FetchOption{}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchDateRange end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, nil
}

func processQueryRS(data *model.QueryRS) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), cursor *string, err error){
	tools.DebugLogPrintf("Debug: processQueryRS start, data: %#v\n", data)
	if data == nil {
		err = errors.New("queryRS is nil")
		log.Printf("Error: %+v\n", err)
		return nil, nil, nil, err
	}
	if data.Transport == nil {
		err = errors.New("transport is nil")
		log.Printf("Error: %+v\n", err)
		return nil, nil, nil, err
	}
	if data.Transport.Body == nil {
		return nil, nil, nil, nil
	}

	// BodyにPointSetが返っていれば、それを処理する
	tools.DebugLogPrintf("Debug: processQueryRS, data.Transport.Body.PointSet: %#v\n", data.Transport.Body.PointSet)
	if data.Transport.Body.PointSet != nil {
		tools.DebugLogPrintln("Debug: processQueryRS, pointSet is not nil")
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
		tools.DebugLogPrintln("Debug: processQueryRS, pointSet is nil")
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
		cursorStr := string(*data.Transport.Header.Query.Cursor)
		cursor = &cursorStr
	} else {
		cursor = nil
	}

	if cursor != nil {
		tools.DebugLogPrintf("Debug: processQueryRS end, pointSets: %v, points: %v, cursor: %v\n", pointSets, points, *cursor)
	} else {
		tools.DebugLogPrintf("Debug: processQueryRS end with no cursor, pointSets: %v, points: %v\n", pointSets, points)
	}
	return pointSets, points, cursor, nil
}