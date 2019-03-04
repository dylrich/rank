package glicko

import (
	"math"
	"testing"
)

var (
	p1 = NewPlayer(Parameters{
		InitialDeviation: 200,
		InitialRating:    1500,
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

func TestToG(t *testing.T) {
	g := toG(30.0)
	if math.Abs(g-.9955) > .0001 {
		t.Log(g)
		t.Fail()
	}

	g = toG(100.0)
	if math.Abs(g-.9531) > .0001 {
		t.Log(g)
		t.Fail()
	}

	g = toG(300.0)
	if math.Abs(g-.7242) > .0001 {
		t.Log(g)
		t.Fail()
	}
}

func TestToE(t *testing.T) {
	e := toE(1500, 1400, .9955)
	if math.Abs(e-.639) > .001 {
		t.Log(e)
		t.Fail()
	}

	e = toE(1500, 1550, .9531)
	if math.Abs(e-.432) > .001 {
		t.Log(e)
		t.Fail()
	}

	e = toE(1500, 1700, .7242)
	if math.Abs(e-.303) > .001 {
		t.Log(e)
		t.Fail()
	}
}

func TestDSquared(t *testing.T) {
	p1.Reset()
	p1.addResult(p2, 1)
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)
	ds := dsquared(&p1.History)
	if math.Abs(ds-53685.74) > 0.01 {
		t.Log(ds)
		t.Fail()
	}
}

func TestGlicko(t *testing.T) {
	p1.Reset()
	p1.Win(p2)
	p1.Lose(p3)
	outcome := p1.Lose(p4)
	if math.Abs(outcome.Rating-1464.1) > 0.1 {
		t.Log(outcome.Rating, outcome.Deviation)
		t.Fail()
	}
	if math.Abs(outcome.Deviation-151.4) > 0.1 {
		t.Log(outcome.Rating, outcome.Deviation)
		t.Fail()
	}
}
