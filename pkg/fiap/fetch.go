package fiap

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"
	"github.com/google/uuid"
	"github.com/hooklift/gowsdl/soap"
)

// TODO この関数は、fiapservice.go の動作テストのために作成したものです。後で削除します。
func TestGoWsdl() {
	client := soap.NewClient("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage")
	service := NewFIAPServiceSoap(client)

	var uuidObj uuid.UUID
	uuidObj, _ = uuid.NewRandom()
	var myUuid Uuid = Uuid(uuidObj.String())

	var storage QueryType = QueryTypeStorage
	var val PositiveInteger = 1000
	var attrTime AttrNameType = AttrNameTypeTime
	var selectType SelectType = SelectTypeMaximum

	header := &Header{
		Query: &Query{
			Id: &myUuid,
			AcceptableSize: &val,
			Type: &storage,
			Key: []*Key{
				{
					Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					AttrName: &attrTime,
					Select: &selectType,
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

// TODO option構造体の使い方について相談し、定義を確定する
func Fetch(connectionURL string, keys []UserInputKey, option *FetchOption) (pointSets map[string]([]ProcessedPointSet), points map[string]([]ProcessedPoint), err error) {
	// ...
	return
}

func FetchRaw(connectionURL string, keys []UserInputKey, option *FetchOption) (raw string, err error) {
	// ...
	return
}

func FetchOnce(connectionURL string, keys []UserInputKey, option *FetchOnceOption) (pointSets map[string]([]ProcessedPointSet), points map[string]([]ProcessedPoint), cursor string, err error) {
	// ...
	return
}

func FetchRawOnce(connectionURL string, keys []UserInputKey, option *FetchOnceOption) (raw string, cursor string, err error) {
	// ...
	return
}

func FetchLatest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	// ...
	return
}

func FetchOldest(connectionURL string, ids ...string) (datas map[string]string, err error) {
	// ...
	return
}

func FetchDateRange(connectionURL string, fromDate time.Time, untilDate time.Time, option *FetchOption) (pointSets map[string]([]ProcessedPointSet), points map[string]([]ProcessedPoint), err error) {
	// ...
	return
}

func FetchByIdsWithKey(connectionURL string, key UserInputKey, option *FetchOption, ids ...string) (pointSets map[string]([]ProcessedPointSet), points map[string]([]ProcessedPoint), err error) {
	// ...
	return
}