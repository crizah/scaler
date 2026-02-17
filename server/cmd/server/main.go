package main

import (
	"log" // blank import registers methods
	"server/internal/auth"
	"server/internal/quiz"
	"server/internal/server"

	"github.com/gin-gonic/gin"
)

// main.go
func main() {
	base, err := server.InitialiseServer()
	if err != nil {
		log.Fatal(err)
	}

	authServer := auth.NewAuthServer(base)
	quizServer := quiz.NewQuizServer(base)

	r := gin.Default()

	v1 := r.Group("/v1")
	v1.POST("/auth/register", authServer.RegisterUser)
	v1.POST("/auth/session", authServer.Session)

	protected := v1.Group("/")
	protected.Use(authServer.AuthMiddleware())
	protected.GET("/quiz/next", quizServer.HandleNextQuestion)
	protected.POST("/quiz/answer", quizServer.SubmitAnswer)
	// protected.GET("/quiz/metrics", quizServer.GetMetrics)
	// protected.GET("/leaderboard/score", quizServer.LeaderboardScore)
	// protected.GET("/leaderboard/streak", quizServer.LeaderboardStreak)

	r.Run(":8080")
}
