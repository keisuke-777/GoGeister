package main

import (
	"fmt"
	"./pkg/action"
	"./pkg/board"
)

func main() {
	test := action.RandomAction(4, 5)
	fmt.Println(test)
	b := board.Board{}
	b.Depth = 0
	b = board.InitBoard(b)
	fmt.Println(b)
	board.DispBoard(b)
}