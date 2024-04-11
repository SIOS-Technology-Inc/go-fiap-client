package model

import (
	"encoding/xml"
	"time"
)

type SelectType string

const (
	SelectTypeMinimum SelectType = "minimum"

	SelectTypeMaximum SelectType = "maximum"

	SelectTypeNone SelectType = ""
)

type Key struct {
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`

	AttrName string `xml:"attrName,attr,omitempty" json:"attr_name,omitempty"`

	Eq string `xml:"eq,attr,omitempty" json:"eq,omitempty"`

	Neq string `xml:"neq,attr,omitempty" json:"neq,omitempty"`

	Lt string `xml:"lt,attr,omitempty" json:"lt,omitempty"`

	Gt string `xml:"gt,attr,omitempty" json:"gt,omitempty"`

	Lteq string `xml:"lteq,attr,omitempty" json:"lteq,omitempty"`

	Gteq string `xml:"gteq,attr,omitempty" json:"gteq,omitempty"`

	Select SelectType `xml:"select,attr,omitempty" json:"select,omitempty"`
}

type Query struct {
	Key []Key `xml:"key,omitempty" json:"key,omitempty"`

	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`

	Type string `xml:"type,attr,omitempty" json:"type,omitempty"`

	Cursor string `xml:"cursor,attr,omitempty" json:"cursor,omitempty"`

	AcceptableSize uint `xml:"acceptableSize,attr,omitempty" json:"acceptable_size,omitempty"`
}

type Error struct {
	Type string `xml:"type,attr,omitempty" json:"type,omitempty"`
}

type OK struct {
}

type Header struct {
	OK *OK `xml:"OK,omitempty" json:"ok,omitempty"`

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
	PointSetId []*string `json:"point_set_id,omitempty"`

	PointId []*string `json:"point_id,omitempty"`

	Id string `json:"id,omitempty"`
}

func (p *PointSet) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	aux := &OriginalPointSet{}

	if err := d.DecodeElement(aux, &start); err != nil {
		return err
	}
	p.Id = aux.Id
	
	//pointSetがnilでない場合、pointSetIDとpointIDをstringの配列として格納
	if aux.PointSet != nil {
		// 受け取ったポイントセットをループさせる
		for _, pointSet := range aux.PointSet {
			// ポイントセット直下のpointSetをループ処理
			for _, pointSetSecondLayerPointset := range pointSet.PointSet {
				s := pointSetSecondLayerPointset.Id
				p.PointSetId = append(p.PointSetId, &s)
			}
			// ポイントセット直下のpointをループ処理
			for _, pointSetSecondLayerPoint := range pointSet.Point {
				s := pointSetSecondLayerPoint.Id
				p.PointId = append(p.PointId, &s)
			}
		}
	}
	return nil
}

type OriginalPointSet struct {
	PointSet []*OriginalPointSet `xml:"pointSet,omitempty" json:"point_set,omitempty"`

	Point []*Point `xml:"point,omitempty" json:"point,omitempty"`

	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
}

type Body struct {
	PointSet []*PointSet `xml:"pointSet,omitempty" json:"point_set,omitempty"`

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