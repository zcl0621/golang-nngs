package wq

import (
	"encoding/json"
	"fmt"
)

// Stones returns all stones in the group at point p, in arbitrary order. The
// argument should be an SGF coordinate, e.g. "dd".
func (self *Board) Stones(p string) []string {

	colour := self.Get(p)
	if colour == EMPTY { // true also if offboard / invalid
		return nil
	}

	touched := make(map[string]bool)
	return self.stones_recurse(p, colour, touched, nil)
}

func (self *Board) stones_recurse(p string, colour Colour, touched map[string]bool, ret []string) []string {

	// Note the constant returning and updating of ret since appends are not visible to caller otherwise.

	touched[p] = true
	ret = append(ret, p)

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.get_fast(a) == colour {
			if touched[a] == false {
				ret = self.stones_recurse(a, colour, touched, ret)
			}
		}
	}

	return ret
}

// HasLiberties checks whether the group at point p has any liberties. The
// argument should be an SGF coordinate, e.g. "dd". For groups of stones on
// normal boards, this is always true, but can be false if the calling program
// is manipulating the board directly.
//
// If the point p is empty, returns false.
func (self *Board) HasLiberties(p string) bool {

	colour := self.Get(p)
	if colour == EMPTY { // true also if offboard / invalid
		return false
	}

	touched := make(map[string]bool)
	return self.has_liberties_recurse(p, colour, touched)
}

func (self *Board) has_liberties_recurse(p string, colour Colour, touched map[string]bool) bool {

	touched[p] = true

	for _, a := range AdjacentPoints(p, self.Size) {
		a_colour := self.get_fast(a)
		if a_colour == EMPTY {
			return true
		} else if a_colour == colour {
			if touched[a] == false {
				if self.has_liberties_recurse(a, colour, touched) {
					return true
				}
			}
		}
	}

	return false
}

// Liberties returns the liberties of the group at point p, in arbitrary order.
// The argument should be an SGF coordinate, e.g. "dd".
func (self *Board) Liberties(p string) []string {

	colour := self.Get(p)
	if colour == EMPTY { // true also if offboard / invalid
		return nil
	}

	touched := make(map[string]bool)
	touched[p] = true // Note this - slightly different setup than the other functions
	return self.liberties_recurse(p, colour, touched, nil)
}

func (self *Board) liberties_recurse(p string, colour Colour, touched map[string]bool, ret []string) []string {

	// Note that this function uses the touched map in a different way from others.
	// Note the constant returning and updating of ret since appends are not visible to caller otherwise.

	for _, a := range AdjacentPoints(p, self.Size) {
		t := touched[a]
		if t == false {
			touched[a] = true
			a_colour := self.get_fast(a)
			if a_colour == EMPTY {
				ret = append(ret, a)
			} else if a_colour == colour {
				ret = self.liberties_recurse(a, colour, touched, ret)
			}
		}
	}

	return ret
}

// Singleton returns true if the specified stone is a group of size 1. The
// argument should be an SGF coordinate, e.g. "dd".
func (self *Board) Singleton(p string) bool {

	colour := self.Get(p)
	if colour == EMPTY { // true also if offboard / invalid
		return false
	}

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.get_fast(a) == colour {
			return false
		}
	}

	return true
}

// Legal returns true if a play at point p would be legal. The argument should
// be an SGF coordinate, e.g. "dd". The colour is determined intelligently. The
// board is not changed. If false, the reason is given in the error.
func (self *Board) Legal(p string) (bool, error) {
	return self.LegalColour(p, self.Player)
}

