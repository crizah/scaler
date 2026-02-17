# run

---

please update .env.example before running

```
git clone https://github.com/crizah/scaler.git
cd scaler
cp .env.example .env
docker compose up --build
```
# video**

---

https://github.com/user-attachments/assets/39c4f8c1-e17e-4d45-8571-aa784f09de0e

# stack

---

* the backend is written in golang (gin framework)
* frontend in react.js (has responsive design and light and dark mode)
* database is mongodb
* caching of user state is done with redis

# algorithm

---

* difficulty ranges from 1 to 10
* 2 consecutive correct answers required to increase difficulty (hysteresis)
* 1 wrong answer decreases difficulty
* last 5 answers stored in CorrectWindow
* MomentumScore = (correct answers in window) / 5
* difficulty increases only if:
  - ConsecutiveUp ≥ 2
  - MomentumScore ≥ 0.6 (60%)
* requiring 2 consecutive correct answers prevents ping-pong oscillation
* score is calculated only for correct answers

```
base        = difficulty * 10
multiplier  = min(1 + (streak * 0.1), 5)
scoreDelta  = base * multiplier
```


# data model

---

can be found in ./server/models
```
type Users struct {
    Username string `bson:"_id" 
    CreatedAt time.Time `bson:"createdAt" 
}
```

```
type Question struct {
	Id            string   `bson:"_id"            json:"questionId"`
	Difficulty    int      `bson:"difficulty"            json:"difficulty"`
	Prompt        string   `bson:"prompt"            json:"prompt"`
	Choices       []string `bson:"choices"            json:"choices"`
	CorrectAnswer string   `bson:"correctans"            json:"correctans"`
}

```


```

type UserState struct {
	Username          string    `bson:"_id"            json:"userId"`
	CurrentDifficulty int       `bson:"currentDifficulty" json:"currentDifficulty"`
	Streak            int       `bson:"streak"            json:"streak"`
	MaxStreak         int       `bson:"maxStreak"         json:"maxStreak"`
	TotalScore        float64   `bson:"totalScore"        json:"totalScore"`
	TotalAnswered     float64   `bson:"totalAnswered"        json:"totalAnswered"`
	TotalCorrect      float64   `bson:"totalCorrect"        json:"totalCorrect"`
	LastQuestionID    string    `bson:"lastQuestionId"    json:"lastQuestionId"`
	LastAnswerAt      time.Time `bson:"lastAnswerAt"      json:"lastAnswerAt"`
	StateVersion      int       `bson:"stateVersion"      json:"stateVersion"`
	// algorithm state
	CorrectWindow   []bool  `bson:"correctWindow"     json:"correctWindow"` 
	MomentumScore   float64 `bson:"momentumScore"     json:"momentumScore"` 
	ConsecutiveUp   int     `bson:"consecutiveUp"     json:"consecutiveUp"`
	ConsecutiveDown int     `bson:"consecutiveDown"     json:"consecutiveDown"`
}

has index at totalScore and maxStreak for faster search


```

```
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

has an index at ikey for faster search
```

# api structure

---


can be found in ./server/cmd/main.go

Protected routes require a session token via Authorization header.


```
POST /auth/register
Request: username
Response: username, sessionToken
```
```
GET /v1/quiz/next 
Request: userId, sessionId (optional)
Response: questionId, difficulty, prompt, choices, sessionId, stateVersion, currentScore, currentStreak


POST /v1/quiz/answer 
Request: userId, sessionId, questionId, answer, stateVersion, answerIdempotencyKey
Response: correct, newDifficulty, newStreak, scoreDelta, totalScore, stateVersion, leaderboardRankScore, leaderboardRankStreak


GET /v1/quiz/metrics 
Response: currentDifficulty, streak, maxStreak, totalScore, accuracy, difficultyHistogram, recentPerformance
```
```
GET /v1/leaderboard/score 
Response: top 5 users by total score (rank, username, value, currentUser)


GET /v1/leaderboard/streak 
Response: top 5 users by max streak (rank, username, value, currentUser)
```


# real time

---

* updates to leaderboard, score and  streaks is done in real time

# edge case handeling

---

* streak gets reset on every wrong answer
* state version checked, stale states are discarded
* duplicate submissions dont update streak because of a check with the answer log (idempotency)


# docker

---

* frontend and backend both have Dockerfiles
* docker compose to get both running
* envionment variables injected into React using an entrypoint scrip found at ./web/entrypoint.sh










