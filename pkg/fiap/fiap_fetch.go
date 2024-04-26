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

var regexpURL = regexp.MustCompile(`^https?://`)

func fiapFetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (httpResponse *http.Response, resBody *model.QueryRS, err error) {
	tools.DebugLogPrintf("Debug: fiapFetch start, connectionURL: %s, keys: %v, option: %v\n", connectionURL, keys, option)

	// 入力チェック
	// if connectionURL == "" {
	// 	err = errors.Newf("connectionURL is empty, %s",	connectionURL)
	// 	log.Printf("Error: %+v\n", err)
	// 	return nil, nil, err
	// }
	if !regexpURL.Match([]byte(connectionURL)) {
		err = errors.Newf("invalid connectionURL: %s", connectionURL)
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
			err = errors.Newf("keys.ID is empty, key: %#v", keys)
			log.Printf("Error: %+v\n", err)
			return nil, nil, err
		}
	}

	client := soap.NewClient(connectionURL, nil)

	// クエリを作成
	queryRQ := newQueryRQ(option, keys)
	resBody = &model.QueryRS{}

	// クエリを実行
	tools.DebugLogPrintf("Debug: fiapFetch, client.Call start, queryRQ: %#v\n", queryRQ)
	httpResponse, err = client.Call(context.Background(), "http://soap.fiap.org/query", queryRQ, resBody)
	tools.DebugLogPrintf("Debug: fiapFetch, client.Call end, httpResponse: %#v, resBody: %#v\n", httpResponse, resBody)

	if err != nil {
		err = errors.Wrap(err, "client.Call error")
		log.Printf("Error: %+v\n", err)
		return nil, nil, err
	}

	tools.DebugLogPrintf("Debug: fiapFetch end, resBody: %#v\n", resBody)
	return httpResponse, resBody, nil
}

func newQueryRQ(option *model.FetchOnceOption, keys []model.UserInputKey) *model.QueryRQ {
	tools.DebugLogPrintf("Debug: CreateFetchQueryRQ start, option: %v, keys: %v\n", option, keys)

	// デフォルト値の設定
	if option == nil {
		option = &model.FetchOnceOption{}
	}

	var uuidObj uuid.UUID
	uuidObj, _ = uuid.NewRandom()

	queryRQ := &model.QueryRQ{
		Transport: &model.Transport{
			Header: &model.Header{
				Query: &model.Query{
					Id:             uuidObj.String(),
					AcceptableSize: option.AcceptableSize,
					Type:           "storage",
					Cursor:         option.Cursor,
					Key:            tools.UserInputKeysToKeys(keys),
				},
			},
		},
	}
	tools.DebugLogPrintf("Debug: CreateFetchQueryRQ end, queryRQ: %#v\n", queryRQ)
	return queryRQ
}
