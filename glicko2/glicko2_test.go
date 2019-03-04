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

func TestToMu(t *testing.T) {
	mu := toMu(1400)
	if math.Abs(mu - -0.5756) > .0001 {
		t.Log(mu)
		t.Fail()
	}

	mu = toMu(1550)
	if math.Abs(mu-.2878) > .0001 {
		t.Log(mu)
		t.Fail()
	}

	mu = toMu(1700)
	if math.Abs(mu-1.1513) > .0001 {
		t.Log(mu)
		t.Fail()
	}

}

func TestToPhi(t *testing.T) {
	phi := toPhi(30)
	if math.Abs(phi-.1727) > .0001 {
		t.Log(phi)
		t.Fail()
	}

	phi = toPhi(100)
	if math.Abs(phi-.5756) > .0001 {
		t.Log(phi)
		t.Fail()
	}

	phi = toPhi(300)
	if math.Abs(phi-1.7269) > .0001 {
		t.Log(phi)
		t.Fail()
	}
}

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

func TestVariance(t *testing.T) {
	totalImpact := 0.5621
	v := variance(totalImpact)
	if math.Abs(v-1.7790) > .0001 {
		t.Log(v)
		t.Fail()
	}
}

func TestImpact(t *testing.T) {
	var i, g, e float64
	g = 0.9955
	e = 0.639
	i = impact(g, e)
	if math.Abs(i-0.228607) > .000001 {
		t.Log(i)
		t.Fail()
	}

	g = 0.9531
	e = 0.432
	i = impact(g, e)
	if math.Abs(i-0.222899) > .000001 {
		t.Log(i)
		t.Fail()
	}

	g = 0.7242
	e = 0.303
	i = impact(g, e)
	if math.Abs(i-0.110762) > .000001 {
		t.Log(i)
		t.Fail()
	}
}

func TestTotalImpact(t *testing.T) {
	p1.Reset()
	p1.addResult(p2, 1)
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)
	ti := totalImpact(&p1.History)

	if math.Abs(ti-0.5621) > .0001 {
		t.Log(ti)
		t.Fail()
	}
}

func TestTotalResultScore(t *testing.T) {
	p1.Reset()
	p1.addResult(p2, 1)
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)
	rs := totalResultScore(&p1.History)
	if math.Abs(rs - -0.2720) > .0001 {
		t.Log(rs)
		t.Fail()
	}
}

func TestAlpha(t *testing.T) {
	a := toAlpha(0.06)
	if math.Abs(a - -5.62682) > .00001 {
		t.Log(a)
		t.Fail()
	}
}

func TestDelta(t *testing.T) {
	variance := 1.7785
	rs := -0.2720
	d := delta(variance, rs)
	if math.Abs(d - -.4837) > .0001 {
		t.Log(d)
		t.Fail()
	}
}

func TestIllinois(t *testing.T) {
	SystemConstant = 0.5
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
}

func TestInitialize(t *testing.T) {
	SystemConstant = 0.5
	sigma := 0.06
	phi := 1.1513
	variance := 1.7785
	delta := -0.4834
	a := -5.62682
	A, B := initializeComparison(sigma, variance, phi, delta, a)
	if math.Abs(A - -5.62682) > .00001 {
		t.Log(A)
		t.Fail()
	}

	if math.Abs(B - -6.12682) > .00001 {
		t.Log(B)
		t.Fail()
	}
}

func TestVolatility(t *testing.T) {
	SystemConstant = 0.5
	sigma := 0.06
	phi := 1.1513
	variance := 1.7785
	delta := -0.4834
	v := volatility(sigma, variance, phi, delta)
	if math.Abs(v-0.05999) > .00001 {
		t.Log(v)
		t.Fail()
	}
}

func TestPhiPrime(t *testing.T) {
	phi := 1.1513
	variance := 1.7785
	volatility := 0.05999
	pp := phiPrime(rd(phi, volatility), variance)
	if math.Abs(pp-0.8722) > .0001 {
		t.Log(pp)
		t.Fail()
	}
}

func TestMuPrime(t *testing.T) {
	ti := -0.272
	mu := 0.0
	pp := 0.8722
	mp := muPrime(mu, pp, ti)
	if math.Abs(mp - -0.2069) > .0001 {
		t.Log(mp)
		t.Fail()
	}
}

func TestFromMu(t *testing.T) {
	mu := -0.2069
	rating := fromMu(mu)
	if math.Abs(rating-1464.06) > .01 {
		t.Log(rating)
		t.Fail()
	}
}

func TestFromPhi(t *testing.T) {
	phi := 0.8722
	deviation := fromPhi(phi)
	if math.Abs(deviation-151.5) > .1 {
		t.Log(deviation)
		t.Fail()
	}
}

func TestGlicko2(t *testing.T) {
	p1.Reset()
	SystemConstant = 0.5
	p1.Win(p2)
	p1.Lose(p3)
	outcome := p1.Lose(p4)
	if math.Abs(outcome.Rating-1464.06) > .01 {
		t.Log(outcome)
		t.Fail()
	}

	if math.Abs(outcome.Volatility-.05999) > .00001 {
		t.Log(outcome)
		t.Fail()
	}

	if math.Abs(outcome.Deviation-151.52) > .01 {
		t.Log(outcome.Deviation)
		t.Fail()
	}
}
