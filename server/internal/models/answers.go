package models

import "time"

type AnswerLog struct {
	Id             string    `bson:"_id"           json:"Id"`
	Username       string    `bson:"username"            json:"username"`
	QuestionID     string    `bson:"questionId"            json:"questionId"`
	Difficulty     int       `bson:"difficulty"            json:"difficulty"`
	Answer         string    `bson:"answer"            json:"answer"`
	Correct        bool      `bson:"correct"            json:"correct"`
	ScoreDelta     float64   `bson:"score"            json:"score"`
	StreakAtAnswer int       `bson:"streak"            json:"streak"`
	IdempotencyKey string    `bson:"ikey"            json:"ikey"`
	AnsweredAt     time.Time `bson:"answeredAt"            json:"answeredAt"`
}
