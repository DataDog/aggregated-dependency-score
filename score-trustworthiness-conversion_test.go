package aggregdepscore

import (
	"math"
	"testing"
)

func TestDefaultScoreTrustworthinessConverter(t *testing.T) {
	c := DefaultScoreTrustworthinessConverter{}
	allowedError := 1e-13

	for score := 0; score <= 10; score++ {
		x := float64(score)
		y := c.ScoreFromTrustworthiness(c.TrustworthinessFromScore(x))

		if math.Abs(y-x) > allowedError {
			t.Fatalf("failed inverse test for %d: got %g", score, y)
		}
	}
}
