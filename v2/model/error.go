package model

import "github.com/champon1020/argus/v2"

var (
	connectError = &argus.ErrorType{Name: "DbConnectFailed", Msg: "Database cannot be connected"}
)
