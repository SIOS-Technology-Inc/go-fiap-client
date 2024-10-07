package model

import (
	"time"
)

/*
UserInputKey holds information about keys specified by the user.

UserInputKey は、ユーザーが指定するキーの情報を保持する型です。

この型は、Fetch、FetchOnceを呼び出す際に引数に使用されます。
*/
type UserInputKey struct {
	ID              string
	Eq              *time.Time
	Neq             *time.Time
	Lt              *time.Time
	Gt              *time.Time
	Lteq            *time.Time
	Gteq            *time.Time
	MinMaxIndicator SelectType
}

/*
UserInputKeyNoID holds information about keys specified by the user. (without ID)

UserInputKeyNoID は、ユーザーが指定するキーの情報を保持する型です。(IDなし)

この型は、FetchByIdsWithKeyを呼び出す際に引数に使用されます。
*/
type UserInputKeyNoID struct {
	Eq              *time.Time
	Neq             *time.Time
	Lt              *time.Time
	Gt              *time.Time
	Lteq            *time.Time
	Gteq            *time.Time
	MinMaxIndicator SelectType
}
