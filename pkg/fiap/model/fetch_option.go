package model

/*
FetchOption is type for Fetch option.

FetchOptionは、Fetchのオプションの型です。Fetch関数のoptionの型として使用します。

AccetableSizeは、fiapのqueryクラス内のacceptableSizeに対応し、一度に受信可能なValueオブジェクトの数を表します。
*/
type FetchOption struct {
	AcceptableSize uint
}
