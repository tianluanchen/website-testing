package services

import (
	"website-testing/pkg"
)

var logger = pkg.NewLogger()

func serializeWithPanic(v any) []byte {
	bs, err := pkg.Serialize(v)
	if err != nil {
		panic(err)
	}
	return bs
}
