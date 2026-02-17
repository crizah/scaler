package main

import (
	"log" // blank import registers methods
	"server/internal/auth"
	"server/internal/quiz"
	"server/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// main.go
func main() {
	godotenv.Load()
	base, err := server.InitialiseServer()
	if err != nil {
		log.Fatal(err)
	}
	// base.PopulateQuestions()

	authServer := auth.NewAuthServer(base)
	quizServer := quiz.NewQuizServer(base)

	r := gin.Default()

	r.Use(auth.CORSMiddleware())

	v1 := r.Group("/v1")
	v1.POST("/auth/register", authServer.RegisterUser) // works
	v1.POST("/auth/session", authServer.Session)       // works

	protected := v1.Group("/")
	protected.Use(authServer.AuthMiddleware())
	protected.GET("/quiz/next", quizServer.HandleNextQuestion) // working
	protected.POST("/quiz/answer", quizServer.SubmitAnswer)    // does not work yet
	// protected.GET("/quiz/metrics", quizServer.GetMetrics)
	// protected.GET("/leaderboard/score", quizServer.LeaderboardScore)
	// protected.GET("/leaderboard/streak", quizServer.LeaderboardStreak)

	r.Run(":8081")
}
