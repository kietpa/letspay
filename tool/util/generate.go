package util

import "github.com/oklog/ulid/v2"

func GenerateReferenceId() string {
	return ulid.Make().String()
}
