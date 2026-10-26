package app

import (
	stdjson "encoding/json"
	"math"
	"testing"
)

func TestStatsDependencyUpgrade(t *testing.T) {
	decode := func(s string) []SingleReviewMiningResult {
		t.Helper()
		var r []SingleReviewMiningResult
		if err := stdjson.Unmarshal([]byte(s), &r); err != nil {
			t.Fatal(err)
		}
		return r
	}
	t.Run("regression", func(t *testing.T) {
		r := decode(`[{"reviewRating":1,"miningResults":[{"attribute":"taste","score":1}]},{"reviewRating":2,"miningResults":[{"attribute":"taste","score":3}]},{"reviewRating":3,"miningResults":[{"attribute":"taste","score":4}]},{"reviewRating":4,"miningResults":[{"attribute":"taste","score":7}]}]`)
		got := simpleLinearRegressEachAttributeScoreAndRatings([]string{"taste", "missing"}, r)
		if len(got) != 1 || math.Abs(got["taste"].Slope-1.9) > 1e-10 {
			t.Fatalf("unexpected regression: %+v", got)
		}
	})
	t.Run("insufficient samples", func(t *testing.T) {
		if got := simpleLinearRegressEachAttributeScoreAndRatings([]string{"taste"}, nil); len(got) != 0 {
			t.Fatal(got)
		}
		if got := reviewsAttributeTTest([]string{"taste"}, nil); got != nil {
			t.Fatal(got)
		}
	})
	t.Run("Welch fallback after variance test error", func(t *testing.T) {
		r := decode(`[{"reviewRating":2},{"reviewRating":2},{"reviewRating":1,"miningResults":[{"attribute":"taste","score":3}]},{"reviewRating":3,"miningResults":[{"attribute":"taste","score":3}]},{"reviewRating":5,"miningResults":[{"attribute":"taste","score":3}]}]`)
		got := reviewsAttributeTTest([]string{"taste"}, r)
		v, ok := got["taste"]
		if !ok || math.IsNaN(v.PValue) || math.IsInf(v.TValue, 0) {
			t.Fatalf("expected valid Welch result: %+v", got)
		}
	})
}
