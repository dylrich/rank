package elo

import (
	"math"
	"testing"
)

var (
	p1 = NewPlayer(Parameters{
		InitialRating: 1500,
	})
	p2 = NewPlayer(Parameters{
		InitialRating: 1500,
	})
)

func TestElo(t *testing.T) {
	p1.Win(p2.Rating)
	p2.Lose(p1.Rating)
	if math.Abs(p1.Rating-1516) > 1 || math.Abs(p2.Rating-1484) > 1 {
		t.Log(p1)
		t.Log(p2)
		t.Fail()
	}
}
