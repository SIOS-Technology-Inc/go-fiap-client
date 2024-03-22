package fiap

import (
	"context"
	"net/http"
	"regexp"

	"github.com/google/uuid"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/globusdigital/soap"
	"github.com/cockroachdb/errors"
)

func fiapFetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (httpResponse *http.Response, resBody *model.QueryRS, err error) {
	tools.DebugLogPrintf("fiapFetch start, connectionURL: %s, keys: %v, option: %v\n", connectionURL, keys, option)
	
	client := soap.NewClient(connectionURL, nil)

	// デフォルト値の設定
	if option == nil {
		option = &model.FetchOnceOption{}
	}
	if option.AcceptableSize == nil {
		*option.AcceptableSize = 1000
	}

	// 入力チェック
	if connectionURL == "" {
		err = errors.New("connectionURL is empty")
		return nil, nil, err
	}
	if !regexp.MustCompile(`^https?://`).Match([]byte(connectionURL)) {
		err = errors.New("invalid connectionURL")
		return nil, nil, err
	}
	if len(keys) == 0 {
		err = errors.New("keys is empty")
		return nil, nil, err
	}
	for _, key := range keys {
		if key.ID == "" {
			err = errors.New("keys.ID is empty")
			return nil, nil, err
		}
	}
	if option.Cursor != nil && !tools.IsUUID(option.Cursor) {
		err = errors.New("cursor must be entered in UUID format. example: '123e4567-e89b-12d3-a456-426614174000'")
		return nil, nil, err
	}

	// クエリを作成
	queryRQ := CreateFetchQueryRQ(option, keys)
	resBody = &model.QueryRS{}

	// クエリを実行
	tools.DebugLogPrintf("fiapFetch, client.Call start, queryRQ: %#v\n", queryRQ)
	httpResponse, err = client.Call(context.Background(), "http://soap.fiap.org/query", queryRQ, resBody)
	tools.DebugLogPrintf("fiapFetch, client.Call end, httpResponse: %#v\n", httpResponse)

	if err != nil {
		err = errors.Wrap(err, "fiapFetch, client.Call error")
		return nil, nil, err
	}

	tools.DebugLogPrintf("fiapFetch end, resBody: %#v\n", resBody)
	return httpResponse, resBody, nil
}

func CreateFetchQueryRQ (option *model.FetchOnceOption, keys []model.UserInputKey) *model.QueryRQ {
	tools.DebugLogPrintf("CreateFetchQueryRQ start, option: %v, keys: %v\n", option, keys)
	var uuidObj uuid.UUID
	uuidObj, _ = uuid.NewRandom()

	val := model.PositiveInteger(*option.AcceptableSize)
	
	queryRQ := &model.QueryRQ{
		Transport: &model.Transport{
			Header: &model.Header{
				Query: &model.Query{
					Id: tools.GoogleUuidToUuidp(uuidObj),
					AcceptableSize: tools.AcceptableSizep(val),
					Type: tools.QueryTypep(model.QueryTypeStorage),
					Cursor: tools.CursorStrpToUuidp(option.Cursor),
					Key: tools.UserInputKeyspToKeysp(keys),
				},
			},
		},
	}
	tools.DebugLogPrintf("CreateFetchQueryRQ end, queryRQ: %#v\n", queryRQ)
	return queryRQ
}
