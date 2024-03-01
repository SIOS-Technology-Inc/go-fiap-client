package fiap

import (
	"github.com/google/uuid"
)

func GoogleUuidToUuidp(uu uuid.UUID) *Uuid {
	us := Uuid(uu.String())
	return &us
}

func CursorStrpToUuidp(s *string) *Uuid {
	if s == nil {
		return nil
	}
	u := Uuid(*s)
	return &u
}
