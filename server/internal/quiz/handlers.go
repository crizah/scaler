package quiz

import (
	"context"
	"math/rand"
	"net/http"
	"server/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type NextQuestionRes struct {
	QuestionID    string   `json:"questionId"`
	Difficulty    int      `json:"difficulty"`
	Prompt        string   `json:"prompt"`
	Choices       []string `json:"choices"`
	StateVersion  int64    `json:"stateVersion"`
	CurrentScore  float64  `json:"currentScore"`
	CurrentStreak int      `json:"currentStreak"`
}

type SubmitAnswerReq struct {
	QuestionID           string `json:"questionId"         binding:"required"`
	Answer               string `json:"answer"             binding:"required"`
	StateVersion         int64  `json:"stateVersion"       binding:"required"`
	AnswerIdempotencyKey string `json:"answerIdempotencyKey" binding:"required"`
}

type SubmitAnswerRes struct {
	Correct               bool    `json:"correct"`
	NewDifficulty         int     `json:"newDifficulty"`
	NewStreak             int     `json:"newStreak"`
	ScoreDelta            float64 `json:"scoreDelta"`
	TotalScore            float64 `json:"totalScore"`
	StateVersion          int64   `json:"stateVersion"`
	LeaderboardRankScore  int     `json:"leaderboardRankScore"`
	LeaderboardRankStreak int     `json:"leaderboardRankStreak"`
}

