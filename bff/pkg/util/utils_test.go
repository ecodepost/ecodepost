package util

import (
	"fmt"
	"testing"

	"github.com/gotomicro/ego/core/util/xtime"
)

func TestCheckPhone(t *testing.T) {
	res := CheckMobile("18827073672")
	fmt.Printf("res--------------->"+"%+v\n", res)
}

func TestDuration(t *testing.T) {
	res := xtime.Duration("61s").Seconds()
	fmt.Printf("res--------------->"+"%+v\n", res)
}
