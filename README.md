# Rating

A collection of Go packages implementing various rating system measures.

## Code status

Elo and Glicko are probably fine to use, but the code will undergo many improvements in the future, including more and better tests as well as better internal structuring. Glicko2 is not yet implemented.

## A note on concurrency

This library will not protect against race conditions and assumes that player data is only accessed one at a time. If you need to support concurrent writes to player data (e.g. two different results occurred at the same time), you will need to implement a mutex in your own application.

## System descriptions

### Elo

Elo is one of the most widely implemented and studied rating systems that exists. Many games use Elo under the hood for their ranking systems - notably chess, for which the system was originally designed. The [Wikipedia article on Elo](https://en.wikipedia.org/wiki/Elo_rating_system) has a ton of useful information, including an example implementation.

### Glicko

The Glicko system was created by [Mark Glickman](http://www.glicko.net/) to be an improvement over Elo in many situations. In fact, Elo is just a special case of the Glicko system. The motivation and background can be found on [Glickman's website](http://www.glicko.net/glicko/glicko.pdf). The implemtation in this repository is based on that paper, a [1999 paper by Glickman](http://www.glicko.net/research/gdescrip.pdf), as well as [the original publication](http://www.glicko.net/research/glicko.pdf) in Applied Statistics.

### Glicko2

Under development.

## Development

#### Install

```bash
git clone git@github.com:dylrich/rating.git && cd rating
go get -u ./...
```

#### Run tests

```bash
mage -v test
```

#### Build

```
mage build
```
