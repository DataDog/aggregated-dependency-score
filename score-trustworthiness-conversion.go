package aggregdepscore

type ScoreTrustworthinessConverter interface {
	ScoreFromTrustworthiness(trustworthiness float64) float64
	TrustworthinessFromScore(score float64) float64
}

type DefaultScoreTrustworthinessConverter struct{}

// compile-time interface checks
var _ ScoreTrustworthinessConverter = &DefaultScoreTrustworthinessConverter{}

func (c *DefaultScoreTrustworthinessConverter) ScoreFromTrustworthiness(trustworthiness float64) float64 {
	// TODO
	return trustworthiness
}

func (c *DefaultScoreTrustworthinessConverter) TrustworthinessFromScore(score float64) float64 {
	// TODO
	return score
}
