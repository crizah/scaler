package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserState struct {
	Username          string             `bson:"_id"            json:"userId"`
	CurrentDifficulty int                `bson:"currentDifficulty" json:"currentDifficulty"`
	Streak            int                `bson:"streak"            json:"streak"`
	MaxStreak         int                `bson:"maxStreak"         json:"maxStreak"`
	TotalScore        float64            `bson:"totalScore"        json:"totalScore"`
	TotalAnswered     float64            `bson:"totalAnswered"        json:"totalAnswered"`
	TotalCorrect      float64            `bson:"totalCorrect"        json:"totalCorrect"`
	LastQuestionID    primitive.ObjectID `bson:"lastQuestionId"    json:"lastQuestionId"`
	LastAnswerAt      time.Time          `bson:"lastAnswerAt"      json:"lastAnswerAt"`
	StateVersion      int64              `bson:"stateVersion"      json:"stateVersion"`
	// Adaptive algorithm state
	CorrectWindow []bool  `bson:"correctWindow"     json:"correctWindow"` // rolling 5-answer window
	MomentumScore float64 `bson:"momentumScore"     json:"momentumScore"` // ping-pong stabilizer
}
