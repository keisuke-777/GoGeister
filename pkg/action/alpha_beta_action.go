package action

import (
	"../board"
	"math"
)

var INFINITY float64 = 10000.0
var ev_table [36]float64 = create_ev_table(val_func_selector(4))

// アルファベータ法で状態価値を計算
func AlphaBetaAction(bs [36]int) (int) {
    ib := board.IncompBoard{}
	ib = board.CreateIncompBoardFromBoardState(bs)
	best_action := -1
	alpha := -INFINITY
	max_depth := 5

	legal_actions := board.IncompGetLegalActions(ib.BoardState)
	for _, action := range legal_actions {
		next_ib := board.IncompNext(ib, action)
		score := alpha_beta(next_ib, -INFINITY, -alpha, 0, max_depth)
		if score > alpha {
			best_action = action
			alpha = score
		}
	}
	return best_action
}

// TODO:実験的に評価関数を変更する場合はExpAlphaBetaActionみたいなのを作って、
// 引数に評価関数あげちゃえば良いのでは？？？


func alpha_beta(ib board.IncompBoard, alpha float64, beta float64, search_depth int, max_depth int) (float64){
	// ゲーム終了時の処理
	if board.IncompIsDone(ib) {
		if ib.Winner == 0 {
			return 0.0
		} else if ib.TurnPlayer == ib.Winner {
			return INFINITY - float64(search_depth)
		} else {
			return -INFINITY + float64(search_depth)
		}
	}

	// 規定の深さまで来たら、探索を打ち切り状態を評価
	if search_depth == max_depth{
		return evaluate_board_state(ib.BoardState)
	}

	dynamic_alpha := alpha // ここいる？（多分値渡しなので引数を直で変更してもいい気がする）
	legal_actions := board.IncompGetLegalActions(ib.BoardState)
	for _, action := range legal_actions {
		next_ib := board.IncompNext(ib, action)
		score := alpha_beta(next_ib, -beta, -dynamic_alpha, search_depth+1, max_depth)
		if score > dynamic_alpha { dynamic_alpha = score }
		// 現ノードのベストスコアが親ノードを超えたら探索終了（カット）
		if dynamic_alpha >= beta { return dynamic_alpha }
	}
	return dynamic_alpha
}

// 評価関数を決める
func val_func_selector(id int) (func(float64) (float64)){
	switch id {
		case 0:
			return func(r float64) float64 { return 1.0 - r }
		case 1:
			return func(r float64) float64 { return 1.0 / math.Pow(r, 1.0) }
		case 2:
			return func(r float64) float64 { return 1.0 / math.Pow(r, 2.0) }
		case 3:
			return func(r float64) float64 { return 1.0 / math.Pow(r, 3.0) }
		case 4:
			return func(r float64) float64 { return 1.0 / math.Pow(r, 4.0) }
		default: // 現状良さげなやつをデフォルトにしておく
			return func(r float64) float64 { return 1.0 / math.Pow(r, 4.0) }
	}
}

func create_ev_table(val_func func(float64) (float64)) ([36]float64) {
	var x, y, r int
	ev_table := [36]float64{}
	for index, _ := range ev_table {
		x = index % 6
		y = index / 6
		if x < 3{
			r = x + 1 + y
		} else {
			r = 6 - x + y
		}
		ev_table[index] = val_func(float64(r))
	}
	return ev_table
}

// 盤面の状態を評価
func evaluate_board_state(bs [36]int) (float64) {
	value := 0.0
	for index, piece := range bs {
		switch piece {
			case 1, 3: // 自分の青駒がどれだけ脱出口に近いか
				value += ev_table[index]
			case -1, -3: // 敵の青駒がどれだけ自分の脱出口に近いか
				value -= ev_table[35 - index]
		}
	}
	// TODO:Pythonの実装では自分のターンか否かで符号を反転させていたので、動作確認が必須
	return value
}