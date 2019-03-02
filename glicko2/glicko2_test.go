package glicko2

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

func TestPhi(t *testing.T) {

	phi := calcPhi(p2.Deviation)
	if math.Abs(phi-.1727) > .0001 {
		t.Log(phi)
		t.Fail()
	}

	phi = calcPhi(p3.Deviation)
	if math.Abs(phi-.5756) > .0001 {
		t.Log(phi)
		t.Fail()
	}

	phi = calcPhi(p4.Deviation)
	if math.Abs(phi-1.7269) > .0001 {
		t.Log(phi)
		t.Fail()
	}
}

func TestG(t *testing.T) {

	g := p2.g()
	if math.Abs(g-.9955) > .0001 {
		t.Log(g)
		t.Fail()
	}

	g = p3.g()
	if math.Abs(g-.9531) > .0001 {
		t.Log(g)
		t.Fail()
	}

	g = p4.g()
	if math.Abs(g-.7242) > .0001 {
		t.Log(g)
		t.Fail()
	}
}

func TestE(t *testing.T) {

	e := p1.e(p2)
	if math.Abs(e-.639) > .001 {
		t.Log(e)
		t.Fail()
	}

	e = p1.e(p3)
	if math.Abs(e-.432) > .001 {
		t.Log(e)
		t.Fail()
	}

	e = p1.e(p4)
	if math.Abs(e-.303) > .001 {
		t.Log(e)
		t.Fail()
	}
}

func TestVariation(t *testing.T) {
	p1.addResult(p2, 1)
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)

	v := p1.variation()

	if math.Abs(v-1.7785) > .001 {
		t.Log(v)
		t.Fail()
	}
}
