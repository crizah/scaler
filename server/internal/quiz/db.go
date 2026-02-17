package quiz

import (
	"context"
	"errors"
	"server/internal/auth"
	"server/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	NO_QUESTIONS     = "no questions at this difficulty"
	VERSION_CONFLICT = "version conflict"
)

func (s *Server) getUserState(username string) (*models.UserState, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result models.UserState

	err := s.CollUserState.FindOne(ctx, bson.M{"_id": username}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New(auth.USER_NOT_FOUND)
	}
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func (s *Server) GetQuestions(diff int) (*[]models.Question, error) {
	var questions []models.Question
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.CollQuestions.Find(ctx, bson.M{
		"difficulty": diff,
	})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &questions); err != nil {
		return nil, err
	}

	return &questions, nil
}

// func (s *Server) updateStreak(username string) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	result, err := s.CollUserState.UpdateOne(
// 		ctx,
// 		bson.M{
// 			"_id":          username,
// 			"stateVersion": expectedVersion,
// 		},
// 		bson.M{"$set": newState},
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	if result.MatchedCount == 0 {
// 		return errors.New("version conflict")
// 	}

// 	return nil

// }

func (s *Server) updateUserState(username string, newState models.UserState, expectedVersion int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.CollUserState.UpdateOne(
		ctx,
		bson.M{
			"_id":          username,
			"stateVersion": expectedVersion,
		},
		bson.M{"$set": newState},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New(VERSION_CONFLICT)
	}
	return nil
}

func (s *Server) updateLeaderboards(username string, state models.UserState) {
	// Upsert score leaderboard
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.CollUserState.UpdateOne(ctx,
		bson.M{"_id": username},
		bson.M{"$set": bson.M{
			"totalScore": state.TotalScore,
			"maxStreak":  state.MaxStreak,
			"updatedAt":  time.Now().UTC(),
		}},

		options.UpdateOne().SetUpsert(true),
	)

}

// getLeaderboardRanks returns (scoreRank, streakRank) for the given username
// Rank = count of users with strictly higher value + 1
func (s *Server) getLeaderboardRanks(username string) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	state, err := s.getUserState(username)
	if err != nil {
		return 0, 0, err
	}

	scoreRank, _ := s.CollUserState.CountDocuments(ctx, bson.M{
		"totalScore": bson.M{"$gt": state.TotalScore},
	})
	streakRank, _ := s.CollUserState.CountDocuments(ctx, bson.M{
		"maxStreak": bson.M{"$gt": state.MaxStreak},
	})

	return int(scoreRank) + 1, int(streakRank) + 1, nil
}
