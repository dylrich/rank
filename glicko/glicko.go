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
	Rating, Delta float64
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
	var outcome Outcome
	return &outcome
}

func (p *Player) addResult(o *Player) {
	var r Result
	r.RD = o.RD
	r.Rating = o.Rating
}

func (p *Player) dsquared() float64 {
	math.Pow(q, 2)
}

func impact(grd, e float64) float64 {

}

func (p *Player) gRD() float64 {
	return 1 / math.Sqrt(1+(3*math.Pow(q, 2)*math.Pow(p.RD, 2)/math.Pow(math.Pi, 2)))
}

func (p *Player) ratingsDeviation() float64 {
	return math.Min(math.Sqrt(math.Pow(p.RD, 2)+math.Pow(p.Parameters.C, 2)), p.Parameters.InitialRD)
}
