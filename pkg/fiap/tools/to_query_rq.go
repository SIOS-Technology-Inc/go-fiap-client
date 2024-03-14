package tools

import (
	"github.com/google/uuid"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)

func CreateQueryRQ (val model.PositiveInteger, option *model.FetchOnceOption, keys []model.UserInputKey) *model.QueryRQ {
	var uuidObj uuid.UUID
	uuidObj, _ = uuid.NewRandom()
	
	queryRQ := &model.QueryRQ{
		Transport: &model.Transport{
			Header: &model.Header{
				Query: &model.Query{
					Id: GoogleUuidToUuidp(uuidObj),
					AcceptableSize: AcceptableSizep(val),
					Type: QueryTypep(model.QueryTypeStorage),
					Cursor: CursorStrpToUuidp(option.Cursor),
					Key: UserInputKeysToKeysp(keys),
				},
			},
		},
	}
	return queryRQ
}