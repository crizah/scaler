// internal/quiz/adaptive.go
package quiz

import "server/internal/models"

const (
	minDifficulty       = 1
	maxDifficulty       = 10
	correctStreakToUp   = 2   // need 2 consecutive correct to go up (hysteresis)
	wrongStreakToDown   = 1   // 1 wrong is enough to go down
	rollingWindowSize   = 5   // last 5 answers for momentum
	momentumThreshold   = 0.6 // 60% correct in window required to increase difficulty
	maxStreakMultiplier = 5
)

// applyAdaptiveAlgorithm returns a mutated copy of state — never modifies in place
func applyAdaptiveAlgorithm(state models.UserState, correct bool) models.UserState {
	s := state // copy

	// update streak
	if correct {
		s.Streak++
		s.ConsecutiveUp++
		s.ConsecutiveDown = 0
		if s.Streak > s.MaxStreak {
			s.MaxStreak = s.Streak
		}
	} else {
		s.Streak = 0
		s.ConsecutiveDown++
		s.ConsecutiveUp = 0
	}

	// rolling window
	s.CorrectWindow = append(s.CorrectWindow, correct)
	if len(s.CorrectWindow) > rollingWindowSize {
		s.CorrectWindow = s.CorrectWindow[len(s.CorrectWindow)-rollingWindowSize:]
	}

	// ── 3. Compute momentum (% correct in window) ─────────────────────────────
	correctCount := 0
	for _, c := range s.CorrectWindow {
		if c {
			correctCount++
		}
	}
	s.MomentumScore = float64(correctCount) / float64(len(s.CorrectWindow))

	// ── 4. Difficulty adjustment with ping-pong stabilizer ────────────────────
	//
	// To go UP:   need consecutiveUp >= correctStreakToUp AND momentum >= threshold
	// To go DOWN: need consecutiveDown >= wrongStreakToDown (fast response to failure)
	//
	// This prevents A/B/A/B oscillation because:
	//   - a single correct after a wrong resets consecutiveUp to 1, not enough to go up
	//   - momentum requires sustained good performance across the window
	//
	if correct {
		if s.ConsecutiveUp >= correctStreakToUp && s.MomentumScore >= momentumThreshold {
			if s.CurrentDifficulty < maxDifficulty {
				s.CurrentDifficulty++
			}
			s.ConsecutiveUp = 0 // reset after adjustment
		}
	} else {
		if s.ConsecutiveDown >= wrongStreakToDown {
			if s.CurrentDifficulty > minDifficulty {
				s.CurrentDifficulty--
			}
			s.ConsecutiveDown = 0
		}
	}

	return s
}

func calculateScore(difficulty int, correct bool, streak int) float64 {
	// calculateScore returns the score delta for a single answer
	//
	// Formula:
	//
	//	base      = difficulty * 10
	//	multiplier = min(1 + (streak * 0.1), maxStreakMultiplier)  → caps at 5x
	//	delta     = base * multiplier  (0 if wrong)
	if !correct {
		return 0
	}
	base := float64(difficulty) * 10.0
	multiplier := 1.0 + float64(streak)*0.1
	if multiplier > maxStreakMultiplier {
		multiplier = maxStreakMultiplier
	}
	return base * multiplier
}
