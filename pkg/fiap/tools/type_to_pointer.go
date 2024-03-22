package tools

import "github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"

func Intp(i int) *int {
	return &i
}

func Stringp(s string) *string {
	return &s
}

func QueryTypep(q model.QueryType) *model.QueryType {
	return &q
}

func AcceptableSizep(p model.PositiveInteger) *model.PositiveInteger {
	return &p
}

func AttrNameTypep(a model.AttrNameType) *model.AttrNameType {
	return &a
}

func SelectTypep(s model.SelectType) *model.SelectType {
	return &s
}