package fiap

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"
	"github.com/google/uuid"
	"github.com/hooklift/gowsdl/soap"
)

type UserInputKey struct {
	ID 				  		string
	Eq							time.Time
	Neq   					time.Time
	Lt    					time.Time
	Gt							time.Time
	Lteq						time.Time
	Gteq						time.Time
	MinMaxIndicator	string
}


type ProcessedValue struct {
	Time  					time.Time			`json:"time"`
	Value						string				`json:"value"`
}

type ProcessedPoint struct {
	Values					[]Value			`json:"values"`
}

type ProcessedPointSet struct {
	PointSetID			string			`json:"point_set_id"`
	PointID					string			`json:"point_id"`
}

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

func Fetch() {
	// ...
}

func FetchRaw() {
	// ...
}

func FetchOnce(acceptableSize *int, cursor Uuid, keys []UserInputKey) (raw string, err error){
	// ...
	return
}

func FetchRawOnce() {
	// ...
}

func FetchLatest() {
	// ...
}

func FetchOldest() {
	// ...
}

func FetchDateRange() {
	// ...
}

func FetchByIdsWithKey() {
	// ...
}