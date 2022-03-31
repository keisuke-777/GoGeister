package board

import (
	"math/rand"
)

type Board struct {
	BoardState [36]int
	Depth int
	TurnPlayer int // 1→先手のターン、-1→後手のターン
	Winner int // 0→勝者なし、1→先手の勝利、2→後手の勝利
}

func InitBoard(b Board) (Board){
	// 盤面の初期化
	enemy_pieces := []int{-1,-1,-1,-1,-2,-2,-2,-2}
	shuffle(enemy_pieces)
	b.BoardState[1] = enemy_pieces[0]
	b.BoardState[2] = enemy_pieces[1]
	b.BoardState[3] = enemy_pieces[2]
	b.BoardState[4] = enemy_pieces[3]
	b.BoardState[7] = enemy_pieces[4]
	b.BoardState[8] = enemy_pieces[5]
	b.BoardState[9] = enemy_pieces[6]
	b.BoardState[10] = enemy_pieces[7]
	my_pieces := []int{1,1,1,1,2,2,2,2}
	shuffle(my_pieces)
	b.BoardState[25] = my_pieces[0]
	b.BoardState[26] = my_pieces[1]
	b.BoardState[27] = my_pieces[2]
	b.BoardState[28] = my_pieces[3]
	b.BoardState[31] = my_pieces[4]
	b.BoardState[32] = my_pieces[5]
	b.BoardState[33] = my_pieces[6]
	b.BoardState[34] = my_pieces[7]
	
	// ターンプレイヤの初期化
	b.Winner = 1

	// ターン数、勝者は0が初期値のため触らなくて良い
	return b
}

// スライスをシャッフルする
func shuffle(list []int){
	for i := len(list); i > 1; i-- {
		j := rand.Intn(i) // 0~i-1 の範囲で乱数生成
		list[i - 1], list[j] = list[j], list[i - 1]
	}
}