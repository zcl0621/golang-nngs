package gameBus

import (
	"fmt"
	"testing"
)

func TestFuck(t *testing.T) {
	areaScore := 184.0 - 177.0
	KM := 7.5
	var win int
	var winResult string
	if areaScore > KM {
		win = 1
		winResult = fmt.Sprintf("B+%v", areaScore-KM)
	} else if areaScore < KM {
		win = 2
		winResult = fmt.Sprintf("W+%v", KM-areaScore)
	} else {
		win = 3
		winResult = fmt.Sprintf("Draw")
	}
	fmt.Println(win)
	fmt.Println(winResult)
}
