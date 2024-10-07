package model

import (
	"encoding/xml"
	"time"
)

/*
SelectType is a type for the select attribute of the Key.

SelectType は Key の select 属性のための型です。

この型は、SelectTypeMinimum、SelectTypeMaximum、SelectTypeNoneの3つの定数を持ちます。
この型の値を指定する場合は、それらの定数を使用してください。
*/
type SelectType string

/*
SelectTypeMinimum is a constant of SelectType.

SelectTypeMinimum は SelectType型の定数です。
	
Keyのselect属性に最小値を指定する場合は、SelectTypeMinimumを使用してください。
*/
const SelectTypeMinimum SelectType = "minimum"

/*
SelectTypeMaximum is a constant of SelectType.

SelectTypeMaximum は SelectType型の定数です。

Keyのselect属性に最大値を指定する場合は、SelectTypeMaximumを使用してください。
*/
const SelectTypeMaximum SelectType = "maximum"

/*
SelectTypeNone is a constant of SelectType.

SelectTypeNone は SelectType型の定数です。

Keyのselect属性に何も指定しない場合は、SelectTypeNoneを使用してください。
*/
const SelectTypeNone SelectType = ""


/*
Key is a type used for the Key attribute of Query.

Key は Query の Key 属性に使われる型です。

この型は、FIAPで取得するデータの範囲条件を指定するために使用します。
型内の各フィールドは、基本的にFIAPのkeyクラスの属性に対応していますが、TRAP手順でのみ使用されるtrap属性は省略されています。
*/
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

