package wq

// Note: boards are created only as needed, and some SGF manipulation
// can be done creating no boards whatsoever.

var mutors = []string{"B", "W", "AB", "AW", "AE", "PL", "SZ"}

var total_board_updates int // For debugging.

// -----------------------------------------------------------------------------------------------
// clear_board_cache_recursive() needs to be called whenever a node's board cache becomes invalid.
// This can be due to:
//
//		* Changing a board-altering property.
//		* Changing the identity of its parent.

func (self *Node) clear_board_cache_recursive() {
	if self.__board_cache == nil { // If nil, all descendent caches are nil also.
		return // See note in the Node struct about this.
	}
	self.__board_cache = nil
	for _, child := range self.children {
		child.clear_board_cache_recursive()
	}
}

func (self *Node) mutor_check(key string) {

	// If the key changes the board, all descendent boards are also invalid.

	for _, s := range mutors {
		if key == s {
			self.clear_board_cache_recursive()
			break
		}
	}
}

// -----------------------------------------------------------------------------------------------

// Board uses the entire history of the tree up to this point to return a board.
// A copy of the result is cached intelligently; the cached board is also purged
// automatically if it becomes invalid (e.g. because a board-altering property
// changed in a relevant part of the SGF tree). Note that modifying a board has
// no effect on the SGF node which created it.
func (self *Node) Board() *Board {

	// Return cache if it exists...

	if self.__board_cache != nil {
		return self.__board_cache.Copy()
	}

	// Generate without recursion... also filling in any empty ancestor caches on the way.
	// This is essential, see note in Node struct about this.

	line := self.GetLine()
	var initial, work *Board

	for _, node := range line {
		if node.__board_cache != nil {
			initial = node.__board_cache // Care: points to the real thing, not a copy!
			continue
		}

		if work == nil {
			if initial == nil {
				work = NewBoard(node.RootBoardSize())
			} else {
				work = initial.Copy() // MUST COPY
			}
		}
		for i := range node.props {
			if node.props[i][0] == "KM" {
				work.KM = stringToFloat64(node.props[i][1])
			}
		}

		work.update_from_node(node)

		node.__board_cache = work.Copy()
	}

	// At this point, work is never nil. It is safe to return work itself since we
	// only stored copies of it in the cache.

	return work
}

func (self *Board) addMove(p string, c Colour, mType int) {
	if mType != 3 {
		x, y, isOnBoard := ParsePoint(p, self.Size)
		if isOnBoard {
			self.Move = append(self.Move, BoardMove{
				X:    x,
				Y:    y,
				C:    c,
				Type: mType,
			})
		}
	} else {
		self.Move = append(self.Move, BoardMove{
			C:    c,
			Type: mType,
		})
	}
}

func (self *Board) addMoves(p string, c Colour, mType int) {
	points := ParsePointList(p, self.Size)
	for _, point := range points {
		self.addMove(point, c, mType)
	}
}

func (self *Board) update_from_node(node *Node) {

	total_board_updates++

	// AB, AW, and AE are updated with AddStone() or AddList() which can create illegal
	// positions; this is normal according to the specs. Ko is cleared, next player is updated.

	for _, p := range node.AllValues("AB") {
		if len(p) == 5 && p[2] == ':' {
			self.AddList(p, BLACK)
			self.addMoves(p, BLACK, 2)
		} else {
			self.AddStone(p, BLACK)
			self.addMove(p, BLACK, 2)
		}
	}

	for _, p := range node.AllValues("AW") {
		if len(p) == 5 && p[2] == ':' {
			self.AddList(p, WHITE)
			self.addMoves(p, WHITE, 2)
		} else {
			self.AddStone(p, WHITE)
			self.addMove(p, WHITE, 2)
		}
	}

	for _, p := range node.AllValues("AE") {
		if len(p) == 5 && p[2] == ':' {
			self.AddList(p, EMPTY)
		} else {
			self.AddStone(p, EMPTY)
		}
	}

	// B and W are updated with ForceStone(), which has no legality checks but does
	// perform captures, as well as swapping the next player and setting the ko square.

	for _, p := range node.AllValues("B") {
		self.ForceStone(p, BLACK)
		self.addMove(p, BLACK, 1)
		self.Step++
	}

	for _, p := range node.AllValues("W") {
		self.ForceStone(p, WHITE)
		self.addMove(p, WHITE, 1)
		self.Step++
	}

	// Respect PL property

	pl, _ := node.GetValue("PL")
	if pl == "B" || pl == "b" {
		self.Player = BLACK
	}
	if pl == "W" || pl == "w" {
		self.Player = WHITE
	}
}
