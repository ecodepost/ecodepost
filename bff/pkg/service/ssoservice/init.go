package ssoservice

import (
	"ecodepost/bff/pkg/invoker"
)

var (
	User *user
	Code *code
)

func Init() error {
	User = InitUser()
	Code = InitCode(invoker.Redis)
	return nil
}
