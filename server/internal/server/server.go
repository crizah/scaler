package server

import (
	"context"
	"os"
	"server/internal/models"
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
}

func InitialiseServer() (*Server, error) {

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

	return &Server{MongoClient: client, CollUsers: u, CollUserState: p, CollQuestions: q, JwtSecret: token}, nil
}

func (s *Server) PopulateQuestions() {
	questions := []interface{}{
		models.Question{
			Id:            "1",
			Difficulty:    1,
			Prompt:        "meow meow?",
			Choices:       []string{"yes", "no", "lol idk", "haha"},
			CorrectAnswer: "yes",
		},
		models.Question{
			Id:            "11",
			Difficulty:    1,
			Prompt:        "what happens when i turn my headlights on?",
			Choices:       []string{"suddenly i can see", "ive got tunnel vision", "im an idiot", "i learn of right and wrong"},
			CorrectAnswer: "suddenly i can see",
		},
		models.Question{
			Id:            "2",
			Difficulty:    4,
			Prompt:        "do u have a cloak?",
			Choices:       []string{"yeah, its a bit of a joke", "no i dont wear clothes", "why are u asking me this", "im going to shoot myself"},
			CorrectAnswer: "yeah, its a bit of a joke",
		},
		models.Question{
			Id:            "3",
			Difficulty:    2,
			Prompt:        "do u want to go home?",
			Choices:       []string{"i want to leave the show", "yeah and take off this uniform", "the worms have entered my brain", "no lol i enjoy being in hell haha"},
			CorrectAnswer: "no lol i enjoy being in hell haha",
		},
		models.Question{
			Id:            "4",
			Difficulty:    2,
			Prompt:        "would somebody care if u stayed with me?",
			Choices:       []string{"baby u can stay and nobody would care", "just pretend im not there", "arnav bought me a gun i will use it one day, yall watch", "u can change"},
			CorrectAnswer: "baby u can stay and nobody would care",
		},
		models.Question{
			Id:            "5",
			Difficulty:    3,
			Prompt:        "what will ur organs do?",
			Choices:       []string{"soon ur organs will grow little mouths", "they will speak for themselves", "soon they will refuse to hold u up, so embarrased to bear your name", "idek what ur talking about"},
			CorrectAnswer: "soon they will refuse to hold u up, so embarrased to bear your name",
		},
		models.Question{
			Id:            "6",
			Difficulty:    4,
			Prompt:        "what kind on sauce do u add?",
			Choices:       []string{"its ltr just sauce", "awesome sauce", "idk i think sauce is a very broiad term", "add where? whatr are u even talking about? what are these questions i want a refund"},
			CorrectAnswer: "awesome sauce",
		},
		models.Question{
			Id:            "7",
			Difficulty:    1,
			Prompt:        "what does she move with",
			Choices:       []string{"im tired grandpa", "lmao wasting time rn but its okay", "a purpouse", "im dumb but happy"},
			CorrectAnswer: "a purpouse",
		},
		models.Question{
			Id:            "8",
			Difficulty:    2,
			Prompt:        "what do u wish for?",
			Choices:       []string{"late at night when im driving", "i wish that they would swoop down in a country lane", "take me on board their beautiful shit", "show me the world as id love to see it"},
			CorrectAnswer: "i wish that they would swoop down in a country lane",
		},
		models.Question{
			Id:            "9",
			Difficulty:    3,
			Prompt:        "what has the bike got?",
			Choices:       []string{"its got a basket a bell and things that make it look good", "idk haha what even is a bike", "i want a refund", "is this a pink floyd reference??"},
			CorrectAnswer: "is this a pink floyd reference??",
		},
		models.Question{
			Id:            "10",
			Difficulty:    5,
			Prompt:        "what do u feel like?",
			Choices:       []string{"i feel like squished face, slick pig living in a smokey city", "wait are all of these songs??? what is wrong with u", "guys, stop", "the worms have enetred my brain"},
			CorrectAnswer: "i feel like squished face, slick pig living in a smokey city",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.CollQuestions.InsertMany(ctx, questions)

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
