package fiap

import (
	"net/http"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/cockroachdb/errors"
)

func Fetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.DebugLogPrintf("Fetch start, connectionURL: %s, keys: %v, option: %#v\n", connectionURL, keys, option)

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
			tools.ErrorLogPrintf("%+v\n", err)
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
	tools.DebugLogPrintf("Fetch end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, err
}

func FetchOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), cursor string, fiapErr *model.Error ,err error) {
	tools.DebugLogPrintf("FetchOnce start, connectionURL: %s, keys: %v, option: %#v\n", connectionURL, keys, option)

	httpResponse, body, err := fiapFetch(connectionURL, keys, option)
	if err != nil {
		err = errors.Wrap(err, "fiapFetch error")
		tools.ErrorLogPrintf("%+v\n", err)
		return nil, nil, "", nil ,err
	}

	pointSets, points, cursor, fiapErr, err = processQueryRS(httpResponse,body)
	if err != nil {
		err = errors.Wrap(err, "processQueryRS error")
		tools.ErrorLogPrintf("%+v\n", err)
		return nil, nil, "", nil, err
	}
	tools.DebugLogPrintf("FetchOnce end, pointSets: %v, points: %v, cursor: %v\n", pointSets, points, cursor)
	return pointSets, points, cursor, fiapErr, nil
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKeyNoID, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error , err error) {
	tools.DebugLogPrintf("FetchByIdsWithKey start, connectionURL: %s, key: %#v, ids: %v\n", connectionURL, key, ids)
	if len(ids) == 0 {
		err = errors.New("ids is empty, set at least one id")
		tools.ErrorLogPrintf("%+v\n", err)
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
		tools.ErrorLogPrintf("%+v\n", err)
		return nil, nil, nil, err
	}
	tools.DebugLogPrintf("FetchByIdsWithKey end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

func FetchLatest(connectionURL string, fromDate *time.Time, untilDate *time.Time ,ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.DebugLogPrintf("FetchLatest start connectionURL: %s, fromDate: %v, untilDate: %v, ids: %v\n", connectionURL, fromDate, untilDate, ids)
	pointSets, points, fiapErr, err = FetchByIdsWithKey(connectionURL, model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMaximum,
		Gteq:            fromDate,
		Lteq:            untilDate,
	}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		tools.ErrorLogPrintf("%+v\n", err)
		return nil, nil, nil, err
	}
	tools.DebugLogPrintf("FetchLatest end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

func FetchOldest(connectionURL string, fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value),fiapErr *model.Error, err error) {
	tools.DebugLogPrintf("FetchOldest start connectionURL: %s, fromDate: %v, untilDate: %v, ids: %v\n", connectionURL, fromDate, untilDate, ids)
	pointSets, points, fiapErr, err = FetchByIdsWithKey(connectionURL, model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMinimum,
		Gteq:            fromDate,
		Lteq:            untilDate,
	}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		tools.ErrorLogPrintf("%+v\n", err)
		return nil, nil, nil, err
	}
	tools.DebugLogPrintf("FetchOldest end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

func FetchDateRange(connectionURL string, fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.DebugLogPrintf("FetchDateRange start, connectionURL: %s, fromDate: %v, untilDate: %v,  ids: %v\n", connectionURL, fromDate, untilDate, ids)
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
		tools.ErrorLogPrintf("%+v\n", err)
		return nil, nil, nil, err
	}
	tools.DebugLogPrintf("FetchDateRange end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

// processQueryRS はQueryRSを処理し、IDをキーとしたPointSetとPointのmapを返す
func processQueryRS(httpResponse *http.Response, queryRS *model.QueryRS) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), cursor string, fiapErr *model.Error, err error) {
	tools.DebugLogPrintf("processQueryRS start, data: %#v\n", queryRS)
	if queryRS.Transport == nil {
		err = errors.Newf("queryRS.Transport is nil, http status: %d", httpResponse.StatusCode)
		tools.ErrorLogPrintf("aaa", err)
		return nil, nil, "", nil, err
	}
	if queryRS.Transport.Header == nil {
		err = errors.Newf("queryRS.Transport.Header is nil, http status: %d", httpResponse.StatusCode)
		tools.ErrorLogPrintf("%+v\n", err)
		return nil, nil, "", nil, err
	}
	if queryRS.Transport.Header.OK != nil &&
		queryRS.Transport.Body == nil {
		err = errors.Newf("queryRS.Transport.Body is nil, http status: %d", httpResponse.StatusCode)
		tools.ErrorLogPrintf("%+v\n", err)
		return nil, nil, "", nil, err
	}
	if queryRS.Transport.Header.Error != nil {
		fiapErr = queryRS.Transport.Header.Error
		return nil, nil, "", fiapErr, nil
	}

	// mapの初期化
	pointSets = make(map[string](model.ProcessedPointSet))
	points = make(map[string]([]model.Value))
	
	// BodyにPointSetが返っていれば、それを処理する
	if queryRS.Transport.Body.PointSet != nil {
		tools.DebugLogPrintf("processQueryRS, pointSet is not nil")
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
		// Pointの数だけ処理を繰り返す
		for _, p := range queryRS.Transport.Body.Point {
			tempValues := p.Value
			if tempValues == nil {
				tempValues = []model.Value{}
			}
			// pointsのkeyが設定されていた場合はデータを追加する
			if existingValues, ok := points[p.Id]; ok {
				tempValues = append(existingValues, tempValues...)
			}
			// pointsのkeyが設定されていない場合にはデータを加工せず代入する
			points[p.Id] = tempValues
		}
	}

	cursor = queryRS.Transport.Header.Query.Cursor

	tools.DebugLogPrintf("processQueryRS end, pointSets: %v, points: %v, cursor: %s\n", pointSets, points, cursor)
	return pointSets, points, cursor, nil, nil
}