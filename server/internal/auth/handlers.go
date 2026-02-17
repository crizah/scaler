package auth

import (
	"net/http"
	"server/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=1,max=10"`
}

type RegisterRes struct {
	Username     string `json:"username"`
	SessionToken string `json:"sessionToken"`
}

func (s *Server) RegisterUser(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username bw 1-10 chars"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)

	// put into db

	err := s.PutUserIntoDb(req.Username)
	if err != nil {
		if err.Error() == USER_EXISTS {
			c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return

	}

	// start state

	state := models.UserState{
		Username:          req.Username,
		CurrentDifficulty: 3, // 1-10 scale
		Streak:            0,
		MaxStreak:         0,
		TotalScore:        0,
		// StateVersion:      1,
		// CorrectWindow:     []bool{},
		// MomentumScore:     0.5, // neutral starting momentum
	}

	err = s.PutIntoUserStateDB(state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return

	}

	// generate jwt token
	token, err := s.GenerateJWT(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	// this needs to set cookies right?

	c.JSON(http.StatusCreated, RegisterRes{
		Username:     req.Username,
		SessionToken: token,
	})

}

func (s *Server) Session(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}

	// var user User
	// err := h.db.Collection("users").FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	// if err == mongo.ErrNoDocuments {
	//     c.JSON(http.StatusNotFound, gin.H{"error": "user not found — register first"})
	//     return
	// }
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
	//     return
	// }
	if ok, err := s.FindInUsersTable(req.Username); !ok {
		if err.Error() == USER_NOT_FOUND {
			c.JSON(http.StatusNotFound, gin.H{"error": "user doenst exist"})
			return

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "sum error finding user"})
			return
		}

	}

	token, err := s.GenerateJWT(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant generate token"})
		return
	}

	c.JSON(http.StatusOK, RegisterRes{

		Username:     req.Username,
		SessionToken: token,
	})
}
