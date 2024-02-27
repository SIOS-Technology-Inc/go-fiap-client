package fiap

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hooklift/gowsdl/soap"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
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
			Id:             &myUuid,
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

func FetchDateRange(connectionURL string, fromDate time.Time, untilDate time.Time, option *model.FetchOption) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), err error) {
	// ...
	return
}

func FetchByIdsWithKey(connectionURL string, key model.UserInputKey, option *model.FetchOption, ids ...string) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), err error) {
	// ...
	return
}
