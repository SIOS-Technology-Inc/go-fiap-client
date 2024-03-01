package fiap

import (
	"github.com/google/uuid"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)

func CreateQueryRQ (val PositiveInteger, option *model.FetchOnceOption, keys []model.UserInputKey) *QueryRQ {
	var uuidObj uuid.UUID
	uuidObj, _ = uuid.NewRandom()
	
	queryRQ := &QueryRQ{
		Transport: &Transport{
			Header: &Header{
				Query: &Query{
					Id: GoogleUuidToUuidp(uuidObj),
					AcceptableSize: AcceptableSizep(val),
					Type: QueryTypep(QueryTypeStorage),
					Cursor: CursorStrpToUuidp(option.Cursor),
					Key: UserInputKeysToKeysp(keys),
				},
			},
		},
	}
	return queryRQ
}