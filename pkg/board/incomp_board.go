package board

import (
	"fmt"
)

type IncompBoard struct {
	BoardState [36]int
	Depth int
	TurnPlayer int // 1→先手のターン、-1→後手のターン
	Winner int // 0→勝者なし、1→先手の勝利、-1→後手の勝利
	MyLeftBluePiece int
	MyLeftRedPiece int
	EnLeftBluePiece int
	EnLeftRedPiece int
}

func CreateIncompBoardFromBoardState(bs [36]int) (IncompBoard) {
	ib := IncompBoard{}
	ib.Depth = 0
	ib.TurnPlayer = 1
	ib.Winner = 0
	// 相手の駒の個数を計測
	for i, piece := range bs {
		if piece == 1 {
			ib.MyLeftBluePiece++
		} else if piece == 2 {
			ib.MyLeftRedPiece++
		} else if piece == -1 {
			ib.EnLeftBluePiece++
			ib.BoardState[i] = -3 // 相手の駒は紫駒にする
		} else if piece == -2 {
			ib.EnLeftRedPiece++
			ib.BoardState[i] = -3 // 相手の駒は紫駒にする
		}
	}
	return ib
}

// 勝者を返す、勝者がいない場合は0を返す
// 対戦相手の駒の色がわからない前提なので、BoardStateからは情報を取らない
func incomp_check_winner(ib IncompBoard) (int) {
	// 敵の赤駒を全て取った場合負け
	if ib.EnLeftRedPiece == 0 { return -ib.TurnPlayer }
	// 敵の青駒を全て取った場合勝ち
	if ib.EnLeftBluePiece == 0 { return ib.TurnPlayer }
	// 自分の青駒が全て取られた場合負け
	if ib.MyLeftBluePiece == 0 { return -ib.TurnPlayer }
	 // 自分の赤駒が全て取られた場合勝ち
	if ib.MyLeftRedPiece == 0 { return ib.TurnPlayer }
	// 勝者なし
	return 0
}

func IncompIsDone(ib IncompBoard) (bool) {
	if ib.Winner != 0{
		return true
	}
	if ib.Depth >= 200{
		return true
	}
	return false
}

// 行動を受けて盤面を更新
func IncompNext(ib IncompBoard, action int) (IncompBoard) {
	if action == 2 || action == 22 { // ゴール行動だけ特殊処理
		ib.Winner = ib.TurnPlayer
		return ib
	}
	before_pos, after_pos := action_to_position(action) // 座標取得
	// 駒が消える場合の処理（IncompNext限定）
	if ib.BoardState[after_pos] != 0 {
		switch ib.BoardState[after_pos] {
			case -1:
				ib.EnLeftBluePiece--
			case -2:
				ib.EnLeftRedPiece--
			case -3: // 紫駒
				ib.EnLeftRedPiece--
		}
	}
	ib.BoardState[before_pos], ib.BoardState[after_pos] = 0, ib.BoardState[before_pos] // 駒を動かす
	ib.Winner = incomp_check_winner(ib) // 勝敗の判定（駒を取った場合、勝敗が確定する場合がある）
	// 相手のターンにする
	ib.EnLeftBluePiece, ib.MyLeftBluePiece = ib.MyLeftBluePiece, ib.EnLeftBluePiece // 自青<->敵青
	ib.EnLeftRedPiece, ib.MyLeftRedPiece = ib.MyLeftRedPiece, ib.EnLeftRedPiece // 自赤<->敵赤
	ReverseBoardState(ib.BoardState)
	ib.Depth++
	ib.TurnPlayer = -ib.TurnPlayer
	return ib
}


// 合法手の取得
func IncompGetLegalActions(bs [36]int) ([]int) {
	legal_actions := make([]int, 0, 32) // 合法手の上限は32
	// ゴール行動だけは可能であれば最初に追加
	// 相手の駒は色が不明なため全てゴール可能にする
	if bs[0] == 1 || bs[0] == 3 {
		legal_actions = append(legal_actions, 2)
	}
	if bs[5] == 1 || bs[5] == 3 {
		legal_actions = append(legal_actions, 22)
	}

	// 合法手を探す
	for position, piece := range bs {
		if piece == 1 || piece == 2 {
			x := position % 6
        	y := position / 6
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
			if x != 5 { // 右端でない
				if bs[position + 1] != 1 && bs[position + 1] != 2{ // 右に自分の駒がいない
					legal_actions = append(legal_actions, pos_and_dir_to_action(position, 3))
				}
			}
		}
	}
	return legal_actions
}

// 盤面の状態を表示
func IncompDispBoard(ib IncompBoard){
	// 相手のターンでも自分目線で表示する
	// これbが参照渡しだとバグるので確認する
	if ib.TurnPlayer == -1 { ReverseBoardState(ib.BoardState) }
	fmt.Println("ーーーーーーーーーーーーーーーーーーー")
	for i, piece := range ib.BoardState {
		switch piece {
			case 0:
				fmt.Print("｜　　")
			case -3:
				fmt.Print("｜敵駒")
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
	fmt.Printf("経過ターン数：%d\n", ib.Depth)
	fmt.Println()
}