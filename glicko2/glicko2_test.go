package glicko2

import (
	"math"
	"testing"
)

var (
	p1 = NewPlayer(Parameters{
		InitialDeviation:  200,
		InitialRating:     1500,
		InitialVolatility: 0.06,
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

	mu := toMu(p2.Rating)
	if math.Abs(mu - -0.5756) > .0001 {
		t.Log(mu)
		t.Fail()
	}

	mu = toMu(p3.Rating)
	if math.Abs(mu-.2878) > .0001 {
		t.Log(mu)
		t.Fail()
	}

	mu = toMu(p4.Rating)
	if math.Abs(mu-1.1513) > .0001 {
		t.Log(mu)
		t.Fail()
	}

}

func TestPhi(t *testing.T) {

	phi := toPhi(p2.Deviation)
	if math.Abs(phi-.1727) > .0001 {
		t.Log(phi)
		t.Fail()
	}

	phi = toPhi(p3.Deviation)
	if math.Abs(phi-.5756) > .0001 {
		t.Log(phi)
		t.Fail()
	}

	phi = toPhi(p4.Deviation)
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

func TestVariance(t *testing.T) {
	p1.addResult(p2, 1)
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)

	v := variance(totalImpact(&p1.History))

	if math.Abs(v-1.7789) > .0001 {
		t.Log(v)
		t.Fail()
	}
	p1.Reset()
}

func TestAlpha(t *testing.T) {
	a := toAlpha(p1.Volatility)
	if math.Abs(a - -5.62682) > .00001 {
		t.Log(a)
		t.Fail()
	}
}

func TestDelta(t *testing.T) {
	p1.addResult(p2, 1)
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)

	v := variance(totalImpact(&p1.History))
	d := delta(v, &p1.History)

	if math.Abs(d - -.4839) > .0001 {
		t.Log(d)
		t.Fail()
	}
	p1.Reset()
}

func TestIllinois(t *testing.T) {
	p1.addResult(p2, 1)
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)
	phi := 1.1513
	variance := 1.7785
	delta := -0.4834
	a := -5.62682
	A := -5.62682
	B := -6.12682
	ia := illinois(A, phi, variance, a, delta)
	if math.Abs(ia - -0.00053567) > .00000001 {
		t.Log(ia)
		t.Fail()
	}

	ib := illinois(B, phi, variance, a, delta)
	if math.Abs(ib-1.999675) > .000001 {
		t.Log(ib)
		t.Fail()
	}
	p1.Reset()
}

func TestInitialize(t *testing.T) {
	p1.addResult(p2, 1)
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)
	phi := toPhi(p1.Deviation)
	ti := totalImpact(&p1.History)
	variance := variance(ti)
	delta := delta(variance, &p1.History)
	a := toAlpha(p1.Volatility)
	A, B := initializeComparison(p1.Volatility, variance, phi, delta, a)
	if math.Abs(A - -5.62682) > .00001 {
		t.Log(A)
		t.Fail()
	}

	if math.Abs(B - -6.12682) > .00001 {
		t.Log(B)
		t.Fail()
	}
	p1.Reset()
}

func TestGlicko2(t *testing.T) {
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)
	o := p1.Win(p2)
	if math.Abs(o.Rating-1464.06) > .01 {
		t.Log(o)
		t.Fail()
	}

	if math.Abs(o.Volatility-.05999) > .00001 {
		t.Log(o)
		t.Fail()
	}

	if math.Abs(o.Deviation-151.52) > .01 {
		t.Log(o.Deviation)
		t.Fail()
	}
	p1.Reset()
}
