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
	// SystemConstant is ...
	SystemConstant = 0.6

	// ConverganceTolerance is ...
	ConverganceTolerance = 0.000001
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

	return &Player{Rating: p.InitialRating, Deviation: p.InitialDeviation, Parameters: p}
}

func toPhi(deviation float64) float64 {
	return deviation / 173.7178
}

func toMu(rating float64) float64 {
	return (rating - 1500) / 173.7178
}

func variance(ti float64) float64 {
	return math.Pow(ti, -1)

}

func toG(deviation float64) float64 {
	return 1 / math.Sqrt(1+(3*math.Pow(toPhi(deviation), 2)/math.Pow(math.Pi, 2)))
}

func toE(playerRating, opponentRating, opponentG float64) float64 {
	return 1 / (1 + math.Pow(math.E, -opponentG*(toMu(playerRating)-toMu(opponentRating))))
}

func (p *Player) addResult(o *Player, score float64) {
	var r Result
	r.Deviation = o.Deviation
	r.Rating = o.Rating
	r.Score = score
	g := toG(o.Parameters.InitialDeviation)
	r.G = g
	r.E = toE(p.Parameters.InitialRating, o.Parameters.InitialRating, g)
	p.History = append(p.History, r)
}

// Win is ...
func (p *Player) Win(o *Player) Outcome {
	p.addResult(o, 1)
	outcome := p.getOutcome()
	p.Deviation = outcome.Deviation
	p.Rating = outcome.Rating
	p.Volatility = outcome.Volatility
	return outcome
}

// Lose is ...
func (p *Player) Lose(o *Player) Outcome {
	p.addResult(o, 0)
	outcome := p.getOutcome()
	p.Deviation = outcome.Deviation
	p.Rating = outcome.Rating
	p.Volatility = outcome.Volatility
	return outcome
}

// Draw is ...
func (p *Player) Draw(o *Player) Outcome {
	p.addResult(o, 0.5)
	outcome := p.getOutcome()
	p.Deviation = outcome.Deviation
	p.Rating = outcome.Rating
	p.Volatility = outcome.Volatility
	return outcome
}

func (p *Player) getOutcome() Outcome {
	mu := toMu(p.Parameters.InitialRating)
	phi := toPhi(p.Parameters.InitialDeviation)
	ti := totalImpact(&p.History)
	ts := totalResultScore(&p.History)
	variance := variance(ti)
	delta := delta(variance, ts)
	volatility := volatility(p.Parameters.InitialVolatility, variance, phi, delta)
	pp := phiPrime(rd(phi, volatility), variance)
	deviation := fromPhi(pp)
	rating := fromMu(muPrime(mu, pp, ts))
	return Outcome{
		Rating:          rating,
		RatingDelta:     rating - p.Rating,
		Deviation:       deviation,
		DeviationDelta:  deviation - p.Deviation,
		Volatility:      volatility,
		VolatilityDelta: volatility - p.Volatility,
	}
}

func totalImpact(history *[]Result) float64 {
	tv := 0.0
	for _, result := range *history {
		tv += impact(result.G, result.E)
	}
	return tv
}

func rd(phi, volatility float64) float64 {
	return math.Sqrt(math.Pow(phi, 2) + math.Pow(volatility, 2))
}

func phiPrime(rd, variance float64) float64 {
	return 1 / math.Sqrt((1/math.Pow(rd, 2) + (1 / variance)))
}

func muPrime(mu, phi, ti float64) float64 {
	return mu + math.Pow(phi, 2)*ti
}

// Reset is ...
func (p *Player) Reset() {
	p.History = []Result{}
	p.Deviation = p.Parameters.InitialDeviation
	p.Rating = p.Parameters.InitialRating
	p.Volatility = p.Parameters.InitialVolatility
}

func impact(g, e float64) float64 {
	return math.Pow(g, 2) * e * (1 - e)
}

func delta(variance, resultScore float64) float64 {
	return variance * resultScore
}

func totalResultScore(history *[]Result) float64 {
	ts := 0.0
	for _, result := range *history {
		ts += resultScore(result.G, result.Score, result.E)
	}
	return ts
}

func resultScore(g, s, e float64) float64 {
	return g * (s - e)
}

func volatility(sigma, variance, phi, delta float64) float64 {
	var A, B, C, fa, fb, fc float64
	a := toAlpha(sigma)
	A, B = initializeComparison(sigma, variance, phi, delta, a)
	fa = illinois(A, phi, variance, a, delta)
	fb = illinois(B, phi, variance, a, delta)
	for math.Abs(B-A) > ConverganceTolerance {
		C = A + (A-B)*fa/(fb-fa)
		fc = illinois(C, phi, variance, a, delta)
		if 0 > (fc * fb) {
			A = B
			fa = fb
		} else {
			fa = fa / 2
		}
		B = C
		fb = fc
	}
	return math.Pow(math.E, (A / 2))
}

func initializeComparison(sigma, variance, phi, delta, a float64) (float64, float64) {
	var A, B float64
	A = a
	deltaSquared := math.Pow(delta, 2)
	if deltaSquared > (math.Pow(phi, 2) + variance) {
		B = math.Log(deltaSquared - math.Pow(phi, 2) - variance)
		return A, B
	}
	k := 1.0
	for 0 > illinois(a-k*SystemConstant, phi, variance, a, delta) {
		k++
	}
	B = a - k*SystemConstant
	return A, B
}

func illinois(x, phi, variance, alpha, delta float64) float64 {
	ex := math.Pow(math.E, x)
	phiSquared := math.Pow(phi, 2)
	left := ex * (math.Pow(delta, 2) - phiSquared - variance - ex) / (2 * math.Pow(phiSquared+variance+ex, 2))
	right := (x - alpha) / math.Pow(SystemConstant, 2)
	return left - right
}

func toAlpha(sigma float64) float64 {
	return math.Log(math.Pow(sigma, 2))
}

func fromMu(mu float64) float64 {
	return (173.7178 * mu) + 1500
}

func fromPhi(phi float64) float64 {
	return 173.7178 * phi
}
