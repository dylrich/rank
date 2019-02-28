package glicko

import (
	"math"
)

const (

	//DefaultInitialRD is ...
	DefaultInitialRD = 350

	// DefaultInitialRanking is ...
	DefaultInitialRanking = 1500
)

var (
	q = math.Ln10 / 400
)

// Player is ...
type Player struct {
	Ranking float64
	RD      float64

	Parameters Parameters
}

// Parameters is ...
type Parameters struct {
	InitialRD, InitialRanking, C float64
}

// Outcome is ...
type Outcome struct {
	Ranking, Delta float64
}

// NewPlayer is ...
func NewPlayer(p Parameters) *Player {
	if &p.InitialRD == nil || &p.InitialRanking == nil {
		p.InitialRD = 350
		p.InitialRanking = 1500
	}
	return &Player{Ranking: p.InitialRanking, RD: p.InitialRD, Parameters: p}
}

// Win is ...
func (p *Player) Win(o *Player) *Outcome {
	var outcome Outcome
	return &outcome
}

// func (p *Player) dsquared() float64 {
// 	math.Pow(q, 2)
// }

func (p *Player) gRD() float64 {
	return 1 / math.Sqrt(1+(3*math.Pow(q, 2)*math.Pow(p.RD, 2)/math.Pow(math.Pi, 2)))
}

func (p *Player) ratingsDeviation() float64 {
	return math.Min(math.Sqrt(math.Pow(p.RD, 2)+math.Pow(p.Parameters.C, 2)), p.Parameters.InitialRD)
}
