package model

/*
FetchOnceOption is type for FetchOnce option.
FetchOnceOptionは、FetchOnceのオプションの型です。

FetchOnce関数のoptionの型として使用します。
FetchOptionとの違いは、連続したデータを取得するためのcursorが含まれていることです。
*/

type FetchOnceOption struct {
	AcceptableSize uint
	Cursor         string
}
