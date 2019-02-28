package glicko

import (
	"log"
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
	e = p1.e(p2)
	if math.Abs(e-0.432) > .001 {
		t.Log(e)
		t.Fail()
	}
	e = p1.e(p2)
	if math.Abs(e-0.303) > .001 {
		t.Log(e)
		t.Fail()
	}
}

func TestDSquared(t *testing.T) {
	p1.addResult(p2, 1)
	p1.addResult(p3, 0)
	p1.addResult(p4, 0)
	ds := p1.dsquared()

	log.Println(ds)
	if math.Abs(ds-53670.85) > 0.01 {
		t.Log(ds)
		t.Fail()
	}
}
