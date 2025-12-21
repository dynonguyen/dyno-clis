package utils

import "github.com/rs/xid"

func GenUniqueStr() string {
	return xid.New().String()
}
