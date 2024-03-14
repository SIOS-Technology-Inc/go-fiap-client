package tools

import (
	"github.com/google/uuid"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)

func GoogleUuidToUuidp(uu uuid.UUID) *model.Uuid {
	us := model.Uuid(uu.String())
	return &us
}

func CursorStrpToUuidp(s *string) *model.Uuid {
	if s == nil {
		return nil
	}
	u := model.Uuid(*s)
	return &u
}
