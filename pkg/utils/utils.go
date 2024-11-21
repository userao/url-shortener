package utils

import (
	"strconv"

	"github.com/deatil/go-encoding/base62"
)

func GetHash(seed uint) string {
	encoding := base62.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	base62 := encoding.EncodeToString([]byte(strconv.Itoa(int(seed))))
	return base62
}
