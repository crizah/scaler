package server

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Server struct {
	MongoClient   *mongo.Client
	JwtSecret     []byte
	CollUsers     *mongo.Collection
	CollUserState *mongo.Collection
	CollQuestions *mongo.Collection
	CollAnswerLog *mongo.Collection
}

func InitialiseServer() (*Server, error) {
	// NEED TO GENERATE AND PUT A JWT TOKEN HERE

	uri := os.Getenv("MONGODB_URI")

	token := []byte(os.Getenv("JWT_SECRET"))

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	u := client.Database("scaler").Collection("Users")
	p := client.Database("scaler").Collection("user-state")
	q := client.Database("scaler").Collection("questions")
	a := client.Database("scaler").Collection("answer-logs")

	return &Server{MongoClient: client, CollUsers: u, CollUserState: p, CollQuestions: q, JwtSecret: token,
		CollAnswerLog: a}, nil
}

func (s *Server) GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JwtSecret)

}
