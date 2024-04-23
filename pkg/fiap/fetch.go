package fiap

import (
	"log"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/cockroachdb/errors"
)

func Fetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), err error) {
	tools.DebugLogPrintf("Debug: Fetch start, connectionURL: %s, keys: %v, option: %#v\n", connectionURL, keys, option)

	pointSets = make(map[string](model.ProcessedPointSet))
	points = make(map[string]([]model.Value))

	// cursorの初期化
	cursor := ""

	// 戻り値のcursorが""になるまで、繰り返し処理を行う
	i := 0
	for {
		i++
		// FetchOnceを実行
		fetchOnceOption := &model.FetchOnceOption{AcceptableSize: option.AcceptableSize, Cursor: cursor}
		fetchOncePointSets, fetchOncePoints, newCursor, err := FetchOnce(connectionURL, keys, fetchOnceOption)
		if err != nil {
			err = errors.Wrapf(err, "FetchOnce error on loop iteration %d", i)
			log.Printf("Error: %+v\n", err)
			return nil, nil, err
		}
		// pointSetにデータを追加
		for key, value := range fetchOncePointSets {
			// pointSetsのkeyが設定されていない場合にはデータを代入する
			if existingPointSet, ok := pointSets[key]; !ok {
				pointSets[key] = value
				// pointSetsのkeyが既に設定されていた場合にはデータを上書きせず追加する
			} else {
				existingPointSet.PointSetID = append(existingPointSet.PointSetID, value.PointSetID...)
				existingPointSet.PointID = append(existingPointSet.PointID, value.PointID...)
				pointSets[key] = existingPointSet
			}
		}
		// pointsにデータを追加
		for key, values := range fetchOncePoints {
			// pointsのkeyが設定されていない場合にはデータを代入する
			if existingPoint, ok := points[key]; !ok {
				points[key] = values
				// pointsのkeyが既に設定されていた場合にはデータを上書きせず追加する
			} else {
				points[key] = append(existingPoint, values...)
			}
		}
		if newCursor == "" {
			break
		}
		cursor = newCursor
	}
	tools.DebugLogPrintf("Debug: Fetch end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, nil
}

func FetchOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), cursor string, err error) {
	tools.DebugLogPrintf("Debug: FetchOnce start, connectionURL: %s, keys: %v, option: %#v\n", connectionURL, keys, option)

	_, body, err := fiapFetch(connectionURL, keys, option)
	if err != nil {
		err = errors.Wrap(err, "fiapFetch error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, "", err
	}

	pointSets, points, cursor, err = processQueryRS(body)
	if err != nil {
		err = errors.Wrap(err, "processQueryRS error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, "", err
	} else {
		tools.DebugLogPrintf("Debug: FetchOnce end, pointSets: %v, points: %v, cursor: %v\n", pointSets, points, cursor)
		return pointSets, points, cursor, nil
	}
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKeyNoID, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), err error) {
	tools.DebugLogPrintf("Debug: FetchByIdsWithKey start, connectionURL: %s, key: %#v, ids: %v\n", connectionURL, key, ids)
	if len(ids) == 0 {
		err = errors.New("ids is empty, set at least one id")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	// Fetchのためのキーを作成
	var keys []model.UserInputKey
	for _, id := range ids {
		keys = append(keys, model.UserInputKey{
			ID:              id,
			Eq:              key.Eq,
			Neq:             key.Neq,
			Lt:              key.Lt,
			Gt:              key.Gt,
			Lteq:            key.Lteq,
			Gteq:            key.Gteq,
			MinMaxIndicator: key.MinMaxIndicator,
		})
	}
	// Fetchを実行
	pointSets, points, err = Fetch(connectionURL, keys, &model.FetchOption{})
	if err != nil {
		err = errors.Wrap(err, "Fetch error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchByIdsWithKey end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, err
}

func FetchLatest(connectionURL string, fromDate *time.Time, untilDate *time.Time ,ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), err error) {
	tools.DebugLogPrintf("Debug: FetchLatest start connectionURL: %s, fromDate: %v, untilDate: %v, ids: %v\n", connectionURL, fromDate, untilDate, ids)
	pointSets, points, err = FetchByIdsWithKey(connectionURL, model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMaximum,
		Gteq:            fromDate,
		Lteq:            untilDate,
	}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchLatest end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, nil
}

func FetchOldest(connectionURL string, fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), err error) {
	tools.DebugLogPrintf("Debug: FetchOldest start connectionURL: %s, fromDate: %v, untilDate: %v, ids: %v\n", connectionURL, fromDate, untilDate, ids)
	pointSets, points, err = FetchByIdsWithKey(connectionURL, model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMinimum,
		Gteq:            fromDate,
		Lteq:            untilDate,
	}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchOldest end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, nil
}

func FetchDateRange(connectionURL string, fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), err error) {
	tools.DebugLogPrintf("Debug: FetchDateRange start, connectionURL: %s, fromDate: %v, untilDate: %v,  ids: %v\n", connectionURL, fromDate, untilDate, ids)
	pointSets, points, err = FetchByIdsWithKey(connectionURL, 
		model.UserInputKeyNoID{
			Gteq:            fromDate,
			Lteq:            untilDate,
			MinMaxIndicator: model.SelectTypeNone,
		}, 
		ids...
	)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchDateRange end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, nil
}

// processQueryRS はQueryRSを処理し、IDをキーとしたPointSetとPointのmapを返す
func processQueryRS(queryRS *model.QueryRS) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), cursor string, err error) {
	tools.DebugLogPrintf("Debug: processQueryRS start, data: %#v\n", queryRS)
	if queryRS.Transport == nil {
		err = errors.New("transport is nil")
		log.Printf("Error: %+v\n", err)
		return nil, nil, "", err
	}
	if queryRS.Transport.Body == nil {
		return nil, nil, "", nil
	}

	// BodyにPointSetが返っていれば、それを処理する
	tools.DebugLogPrintf("Debug: processQueryRS, data.Transport.Body.PointSet: %#v\n", queryRS.Transport.Body.PointSet)
	if queryRS.Transport.Body.PointSet != nil {
		tools.DebugLogPrintln("Debug: processQueryRS, pointSet is not nil")
		// pointSetsを初期化
		pointSets = make(map[string](model.ProcessedPointSet))
		// PointSetの数だけ処理を繰り返す
		for _, ps := range queryRS.Transport.Body.PointSet {
			proccessed := model.ProcessedPointSet{}
			for _, id := range ps.PointSetId {
				if id != nil {
					proccessed.PointSetID = append(proccessed.PointSetID, *id)
				}
			}
			for _, id := range ps.PointId {
				if id != nil {
					proccessed.PointID = append(proccessed.PointID, *id)
				}
			}
			// pointSetsのkeyが設定されていない場合にはデータを代入する
			if existingPointSet, ok := pointSets[ps.Id]; !ok {
				pointSets[ps.Id] = proccessed
				// pointSetsのkeyが既に設定されていた場合にはデータを上書きせず追加する
			} else {
				proccessed.PointSetID = append(existingPointSet.PointSetID, proccessed.PointSetID...)
				proccessed.PointID = append(existingPointSet.PointID, proccessed.PointID...)
				pointSets[ps.Id] = proccessed
			}
		}
	} else {
		tools.DebugLogPrintln("Debug: processQueryRS, pointSet is nil")
		pointSets = nil
	}

	// BodyにPointが返っていれば、それを処理する
	if queryRS.Transport.Body.Point != nil {
		// pointsを初期化
		points = make(map[string]([]model.Value))
		// Pointの数だけ処理を繰り返す
		for _, p := range queryRS.Transport.Body.Point {
			values := make([]model.Value, len(p.Value))
			for i, v := range p.Value {
				values[i] = *v
			}
			// pointsのkeyが設定されていない場合にはデータを代入する
			if existingValues, ok := points[p.Id]; !ok {
				points[p.Id] = values
				// pointsのkeyが既に設定されていた場合にはデータを上書きせず追加する
			} else {
				points[p.Id] = append(existingValues, values...)
			}
		}
	} else {
		points = nil
	}

	// QueryクラスにCursorがあれば、それを処理する
	if queryRS.Transport.Header.Query.Cursor != "" {
		cursor = queryRS.Transport.Header.Query.Cursor
	} else {
		cursor = ""
	}

	tools.DebugLogPrintf("Debug: processQueryRS end, pointSets: %v, points: %v, cursor: %s\n", pointSets, points, cursor)
	
	return pointSets, points, cursor, nil
}
