package glicko2

import "math"

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
	Volatility float64
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
