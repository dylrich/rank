package elo

import "testing"

func TestWinLoseDraw(t *testing.T) {
	p := Parameters{
		K: DefaultKFactor,
		D: DefaultDeviation,
	}
	player1 := NewPlayer(p)
	player2 := NewPlayer(p)

	player1.Win(player2)
	player2.Lose(player1)
	if player1.Ranking != 1516 || int(player2.Ranking) != 1484 {
		t.Log(player1)
		t.Log(player2)
		t.Fail()
	}

	player1.Draw(player2)
	if int(player1.Ranking) != 1514 {
		t.Log(player1)
		t.Log(player2)
		t.Fail()
	}
}
