package elo

import (
	"math"
)

const (

	// DefaultInitialRating is the standard value for an initial rating for players that have no result history from the previous rating period.
	DefaultInitialRating = 1500
)

var (
	// KFactor regulates how much new results impact the player's rating.
	KFactor = 32.0

	// D represents the Elo standard deviation value
	D = 400.0
)

// Player represents an individual participant in the competition. The Player struct contains the Rating measure, which is the Elo system's estimation of how skilled that player is. This is a moment-in-time snapshot, and will be updated on any new results for that player. The Parameters attribute contains initial values for that player which can be used to reconstruct the player's current rating from scratch when combined with the History data. Parameters should be altered at the beginning of a new rating period to be the final Rating from the previous period.
type Player struct {
	Rating     float64
	History    []Result
	Parameters Parameters
}

// Parameters contains initial Rating for a player. This is set on instantiation of the player.
type Parameters struct {
	InitialRating float64
}

// Result contains the important information from a match that has occurred. The information is used to calculate new ratings when new results are added.
type Result struct {
	Rating, Score float64
}

// Outcome is a snapshot of the current state for a player, including the delta value for this result's Rating change. This information can be passed to users to give them an idea of how much the most recent result has impacted their ranking criteria.
type Outcome struct {
	Rating, RatingDelta float64
}

// NewPlayer is used to instantiate a new Player object based on the input parameters. If any of the parameters are nil, they will be automatically populated with the default values.
func NewPlayer(p Parameters) *Player {
	if &p.InitialRating == nil {
		p.InitialRating = DefaultInitialRating
	}
	return &Player{Rating: p.InitialRating, Parameters: p}
}

// Win is called when a player has won a match against another player, earning an Elo score of 1. This function will handle updating the calling Player only. To add the loss to the opponent's rating, call Opponent.Lose(Player) as appropriate.
func (p *Player) Win(opponentRating float64) *Outcome {
	p.addResult(1, opponentRating)
	outcome := p.getOutcome(1, opponentRating)
	p.Rating = outcome.Rating
	return &outcome
}

// Lose is called when a player has won a match against another player, earning an Elo score of 0. This function will handle updating the calling Player only. To add the loss to the opponent's rating, call Opponent.Lose(Player) as appropriate.
func (p *Player) Lose(opponentRating float64) *Outcome {
	p.addResult(0, opponentRating)
	outcome := p.getOutcome(0, opponentRating)
	p.Rating = outcome.Rating
	return &outcome
}

// Draw is called when a player has won a match against another player, earning an Elo score of 0.5. This function will handle updating the calling Player only. To add the draw record to the opponent's rating, call Opponent.Draw(Player) as appropriate.
func (p *Player) Draw(opponentRating float64) *Outcome {
	p.addResult(0.5, opponentRating)
	outcome := p.getOutcome(0.5, opponentRating)
	p.Rating = outcome.Rating
	return &outcome
}

// Reset will wipe the calling Player's history completely, and revert the current Rating to its initial value.
func (p *Player) Reset() {
	p.History = []Result{}
	p.Rating = p.Parameters.InitialRating
}

// NewPeriod takes the calling Player's current Rating and sets it as the new initital rating before resetting the player's history to empty.
func (p *Player) NewPeriod() {
	p.Parameters.InitialRating = p.Rating
	p.Reset()
}

func (p *Player) getOutcome(score, opponentRating float64) Outcome {
	rd := ratingDelta(score, expectation(transform(p.Rating), transform(opponentRating)))
	return Outcome{
		Rating:      rd + p.Rating,
		RatingDelta: rd,
	}
}

func (p *Player) addResult(score, rating float64) {
	var r Result
	r.Rating = rating
	r.Score = score
	p.History = append(p.History, r)
}

func ratingDelta(score, expectation float64) float64 {
	return KFactor * (score - expectation)
}

func transform(rating float64) float64 {
	return math.Pow(10, (rating / D))
}

func expectation(t1, t2 float64) float64 {
	return t1 / (t1 + t2)
}
