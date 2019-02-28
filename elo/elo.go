package elo

import (
	"math"
)

const (

	// DefaultKFactor is ...
	DefaultKFactor = 32

	// DefaultDeviation is ...
	DefaultDeviation = 400

	// DefaultInitialRating is ...
	DefaultInitialRating = 1500
)

// Player is ...
type Player struct {
	Rating     float64
	Parameters Parameters
}

// Parameters is ...
type Parameters struct {
	K, D, InitialRating float64
}

// Outcome is ...
type Outcome struct {
	Rating, Delta float64
}

// NewPlayer is ...
func NewPlayer(p Parameters) *Player {
	return &Player{Rating: p.InitialRating, Parameters: p}
}

// Win is ...
func (p *Player) Win(o *Player) *Outcome {
	t := p.transform()
	e := p.expectation(t, o.transform())
	delta := p.delta(e, 1)
	p.Rating = delta + p.Rating
	return &Outcome{
		Rating: p.Rating,
		Delta:  delta,
	}
}

// Lose is ...
func (p *Player) Lose(o *Player) *Outcome {
	t := p.transform()
	e := p.expectation(t, o.transform())
	delta := p.delta(e, 0)
	p.Rating = delta + p.Rating
	return &Outcome{
		Rating: p.Rating,
		Delta:  delta,
	}
}

// Draw is ...
func (p *Player) Draw(o *Player) *Outcome {
	t := p.transform()
	e := p.expectation(t, o.transform())
	delta := p.delta(e, .5)
	p.Rating = delta + p.Rating
	return &Outcome{
		Rating: p.Rating,
		Delta:  delta,
	}
}

func (p *Player) delta(e, s float64) float64 {
	return p.Parameters.K * (s - e)
}

func (p *Player) transform() float64 {
	return math.Pow(10, (p.Rating / p.Parameters.D))
}

func (p *Player) expectation(t1, t2 float64) float64 {
	return t1 / (t1 + t2)
}
