package model

import (
	"encoding/xml"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name


type SelectType string

const (
	SelectTypeMinimum SelectType = "minimum"

	SelectTypeMaximum SelectType = "maximum"

	SelectTypeNone SelectType = ""
)

type TrapType string

type Key struct {
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`

	AttrName string `xml:"attrName,attr,omitempty" json:"attrName,omitempty"`

	Eq string `xml:"eq,attr,omitempty" json:"eq,omitempty"`

	Neq string `xml:"neq,attr,omitempty" json:"neq,omitempty"`

	Lt string `xml:"lt,attr,omitempty" json:"lt,omitempty"`

	Gt string `xml:"gt,attr,omitempty" json:"gt,omitempty"`

	Lteq string `xml:"lteq,attr,omitempty" json:"lteq,omitempty"`

	Gteq string `xml:"gteq,attr,omitempty" json:"gteq,omitempty"`

	Select SelectType `xml:"select,attr,omitempty" json:"select,omitempty"`
}

type NonNegativeInteger uint

// TODO: uintは符号なし整数を表すため、0を許容する
// 0を許容しない場合は、関数の処理に注意が必要
type PositiveInteger uint

type Query struct {
	Key []Key `xml:"key,omitempty" json:"key,omitempty"`

	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`

	Type string `xml:"type,attr,omitempty" json:"type,omitempty"`

	Cursor string `xml:"cursor,attr,omitempty" json:"cursor,omitempty"`

	AcceptableSize uint `xml:"acceptableSize,attr,omitempty" json:"acceptableSize,omitempty"`
}

type Error struct {
	Type string `xml:"type,attr,omitempty" json:"type,omitempty"`
}

type OK struct {
}

type Header struct {
	OK *OK `xml:"OK,omitempty" json:"OK,omitempty"`

	Error *Error `xml:"error,omitempty" json:"error,omitempty"`

	Query *Query `xml:"query,omitempty" json:"query,omitempty"`
}

type Value struct {
	Time time.Time `xml:"time,attr,omitempty" json:"time,omitempty"`

	Value string `xml:",chardata" json:"value"`
}

type Point struct {
	Value []*Value `xml:"value,omitempty" json:"value,omitempty"`

	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
}

// TODO PointSet内のPointSetの型を必要に応じて変更する（stringの配列にする?）
type PointSet struct {
	PointSet []*PointSet `xml:"pointSet,omitempty" json:"pointSet,omitempty"`

	Point []*Point `xml:"point,omitempty" json:"point,omitempty"`

	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
}

type Body struct {
	PointSet []*PointSet `xml:"pointSet,omitempty" json:"pointSet,omitempty"`

	Point []*Point `xml:"point,omitempty" json:"point,omitempty"`
}

type Transport struct {
	Header *Header `xml:"header,omitempty" json:"header,omitempty"`

	Body *Body `xml:"body,omitempty" json:"body,omitempty"`
}

type QueryRQ struct {
	Transport *Transport `xml:"transport,omitempty" json:"transport,omitempty"`
}

type QueryRS struct {
	Transport *Transport `xml:"transport,omitempty" json:"transport,omitempty"`
}