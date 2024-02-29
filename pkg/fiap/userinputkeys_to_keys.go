package fiap

import (
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"	
)

func UserInputKeyToKeyp(uk model.UserInputKey) *Key {
	k := Key{
		Id: AnyURI(uk.ID),
		AttrName: AttrNameTypep(AttrNameTypeTime),
		Eq: CheckAndConvertTime(uk.Eq),
		Neq: CheckAndConvertTime(uk.Neq),
		Lt: CheckAndConvertTime(uk.Lt),
		Gt: CheckAndConvertTime(uk.Gt),
		Lteq: CheckAndConvertTime(uk.Lteq),
		Gteq: CheckAndConvertTime(uk.Gteq),
		Select: SelectTypep(SelectType(uk.MinMaxIndicator)),
	}
	return &k
}

func UserInputKeysToKeysp(uk []model.UserInputKey) []*Key {
	var keys []*Key
	for _, k := range uk {
		keys = append(keys, UserInputKeyToKeyp(k))
	}
	return keys
}