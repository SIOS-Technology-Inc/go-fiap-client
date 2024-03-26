package fiap

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"github.com/google/uuid"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
	"github.com/cockroachdb/errors"
	"github.com/globusdigital/soap"
)

func fiapFetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (httpResponse *http.Response, resBody *model.QueryRS, err error) {
	tools.DebugLogPrintf("Debug: fiapFetch start, connectionURL: %s, keys: %v, option: %v\n", connectionURL, keys, option)
	
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
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	if !regexp.MustCompile(`^https?://`).Match([]byte(connectionURL)) {
		err = errors.New("invalid connectionURL")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	if len(keys) == 0 {
		err = errors.New("keys is empty")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}
	for _, key := range keys {
		if key.ID == "" {
			err = errors.New("keys.ID is empty")
			log.Printf("Error: %+v\n", err)
			return nil, nil, err
		}
	}
	if option.Cursor != nil && !tools.IsUUID(option.Cursor) {
		tools.DebugLogPrintf("Debug: option.Cursor: %#v\n", *option.Cursor)
		err = errors.New("cursor must be entered in UUID format. example: '123e4567-e89b-12d3-a456-426614174000'")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}

	// クエリを作成
	queryRQ := CreateFetchQueryRQ(option, keys)
	resBody = &model.QueryRS{}

	// クエリを実行
	tools.DebugLogPrintf("Debug: fiapFetch, client.Call start, queryRQ: %#v\n", queryRQ)
	httpResponse, err = client.Call(context.Background(), "http://soap.fiap.org/query", queryRQ, resBody)
	tools.DebugLogPrintf("Debug: fiapFetch, client.Call end, httpResponse: %#v, resBody: %#v\n", httpResponse, resBody)

	if err != nil {
		err = errors.Wrap(err, "fiapFetch, client.Call error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}

	tools.DebugLogPrintf("Debug: fiapFetch end, resBody: %#v\n", resBody)
	return httpResponse, resBody, nil
}

func CreateFetchQueryRQ (option *model.FetchOnceOption, keys []model.UserInputKey) *model.QueryRQ {
	tools.DebugLogPrintf("Debug: CreateFetchQueryRQ start, option: %v, keys: %v\n", option, keys)
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
	tools.DebugLogPrintf("Debug: CreateFetchQueryRQ end, queryRQ: %#v\n", queryRQ)
	return queryRQ
}
