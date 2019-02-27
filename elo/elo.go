package elo

import (
	"fmt"
	"math"
)

const (
	defaultDeviation = 400
	defaultKFactor   = 32
)

// Player is ...
type Player struct {
	Ranking float64
}

// Parameters is ...
type Parameters struct {
	K, D float64
}

// Outcome is ...
type Outcome struct {
	Ranking, Delta float64
}

// Rank is ...
func Rank(player1, player2 Player, result int, params Parameters) (Outcome, Outcome, error) {
	var o1, o2 Outcome
	if result > 2 || result < 0 {
		return o1, o2, fmt.Errorf("Error - result value is out of range. Valid values are 0 a player 1 win, 1 for a player 2 win, or 2 for a draw")
	}
	t1 := transform(player1, params.D)
	t2 := transform(player2, params.D)

	e1 := expectation(player1, player2)
	e2 := expectation(player2, player1)

	s1 := score(0, result)
	s2 := score(1, result)

	r1 := rank(t1, e1, s1, params.K)
	r2 := rank(t2, e2, s2, params.K)

	o1 = Outcome{
		Ranking: r1,
		Delta:   r1 - player1.Ranking,
	}

	o2 = Outcome{
		Ranking: r2,
		Delta:   r2 - player1.Ranking,
	}

	return o1, o2, nil

}

func rank(t, e, s, k float64) float64 {
	return t + k*(s-e)
}

func transform(player Player, d float64) float64 {
	return math.Pow(10, (player.Ranking - d))
}

func expectation(player1, player2 Player) float64 {
	return player1.Ranking / (player1.Ranking + player2.Ranking)
}

func score(playerNum, result int) float64 {
	if result == 2 {
		return 0.5
	}
	if playerNum == 0 {
		if result == 0 {
			return 1.0
		}
		return 0
	}
	if result == 0 {
		return 0
	}
	return 1.0

}
