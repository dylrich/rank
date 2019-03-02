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

	// ConverganceTolerance is ...
	ConverganceTolerance = 0.000001

	// SystemConstant is ...
	SystemConstant = 0.6
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

func (p *Player) g() float64 {
	return 1 / math.Sqrt(1+(3*math.Pow(toPhi(p.Deviation), 2)/math.Pow(math.Pi, 2)))
}

func (p *Player) e(o *Player) float64 {
	return 1 / (1 + math.Pow(math.E, -o.g()*(toMu(p.Rating)-toMu(o.Rating))))
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

// Win is ...
func (p *Player) Win(o *Player) Outcome {
	p.addResult(o, 1)
	outcome := p.getOutcome()
	return outcome
}

func (p *Player) getOutcome() Outcome {
	mu := toMu(p.Rating)
	phi := toPhi(p.Deviation)
	ti := totalImpact(&p.History)
	variance := variance(ti)
	// delta := delta(variance, &p.History)
	volatility := volatility(p.Volatility, variance, phi, &p.History)
	pp := phiPrime(rd(phi, volatility), variance)
	deviation := fromPhi(pp)
	rating := fromMu(muPrime(mu, pp, ti))
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

func delta(variance float64, history *[]Result) float64 {
	td := 0.0
	for _, result := range *history {
		td += resultScore(result.G, result.Score, result.E)
	}

	return variance * td
}

func resultScore(g, s, e float64) float64 {
	return g * (s - e)
}

func volatility(sigma, variance, phi float64, history *[]Result) float64 {
	var b float64
	var a float64
	a, b = initializeComparison(sigma, variance, phi, history)
	for math.Abs(b-a) > ConverganceTolerance {

	}
	return math.Pow(math.E, (a / 2))
}

func initializeComparison(sigma, variance, phi float64, history *[]Result) (float64, float64) {
	var b float64
	var a float64
	a = alpha(sigma)
	deltaSquared := math.Pow(delta(variance, history), 2)
	if deltaSquared > math.Pow(phi, 2)+variance {
		b = math.Log(deltaSquared - math.Pow(phi, 2) - variance)
	}

	return a, b

}

func alpha(sigma float64) float64 {
	return math.Log(math.Pow(sigma, 2))
}

func fromMu(mu float64) float64 {
	return (173.7178 * mu) + 1500
}

func fromPhi(phi float64) float64 {
	return 173.7178 * phi
}
