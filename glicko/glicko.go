package glicko

import (
	"math"
)

const (

	//DefaultInitialRD is ...
	DefaultInitialRD = 350

	// DefaultInitialRating is ...
	DefaultInitialRating = 1500
)

var (
	q = math.Ln10 / 400
)

// Result is ...
type Result struct {
	Rating, RD, GRD, E, Score float64
}

// Player is ...
type Player struct {
	Rating     float64
	RD         float64
	History    []Result
	Parameters Parameters
}

// Parameters is ...
type Parameters struct {
	InitialRD, InitialRating, C float64
}

// Outcome is ...
type Outcome struct {
	Rating, RatingDelta, RD, RDDelta float64
}

// NewPlayer is ...
func NewPlayer(p Parameters) *Player {
	if &p.InitialRD == nil || &p.InitialRating == nil {
		p.InitialRD = 350
		p.InitialRating = 1500
	}
	return &Player{Rating: p.InitialRating, RD: p.InitialRD, Parameters: p}
}

// Win is ...
func (p *Player) Win(o *Player) *Outcome {
	p.addResult(o, 1)
	ratingPrime := p.ratingPrime()
	rdPrime := p.rdPrime()
	ratingDelta := ratingPrime - p.Rating
	rdDelta := rdPrime - p.RD
	p.Rating = ratingPrime
	p.RD = rdPrime
	return &Outcome{
		Rating:      p.Rating,
		RatingDelta: ratingDelta,
		RD:          rdPrime,
		RDDelta:     rdDelta,
	}
}

func (p *Player) addResult(o *Player, score float64) {
	var r Result
	r.RD = o.RD
	r.Rating = o.Rating
	r.Score = score
	r.GRD = o.gRD()
	r.E = p.e(o)
	p.History = append(p.History, r)
}

func (p *Player) e(o *Player) float64 {
	return 1 / (1 + math.Pow(10, (-o.gRD()*ratingDelta(p.Rating, o.Rating)/400)))
}

func ratingDelta(r1, r2 float64) float64 {
	return r1 - r2
}

func (p *Player) dsquared() float64 {
	ti := 0.0
	for _, r := range p.History {
		ti += impact(r.GRD, r.E)
	}
	return math.Pow(math.Pow(q, 2)*ti, -1)
}

func impact(grd, e float64) float64 {
	return math.Pow(grd, 2) * e * (1 - e)
}

func (p *Player) gRD() float64 {
	return 1 / math.Sqrt(1+(3*math.Pow(q, 2)*math.Pow(p.RD, 2)/math.Pow(math.Pi, 2)))
}

func (p *Player) ratingsDeviation() float64 {
	return math.Min(math.Sqrt(math.Pow(p.RD, 2)+math.Pow(p.Parameters.C, 2)), p.Parameters.InitialRD)
}

func (p *Player) ratingPrime() float64 {
	adjustment := 0.0
	for _, r := range p.History {
		adjustment += adjust(r.GRD, r.E, r.Score)
	}
	return p.Rating + (q/p.deviationAdjustment())*adjustment
}

func (p *Player) rdPrime() float64 {
	return math.Sqrt(math.Pow(p.deviationAdjustment(), -1))
}

func (p *Player) deviationAdjustment() float64 {
	return (1 / math.Pow(p.RD, 2)) + (1 / p.dsquared())
}

func adjust(grd, e, score float64) float64 {
	return grd * (score - e)
}
