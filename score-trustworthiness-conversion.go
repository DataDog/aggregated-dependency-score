package aggregdepscore

import "math"

const (
	minTrustworthiness = 0.8
	// trustworthinessOffset is noted as "k" in the design paper
	trustworthinessOffset = 60.0
)

type ScoreTrustworthinessConverter interface {
	ScoreFromTrustworthiness(trustworthiness float64) float64
	TrustworthinessFromScore(score float64) float64
}

type DefaultScoreTrustworthinessConverter struct{}

// compile-time interface checks
var _ ScoreTrustworthinessConverter = &DefaultScoreTrustworthinessConverter{}

func (c *DefaultScoreTrustworthinessConverter) ScoreFromTrustworthiness(trustworthiness float64) float64 {
	t := trustworthiness
	k := trustworthinessOffset

	if t < minTrustworthiness {
		return 0
	}

	return ((1 - math.Pow(k, (1-(1-t)/0.2))) /
		(1 - k))
}

func (c *DefaultScoreTrustworthinessConverter) TrustworthinessFromScore(score float64) float64 {
	k := trustworthinessOffset
	min := minTrustworthiness

	return 1 - (1-min)*(1-math.Log(1+(k-1)*score)/
		math.Log(k))
}
