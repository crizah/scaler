package main

import (
	"server/internal/auth"

	"github.com/gin-gonic/gin"
)

// main.go (route registration)
func setupRoutes(r *gin.Engine) {
	// initialise server

	s, _ := auth.InitialiseServer()

	v1 := r.Group("/v1")
	{
		// public
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.RegisterUser)
			auth.POST("/session", s.Session)
		}

		// // protected
		// quiz := v1.Group("/quiz", auth.AuthMiddleware(secret))
		// {
		// 	quiz.GET("/next", quizH.NextQuestion)
		// 	quiz.POST("/answer", quizH.SubmitAnswer)
		// 	quiz.GET("/metrics", quizH.GetMetrics)
		// }

		// lb := v1.Group("/leaderboard", auth.AuthMiddleware(secret))
		// {
		// 	lb.GET("/score", lbH.TopScore)
		// 	lb.GET("/streak", lbH.TopStreak)
		// }
	}
}
