package fiap

import (
	"github.com/google/uuid"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
)

func CreateQueryRQ (val model.PositiveInteger, option *model.FetchOnceOption, keys []model.UserInputKey) *model.QueryRQ {
	var uuidObj uuid.UUID
	uuidObj, _ = uuid.NewRandom()
	
	queryRQ := &model.QueryRQ{
		Transport: &model.Transport{
			Header: &model.Header{
				Query: &model.Query{
					Id: tools.GoogleUuidToUuidp(uuidObj),
					AcceptableSize: tools.AcceptableSizep(val),
					Type: tools.QueryTypep(model.QueryTypeStorage),
					Cursor: tools.CursorStrpToUuidp(option.Cursor),
					Key: tools.UserInputKeysToKeysp(keys),
				},
			},
		},
	}
	return queryRQ
}