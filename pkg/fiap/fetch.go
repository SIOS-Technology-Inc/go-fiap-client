package fiap

import (
	"net/http"
	"time"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/cockroachdb/errors"
)

/*
Fetcher is an interface for fetching data from the FIAP server.

FetcherはFIAPサーバからデータを取得するためのインターフェースです。

Fetch: 与えられたキーとオプションを使用して、FIAPサーバからデータを取得します。

FetchOnce: Fetchメソッドと同様に、与えられたキーとオプションを使用してデータを取得しますが、一度だけ取得します。また、後続のfetchのためのカーソルも返します。

FetchByIdsWithKey: 指定したキーとIDのセットを使用してデータを取得します。

FetchLatest: 指定された日付範囲とIDセット内の最新データを取得します。

FetchOldest: 指定された日付範囲とIDセット内の最古のデータを取得します。

FetchDateRange: 指定された日付範囲とIDセット内のデータを取得します。
*/
type Fetcher interface {
	Fetch(keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error)
	FetchOnce(keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), cursor string, fiapErr *model.Error, err error)
	FetchByIdsWithKey(key model.UserInputKeyNoID, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error)
	FetchLatest(fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error)
	FetchOldest(fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error)
	FetchDateRange(fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error)
}

/* 
FetchClient is a client struct for fetching data from a FIAP server.

FetchClientはFIAPサーバからデータを取得するためのクライアント構造体です。
*/
type FetchClient struct {
	ConnectionURL string
}

