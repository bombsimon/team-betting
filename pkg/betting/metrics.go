package betting

import (
	"context"

	"github.com/bombsimon/team-betting/pkg"
)

type metricCompareType int

const (
	metricGt metricCompareType = iota
	metricLt
)

type betterBets struct {
	better       *pkg.Better
	scores       []int64
	longestNote  string
	shortestNote string
}

//GetCompetitionMetrics returns metrics for a competition.
func (s *Service) GetCompetitionMetrics(ctx context.Context, id int) (*pkg.CompetitionMetrics, error) {
	competition, err := s.GetCompetition(ctx, id)
	if err != nil {
		return nil, err
	}

	var (
		cm                = &pkg.CompetitionMetrics{}
		totalTopScores    = 0
		totalBottomScores = 0
		totalScore        = 0
		totalBets         = 0
	)

	for _, us := range mapBetsToUserID(competition.Bets) {
		var (
			topScores        = 0
			bottomScores     = 0
			totalBetterScore = 0
		)

		for _, v := range us.scores {
			totalBetterScore += int(v)

			if int(v) == competition.MaxScore {
				topScores++
			}

			if int(v) == competition.MinScore {
				bottomScores++
			}
		}

		// Update all global values to calculate group metrics.
		totalBets += len(us.scores)
		totalTopScores += topScores
		totalBottomScores += bottomScores
		totalScore += totalBetterScore

		// Check better values to update record metrics.
		averageScore := float64(totalBetterScore) / float64(len(us.scores))

		if mv, ok := shouldUpdateMetric(averageScore, metricGt, cm.HighestAverageBetter.Value, us.better); ok {
			cm.HighestAverageBetter = *mv
		}

		if mv, ok := shouldUpdateMetric(averageScore, metricLt, cm.LowestAverageBetter.Value, us.better); ok {
			cm.LowestAverageBetter = *mv
		}

		if mv, ok := shouldUpdateMetric(topScores, metricGt, cm.MostTopScores.Value, us.better); ok {
			cm.MostTopScores = *mv
		}

		if mv, ok := shouldUpdateMetric(bottomScores, metricGt, cm.MostBottomScores.Value, us.better); ok {
			cm.MostBottomScores = *mv
		}

		if mv, ok := shouldUpdateMetric(us.longestNote, metricGt, cm.LongestNote.Value, us.better); ok {
			cm.LongestNote = *mv
		}

		if mv, ok := shouldUpdateMetric(us.shortestNote, metricLt, cm.ShortestNote.Value, us.better); ok {
			cm.ShortestNote = *mv
		}
	}

	cm.NumberOfBottomScores = totalBottomScores
	cm.NumberOfTopScores = totalTopScores
	cm.GroupAverageScore = float64(totalScore) / float64(totalBets)

	return cm, nil
}

func mapBetsToUserID(bets []*pkg.Bet) map[int]*betterBets {
	var btu = map[int]*betterBets{}

	for _, bet := range bets {
		if _, ok := btu[bet.BetterID]; !ok {
			btu[bet.BetterID] = &betterBets{bet.Better, []int64{}, "", ""}
		}

		btu[bet.BetterID].scores = append(btu[bet.BetterID].scores, bet.Score.ValueOrZero())

		if len(bet.Note.ValueOrZero()) > len(btu[bet.BetterID].longestNote) {
			btu[bet.BetterID].longestNote = bet.Note.ValueOrZero()
		}

		if btu[bet.BetterID].shortestNote == "" {
			btu[bet.BetterID].shortestNote = bet.Note.ValueOrZero()
		}

		if len(bet.Note.ValueOrZero()) < len(btu[bet.BetterID].shortestNote) {
			btu[bet.BetterID].shortestNote = bet.Note.ValueOrZero()
		}
	}

	return btu
}

func shouldUpdateMetric(potentialVal interface{}, ct metricCompareType, currentVal interface{}, who *pkg.Better) (*pkg.MetricValue, bool) {
	// Construct the potential metric.
	r := &pkg.MetricValue{
		Who:   who,
		Value: potentialVal,
	}

	// Based on the metric value type check if we should update.
	switch v := currentVal.(type) {
	case nil:
		// If there's no current value we should always set the potential value.
		return r, true
	case float64:
		switch ct {
		case metricGt:
			if potentialVal.(float64) > v {
				return r, true
			}
		case metricLt:
			if potentialVal.(float64) < v {
				return r, true
			}
		}
	case int:
		switch ct {
		case metricGt:
			if potentialVal.(int) > v {
				return r, true
			}
		case metricLt:
			if potentialVal.(int) < v {
				return r, true
			}
		}
	case string:
		// For strings we want to update the value if the length of the string
		// is shorter or longer and the current value, but only if the potential
		// value isn't empty.
		if v == "" {
			return nil, false
		}

		switch ct {
		case metricGt:
			if len(potentialVal.(string)) > len(v) {
				return r, true
			}
		case metricLt:
			if len(potentialVal.(string)) < len(v) {
				return r, true
			}
		}
	}

	return nil, false
}
