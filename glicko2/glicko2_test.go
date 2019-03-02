package glicko2

import (
	"math"
	"testing"
)

var (
	p1 = NewPlayer(Parameters{
		InitialDeviation: DefaultInitialDeviation,
		InitialRating:    DefaultInitialRating,
	})
	p2 = NewPlayer(Parameters{
		InitialDeviation: 30,
		InitialRating:    1400,
	})
	p3 = NewPlayer(Parameters{
		InitialDeviation: 100,
		InitialRating:    1550,
	})
	p4 = NewPlayer(Parameters{
		InitialDeviation: 300,
		InitialRating:    1700,
	})
)

func TestMu(t *testing.T) {

	mu := calcMu(p2.Rating)
	if math.Abs(mu - -0.5756) > .0001 {
		t.Log(mu)
		t.Fail()
	}

	mu = calcMu(p3.Rating)
	if math.Abs(mu-.2878) > .0001 {
		t.Log(mu)
		t.Fail()
	}

	mu = calcMu(p4.Rating)
	if math.Abs(mu-1.1513) > .0001 {
		t.Log(mu)
		t.Fail()
	}

}
