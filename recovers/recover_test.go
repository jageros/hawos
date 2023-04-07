package recovers

import "testing"

func TestGo(t *testing.T) {
	Go(func() {
		panic("xxxx")
	})
}
