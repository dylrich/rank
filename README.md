# Rating

A collection of Go packages implementing the Elo, Glicko, and Glicko 2 rating systems.

## Install

```bash
go get -u github.com/dylrich/rating
```

## Elo

Elo is one of the most widely implemented and studied rating systems that exists. Many games use Elo under the hood for their ranking systems - notably chess, for which the system was originally designed. The [Wikipedia article on Elo](https://en.wikipedia.org/wiki/Elo_rating_system) has a ton of useful information, including an example implementation.

### Usage

```go
package main

import (
    "fmt"
    "github.com/dylrich/rating/elo"
)

func main(){
    params := elo.Parameters{InitialRating: elo.DefaultInitialRating}
    p1 := elo.NewPlayer(params)
    p2 := elo.NewPlayer(params)

    p1Rating := p1.Rating
    p2Rating := p2.Rating

    p1Outcome := p1.Win(p2Rating)
    p2Outcome := p2.Lose(p1Rating)

    fmt.Printf("Player 1's rating is now %v (%v)", p1Outcome.Rating, p1Outcome.RatingDelta)
    fmt.Printf("Player 2's rating is now %v (%v)", p2Outcome.Rating, p2Outcome.RatingDelta)
}
```

### Status

Elo has been tested against known datasets and should be suitable for use in your application. It is currently missing a few features, such as an automatic KFactor calculator, but these will be implemented in the future.

## Glicko

The Glicko system was created by [Mark Glickman](http://www.glicko.net/) to be an improvement over Elo in many situations. In fact, Elo is just a particular case of the Glicko system. The motivation and background can be found on [Glickman's website](http://www.glicko.net/glicko/glicko.pdf). The implemtation in this repository is based on that paper, a [1999 paper by Glickman](http://www.glicko.net/research/gdescrip.pdf), as well as [the original publication](http://www.glicko.net/research/glicko.pdf) in Applied Statistics.

### Usage

```go
package main

import (
    "fmt"
    "github.com/dylrich/rating/glicko"
)

func main(){
    params := glicko.Parameters{
        InitialRating: glicko.DefaultInitialRating,
        InitialDeviation: glicko.DefaultInitialDeviation,
        }
    p1 := glicko.NewPlayer(params)
    p2 := glicko.NewPlayer(params)

    p1Rating := p1.Rating
    p1Deviation := p1.Deviation

    p2Rating := p2.Rating
    p2Deviation := p2.Deviation

    p1Outcome := p1.Win(p2Rating, p2Deviation)
    p2Outcome := p2.Lose(p1Rating, p2Deviation)

    fmt.Printf("Player 1's rating is now %v (%v) with a deviation of %v (%v)", p1Outcome.Rating, p1Outcome.RatingDelta, p1Outcome.Deviation, p1Outcome.DeviationDelta)
    fmt.Printf("Player 2's rating is now %v (%v) with a deviation of %v (%v)", p2Outcome.Rating, p2Outcome.RatingDelta, p2Outcome.Deviation, p2Outcome.DeviationDelta)
}
```

### Status

Glicko has been tested against known datasets and should be suitable for use in your application. It is currently missing a few features, such as an automatic C value calculator and additional rating reporting utilities, but these will be implemented in the future.

## Glicko2

Glicko2 was also created by Mark Glickman to be an improvement over his original Glicko system. Glicko 2 adds a volatility measure, which provides a range that the actual player rating is anticipated to fluctuate. Volatility will be higher for players who have very inconsistent performances, and lower for players who have a steady history of results. This implementation was based on an [example Glickman produced on his website](http://www.glicko.net/glicko/glicko2.pdf) as well as [the original publication](http://www.glicko.net/research/dpcmsv.pdf) in Applied Statistics

### Usage

```go
package main

import (
    "fmt"
    "github.com/dylrich/rating/glicko2"
)

func main(){
    params := glicko2.Parameters{
        InitialRating: glicko2.DefaultInitialRating,
        InitialDeviation: glicko2.DefaultInitialDeviation,
        InitialVolatility: glicko2.DefaultInitialVolatility,
        }
        
    p1 := glicko2.NewPlayer(params)
    p2 := glicko2.NewPlayer(params)

    p1Rating := p1.Rating
    p1Deviation := p1.Deviation

    p2Rating := p2.Rating
    p2Deviation := p2.Deviation

    p1Outcome := p1.Win(p2Rating, p2Deviation)
    p2Outcome := p2.Lose(p1Rating, p2Deviation)

    fmt.Printf("Player 1's rating is now %v (%v) with a deviation of %v (%v) and volatility of %v (%v)", p1Outcome.Rating, p1Outcome.RatingDelta, p1Outcome.Deviation, p1Outcome.DeviationDelta, p1Outcome.Volatility, p1Outcome.VolatilityDelta)
    fmt.Printf("Player 2's rating is now %v (%v) with a deviation of %v (%v) and volatility of %v (%v)", p2Outcome.Rating, p2Outcome.RatingDelta, p2Outcome.Deviation, p2Outcome.DeviationDelta, p2Outcome.Volatility, p2Outcome.VolatilityDelta)
}
```

### Status

Glicko2 has been tested against known datasets and should be suitable for use in your application. It is currently missing a few features, such as an automatic SystemConstant calculator and additional rating reporting utilities, but these will be implemented in the future.

## A note on concurrency

This library will not protect against race conditions and assumes that player data is only accessed one at a time. If you need to support concurrent writes to player data (e.g. two different results occurred at the same time), you will need to implement a mutex in your own application.

## Development

### Install

```bash
git clone git@github.com:dylrich/rating.git && cd rating
go get -u ./...
```

### Run tests

```bash
mage -v test
```

### Run benchmarks

```bash
mage -v benchmark
```

### Current benchmarks

```bash
BenchmarkGlicko-16         10000            222202 ns/op
BenchmarkGlicko2-16        10000            221210 ns/op
```
