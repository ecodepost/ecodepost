package bffsso_test

import (
	"errors"
	"fmt"
	"testing"

	errcodev1 "ecodepost/pb/errcode/v1"

	"github.com/gotomicro/ego/core/eerrors"
)

func TestErrorWrap(t *testing.T) {
	err := errcodev1.AuthErrBrowserCookieSystemError()
	err2 := fmt.Errorf("something, err: %w", err)

	some := eerrors.FromError(err2)
	fmt.Printf("	errors.Is(some, err)--------------->"+"%+v\n", errors.Is(some, err))
	fmt.Printf("some--------------->"+"%+v\n", err2)
	fmt.Printf("some--------------->"+"%+v\n", some)
}
