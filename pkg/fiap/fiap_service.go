package fiap

import (
	"context"
	"encoding/xml"
	"github.com/hooklift/gowsdl/soap"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type AnyType struct {
	InnerXML string `xml:",innerxml"`
}

type AnyURI string

type NCName string

type Uuid string

type QueryType string

const (
	QueryTypeStorage QueryType = "storage"

	QueryTypeStream QueryType = "stream"
)

type AttrNameType string

const (
	AttrNameTypeTime AttrNameType = "time"

	AttrNameTypeValue AttrNameType = "value"
)

type SelectType string

const (
	SelectTypeMinimum SelectType = "minimum"

	SelectTypeMaximum SelectType = "maximum"
)

type TrapType string

const (
	TrapTypeChanged TrapType = "changed"
)

type Key struct {
	XMLName xml.Name `xml:"http://gutp.jp/fiap/2009/11/ key"`

	Id AnyURI `xml:"id,attr,omitempty" json:"id,omitempty"`

	AttrName *AttrNameType `xml:"attrName,attr,omitempty" json:"attrName,omitempty"`

	Eq string `xml:"eq,attr,omitempty" json:"eq,omitempty"`

	Neq string `xml:"neq,attr,omitempty" json:"neq,omitempty"`

	Lt string `xml:"lt,attr,omitempty" json:"lt,omitempty"`

	Gt string `xml:"gt,attr,omitempty" json:"gt,omitempty"`

	Lteq string `xml:"lteq,attr,omitempty" json:"lteq,omitempty"`

	Gteq string `xml:"gteq,attr,omitempty" json:"gteq,omitempty"`

	Select *SelectType `xml:"select,attr,omitempty" json:"select,omitempty"`

	Trap *TrapType `xml:"trap,attr,omitempty" json:"trap,omitempty"`
}

type NonNegativeInteger uint

// TODO: uintは符号なし整数を表すため、0を許容する
// 0を許容しない場合は、関数の処理に注意が必要
type PositiveInteger uint

type Query struct {
	XMLName xml.Name `xml:"http://gutp.jp/fiap/2009/11/ query"`

	Key []*Key `xml:"key,omitempty" json:"key,omitempty"`

	Id *Uuid `xml:"id,attr,omitempty" json:"id,omitempty"`

	Type *QueryType `xml:"type,attr,omitempty" json:"type,omitempty"`

	Cursor *Uuid `xml:"cursor,attr,omitempty" json:"cursor,omitempty"`

	Ttl *NonNegativeInteger `xml:"ttl,attr,omitempty" json:"ttl,omitempty"`

	AcceptableSize *PositiveInteger `xml:"acceptableSize,attr,omitempty" json:"acceptableSize,omitempty"`

	CallbackData AnyURI `xml:"callbackData,attr,omitempty" json:"callbackData,omitempty"`

	CallbackControl AnyURI `xml:"callbackControl,attr,omitempty" json:"callbackControl,omitempty"`
}

type Error struct {
	Type *string `xml:"type,attr,omitempty" json:"type,omitempty"`
}

type OK struct {
}

type Header struct {
	XMLName xml.Name `xml:"http://gutp.jp/fiap/2009/11/ header"`

	OK *OK `xml:"OK,omitempty" json:"OK,omitempty"`

	Error *Error `xml:"error,omitempty" json:"error,omitempty"`

	Query *Query `xml:"query,omitempty" json:"query,omitempty"`
}

type Value struct {
	XMLName xml.Name `xml:"http://gutp.jp/fiap/2009/11/ value"`

	Time *time.Time `xml:"time,attr,omitempty" json:"time,omitempty"`

	Value string `xml:",chardata" json:"value"`
}

type Point struct {
	XMLName xml.Name `xml:"http://gutp.jp/fiap/2009/11/ point"`

	Value []*Value `xml:"value,omitempty" json:"value,omitempty"`

	Id AnyURI `xml:"id,attr,omitempty" json:"id,omitempty"`
}

type PointSet struct {
	XMLName xml.Name `xml:"http://gutp.jp/fiap/2009/11/ pointSet"`

	PointSet []*PointSet `xml:"pointSet,omitempty" json:"pointSet,omitempty"`

	Point []*Point `xml:"point,omitempty" json:"point,omitempty"`

	Id AnyURI `xml:"id,attr,omitempty" json:"id,omitempty"`
}

type Body struct {
	XMLName xml.Name `xml:"http://gutp.jp/fiap/2009/11/ body"`

	PointSet []*PointSet `xml:"pointSet,omitempty" json:"pointSet,omitempty"`

	Point []*Point `xml:"point,omitempty" json:"point,omitempty"`
}

type Transport struct {
	XMLName xml.Name `xml:"http://gutp.jp/fiap/2009/11/ transport"`

	Header *Header `xml:"header,omitempty" json:"header,omitempty"`

	Body *Body `xml:"body,omitempty" json:"body,omitempty"`
}

type QueryRQ struct {
	XMLName xml.Name `xml:"http://soap.fiap.org/ queryRQ"`

	Transport *Transport `xml:"transport,omitempty" json:"transport,omitempty"`
}

type QueryRS struct {
	XMLName xml.Name `xml:"http://soap.fiap.org/ queryRS"`

	Transport *Transport `xml:"transport,omitempty" json:"transport,omitempty"`
}

type DataRQ struct {
	XMLName xml.Name `xml:"http://soap.fiap.org/ dataRQ"`

	Transport *Transport `xml:"transport,omitempty" json:"transport,omitempty"`
}

type DataRS struct {
	XMLName xml.Name `xml:"http://soap.fiap.org/ dataRS"`

	Transport *Transport `xml:"transport,omitempty" json:"transport,omitempty"`
}

type FIAPServiceSoap interface {
	Query(request *QueryRQ) (*QueryRS, error)

	QueryContext(ctx context.Context, request *QueryRQ) (*QueryRS, error)

	Data(request *DataRQ) (*DataRS, error)

	DataContext(ctx context.Context, request *DataRQ) (*DataRS, error)
}

type fIAPServiceSoap struct {
	client *soap.Client
}

func NewFIAPServiceSoap(client *soap.Client) FIAPServiceSoap {
	return &fIAPServiceSoap{
		client: client,
	}
}

func (service *fIAPServiceSoap) QueryContext(ctx context.Context, request *QueryRQ) (*QueryRS, error) {
	response := new(QueryRS)
	err := service.client.CallContext(ctx, "http://soap.fiap.org/query", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *fIAPServiceSoap) Query(request *QueryRQ) (*QueryRS, error) {
	return service.QueryContext(
		context.Background(),
		request,
	)
}

func (service *fIAPServiceSoap) DataContext(ctx context.Context, request *DataRQ) (*DataRS, error) {
	response := new(DataRS)
	err := service.client.CallContext(ctx, "http://soap.fiap.org/data", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *fIAPServiceSoap) Data(request *DataRQ) (*DataRS, error) {
	return service.DataContext(
		context.Background(),
		request,
	)
}
