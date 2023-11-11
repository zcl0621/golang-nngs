package ai

import (
	"fmt"
	"testing"
)

func TestBBB(t *testing.T) {
	a := 0
	for {
		switch a {
		case 100:
			a++
			break
		case 200:
			fmt.Println(a)
			return
		default:
			fmt.Println(a)
			a++
		}
	}

}
