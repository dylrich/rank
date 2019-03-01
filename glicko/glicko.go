package glicko

import (
	"math"
)

const (

	//DefaultInitialDeviation is ...
	DefaultInitialDeviation = 350

	// DefaultInitialRating is ...
	DefaultInitialRating = 1500
)

var (
	q = math.Ln10 / 400
)

// Result is ...
type Result struct {
	Rating, Deviation, GDeviation, E, Score float64
}

// Player is ...
type Player struct {
	Rating     float64
	Deviation  float64
	History    []Result
	Parameters Parameters
}

// Parameters is ...
type Parameters struct {
	InitialDeviation, InitialRating, C float64
}

// Outcome is ...
type Outcome struct {
	Rating, RatingDelta, Deviation, DeviationDelta float64
}

// NewPlayer is ...
func NewPlayer(p Parameters) *Player {
	if &p.InitialDeviation == nil || &p.InitialRating == nil {
		p.InitialDeviation = 350
		p.InitialRating = 1500
	}
	return &Player{Rating: p.InitialRating, Deviation: p.InitialDeviation, Parameters: p}
}

// Win is ...
func (p *Player) Win(o *Player) *Outcome {
	p.addResult(o, 1)
	ratingPrime := p.ratingPrime()
	DeviationPrime := p.deviationPrime()
	ratingDelta := ratingPrime - p.Rating
	DeviationDelta := DeviationPrime - p.Deviation
	p.Rating = ratingPrime
	p.Deviation = DeviationPrime
	return &Outcome{
		Rating:         p.Rating,
		RatingDelta:    ratingDelta,
		Deviation:      DeviationPrime,
		DeviationDelta: DeviationDelta,
	}
}

// Loss is ...
func (p *Player) Loss(o *Player) *Outcome {
	p.addResult(o, 0)
	ratingPrime := p.ratingPrime()
	DeviationPrime := p.deviationPrime()
	ratingDelta := ratingPrime - p.Rating
	DeviationDelta := DeviationPrime - p.Deviation
	p.Rating = ratingPrime
	p.Deviation = DeviationPrime
	return &Outcome{
		Rating:         p.Rating,
		RatingDelta:    ratingDelta,
		Deviation:      DeviationPrime,
		DeviationDelta: DeviationDelta,
	}
}

// Draw is ...
func (p *Player) Draw(o *Player) *Outcome {
	p.addResult(o, 0.5)
	ratingPrime := p.ratingPrime()
	DeviationPrime := p.deviationPrime()
	ratingDelta := ratingPrime - p.Rating
	DeviationDelta := DeviationPrime - p.Deviation
	p.Rating = ratingPrime
	p.Deviation = DeviationPrime
	return &Outcome{
		Rating:         p.Rating,
		RatingDelta:    ratingDelta,
		Deviation:      DeviationPrime,
		DeviationDelta: DeviationDelta,
	}
}

func (p *Player) addResult(o *Player, score float64) {
	var r Result
	r.Deviation = o.Deviation
	r.Rating = o.Rating
	r.Score = score
	r.GDeviation = o.gDeviation()
	r.E = p.e(o)
	p.History = append(p.History, r)
}

func (p *Player) e(o *Player) float64 {
	return 1 / (1 + math.Pow(10, (-o.gDeviation()*ratingDelta(p.Rating, o.Rating)/400)))
}

func ratingDelta(r1, r2 float64) float64 {
	return r1 - r2
}

func (p *Player) dsquared() float64 {
	ti := 0.0
	for _, r := range p.History {
		ti += impact(r.GDeviation, r.E)
	}
	return math.Pow(math.Pow(q, 2)*ti, -1)
}

func impact(gDeviation, e float64) float64 {
	return math.Pow(gDeviation, 2) * e * (1 - e)
}

func (p *Player) gDeviation() float64 {
	return 1 / math.Sqrt(1+(3*math.Pow(q, 2)*math.Pow(p.Deviation, 2)/math.Pow(math.Pi, 2)))
}

func (p *Player) ratingsDeviation() float64 {
	return math.Min(math.Sqrt(math.Pow(p.Deviation, 2)+math.Pow(p.Parameters.C, 2)), p.Parameters.InitialDeviation)
}

func (p *Player) ratingPrime() float64 {
	adjustment := 0.0
	for _, r := range p.History {
		adjustment += adjust(r.GDeviation, r.E, r.Score)
	}
	return p.Rating + (q/p.deviationAdjustment())*adjustment
}

func (p *Player) deviationPrime() float64 {
	return math.Sqrt(math.Pow(p.deviationAdjustment(), -1))
}

func (p *Player) deviationAdjustment() float64 {
	return (1 / math.Pow(p.Deviation, 2)) + (1 / p.dsquared())
}

func adjust(gDeviation, e, score float64) float64 {
	return gDeviation * (score - e)
}