func (s *Server) HandleNextQuestion(c *gin.Context) {
	username := c.GetString("username")

	// get the users state
	state, err := s.getUserState(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load state"})
		return
	}

	// get all the questions at current difficulty
	diff := state.CurrentDifficulty
	questions, err := s.GetQuestions(diff)
	if err != nil {
		if err.Error() == NO_QUESTIONS {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no questions rn"})
			return

		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant get questions"})
		return
	}

	// pick a random guy, not the lasr asked question tho
	q := pickQuestion(*questions, state.LastQuestionID)

	c.JSON(http.StatusOK, NextQuestionRes{
		QuestionID:    q.Id,
		Difficulty:    q.Difficulty,
		Prompt:        q.Prompt,
		Choices:       q.Choices,
		StateVersion:  state.StateVersion,
		CurrentScore:  state.TotalScore,
		CurrentStreak: state.Streak,
	})

}

func (s *Server) SubmitAnswer(c *gin.Context) {
	username := c.GetString("username")

	var req SubmitAnswerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// reject duplicate submissions
	var existing models.AnswerLog
	err := s.CollAnswerLog.FindOne(ctx, bson.M{"idempotencyKey": req.AnswerIdempotencyKey}).Decode(&existing)
	if err == nil {
		// Already processed — return the stored result idempotently
		rankScore, rankStreak := s.getLeaderboardRanks(username)
		c.JSON(http.StatusOK, SubmitAnswerRes{
			Correct:               existing.Correct,
			NewDifficulty:         0, // caller should use stateVersion to re-fetch
			NewStreak:             existing.StreakAtAnswer,
			ScoreDelta:            existing.ScoreDelta,
			TotalScore:            0,
			StateVersion:          req.StateVersion,
			LeaderboardRankScore:  rankScore,
			LeaderboardRankStreak: rankStreak,
		})
		return
	}

	// versin check and get state
	state, err := s.getUserState(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load state"})
		return
	}
	if state.StateVersion != req.StateVersion {
		c.JSON(http.StatusConflict, gin.H{"error": "stale state"})
		return
	}

	// get the question
	var q models.Question
	if err := s.CollQuestions.FindOne(ctx, bson.M{"_id": req.QuestionID}).Decode(&q); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
		return
	}

	// correctness
	correct := req.Answer == q.CorrectAnswer

	// new difficulty + updated state
	newState := applyAdaptiveAlgorithm(*state, correct)

	//score delta
	scoreDelta := calculateScore(q.Difficulty, correct, newState.Streak)
	newState.TotalScore += scoreDelta
	newState.LastQuestionID = req.QuestionID
	newState.LastAnswerAt = time.Now().UTC()
	newState.StateVersion = state.StateVersion + 1
	newState.TotalAnswered++
	if correct {
		newState.TotalCorrect++
	}
	if newState.TotalAnswered > 0 {
		newState.TotalScore = float64(newState.TotalCorrect) / float64(newState.TotalAnswered)
	}

	// 7. Persist updated state (optimistic concurrency — only update if version matches)
	err = s.updateUserState(username, newState, state.StateVersion)
	if err != nil {
		if err.Error() == "version conflict" {
			c.JSON(http.StatusConflict, gin.H{"error": "concurrent update — fetch /next again"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save state"})
		return
	}

	// answers log
	log := models.AnswerLog{
		Id:             uuid.NewString(),
		Username:       username,
		QuestionID:     req.QuestionID,
		Difficulty:     q.Difficulty,
		Answer:         req.Answer,
		Correct:        correct,
		ScoreDelta:     scoreDelta,
		StreakAtAnswer: newState.Streak,
		IdempotencyKey: req.AnswerIdempotencyKey,
		AnsweredAt:     time.Now().UTC(),
	}

	s.CollAnswerLog.InsertOne(ctx, log) // write to db

	//update leaderboa5rd
	s.updateLeaderboards(username, newState)

	// get rank
	rankScore, rankStreak := s.getLeaderboardRanks(username)

	c.JSON(http.StatusOK, SubmitAnswerRes{
		Correct:               correct,
		NewDifficulty:         newState.CurrentDifficulty,
		NewStreak:             newState.Streak,
		ScoreDelta:            scoreDelta,
		TotalScore:            newState.TotalScore,
		StateVersion:          newState.StateVersion,
		LeaderboardRankScore:  rankScore,
		LeaderboardRankStreak: rankStreak,
	})
}

// func (s *Server) GetMetrics(c *gin.Context) {
// 	username := c.GetString("username")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	state, err := s.getUserState(username)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load state"})
// 		return
// 	}

// 	// Pull last N answers for recentPerformance + difficultyHistogram
// 	opts := options.Find().
// 		SetSort(bson.D{{Key: "answeredAt", Value: -1}}).
// 		SetLimit(recentWindowSize)

// 	cursor, err := s.collAnswerLog.Find(ctx, bson.M{"username": username}, opts)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
// 		return
// 	}

// 	var logs []AnswerLog
// 	cursor.All(ctx, &logs)

// 	// Build recent performance slice + difficulty histogram
// 	recent := make([]RecentAnswer, 0, len(logs))
// 	histogram := make(map[string]int)

// 	for _, l := range logs {
// 		recent = append(recent, RecentAnswer{
// 			QuestionID: l.QuestionID,
// 			Correct:    l.Correct,
// 			Difficulty: l.Difficulty,
// 			ScoreDelta: l.ScoreDelta,
// 			AnsweredAt: l.AnsweredAt.Format(time.RFC3339),
// 		})
// 		key := difficultyLabel(l.Difficulty)
// 		histogram[key]++
// 	}

// 	c.JSON(http.StatusOK, MetricsRes{
// 		CurrentDifficulty:   state.CurrentDifficulty,
// 		Streak:              state.Streak,
// 		MaxStreak:           state.MaxStreak,
// 		TotalScore:          state.TotalScore,
// 		Accuracy:            state.Accuracy,
// 		DifficultyHistogram: histogram,
// 		RecentPerformance:   recent,
// 	})
// }

func pickQuestion(q []models.Question, last string) models.Question {
	// Filter out the last asked question to avoid immediate repeats
	filtered := make([]models.Question, 0, len(q))
	for _, qu := range q {
		if qu.Id != last {
			filtered = append(filtered, qu)
		}
	}
	if len(filtered) == 0 {
		filtered = q
	}
	return filtered[rand.Intn(len(filtered))]

}