/*
Query is a type used for the Query attribute of Header.

Query は Header の Query 属性に使われる型です。

この型は、FIAPで行うクエリの内容を表現するために使用します。
型内の各フィールドは、基本的にFIAPのqueryクラスの属性に対応していますが、TRAP手順を利用するための属性は省略されています。

省略されている属性: type, ttl, callbackData, callbackControl
*/
type Query struct {
	Key []Key `xml:"key,omitempty" json:"key,omitempty"`

	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`

	Type string `xml:"type,attr,omitempty" json:"type,omitempty"`

	Cursor string `xml:"cursor,attr,omitempty" json:"cursor,omitempty"`

	AcceptableSize uint `xml:"acceptableSize,attr,omitempty" json:"acceptable_size,omitempty"`
}

/*
Error is a type used for the Error attribute of Header.

Error は Header の Error 属性に使われる型です。

この型は、FIAPでエラーが発生した際にエラーの内容を受け取るために使用します。
型内の各フィールドは、FIAPのerrorクラスの属性に対応しています。
*/
type Error struct {
	Type string `xml:"type,attr,omitempty" json:"type,omitempty"`

	Value string `xml:",chardata" json:"value"`
}

/*
OK is a type used for the OK attribute of Header.

OK は Header の OK 属性に使われる型です。

この型は、FIAP通信が正常に処理された場合に、成功したことを示すために使用します。FIAPの仕様に従い、この型は空の構造体です。
*/
type OK struct {
}

/*
Header is a type used for the Header attribute of Transport.

Header は Transport の Header 属性に使われる型です。

この型は、FIAP通信のヘッダ部をまとめるために使用します。型内の各フィールドは、FIAPのheaderクラスの属性に対応しています。
*/
type Header struct {
	OK *OK `xml:"OK,omitempty" json:"ok,omitempty"`

	Error *Error `xml:"error,omitempty" json:"error,omitempty"`

	Query *Query `xml:"query,omitempty" json:"query,omitempty"`
}

/*
Value is a type used for the Value attribute of Point.

Value は Point の Value 属性に使われる型です。

この型は、時系列データの値を表現するために使用します。型内の各フィールドは、FIAPのvalueクラスの属性に対応しています。
*/
type Value struct {
	Time time.Time `xml:"time,attr,omitempty" json:"time,omitempty"`

	Value string `xml:",chardata" json:"value"`
}

/*
Point is a type used for the Point attribute of Body.

Point は Body の Point 属性に使われる型です。

この型は、ポイントを表現するために使用します。型内の各フィールドは、FIAPのpointクラスの属性に対応しています。
*/
type Point struct {
	Value []Value `xml:"value,omitempty" json:"value,omitempty"`

	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
}

/*
PointSet is a type representing a point set.

PointSet はポイントセットを表す型です。

この型は、FIAPのpointSetクラスに対応した型を、PointSetIdとPointIdの2つのフィールドに分割して扱いやすくしたものです。
*/
type PointSet struct {
	PointSetId []string `json:"point_set_id,omitempty"`

	PointId []string `json:"point_id,omitempty"`

	Id string `json:"id,omitempty"`
}

/*
UnmarshalXML is a custom unmarshal function for PointSet.

PointSet に対するカスタムのUnmarshal関数です。

FIAP通信でOriginalPointSet型として受け取ったデータをPointSet型に変換します。
*/
func (p *PointSet) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	aux := new(OriginalPointSet)

	if err := d.DecodeElement(&aux, &start); err != nil {
		return err
	}
	p.Id = aux.Id

	// pointSetIdとpointIdを初期化
	p.PointSetId = make([]string, 0)
	p.PointId = make([]string, 0)

	//pointSetがnilでない場合、pointSetIDとpointIDをstringの配列として格納
	if aux.PointSet != nil {
		// 受け取ったポイントセットをループさせる
		for _, pointSet := range aux.PointSet {
			p.PointSetId = append(p.PointSetId, pointSet.Id)
		}
	}
	if aux.Point != nil {
		// 受け取ったポイントをループさせる
		for _, point := range aux.Point {
			p.PointId = append(p.PointId, point.Id)
		}
	}
	return nil
}

/*
OriginalPointSet is a type used for the PointSet attribute of Body.

OriginalPointSet は Body の PointSet 属性に使われる型です。

この型は、pointSetをFIAP通信で扱うために使用します。型内の各フィールドは、FIAPのpointSetクラスの属性に対応しています。
*/
type OriginalPointSet struct {
	PointSet []*OriginalPointSet `xml:"pointSet,omitempty" json:"point_set,omitempty"`

	Point []*Point `xml:"point,omitempty" json:"point,omitempty"`

	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
}

/*
Body is a type used for the Body attribute of Transport.

Body は Transport の Body 属性に使われる型です。

この型は、FIAP通信のボディ部をまとめるために使用します。型内の各フィールドは、FIAPのbodyクラスの属性に対応しています。
*/
type Body struct {
	PointSet []*PointSet `xml:"pointSet,omitempty" json:"point_set,omitempty"`

	Point []*Point `xml:"point,omitempty" json:"point,omitempty"`
}

/*
Transport represents a type for Transport.

Transport は トランスポート部 を表すための型です。

この型は、FIAP通信のヘッダ部とボディ部をトランスポート部としてまとめるために使用します。
型内の各フィールドは、FIAPのtransportクラスの属性に対応しています。
*/
type Transport struct {
	XMLName xml.Name `xml:"http://gutp.jp/fiap/2009/11/ transport"`

	Header *Header `xml:"header,omitempty" json:"header,omitempty"`

	Body *Body `xml:"body,omitempty" json:"body,omitempty"`
}

/*
QueryRQ is a type used for sending requests with the soap package.

QueryRQは、soapパッケージでリクエストを送信する際に使用する型です。
*/
type QueryRQ struct {
	XMLName xml.Name `xml:"http://soap.fiap.org/ queryRQ"`

	Transport *Transport `xml:"transport,omitempty" json:"transport,omitempty"`
}

/*
QueryRS is a type used for receiving responses with the soap package.

QueryRSは、soapパッケージでレスポンスを受信する際に使用する型です。
*/
type QueryRS struct {
	XMLName xml.Name `xml:"http://soap.fiap.org/ queryRS"`

	Transport *Transport `xml:"transport,omitempty" json:"transport,omitempty"`
}