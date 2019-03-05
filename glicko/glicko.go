package glicko

import (
	"math"
)

const (

	// DefaultInitialDeviation is the standard value for an initial deviation for players that have no result history from the previous rating period.
	DefaultInitialDeviation = 350

	// DefaultInitialRating is the standard value for an initial rating for players that have no result history from the previous rating period.
	DefaultInitialRating = 1500
)

var (
	// C is a constant that governs the increase in uncertainty between rating periods.
	C = 40
	q = math.Ln10 / 400
)

// Player represents an individual participant in the competition. The Player struct contains the Rating and Deviation measures which all compose the Glicko system's estimation of how skilled that player is as well as how reliable that estimation is. These values are all moment-in-time snapshots, and will be updated on any new results for that player. The Parameters attribute contains initial values for that player which can be used to reconstruct the player's current rating from scratch when combined with the History data. Parameters should be altered at the beginning of a new rating period to be the final Rating and Deviation values of the previous period.
type Player struct {
	Rating     float64
	Deviation  float64
	History    []Result
	Parameters Parameters
}

// Parameters contains initial values for a player. These are set on instantiation of the player, and can be altered later by using the Player.NewPeriod() method.
type Parameters struct {
	InitialDeviation, InitialRating float64
}

// Result contains the important information from a match that has occurred. The information is used to calculate new ratings when new results are added.
type Result struct {
	Rating, Deviation, G, E, Score float64
}

// Outcome is a snapshot of the current state for a player, including delta values for each Deviation and Rating change. This information can be passed to users to give them an idea of how much the most recent result has impacted their ranking criteria.
type Outcome struct {
	Rating, RatingDelta, Deviation, DeviationDelta float64
}

// NewPlayer is used to instantiate a new Player object based on the input parameters. If any of the parameters are nil, they will be automatically populated with the default values.
func NewPlayer(p Parameters) *Player {
	if &p.InitialDeviation == nil {
		p.InitialDeviation = DefaultInitialDeviation
	}
	if &p.InitialRating == nil {
		p.InitialRating = DefaultInitialRating
	}
	return &Player{Rating: p.InitialRating, Deviation: p.InitialDeviation, Parameters: p}
}

// Win is called when a player has won a match against another player, earning a Glicko score of 1. This function will handle adding the result to the history of the player who wins only. To add the loss record to the opponent's history, call Opponent.Lose(Player) as appropriate.
func (p *Player) Win(rating, deviation float64) Outcome {
	p.addResult(rating, deviation, 1)
	outcome := p.getOutcome()
	p.Rating = outcome.Rating
	p.Deviation = outcome.Deviation
	return outcome
}

// Lose is called when a player has won a match against another player, earning a Glicko score of 0. This function will handle adding the result to the history of the player who loses only. To add the win record to the opponent's history, call Opponent.Win(Player) as appropriate.
func (p *Player) Lose(rating, deviation float64) Outcome {
	p.addResult(rating, deviation, 0)
	outcome := p.getOutcome()
	p.Rating = outcome.Rating
	p.Deviation = outcome.Deviation
	return outcome
}

// Draw is called when a player has tied in a match against another player, earning a Glicko score of 0.5. This function will handle adding the result to the history of the player this method is called on only. To add the draw record to the opponent's history, call Opponent.Draw(Player) as appropriate.
func (p *Player) Draw(rating, deviation float64) Outcome {
	p.addResult(rating, deviation, 0.5)
	outcome := p.getOutcome()
	p.Rating = outcome.Rating
	p.Deviation = outcome.Deviation
	return outcome
}

// Reset will wipe the calling Player's history completely, and revert the current Rating and Deviation to the initial values.
func (p *Player) Reset() {
	p.History = []Result{}
	p.Deviation = p.Parameters.InitialDeviation
	p.Rating = p.Parameters.InitialRating
}

// NewPeriod takes the calling Player's current Rating and Deviation, and sets them as the new initital values before resetting the player's history to empty.
func (p *Player) NewPeriod() {
	p.Parameters.InitialDeviation = p.Deviation
	p.Parameters.InitialRating = p.Rating
	p.Reset()
}

func (p *Player) getOutcome() Outcome {
	ds := deviationScore(p.Parameters.InitialDeviation, &p.History)
	rp := ratingPrime(p.Parameters.InitialRating, ds, &p.History)
	dp := deviationPrime(ds)
	return Outcome{
		Rating:         rp,
		RatingDelta:    rp - p.Rating,
		Deviation:      dp,
		DeviationDelta: dp - p.Deviation,
	}
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

func toG(deviation float64) float64 {
	return 1 / math.Sqrt(1+(3*math.Pow(q, 2)*math.Pow(deviation, 2)/math.Pow(math.Pi, 2)))
}

func toE(playerRating, opponentRating, opponentG float64) float64 {
	return 1 / (1 + math.Pow(10, -opponentG*(playerRating-opponentRating)/400))
}

func dsquared(history *[]Result) float64 {
	return math.Pow(math.Pow(q, 2)*totalImpact(history), -1)
}

func ratingPrime(rating, deviationScore float64, history *[]Result) float64 {
	return rating + (q/deviationScore)*totalResultScore(history)
}

func deviationPrime(deviationScore float64) float64 {
	return math.Sqrt(math.Pow(deviationScore, -1))
}

func deviationScore(deviation float64, history *[]Result) float64 {
	return (1 / math.Pow(deviation, 2)) + (1 / dsquared(history))
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
