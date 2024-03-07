package fiap

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/globusdigital/soap"
)

func fiapFetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (httpResponse *http.Response, resBody *QueryRS, err error) {
	// クライアントを作成
	client := soap.NewClient(connectionURL, nil)

	// デフォルト値の設定
	if option == nil {
		option = &model.FetchOnceOption{}
	}
	if option.AcceptableSize == nil {
		*option.AcceptableSize = 1000
	}

	var val PositiveInteger = PositiveInteger(*option.AcceptableSize)

	// 入力チェック
	if connectionURL == "" {
		return nil, nil, fmt.Errorf("connectionURL is empty")
	}
	if !regexp.MustCompile(`^https?://`).Match([]byte(connectionURL)) {
		return nil, nil, fmt.Errorf("invalid connectionURL: %s", connectionURL)
	}
	if len(keys) == 0 {
		return nil, nil, fmt.Errorf("keys is empty")
	}
	for _, key := range keys {
		if key.ID == "" {
			return nil, nil, fmt.Errorf("keys.ID is empty")
		}
	}
	if option.Cursor != nil && !IsUUID(option.Cursor) {
		return nil, nil, fmt.Errorf("cursor must be entered in UUID format. example: '123e4567-e89b-12d3-a456-426614174000'")
	}

	// クエリを作成
	queryRQ := CreateQueryRQ(val, option, keys)
	response := &QueryRS{}

	// クエリを実行
	httpResponse, err = client.Call(context.Background(), "http://soap.fiap.org/query", queryRQ, response)

	// エラーがあればエラーを返す
	if err != nil {
		return nil, nil, err
	}

	// エラーがなければ結果を返す
	return httpResponse, resBody, nil
}