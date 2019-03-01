package glicko

import (
	"math"
	"testing"
)

var (
	p1 = NewPlayer(Parameters{
		InitialRD:     DefaultInitialRD,
		InitialRating: DefaultInitialRating,
	})
	p2 = NewPlayer(Parameters{
		InitialRD:     30,
		InitialRating: 1400,
	})
	p3 = NewPlayer(Parameters{
		InitialRD:     100,
		InitialRating: 1550,
	})
	p4 = NewPlayer(Parameters{
		InitialRD:     300,
		InitialRating: 1700,
	})
	p5 = NewPlayer(Parameters{
		InitialRD:     200,
		InitialRating: 1500,
	})
	p6 = NewPlayer(Parameters{
		InitialRD:     200,
		InitialRating: 1500,
	})
	p7 = NewPlayer(Parameters{
		InitialRD:     200,
		InitialRating: 1500,
	})
	p8 = NewPlayer(Parameters{
		InitialRD:     200,
		InitialRating: 1500,
	})
)

func TestGRD(t *testing.T) {
	grd := p2.gRD()
	if math.Abs(grd-0.9955) > .0001 {
		t.Log(grd)
		t.Fail()
	}
	grd = p3.gRD()
	if math.Abs(grd-0.9531) > .0001 {
		t.Log(grd)
		t.Fail()
	}
	grd = p4.gRD()
	if math.Abs(grd-0.7242) > .0001 {
		t.Log(grd)
		t.Fail()
	}
}

func TestE(t *testing.T) {
	e := p1.e(p2)
	if math.Abs(e-0.639) > .001 {
		t.Log(e)
		t.Fail()
	}
	e = p1.e(p3)
	if math.Abs(e-0.432) > .001 {
		t.Log(e)
		t.Fail()
	}
	e = p1.e(p4)
	if math.Abs(e-0.303) > .001 {
		t.Log(e)
		t.Fail()
	}
}

func TestDSquared(t *testing.T) {
	p5.addResult(p2, 1)
	p5.addResult(p3, 0)
	p5.addResult(p4, 0)
	ds := p5.dsquared()
	if math.Abs(ds-53685.74) > 0.01 {
		t.Log(ds)
		t.Fail()
	}
}

func TestWin(t *testing.T) {
	p6.addResult(p3, 0)
	p6.addResult(p4, 0)
	p6.Win(p2)
	if math.Abs(p6.Rating-1464.1) > 0.1 {
		t.Log(p6.Rating, p6.RD)
		t.Fail()
	}
	if math.Abs(p6.RD-151.4) > 0.1 {
		t.Log(p6.Rating, p6.RD)
		t.Fail()
	}
}

func TestLoss(t *testing.T) {
	p7.addResult(p3, 0)
	p7.addResult(p4, 0)
	p7.Loss(p2)
	if math.Abs(p7.Rating-1332.7) > 0.1 {
		t.Log(p7.Rating, p7.RD)
		t.Fail()
	}
	if math.Abs(p7.RD-151.4) > 0.1 {
		t.Log(p7.Rating, p7.RD)
		t.Fail()
	}
}

func TestDraw(t *testing.T) {
	p8.addResult(p3, 0)
	p8.addResult(p4, 0)
	p8.Draw(p2)
	if math.Abs(p8.Rating-1398.4) > 0.1 {
		t.Log(p8.Rating, p8.RD)
		t.Fail()
	}
	if math.Abs(p8.RD-151.4) > 0.1 {
		t.Log(p8.Rating, p8.RD)
		t.Fail()
	}
}