// LegalColour is like Legal, except the colour is specified rather than being
// automatically determined.
func (self *Board) LegalColour(p string, colour Colour) (bool, error) {

	if colour != BLACK && colour != WHITE {
		return false, fmt.Errorf("colour not BLACK or WHITE")
	}

	x, y, onboard := ParsePoint(p, self.Size)

	if onboard == false {
		return false, fmt.Errorf("invalid or off-board string %q", p)
	}

	if self.State[x][y] != EMPTY {
		return false, fmt.Errorf("point %q (%v,%v) was not empty", p, x, y)
	}

	if self.Ko == p {
		if colour == self.Player { // i.e. we've not forced a move by the wrong colour.
			return false, fmt.Errorf("ko recapture at %q (%v,%v) forbidden", p, x, y)
		}
	}

	has_own_liberties := false
	for _, a := range AdjacentPoints(p, self.Size) {
		if self.get_fast(a) == EMPTY {
			has_own_liberties = true
			break
		}
	}

	if has_own_liberties == false {

		// The move we are playing will have no liberties of its own.
		// Therefore, it will be legal iff it has a neighbour which:
		//
		//		- Is an enemy group with 1 liberty, or
		//		- Is a friendly group with 2 or more liberties.

		allowed := false

		for _, a := range AdjacentPoints(p, self.Size) {
			if self.get_fast(a) == colour.Opposite() {
				if len(self.Liberties(a)) == 1 {
					allowed = true
					break
				}
			} else if self.get_fast(a) == colour {
				if len(self.Liberties(a)) >= 2 {
					allowed = true
					break
				}
			} else {
				panic("wat")
			}
		}

		if allowed == false {
			return false, fmt.Errorf("suicide at %q (%v,%v) forbidden", p, x, y)
		}
	}

	// The move is legal!

	return true, nil
}

type KataGoAnalysis struct {
	Id               string     `json:"id"`
	InitialStones    [][]string `json:"initialStones,omitempty"`
	Moves            [][]string `json:"moves,omitempty"`
	Rules            string     `json:"rules"`
	Komi             float64    `json:"komi"`
	BoardXSize       int        `json:"boardXSize"`
	BoardYSize       int        `json:"boardYSize"`
	IncludeOwnership bool       `json:"includeOwnership"`
}

func (a *KataGoAnalysis) ToString() string {
	d, _ := json.Marshal(a)
	return string(d)
}

func (self *Board) GetKataGoAnalysisData(gameId uint) *KataGoAnalysis {
	var res KataGoAnalysis
	res.Id = fmt.Sprintf("%d", gameId)
	for i := range self.Move {
		switch self.Move[i].Type {
		case 1:
			var move []string
			switch self.Move[i].C {
			case BLACK:
				move = append(move, "B")
			case WHITE:
				move = append(move, "W")
			default:
				break
			}
			move = append(move, XYToGTP(self.Move[i].X, self.Move[i].Y))
			res.Moves = append(res.Moves, move)
		case 2:
			var move []string
			switch self.Move[i].C {
			case BLACK:
				move = append(move, "B")
			case WHITE:
				move = append(move, "W")
			default:
				break
			}
			move = append(move, XYToGTP(self.Move[i].X, self.Move[i].Y))
			res.InitialStones = append(res.InitialStones, move)
		case 3:
			var move []string
			switch self.Move[i].C {
			case BLACK:
				move = append(move, "B")
			case WHITE:
				move = append(move, "W")
			default:
				break
			}
			move = append(move, "pass")
			res.Moves = append(res.Moves, move)
		}
	}
	res.Komi = self.KM
	res.Rules = "chinese"
	res.BoardXSize = self.Size
	res.BoardYSize = self.Size
	res.IncludeOwnership = true
	return &res
}

type OwnerShip struct {
	X    int     `json:"x"`
	Y    int     `json:"y"`
	C    int     `json:"c"`
	Size float64 `json:"size"`
}

func GetColorOwnerShip(ownership []float64, boardSize int) *[]OwnerShip {
	distOwnerShip := make([][]float64, boardSize)
	for i := range distOwnerShip {
		distOwnerShip[i] = make([]float64, boardSize)
	}

	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			index := i*boardSize + j
			distOwnerShip[j][i] = (ownership)[index]
		}
	}
	swapArray(distOwnerShip)
	var res []OwnerShip
	for x := range distOwnerShip {
		for y := range distOwnerShip[x] {
			d := OwnerShip{
				X:    x,
				Y:    y,
				Size: distOwnerShip[x][y],
			}
			if d.Size > 0 {
				d.C = 1
			} else if d.Size < 0 {
				d.C = 2
			}
			res = append(res, d)
		}
	}
	return &res
}

func swap(arr [][]float64, i1, j1, i2, j2 int) {
	arr[i1][j1], arr[i2][j2] = arr[i2][j2], arr[i1][j1]
}

func swapArray(arr [][]float64) {
	rows, cols := len(arr), len(arr[0])
	for i := 0; i < rows; i++ {
		for j := 0; j < cols/2; j++ {
			swap(arr, i, j, i, cols-j-1)
		}
	}
}
