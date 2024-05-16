package fiap

import (
	"context"
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
	tools.LogPrintf(tools.LogLevelDebug, "fiapFetch start, connectionURL: %s, keys: %v, option: %v\n", connectionURL, keys, option)

	if !regexpURL.Match([]byte(connectionURL)) {
		err = errors.Newf("invalid connectionURL: %s", connectionURL)
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, err
	}
	if len(keys) == 0 {
		err = errors.New("keys is empty")
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, err
	}
	for _, key := range keys {
		if key.ID == "" {
			err = errors.Newf("keys.ID is empty, key: %#v", keys)
			tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
			return nil, nil, err
		}
	}

	client := soap.NewClient(connectionURL, nil)

	// クエリを作成
	queryRQ := newQueryRQ(option, keys)
	resBody = &model.QueryRS{}

	// クエリを実行
	tools.LogPrintf(tools.LogLevelDebug, "fiapFetch, client.Call start, queryRQ: %#v\n", queryRQ)
	httpResponse, err = client.Call(context.Background(), "http://soap.fiap.org/query", queryRQ, resBody)
	tools.LogPrintf(tools.LogLevelDebug, "fiapFetch, client.Call end, httpResponse: %#v, resBody: %#v\n", httpResponse, resBody)

	if err != nil {
		err = errors.Wrap(err, "client.Call error")
		tools.LogPrintf(tools.LogLevelError, "%+v\n", err)
		return nil, nil, err
	}

	tools.LogPrintf(tools.LogLevelDebug, "fiapFetch end, resBody: %#v\n", resBody)
	return httpResponse, resBody, nil
}

func newQueryRQ(option *model.FetchOnceOption, keys []model.UserInputKey) *model.QueryRQ {
	tools.LogPrintf(tools.LogLevelDebug, "CreateFetchQueryRQ start, option: %v, keys: %v\n", option, keys)

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
	tools.LogPrintf(tools.LogLevelDebug, "CreateFetchQueryRQ end, queryRQ: %#v\n", queryRQ)
	return queryRQ
}
