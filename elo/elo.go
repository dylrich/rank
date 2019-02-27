package elo

import (
	"math"
)

const (

	// DefaultKFactor is ...
	DefaultKFactor = 32

	// DefaultDeviation is ...
	DefaultDeviation = 400

	// DefaultInitialRanking is ...
	DefaultInitialRanking = 1500
)

// Player is ...
type Player struct {
	Ranking    float64
	Parameters Parameters
}

// Parameters is ...
type Parameters struct {
	K, D, InitialRanking float64
}

// Outcome is ...
type Outcome struct {
	Ranking, Delta float64
}

// NewPlayer is ...
func NewPlayer(p Parameters) *Player {
	return &Player{Ranking: p.InitialRanking, Parameters: p}
}

// Win is ...
func (p *Player) Win(o *Player) *Outcome {
	t := p.transform()
	e := p.expectation(t, o.transform())
	delta := p.delta(e, 1)
	p.Ranking = delta + p.Ranking
	return &Outcome{
		Ranking: p.Ranking,
		Delta:   delta,
	}
}

// Lose is ...
func (p *Player) Lose(o *Player) *Outcome {
	t := p.transform()
	e := p.expectation(t, o.transform())
	delta := p.delta(e, 0)
	p.Ranking = delta + p.Ranking
	return &Outcome{
		Ranking: p.Ranking,
		Delta:   delta,
	}
}

// Draw is ...
func (p *Player) Draw(o *Player) *Outcome {
	t := p.transform()
	e := p.expectation(t, o.transform())
	delta := p.delta(e, .5)
	p.Ranking = delta + p.Ranking
	return &Outcome{
		Ranking: p.Ranking,
		Delta:   delta,
	}
}

func (p *Player) delta(e, s float64) float64 {
	return p.Parameters.K * (s - e)
}

func (p *Player) transform() float64 {
	return math.Pow(10, (p.Ranking / p.Parameters.D))
}

func (p *Player) expectation(t1, t2 float64) float64 {
	return t1 / (t1 + t2)
}
