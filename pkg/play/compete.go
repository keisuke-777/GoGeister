package play

import (
	"fmt"
	"math/rand"
	"time"
	"../action"
	"../board"
)


func compete(fst_act_fn func(bs [36]int) (int), sec_act_fn func(bs [36]int) (int), seed int64, num_of_games int){
	num_of_wins := 0 // 先手の勝利回数
	num_of_drows := 0 // 引き分けの回数
	for i := 1; i < (num_of_games / 2) + 1; i++ {
		rand.Seed(seed * int64(i))
		b := board.Board{}
		b = board.InitBoard(b)
		for {
			if board.IsDone(b) { // 勝者が決定している場合
				switch b.Winner {
					case 0:
						num_of_drows++
					case 1:
						num_of_wins++
					default:
						// case -1 => 敗北なので何もしない
				}
				break
			}
			// 次の状態の取得
			if (b.Depth % 2) == 0{
				b = board.Next(b, fst_act_fn(b.BoardState))
			} else {
				b = board.Next(b, sec_act_fn(b.BoardState))
			}
		}
		
		// 同一の条件で先手後手を入れ替える
		rand.Seed(seed * int64(i))
		b = board.Board{}
		b = board.InitBoard(b)
		for {
			if board.IsDone(b) {
				switch b.Winner {
					case 0:
						num_of_drows++
					case -1:
						num_of_wins++
				}
				break
			}
			if (b.Depth % 2) == 0{
				b = board.Next(b, sec_act_fn(b.BoardState))
			} else {
				b = board.Next(b, fst_act_fn(b.BoardState))
			}
		}
	}
	// 勝率などを出力
	fmt.Printf("先手勝利回数：%d, 後手勝利回数：%d, 引き分け回数：%d\n", num_of_wins, num_of_games - num_of_wins- num_of_drows, num_of_drows)
}

func RandVsRand(){
	seed := time.Now().UnixNano()
	fmt.Printf("ランダムvsランダム（初期seed:%d）\n", seed)
	compete(action.RandomAction, action.RandomAction, seed, 100)
}