package fiap

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"
	"regexp"

	"github.com/google/uuid"
	"github.com/globusdigital/soap"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)

// TODO この関数は、fiapservice.go の動作テストのために作成したものです。後で削除します。
func TestGoWsdl() {
	client := soap.NewClient("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",nil)
	service := NewFIAPServiceSoap(client)

	var uuidObj uuid.UUID
	uuidObj, _ = uuid.NewRandom()

	var storage QueryType = QueryTypeStorage
	var val PositiveInteger = 1000
	var attrTime AttrNameType = AttrNameTypeTime
	var selectType SelectType = SelectTypeMaximum

	header := &Header{
		Query: &Query{
			Id:             Uuidp(uuidObj),
			AcceptableSize: &val,
			Type:           &storage,
			Key: []*Key{
				{
					Id:       "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					AttrName: &attrTime,
					Select:   &selectType,
				},
			},
		},
	}

	queryRQ := &QueryRQ{
		Transport: &Transport{
			Header: header,
		},
	}

	res, err := service.Query(queryRQ)
	if err != nil {
		log.Fatalf("could't get point data: %v", err)
	}
	log.Println(*res)
	log.Printf("%#v\n", *res)

	xmlBytes, err := xml.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(xmlBytes))
}


func Fetch(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), err error) {
	// ...
	return
}

func FetchRaw(connectionURL string, keys []model.UserInputKey, option *model.FetchOption) (raw string, err error) {
	// ...
	return
}

func FetchOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), cursor string, err error) {
	// ...
	return
}

func FetchRawOnce(connectionURL string, keys []model.UserInputKey, option *model.FetchOnceOption) (raw string, cursor string, err error) {
	// クライアントを作成
	client := soap.NewClient(connectionURL, nil)
	service := NewFIAPServiceSoap(client)

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
		return "", "", fmt.Errorf("connectionURL is empty")
	}
	if !regexp.MustCompile(`^https?://`).Match([]byte(connectionURL)) {
		return "", "", fmt.Errorf("invalid connectionURL: %s", connectionURL)
	}
	if len(keys) == 0 {
		return "", "", fmt.Errorf("keys is empty")
	}
	for _, key := range keys {
		if key.ID == "" {
			return "", "", fmt.Errorf("keys.ID is empty")
		}
	}
	if option.Cursor != nil && !IsUUID(option.Cursor) {
		return "", "", fmt.Errorf("cursor must be entered in UUID format. example: '123e4567-e89b-12d3-a456-426614174000'")
	}

	// クエリを作成
	queryRQ := CreateQueryRQ(val, option, keys)

	res, err := service.Query(queryRQ)
	
	// エラーがあればログを出力して終了
	if err != nil {
		log.Fatalf("couldn't get point data: %v", err)
		return "", "", err
	}

	xmlBytes, err := xml.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatalf("couldn't parse xml data: %v", err)
		return "", "", err
	}

	// デバック用に出力
	fmt.Println(string(xmlBytes))

	// カーソルがある場合はrawと空文字カーソルを返す
	if res.Transport.Header.Query.Cursor != nil {
		return string(xmlBytes), string(*res.Transport.Header.Query.Cursor), err
		// カーソルがない場合はrawと空文字を返す
	} else {
		return string(xmlBytes), "", err
	}
}

func FetchLatest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	// ...
	return
}

func FetchOldest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	// ...
	return
}

func FetchDateRange(connectionURL string, fromDate time.Time, untilDate time.Time, option *model.FetchOption) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), err error) {
	// ...
	return
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKey, option *model.FetchOption, ids ...string) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), err error) {
	// ...
	return
}
