package board

import (
	"fmt"
	"math/rand"
)

type Board struct {
	BoardState [36]int
	Depth int
	TurnPlayer int // 1→先手のターン、-1→後手のターン
	Winner int // 0→勝者なし、1→先手の勝利、-1→後手の勝利
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

// 勝者を返す、勝者がいない場合は0を返す
func check_winner(b Board) (int){
	en_b, en_r, my_b, my_r := 0,0,0,0
	for _, piece := range b.BoardState {
		switch piece {
			case 0:
				continue
			case -2:
				en_r++
			case -1:
				en_b++
			case 1:
				my_b++
			case 2:
				my_r++
			default: // 想定していない駒が紛れている場合
				fmt.Println(piece)
		}
	}
	// 敵の赤駒を全て取った場合負け
	if en_r == 0{return -b.TurnPlayer}
	// 敵の青駒を全て取った場合勝ち
	if en_b == 0{return b.TurnPlayer}
	// 自分の青駒が全て取られた場合負け
	if my_b == 0{return -b.TurnPlayer}
	 // 自分の赤駒が全て取られた場合勝ち
	if my_r == 0{return b.TurnPlayer}
	// 勝者なし
	return 0
}

// 行動を受けて盤面を更新
func Next(b Board, action int) (Board){
	if action == 2 || action == 22 { // ゴール行動だけ特殊処理
		b.Winner = b.TurnPlayer
		return b
	}
	before_pos, after_pos := action_to_position(action) // 座標取得
	b.BoardState[before_pos], b.BoardState[after_pos] = 0, b.BoardState[before_pos] // 駒を動かす
	// 相手のターンにする
	ReverseBoardState(b.BoardState)
	b.Depth++
	b.TurnPlayer = -b.TurnPlayer
	return b
}

func action_to_position(action int) (int, int){
	before_pos, direction := action / 4, action % 4
	var after_pos int
	switch direction {
		case 0: // 下
			after_pos = before_pos + 6
		case 1: // 左
			after_pos = before_pos - 1
		case 2: // 上
			after_pos = before_pos - 6
		case 3: // 右
			after_pos = before_pos + 1
	}
	return before_pos, after_pos
}


// 盤面を反転させる
func ReverseBoardState(bs [36]int){
	n := len(bs)
	for i := 0; i < (n / 2); i++ {
		// 端から符号反転させてひっくり返す
		bs[i], bs[n-i-1] = -bs[n-i-1], -bs[i]
	}
}

// 座標と方角から行動番号を算出
func pos_and_dir_to_action(position int, direction int) (int){
	return position * 4 + direction
}

// 合法手の取得
func GetLegalActions(bs [36]int) ([]int) {
	legal_actions := make([]int, 0, 32) // 合法手の上限は32
	// ゴール行動だけは可能であれば最初に追加
	if bs[0] == 1 {
		legal_actions = append(legal_actions, 2)
	}
	if bs[5] == 1 {
		legal_actions = append(legal_actions, 22)
	}

	// 合法手を探す
	for position, piece := range bs {
		if piece == 1 || piece == 2 {
			x := position % 6
        	y := int(position / 6)
			if y != 5 { // 下端でない
				if bs[position + 6] != 1 && bs[position + 6] != 2{ // 下に自分の駒がいない
					legal_actions = append(legal_actions, pos_and_dir_to_action(position, 0))
				}
			}
			if x != 0 { // 左端でない
				if bs[position - 1] != 1 && bs[position - 1] != 2{ // 左に自分の駒がいない
					legal_actions = append(legal_actions, pos_and_dir_to_action(position, 1))
				}
			}
			if y != 0 { // 上端でない
				if bs[position - 6] != 1 && bs[position - 6] != 2{ // 上に自分の駒がいない
					legal_actions = append(legal_actions, pos_and_dir_to_action(position, 2))
				}
			}
			if y != 0 { // 右端でない
				if bs[position + 1] != 1 && bs[position + 1] != 2{ // 右に自分の駒がいない
					legal_actions = append(legal_actions, pos_and_dir_to_action(position, 3))
				}
			}
		}
	}
	return legal_actions
}

// 盤面の状態を表示
func DispBoard(b Board){
	// 相手のターンでも自分目線で表示する
	// これbが参照渡しだとバグるので確認する
	if b.TurnPlayer == -1 { ReverseBoardState(b.BoardState) }
	fmt.Println("ーーーーーーーーーーーーーーーーーーー")
	for i, piece := range b.BoardState {
		switch piece {
			case 0:
				fmt.Print("｜　　")
			case -2:
				fmt.Print("｜敵赤")
			case -1:
				fmt.Print("｜敵青")
			case 1:
				fmt.Print("｜自青")
			case 2:
				fmt.Print("｜自赤")
			default: // 想定していない駒が紛れている場合
				fmt.Println(piece)
		}
		if (i+1) % 6 == 0{
			fmt.Println("｜")
			fmt.Println("ーーーーーーーーーーーーーーーーーーー")
		}
	}
	fmt.Printf("経過ターン数：%d\n", b.Depth)
	fmt.Println()
}