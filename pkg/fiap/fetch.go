package fiap

import (
	"log"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/cockroachdb/errors"
)

func Fetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
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
		fetchOncePointSets, fetchOncePoints, newCursor ,fiapErr, err := FetchOnce(connectionURL, keys, fetchOnceOption)
		if err != nil {
			err = errors.Wrapf(err, "FetchOnce error on loop iteration %d", i)
			log.Printf("Error: %+v\n", err)
			return nil, nil, nil, err
		}
		if fiapErr != nil {
			return pointSets, points, fiapErr, nil
		}
		
		// pointSetにデータを追加
		for key, value := range fetchOncePointSets {
			tempPointSet := value
			// keyが既に設定されている場合は既存データに追加する
			if existingPointSet, ok := pointSets[key]; ok {
				existingPointSet.PointSetID = append(existingPointSet.PointSetID, value.PointSetID...)
				existingPointSet.PointID = append(existingPointSet.PointID, value.PointID...)
				tempPointSet = existingPointSet
			}
			// keyが設定されていない場合にはデータを加工せず代入する
			pointSets[key] = tempPointSet
		}
		// pointsにデータを追加
		for key, values := range fetchOncePoints {
			tempValues := values
			// pointsのkeyが設定されている場合は既存データに追加する
			if existingPoint, ok := points[key]; ok {
				tempValues = append(existingPoint, values...)
			}
			// pointsのkeyが設定されていない場合にはデータを加工せず代入する
			points[key] = tempValues
		}
		
		if newCursor == "" {
			break
		}
		cursor = newCursor
	}
	tools.DebugLogPrintf("Debug: Fetch end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, err
}

func FetchOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), cursor string, fiapErr *model.Error ,err error) {
	tools.DebugLogPrintf("Debug: FetchOnce start, connectionURL: %s, keys: %v, option: %#v\n", connectionURL, keys, option)

	_, body, err := fiapFetch(connectionURL, keys, option)
	if err != nil {
		err = errors.Wrap(err, "fiapFetch error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, "", nil ,err
	}

	pointSets, points, cursor, fiapErr, err = processQueryRS(body)
	if err != nil {
		err = errors.Wrap(err, "processQueryRS error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, "", nil, err
	}
	tools.DebugLogPrintf("Debug: FetchOnce end, pointSets: %v, points: %v, cursor: %v\n", pointSets, points, cursor)
	return pointSets, points, cursor, fiapErr, nil
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKeyNoID, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error , err error) {
	tools.DebugLogPrintf("Debug: FetchByIdsWithKey start, connectionURL: %s, key: %#v, ids: %v\n", connectionURL, key, ids)
	if len(ids) == 0 {
		err = errors.New("ids is empty, set at least one id")
		log.Printf("Error: %+v\n", err)
		return nil, nil, nil, err
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
	pointSets, points, fiapErr, err = Fetch(connectionURL, keys, &model.FetchOption{})
	if err != nil {
		err = errors.Wrap(err, "Fetch error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchByIdsWithKey end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

func FetchLatest(connectionURL string, fromDate *time.Time, untilDate *time.Time ,ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.DebugLogPrintf("Debug: FetchLatest start connectionURL: %s, fromDate: %v, untilDate: %v, ids: %v\n", connectionURL, fromDate, untilDate, ids)
	pointSets, points, fiapErr, err = FetchByIdsWithKey(connectionURL, model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMaximum,
		Gteq:            fromDate,
		Lteq:            untilDate,
	}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchLatest end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

func FetchOldest(connectionURL string, fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value),fiapErr *model.Error, err error) {
	tools.DebugLogPrintf("Debug: FetchOldest start connectionURL: %s, fromDate: %v, untilDate: %v, ids: %v\n", connectionURL, fromDate, untilDate, ids)
	pointSets, points, fiapErr, err = FetchByIdsWithKey(connectionURL, model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMinimum,
		Gteq:            fromDate,
		Lteq:            untilDate,
	}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchOldest end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

func FetchDateRange(connectionURL string, fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.DebugLogPrintf("Debug: FetchDateRange start, connectionURL: %s, fromDate: %v, untilDate: %v,  ids: %v\n", connectionURL, fromDate, untilDate, ids)
	pointSets, points, fiapErr, err = FetchByIdsWithKey(connectionURL, 
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
		return nil, nil, nil, err
	}
	tools.DebugLogPrintf("Debug: FetchDateRange end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

// processQueryRS はQueryRSを処理し、IDをキーとしたPointSetとPointのmapを返す
func processQueryRS(queryRS *model.QueryRS) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), cursor string, fiapErr *model.Error, err error) {
	tools.DebugLogPrintf("Debug: processQueryRS start, data: %#v\n", queryRS)
	if queryRS.Transport == nil {
		err = errors.New("queryRS.Transport is nil")
		log.Printf("Error: %+v\n", err)
		return nil, nil, "", nil, err
	}
	if queryRS.Transport.Header == nil {
		err = errors.New("queryRS.Transport.Header is nil")
		log.Printf("Error: %+v\n", err)
		return nil, nil, "", nil, nil
	}

	if queryRS.Transport.Header.OK != nil &&
		queryRS.Transport.Body == nil {
		err = errors.New("queryRS.Transport.Body is nil")
		log.Printf("Error: %+v\n", err)
		return nil, nil, "", nil, nil
	}

	if queryRS.Transport.Header.Error != nil {
		fiapErr = queryRS.Transport.Header.Error
		log.Printf("Error: fiap error: %+v\n", fiapErr)
		return nil, nil, "", fiapErr, nil
	}

	// mapの初期化
	pointSets = make(map[string](model.ProcessedPointSet))
	points = make(map[string]([]model.Value))
	
	// BodyにPointSetが返っていれば、それを処理する
	if queryRS.Transport.Body.PointSet != nil {
		tools.DebugLogPrintln("Debug: processQueryRS, pointSet is not nil")
		// pointSetsを初期化
		// PointSetの数だけ処理を繰り返す
		for _, ps := range queryRS.Transport.Body.PointSet {
			proccessed := model.ProcessedPointSet{}
			proccessed.PointSetID = ps.PointSetId
			proccessed.PointID = ps.PointId

			// pointSetsのkeyが既に設定されていた場合はデータを追加する
			if existingPointSet, ok := pointSets[ps.Id]; ok {
				proccessed.PointSetID = append(existingPointSet.PointSetID, proccessed.PointSetID...)
				proccessed.PointID = append(existingPointSet.PointID, proccessed.PointID...)
			}
			// pointSetsのkeyが設定されていない場合にはデータを加工せず代入する
			pointSets[ps.Id] = proccessed
		}
	}

	// BodyにPointが返っていれば、それを処理する
	if queryRS.Transport.Body.Point != nil {
		// pointsを初期化
		points = make(map[string]([]model.Value))
		// Pointの数だけ処理を繰り返す
		for _, p := range queryRS.Transport.Body.Point {
			tempValues := p.Value
			// pointsのkeyが設定されていた場合はデータを追加する
			if existingValues, ok := points[p.Id]; ok {
				tempValues = append(existingValues, p.Value...)
			}
			// pointsのkeyが設定されていない場合にはデータを加工せず代入する
			points[p.Id] = tempValues
		}
	}

	cursor = queryRS.Transport.Header.Query.Cursor

	tools.DebugLogPrintf("Debug: processQueryRS end, pointSets: %v, points: %v, cursor: %s\n", pointSets, points, cursor)
	return pointSets, points, cursor, nil, nil
}