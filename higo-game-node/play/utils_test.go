package play

import (
	"errors"
	"fmt"
	"testing"
)

func AAA() (e error) {
	e = errors.New("111")
	defer func() {
		if e != nil {
			fmt.Print("xxx")
		}
	}()
	return
}

func TestAAA(t *testing.T) {
	e := AAA()
	if e != nil {
		t.Error(e)
	}
}
