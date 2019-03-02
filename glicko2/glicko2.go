package glicko2

import (
	"math"
)

const (

	// DefaultInitialDeviation is ...
	DefaultInitialDeviation = 350

	// DefaultInitialRating is ...
	DefaultInitialRating = 1500

	// DefaultInitialVolatility is ...
	DefaultInitialVolatility = 0.06
)

var (
	q = math.Ln10 / 400
)

// Result is ...
type Result struct {
	Rating, Deviation, G, E, Score float64
}

// Player is ...
type Player struct {
	Rating     float64
	Deviation  float64
	Volatility float64
	mu         float64
	phi        float64
	History    []Result
	Parameters Parameters
}

// Parameters is ...
type Parameters struct {
	InitialDeviation, InitialRating, InitialVolatility, C float64
}

// Outcome is ...
type Outcome struct {
	Rating, RatingDelta, Deviation, DeviationDelta, Volatility, VolatilityDelta float64
}

// NewPlayer is ...
func NewPlayer(p Parameters) *Player {
	if &p.InitialDeviation == nil || &p.InitialRating == nil || &p.InitialVolatility == nil {
		p.InitialDeviation = DefaultInitialDeviation
		p.InitialRating = DefaultInitialRating
		p.InitialVolatility = DefaultInitialVolatility
	}

	return &Player{Rating: p.InitialRating, Deviation: p.InitialDeviation, Parameters: p, mu: calcMu(p.InitialRating), phi: calcPhi(p.InitialDeviation)}
}

func calcPhi(deviation float64) float64 {
	return deviation / 173.7178
}

func calcMu(rating float64) float64 {
	return rating / 173.7178
}

// func (p *Player) variation() float64 {
// }

func (p *Player) g() float64 {
	return 1 / math.Sqrt(1+(3*math.Pow(p.phi, 2)/math.Pow(math.Pi, 2)))
}

func (p *Player) e(o *Player) float64 {
	return 1 / (1 + math.Pow(math.E, -o.g()*(p.mu-o.mu)))
}

func (p *Player) addResult(o *Player, score float64) {
	var r Result
	r.Deviation = o.Deviation
	r.Rating = o.Rating
	r.Score = score
	r.G = o.g()
	r.E = p.e(o)
	p.History = append(p.History, r)
}

// func (p *Player) delta() float64 {

// }

// func (p *Player) volatilityPrime() float64 {

// }

// func (p *Player) phiPrime() float64 {

// }

// func (p *Player) muPrime() float64 {

// }

func (p *Player) convertRating() float64 {
	return (173.7178 * p.mu) + 1500
}

func (p *Player) convertDeviation() float64 {
	return 173.7178 * p.phi
}