/*
Fetch fetches data from the FIAP server using the provided keys and options.

Fetchは、与えられたキーとオプションを使用してFIAPサーバからデータを取得します。

この関数はFetchOnceメソッドを繰り返し呼び出し、取得したデータを結果のmapに追加します。
そして、FetchOnceの戻り値のcursorが""になるまで、繰り返し処理を行います。

※cursorは、連続した大量のデータを一度に取得できる量を制限し分割して取得する場合に、どの位置までデータを取得したかを示すために使用されます。

引数
 - keys: データの範囲を指定するためのkeyの配列。1つのkeyの条件はAND結合です。複数のkeyを指定すると、OR結合になります。
 - option: オプションの指定は任意です。指定しない場合はnilを設定して下さい。

戻り値
 - pointSets: keysの中で指定したIDをキーとして取得したpointSetIDとPointIDのデータのmap
 - points: keysで指定したIDをキーとして取得した時系列データのmap
 - fiapErr: fiap通信の<error>タグを格納する構造体。タグがない場合はnil
 - err: goのエラー情報を格納する構造体。エラーが発生した場合、スタックトレースを含むエラー情報が返される。エラーがない場合はnil。
 
errの発生条件
  - fetchOnceメソッドでエラーが発生した場合
*/
func (f *FetchClient) Fetch(keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.LogPrintf(tools.LogLevelDebug, "Fetch start, connectionURL: %s, keys: %v, option: %#v\n", f.ConnectionURL, keys, option)

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
		fetchOncePointSets, fetchOncePoints, newCursor, fiapErr, err := f.FetchOnce(keys, fetchOnceOption)
		if err != nil {
			err = errors.Wrapf(err, "FetchOnce error on loop iteration %d", i)
			tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
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
	tools.LogPrintf(tools.LogLevelDebug, "Fetch end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, err
}

/*
FetchOnce fetches data from the FIAP server only once using the provided keys and options.

FetchOnceは、与えられたキーとオプションを使用して、FIAPサーバからデータを一度だけ取得します。

この関数はデータをFIAPサーバーから一度だけ取得し、取得したデータを返します。
一度に取得するデータが多すぎる場合は、後続のfetchのためのカーソルが返されます。

以下は、cursorが返ってきた場合に後続のデータを取得するための具体的なコード例
	// FetchOnceの1度目の呼び出し
	pointSets, points, cursor, fiapErr, err := fetchClient.FetchOnce([]model.UserInputKey{
		{
			ID: "id1",
		}
	}, &model.FetchOnceOption{
		AcceptableSize: 1,
	})

	// cursorが返ってきた場合、FetchOnceをもう一度呼び出し、続きのデータを取得する
	if cursor != "" {
		// cursorが返ってきた場合は、続きのデータを取得する
		pointSets, points, cursor, fiapErr, err = fetchClient.FetchOnce([]model.UserInputKey{
			{
				ID: "id1",
			}
		}, &model.FetchOnceOption{
			AcceptableSize: 1,
			// FetchOnceの1度目の呼び出しで取得したcursorを指定する
			Cursor: cursor,
		})
	}
	
	// もしさらにcursorが返ってきた場合は、同様に続きのデータを取得する
		
引数
 - keys: データの範囲を指定するためのkeyの配列。1つのkeyの条件はAND結合です。複数のkeyを指定すると、OR結合になります。
 - option: オプションの指定は任意です。指定しない場合はnilを設定して下さい。

戻り値
 - pointSets: keysの中で指定したIDをキーとして取得したpointSetIDとPointIDのデータのmap
 - points: keysで指定したIDをキーとして取得した時系列データのmap
 - cursor: 後続のfetchのためのカーソル。fetchOnceで取得するデータの量が多すぎる場合にカーソルが返されます。データを最後まで取得できた場合は""が返されます。
 - fiapErr: fiap通信の<error>タグを格納する構造体。タグがない場合はnil
 - err: goのエラー情報を格納する構造体。エラーが発生した場合、スタックトレースを含むエラー情報が返される。エラーがない場合はnil。

fiapErrの発生条件
 - processQueryRSでエラーが発生しqueryRS.Transport.Header.Errorがnilでない場合、SOAP通信は成功したがFIAP通信が失敗したことを示すqueryRS.Transport.Header.Errorの情報をfiapErrとして返す

errの発生条件
 - レシーバfで設定したconnectionURLが http:// または https:// で始まっていない場合(fiapFetch内でエラー)
 - メソッドの引数のkeysの長さが0の場合(fiapFetch内でエラー)
 - メソッドの引数のkeys.IDが空の場合(fiapFetch内でエラー)
 - soap通信を行うclient.Callメソッドでエラーが発生した場合(fiapFetch内でエラー)
 - queryRS.Transportがnilの場合(processQueryRS内でエラー): データが取得できていないためエラーとし、その原因を特定するためにhttp status codeを表示する
 - queryRS.Transport.Headerがnilの場合(processQueryRS内でエラー): SOAP通信に成功した場合はHeader内にokまたはerrorが格納されるためHeaderがnilの場合はエラーとし、その原因を特定するためhttp status codeを表示する
 - queryRS.Transport.Header.OKがnilでなく、queryRS.Transport.Bodyがnilの場合(processQueryRS内でエラー): SOAP通信に成功した場合はBody内にデータが格納されるためBodyがnilの場合はエラーとし、その原因を特定するためにhttp status codeを表示する
*/
func (f *FetchClient) FetchOnce(keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), cursor string, fiapErr *model.Error, err error) {
	tools.LogPrintf(tools.LogLevelDebug, "FetchOnce start, connectionURL: %s, keys: %v, option: %#v\n", f.ConnectionURL, keys, option)

	httpResponse, body, err := fiapFetch(f.ConnectionURL, keys, option)
	if err != nil {
		err = errors.Wrap(err, "fiapFetch error")
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, "", nil, err
	}

	pointSets, points, cursor, fiapErr, err = processQueryRS(httpResponse, body)
	if err != nil {
		err = errors.Wrap(err, "processQueryRS error")
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, "", nil, err
	}
	tools.LogPrintf(tools.LogLevelDebug, "FetchOnce end, pointSets: %v, points: %v, cursor: %v\n", pointSets, points, cursor)
	return pointSets, points, cursor, fiapErr, nil
}

/*
FetchByIdsWithKey fetches data using the specified key and set of IDs.

FetchByIdsWithKeyは、指定されたキーとIDのセットを使用してデータを取得します。

この関数はFetchメソッドを使用し、指定されたキーとIDのセットのデータを取得します。

以下は、FetchByIdsWithKeyの呼び出しの例と、それと同じ結果を返すFetchメソッドの呼び出しの例を示します。
	// FetchByIdsWithKeyの呼び出しの例
	pointSets, points, fiapErr, err := fetchClient.FetchByIdsWithKey(model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeNone,
	}, "id1", "id2", "id3")

	// FetchByIdsWithKeyの呼び出しの例と同じ結果を返すFetchメソッドの呼び出しの例
	pointSets, points, fiapErr, err := fetchClient.Fetch([]model.UserInputKey{
		{ID: "id1"},
		{ID: "id2"},
		{ID: "id3"},
	}, &model.FetchOption{})

引数
 - key: データの範囲を指定するためのkey。idsで指定したすべてのIDに適用される。key内の条件を組み合わせると、AND結合になる。
 - ids: データのID、複数のIDをカンマ区切りで指定可能

戻り値
 - pointSets: keysの中で指定したIDをキーとして取得したpointSetIDとPointIDのデータのmap
 - points: keysで指定したIDをキーとして取得した時系列データのmap
 - fiapErr: fiap通信の<error>タグを格納する構造体。タグがない場合はnil
 - err: goのエラー情報を格納する構造体。エラーが発生した場合、スタックトレースを含むエラー情報が返される。エラーがない場合はnil。

errの発生条件
 - Fetchメソッドでエラーが発生した場合
*/
func (f *FetchClient) FetchByIdsWithKey(key model.UserInputKeyNoID, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.LogPrintf(tools.LogLevelDebug, "FetchByIdsWithKey start, connectionURL: %s, key: %#v, ids: %v\n", f.ConnectionURL, key, ids)
	if len(ids) == 0 {
		err = errors.New("ids is empty, set at least one id")
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
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
	pointSets, points, fiapErr, err = f.Fetch(keys, &model.FetchOption{})
	if err != nil {
		err = errors.Wrap(err, "Fetch error")
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, nil, err
	}
	tools.LogPrintf(tools.LogLevelDebug, "FetchByIdsWithKey end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

/*
FetchLatest fetches the latest data within a specified date range and set of IDs.

FetchLatestは、指定された日付範囲とIDセット内の最新データを取得します。

この関数は、FetchByIdWithKeyを使用して、指定された日付範囲とIDセット内の最新データを取得します。

以下に、FetchLatestの呼び出しの例と、それと同じ結果を返すFetchByIdWithKeyの呼び出しの例を示します。
	// FetchLatestの呼び出しの例
	pointSets, points, fiapErr, err := fetchClient.FetchLatest(&fromDate, &untilDate, "id1", "id2", "id3")

	// FetchLatestの呼び出しの例と同じ結果を返すFetchByIdWithKeyの呼び出しの例
	pointSets, points, fiapErr, err := fetchClient.FetchByIdWithKey(model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMaximum,
		Gteq:            &fromDate,
		Lteq:            &untilDate,
	}, "id1", "id2", "id3") 

引数
 - fromDate: fetchするデータの開始日時
 - untilDate: fetchするデータの終了日時
 - ids: fetchするデータのID、複数のIDをカンマ区切りで指定可能

戻り値
 - pointSets: keysの中で指定したIDをキーとして取得したpointSetIDとPointIDのデータのmap
 - points: keysで指定したIDをキーとして取得した時系列データのmap
 - fiapErr: fiap通信の<error>タグを格納する構造体。タグがない場合はnil
 - err: goのエラー情報を格納する構造体。エラーが発生した場合、スタックトレースを含むエラー情報が返される。エラーがない場合はnil。

errの発生条件
 - FetchByIdWithKeyでエラーが発生した場合
*/
func (f *FetchClient) FetchLatest(fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.LogPrintf(tools.LogLevelDebug, "FetchLatest start connectionURL: %s, fromDate: %v, untilDate: %v, ids: %v\n", f.ConnectionURL, fromDate, untilDate, ids)
	pointSets, points, fiapErr, err = f.FetchByIdsWithKey(model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMaximum,
		Gteq:            fromDate,
		Lteq:            untilDate,
	}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, nil, err
	}
	tools.LogPrintf(tools.LogLevelDebug, "FetchLatest end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

/*
FetchOldest fetches the oldest data within a specified date range and set of IDs.

FetchOldestは、指定された日付範囲とIDセット内の最古のデータを取得します。

この関数は、FetchByIdWithKeyを使用して、指定された日付範囲とIDセット内の最古のデータを取得します。

以下に、FetchOldestの呼び出しの例と、それと同じ結果を返すFetchByIdWithKeyの呼び出しの例を示します。
	// FetchOldestの呼び出しの例
	pointSets, points, fiapErr, err := fetchClient.FetchOldest(&fromDate, &untilDate, "id1", "id2", "id3")

	// FetchOldestの呼び出しの例と同じ結果を返すFetchByIdWithKeyの呼び出しの例
	pointSets, points, fiapErr, err := fetchClient.FetchByIdWithKey(model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMinimum,
		Gteq:            &fromDate,
		Lteq:            &untilDate,
	}, "id1", "id2", "id3")

引数
 - fromDate: fetchするデータの開始日時
 - untilDate: fetchするデータの終了日時
 - ids: fetchするデータのID、複数のIDをカンマ区切りで指定可能

戻り値
 - pointSets: keysの中で指定したIDをキーとして取得したpointSetIDとPointIDのデータのmap
 - points: keysで指定したIDをキーとして取得した時系列データのmap
 - fiapErr: fiap通信の<error>タグを格納する構造体。タグがない場合はnil
 - err: goのエラー情報を格納する構造体。エラーが発生した場合、スタックトレースを含むエラー情報が返される。エラーがない場合はnil。

errの発生条件
 - FetchByIdWithKeyでエラーが発生した場合
*/
func (f *FetchClient) FetchOldest(fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.LogPrintf(tools.LogLevelDebug, "FetchOldest start connectionURL: %s, fromDate: %v, untilDate: %v, ids: %v\n", f.ConnectionURL, fromDate, untilDate, ids)
	pointSets, points, fiapErr, err = f.FetchByIdsWithKey(model.UserInputKeyNoID{
		MinMaxIndicator: model.SelectTypeMinimum,
		Gteq:            fromDate,
		Lteq:            untilDate,
	}, ids...)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, nil, err
	}
	tools.LogPrintf(tools.LogLevelDebug, "FetchOldest end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

/*
FetchDateRange fetches data within a specified date range and set of IDs.

FetchDateRangeは、指定された日付範囲とIDセット内のデータを取得します。

この関数はFetchByIdWithKeyを使用し、指定された日付範囲とIDセット内のデータを取得します。

以下に、FetchDateRangeの呼び出しの例と、それと同じ結果を返すFetchByIdWithKeyの呼び出しの例を示します。
	// FetchDateRangeの呼び出しの例
	pointSets, points, fiapErr, err := fetchClient.FetchDateRange(&fromDate, &untilDate, "id1", "id2", "id3")
	// FetchDateRangeの呼び出しの例と同じ結果を返すFetchByIdWithKeyの呼び出しの例
	pointSets, points, fiapErr, err := fetchClient.FetchByIdWithKey(model.UserInputKeyNoID{
		Gteq:            &fromDate,
		Lteq:            &untilDate,
	}, "id1", "id2", "id3")

引数
 - fromDate: fetchするデータの開始日時
 - untilDate: fetchするデータの終了日時
 - ids: fetchするデータのID、複数のIDをカンマ区切りで指定可能

戻り値
 - pointSets: keysの中で指定したIDをキーとして取得したpointSetIDとPointIDのデータのmap
 - points: keysで指定したIDをキーとして取得した時系列データのmap
 - fiapErr: fiap通信の<error>タグを格納する構造体。タグがない場合はnil
 - err: goのエラー情報を格納する構造体。エラーが発生した場合、スタックトレースを含むエラー情報が返される。エラーがない場合はnil。

errの発生条件
 - FetchByIdWithKeyでエラーが発生した場合
*/
func (f *FetchClient) FetchDateRange(fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), fiapErr *model.Error, err error) {
	tools.LogPrintf(tools.LogLevelDebug, "FetchDateRange start, connectionURL: %s, fromDate: %v, untilDate: %v,  ids: %v\n", f.ConnectionURL, fromDate, untilDate, ids)
	pointSets, points, fiapErr, err = f.FetchByIdsWithKey(
		model.UserInputKeyNoID{
			Gteq:            fromDate,
			Lteq:            untilDate,
			MinMaxIndicator: model.SelectTypeNone,
		},
		ids...,
	)
	if err != nil {
		err = errors.Wrap(err, "FetchByIdsWithKey error")
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, nil, err
	}
	tools.LogPrintf(tools.LogLevelDebug, "FetchDateRange end, pointSets: %v, points: %v\n", pointSets, points)
	return pointSets, points, fiapErr, nil
}

// processQueryRS はQueryRSを処理し、IDをキーとしたPointSetとPointのmapを返す
func processQueryRS(httpResponse *http.Response, queryRS *model.QueryRS) (pointSets map[string](model.ProcessedPointSet), points map[string]([]model.Value), cursor string, fiapErr *model.Error, err error) {
	tools.LogPrintf(tools.LogLevelDebug, "processQueryRS start, data: %#v\n", queryRS)
	if queryRS.Transport == nil {
		err = errors.Newf("queryRS.Transport is nil, http status: %d", httpResponse.StatusCode)
		tools.LogPrintf(tools.LogLevelError, "aaa", err)
		return nil, nil, "", nil, err
	}
	if queryRS.Transport.Header == nil {
		err = errors.Newf("queryRS.Transport.Header is nil, http status: %d", httpResponse.StatusCode)
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, "", nil, err
	}
	if queryRS.Transport.Header.OK != nil &&
		queryRS.Transport.Body == nil {
		err = errors.Newf("queryRS.Transport.Body is nil, http status: %d", httpResponse.StatusCode)
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
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
		tools.LogPrintf(tools.LogLevelDebug, "processQueryRS, pointSet is not nil")
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

	tools.LogPrintf(tools.LogLevelDebug, "processQueryRS end, pointSets: %v, points: %v, cursor: %s\n", pointSets, points, cursor)
	return pointSets, points, cursor, nil, nil
}
