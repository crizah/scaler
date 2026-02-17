package auth

import (
	"context"
	"errors"
	"server/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	USER_EXISTS    = "user already exists"
	USER_NOT_FOUND = "user not found"
)

func (s *Server) PutUserIntoDb(username string) error {
	user := bson.M{
		"_id": username, // PK

		"createdAt": time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := s.CollUsers.InsertOne(ctx, user)
	if err != nil {
		// Duplicate key error (username )
		if mongo.IsDuplicateKeyError(err) {
			return errors.New(USER_EXISTS)
		}
		return err
	}

	return nil

}

func (s *Server) PutIntoUserStateDB(state models.UserState) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := s.CollUserState.InsertOne(ctx, state)
	return err

}

func (s *Server) FindInUsersTable(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.CollUsers.FindOne(ctx, bson.M{"_id": username}).Err()

	if err == mongo.ErrNoDocuments {
		return false, errors.New(USER_NOT_FOUND)
	}
	if err != nil {
		return false, err
	}

	return true, nil

}
