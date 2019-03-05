package glicko2

import (
	"math"
)

const (

	// DefaultInitialDeviation is the standard value for an initial deviation for players that have no result history from the previous rating period.
	DefaultInitialDeviation = 350

	// DefaultInitialRating is the standard value for an initial rating for players that have no result history from the previous rating period.
	DefaultInitialRating = 1500

	// DefaultInitialVolatility is the standard value for an initial volatility for players that have no result history from the previous rating period.
	DefaultInitialVolatility = 0.06
)

var (
	q = math.Ln10 / 400

	// SystemConstant (τ) constrains the change in volatility over time. It  needs to be set  prior  to  application  of  the  system. Reasonable  choices  are  between  0.3  and  1.2 ,though the system should be tested to decide which value results in greatest predictive accuracy. Smaller values of τ prevent the volatility measures from changing by large amounts, which in turn prevent enormous changes in ratings based on very improbable results. If  the  application  of  Glicko2  is  expected  to  involve  extremely  improbable collections of game outcomes, then τ should be set to a small value, even as small as τ = 0.
	SystemConstant = 0.6

	// ConverganceTolerance (ε) is the value that the illinois algorithm uses to detect whether A and B have converged to each other.
	ConverganceTolerance = 0.000001
)

// Player represents an individual participant in the competition. The Player struct contains the Rating, Deviation, and Volatility measures which all compose the Glicko2 system's estimation of how skilled that player is as well as how reliable that estimation is. These values are all moment-in-time snapshots, and will be updated on any new results for that player. The Parameters attribute contains initial values for that player which can be used to reconstruct the player's current rating from scratch when combined with the History data. Parameters should be altered at the beginning of a new rating period to be the final Rating, Deviation, and Volatility values of the previous period.
type Player struct {
	Rating     float64
	Deviation  float64
	Volatility float64
	History    []Result
	Parameters Parameters
}

// Parameters contains initial values for a player. These are set on instantiation of the player, and can be altered later by using the Player.NewPeriod() method.
type Parameters struct {
	InitialDeviation, InitialRating, InitialVolatility float64
}

// Result contains the important information from a match that has occurred. The information is used to calculate new ratings when new results are added.
type Result struct {
	Rating, Deviation, G, E, Score float64
}

// Outcome is a snapshot of the current state for a player, including delta values for each Deviation, Rating, and Volatility change. This information can be passed to users to give them an idea of how much the most recent result has impacted their ranking criteria.
type Outcome struct {
	Rating, RatingDelta, Deviation, DeviationDelta, Volatility, VolatilityDelta float64
}

// NewPlayer is used to instantiate a new Player object based on the input parameters. If any of the parameters are nil, they will be automatically populated with the default values.
func NewPlayer(p Parameters) *Player {
	if &p.InitialDeviation == nil {
		p.InitialDeviation = DefaultInitialDeviation
	}
	if &p.InitialRating == nil {
		p.InitialRating = DefaultInitialRating
	}
	if &p.InitialVolatility == nil {
		p.InitialVolatility = DefaultInitialVolatility
	}

	return &Player{Rating: p.InitialRating, Deviation: p.InitialDeviation, Volatility: p.InitialVolatility, Parameters: p}
}

// Win is called when a player has won a match against another player, earning a Glicko2 score of 1. This function will handle adding the result to the history of the player who wins only. To add the loss record to the opponent's history, call Opponent.Lose(Player) as appropriate.
func (p *Player) Win(rating, deviation float64) Outcome {
	p.addResult(rating, deviation, 1)
	outcome := p.getOutcome()
	p.Deviation = outcome.Deviation
	p.Rating = outcome.Rating
	p.Volatility = outcome.Volatility
	return outcome
}

// Lose is called when a player has won a match against another player, earning a Glicko2 score of 0. This function will handle adding the result to the history of the player who loses only. To add the win record to the opponent's history, call Opponent.Win(Player) as appropriate.
func (p *Player) Lose(rating, deviation float64) Outcome {
	p.addResult(rating, deviation, 0)
	outcome := p.getOutcome()
	p.Deviation = outcome.Deviation
	p.Rating = outcome.Rating
	p.Volatility = outcome.Volatility
	return outcome
}

// Draw is called when a player has tied in a match against another player, earning a Glicko2 score of 0.5. This function will handle adding the result to the history of the player this method is called on only. To add the draw record to the opponent's history, call Opponent.Draw(Player) as appropriate.
func (p *Player) Draw(rating, deviation float64) Outcome {
	p.addResult(rating, deviation, 0.5)
	outcome := p.getOutcome()
	p.Deviation = outcome.Deviation
	p.Rating = outcome.Rating
	p.Volatility = outcome.Volatility
	return outcome
}

// Reset will wipe the calling Player's history completely, and revert the current Rating, Deviation, and Volatility to the initial values.
func (p *Player) Reset() {
	p.History = []Result{}
	p.Deviation = p.Parameters.InitialDeviation
	p.Rating = p.Parameters.InitialRating
	p.Volatility = p.Parameters.InitialVolatility
}

// NewPeriod takes the calling Player's current Rating, Volatility, and Deviation, and sets them as the new initital values before resetting the player's history to empty.
func (p *Player) NewPeriod() {
	p.Parameters.InitialDeviation = p.Deviation
	p.Parameters.InitialRating = p.Rating
	p.Parameters.InitialVolatility = p.Volatility
	p.Reset()
}

func (p *Player) addResult(rating, deviation, score float64) {
	var r Result
	r.Deviation = deviation
	r.Rating = rating
	r.Score = score
	g := toG(deviation)
	r.G = g
	r.E = toE(p.Parameters.InitialRating, rating, g)
	p.History = append(p.History, r)
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

// The Illinois algorithm is a variant of the regula falsi (false position) procedure. The Illinois algorithm is quite stable, reliable, and converges quickly. The algorithm takes advantage of the knowledge that the desired value of σ′ can be sandwiched at the start of the algorithm by the initial choices of A and B.
func illinois(x, phi, variance, alpha, delta float64) float64 {
	ex := math.Pow(math.E, x)
	phiSquared := math.Pow(phi, 2)
	left := ex * (math.Pow(delta, 2) - phiSquared - variance - ex) / (2 * math.Pow(phiSquared+variance+ex, 2))
	right := (x - alpha) / math.Pow(SystemConstant, 2)
	return left - right
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

func delta(variance, resultScore float64) float64 {
	return variance * resultScore
}

func variance(ti float64) float64 {
	return math.Pow(ti, -1)
}

func impact(g, e float64) float64 {
	return math.Pow(g, 2) * e * (1 - e)
}

func totalImpact(history *[]Result) float64 {
	tv := 0.0
	for _, result := range *history {
		tv += impact(result.G, result.E)
	}
	return tv
}

func resultScore(g, s, e float64) float64 {
	return g * (s - e)
}

func totalResultScore(history *[]Result) float64 {
	ts := 0.0
	for _, result := range *history {
		ts += resultScore(result.G, result.Score, result.E)
	}
	return ts
}

func toPhi(deviation float64) float64 {
	return deviation / 173.7178
}

func toMu(rating float64) float64 {
	return (rating - 1500) / 173.7178
}

func toG(deviation float64) float64 {
	return 1 / math.Sqrt(1+(3*math.Pow(toPhi(deviation), 2)/math.Pow(math.Pi, 2)))
}

func toE(playerRating, opponentRating, opponentG float64) float64 {
	return 1 / (1 + math.Pow(math.E, -opponentG*(toMu(playerRating)-toMu(opponentRating))))
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
