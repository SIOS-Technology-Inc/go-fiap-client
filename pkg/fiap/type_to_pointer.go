package fiap

func Intp(i int) *int {
	return &i
}

func QueryTypep(q QueryType) *QueryType {
	return &q
}

func AcceptableSizep(p PositiveInteger) *PositiveInteger {
	return &p
}

func AttrNameTypep(a AttrNameType) *AttrNameType {
	return &a
}

func SelectTypep(s SelectType) *SelectType {
	return &s
}