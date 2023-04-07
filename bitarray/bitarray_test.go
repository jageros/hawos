package bitarray

import (
	"fmt"
	"testing"
)

func TestIsHigh(t *testing.T) {
	var bs []byte
	bs = SetHigh(bs, 30)
	bs = SetHigh(bs, 31)
	fmt.Println(bs)
	bs = SetHigh(bs, 511)
	fmt.Println(bs)
	bs = SetLow(bs, 30)
	fmt.Println(bs)
	bs = SetLow(bs, 31)
	fmt.Println(bs)
	bs = SetLow(bs, 511)
	fmt.Println(bs)

	fmt.Println(IsHigh(bs, 30))
}
