package fiap

import (
	"regexp"
)

func IsUUID(uuid *string) bool {
	r := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return r.MatchString(*uuid)
}
