package glicko

import (
	"math"
	"testing"
)

func TestGRD(t *testing.T) {
	p := Parameters{
		InitialRD:     30,
		InitialRating: DefaultInitialRating,
	}
	player := NewPlayer(p)
	grd := player.gRD()
	if math.Abs(grd-0.9955) > .0001 {
		t.Log(grd)
		t.Fail()
	}

	p = Parameters{
		InitialRD:     100,
		InitialRating: DefaultInitialRating,
	}
	player = NewPlayer(p)
	grd = player.gRD()
	if math.Abs(grd-0.9531) > .0001 {
		t.Log(grd)
		t.Fail()
	}

	p = Parameters{
		InitialRD:     300,
		InitialRating: DefaultInitialRating,
	}
	player = NewPlayer(p)
	grd = player.gRD()
	if math.Abs(grd-0.7242) > .0001 {
		t.Log(grd)
		t.Fail()
	}
}

func TestE(t *testing.T) {
	p := Parameters{
		InitialRD:     350,
		InitialRating: DefaultInitialRating,
	}
	op := Parameters{
		InitialRD:     30,
		InitialRating: 1400,
	}
	player1 := NewPlayer(p)
	player2 := NewPlayer(op)
	e := player1.e(player2)
	if math.Abs(e-0.639) > .001 {
		t.Log(e)
		t.Fail()
	}

	p = Parameters{
		InitialRD:     350,
		InitialRating: DefaultInitialRating,
	}
	op = Parameters{
		InitialRD:     100,
		InitialRating: 1550,
	}
	player1 = NewPlayer(p)
	player2 = NewPlayer(op)
	e = player1.e(player2)
	if math.Abs(e-0.432) > .001 {
		t.Log(e)
		t.Fail()
	}

	p = Parameters{
		InitialRD:     350,
		InitialRating: DefaultInitialRating,
	}
	op = Parameters{
		InitialRD:     300,
		InitialRating: 1700,
	}
	player1 = NewPlayer(p)
	player2 = NewPlayer(op)
	e = player1.e(player2)
	if math.Abs(e-0.303) > .001 {
		t.Log(e)
		t.Fail()
	}
}
