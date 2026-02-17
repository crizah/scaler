package quiz

import (
	"context"
	"net/http"
	"server/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LeaderboardEntry struct {
	Rank     int     `json:"rank"`
	Username string  `json:"username"`
	Value    float64 `json:"value"`
}

type LeaderboardRes struct {
	Entries     []LeaderboardEntry `json:"entries"`
	CurrentUser LeaderboardEntry   `json:"currentUser"`
}

func (s *Server) GetScoreLeaderboard(c *gin.Context) {
	username := c.GetString("username")

	// top 5
	var states []models.UserState
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().
		SetSort(bson.M{"totalScore": -1}).
		SetLimit(5)

	cursor, err := s.CollUserState.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch leaderboard"})
		return
	}
	cursor.All(ctx, &states)

	entries := make([]LeaderboardEntry, 0, len(states))
	for i, st := range states {
		entries = append(entries, LeaderboardEntry{
			Rank:     i + 1,
			Username: st.Username,
			Value:    st.TotalScore,
		})
	}

	rankScore, _, err := s.getLeaderboardRanks(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get rank " + err.Error()})
		return

	}
	state, err := s.getUserState(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get state " + err.Error()})
		return

	}

	c.JSON(http.StatusOK, LeaderboardRes{
		Entries: entries,
		CurrentUser: LeaderboardEntry{
			Rank:     rankScore,
			Username: username,
			Value:    state.TotalScore,
		},
	})
}

func (s *Server) GetStreakLeaderboard(c *gin.Context) {
	username := c.GetString("username")

	var states []models.UserState
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().
		SetSort(bson.M{"totalScore": -1}).
		SetLimit(5)

	cursor, err := s.CollUserState.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch leaderboard " + err.Error()})
		return
	}
	cursor.All(ctx, &states)

	entries := make([]LeaderboardEntry, 0, len(states))
	for i, st := range states {
		entries = append(entries, LeaderboardEntry{
			Rank:     i + 1,
			Username: st.Username,
			Value:    float64(st.MaxStreak),
		})
	}

	_, rankStreak, err := s.getLeaderboardRanks(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get streak " + err.Error()})
		return

	}

	state, err := s.getUserState(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user state " + err.Error()})
		return

	}

	c.JSON(http.StatusOK, LeaderboardRes{
		Entries: entries,
		CurrentUser: LeaderboardEntry{
			Rank:     rankStreak,
			Username: username,
			Value:    float64(state.MaxStreak),
		},
	})
}
