package glicko

import (
	"math"
	"testing"
)

func TestGRD(t *testing.T) {
	p := Parameters{
		InitialRD:      30,
		InitialRanking: DefaultInitialRanking,
	}
	player := NewPlayer(p)
	grd := player.gRD()
	if math.Abs(grd-0.9955) > .0001 {
		t.Log(grd)
		t.Fail()
	}

	p = Parameters{
		InitialRD:      100,
		InitialRanking: DefaultInitialRanking,
	}
	player = NewPlayer(p)
	grd = player.gRD()
	if math.Abs(grd-0.9531) > .0001 {
		t.Log(grd)
		t.Fail()
	}

	p = Parameters{
		InitialRD:      300,
		InitialRanking: DefaultInitialRanking,
	}
	player = NewPlayer(p)
	grd = player.gRD()
	if math.Abs(grd-0.7242) > .0001 {
		t.Log(grd)
		t.Fail()
	}
}
