package tools

import (
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)

func UserInputKeyspToKeysp(ks []model.UserInputKey) []*model.Key {
	var keys []*model.Key
	for _, k := range ks {
		key := model.Key{
			Id: model.AnyURI(k.ID),
			AttrName: AttrNameTypep(model.AttrNameTypeTime),
			Eq: TimeToString(k.Eq),
			Neq: TimeToString(k.Neq),
			Lt: TimeToString(k.Lt),
			Gt: TimeToString(k.Gt),
			Lteq: TimeToString(k.Lteq),
			Gteq: TimeToString(k.Gteq),
			Select: SelectTypep(model.SelectType(*k.MinMaxIndicator)),
		}
		keys = append(keys, &key)
	}
	return keys
}