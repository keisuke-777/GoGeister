package action

import (
	"../board"
    "math/rand"
)

// ランダムな行動を返す
func RandomAction(bs [36]int) (int) {
    legal_actions := board.GetLegalActions(bs)
    return legal_actions[rand.Intn(len(legal_actions))]
}