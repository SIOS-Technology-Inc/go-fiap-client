package model

/*
FetchOnceOption is type for FetchOnce option.

FetchOnceOptionは、FetchOnceのオプションの型です。

FetchOnce関数のoptionの型として使用します。FetchOptionとの違いは、連続したデータを取得するためのcursorが含まれていることです。

AccetableSizeは、fiapのqueryクラス内のacceptableSizeに対応し、一度に受信可能なValueオブジェクトの数を表します。

Cursorは、fiapのqueryクラス内のcursorに対応し、連続したデータを取得するためのポインタを表します。
*/
type FetchOnceOption struct {
	AcceptableSize uint
	Cursor         string
}
