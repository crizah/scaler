package models

type Question struct {
	Id            string   `bson:"_id"            json:"questionId"`
	Difficulty    int      `bson:"difficulty"            json:"difficulty"`
	Prompt        string   `bson:"prompt"            json:"prompt"`
	Choices       []string `bson:"choices"            json:"choices"`
	CorrectAnswer string   `bson:"correctans"            json:"correctans"`
}
