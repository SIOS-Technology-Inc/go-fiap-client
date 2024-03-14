package tools

import (
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)

func UserInputKeyToKeyp(uk model.UserInputKey) *model.Key {
	k := model.Key{
		Id: model.AnyURI(uk.ID),
		AttrName: AttrNameTypep(model.AttrNameTypeTime),
		Eq: TimeToString(uk.Eq),
		Neq: TimeToString(uk.Neq),
		Lt: TimeToString(uk.Lt),
		Gt: TimeToString(uk.Gt),
		Lteq: TimeToString(uk.Lteq),
		Gteq: TimeToString(uk.Gteq),
		Select: SelectTypep(model.SelectType(uk.MinMaxIndicator)),
	}
	return &k
}

func UserInputKeysToKeysp(uk []model.UserInputKey) []*model.Key {
	var keys []*model.Key
	for _, k := range uk {
		keys = append(keys, UserInputKeyToKeyp(k))
	}
	return keys
}