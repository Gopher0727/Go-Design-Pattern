package main

var units = map[int]*ChessPieceUnit{
	1: {
		ID:    1,
		Name:  "车",
		Color: "red",
	},
	2: {
		ID:    2,
		Name:  "马",
		Color: "black",
	},
}

// 棋子享元
type ChessPieceUnit struct {
	ID    uint
	Name  string
	Color string
}

func NewChessPieceUnit(id int) *ChessPieceUnit {
	return units[id]
}

type ChessPiece struct {
	Unit *ChessPieceUnit
	x    int
	y    int
}

type ChessBoard struct {
	chessPieces map[int]*ChessPiece
}

func NewChessBoard() *ChessBoard {
	board := &ChessBoard{
		chessPieces: make(map[int]*ChessPiece),
	}
	for id := range units {
		board.chessPieces[id] = &ChessPiece{
			Unit: NewChessPieceUnit(id),
			x:    0,
			y:    0,
		}
	}
	return board
}

func (c *ChessBoard) Move(id, x, y int) {
	c.chessPieces[id].x = x
	c.chessPieces[id].y = y
}
