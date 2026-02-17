// src/pages/Quiz.jsx
import { useState, useEffect, useCallback } from "react";
import { useAuth } from "../context/AuthContext";
import axios from "axios";
import { useTheme } from "../hooks/useTheme";
import styles from "./Quiz.module.css";

import { Leaderboard } from "./Leaderboard";

// const BASE_URL = process.env.REACT_APP_BACKEND_URL;
const BASE_URL = window.RUNTIME_CONFIG.BACKEND_URL;

export function Quiz() {
  const { user } = useAuth();
  const { theme, toggleTheme } = useTheme();


  const [question, setQuestion]     = useState(null);
  const [selected, setSelected]     = useState(null);
  const [result, setResult]         = useState(null);
  const [stats, setStats]           = useState({ score: 0, streak: 0 });
  const [phase, setPhase]           = useState("loading"); // loading | answering | result
  const [error, setError]           = useState(null);

  const fetchQuestion = useCallback(async () => {
    setPhase("loading");
    setSelected(null);
    setResult(null);
    setError(null);
    try {
      const res = await axios.get(`${BASE_URL}/v1/quiz/next`, {
        headers: { Authorization: `Bearer ${localStorage.getItem("sessionToken")}` },
      });

// type NextQuestionRes struct {
// 	QuestionID    string   `json:"questionId"`
// 	Difficulty    int      `json:"difficulty"`
// 	Prompt        string   `json:"prompt"`
// 	Choices       []string `json:"choices"`
// 	StateVersion  int64    `json:"stateVersion"`
// 	CurrentScore  float64  `json:"currentScore"`
// 	CurrentStreak int      `json:"currentStreak"`
// }
      setQuestion(res.data);
      setStats(s => ({ ...s, score: res.data.currentScore, streak: res.data.currentStreak }));
      setPhase("answering");
    } catch (err) {
      setError(err.response?.data?.error ?? "Failed to load question");
      setPhase("error");
    }
  }, []);

  // load first question on mount
  useEffect(() => { fetchQuestion() }, [fetchQuestion]);

  const submitAnswer = async (choice) => {
    if (phase !== "answering" || selected) return;
    setSelected(choice);
    setPhase("submitting");
    try {
      const res = await axios.post(
        `${BASE_URL}/v1/quiz/answer`,
        {
          questionId:           question.questionId,
          answer:               choice,
          stateVersion:         question.stateVersion,
          answerIdempotencyKey: crypto.randomUUID(),
        },
        { headers: { Authorization: `Bearer ${localStorage.getItem("sessionToken")}` } }
      );
      setResult(res.data);
      setStats({ score: res.data.totalScore, streak: res.data.newStreak });
      setPhase("result");
    } catch (err) {
      setError(err.response?.data?.error ?? "Failed to submit answer");
      setPhase("error");
    }
  };



// replace the outer return with this
return (
  <div className={styles.page}>
    <button type="button" onClick={toggleTheme} className={styles.themeToggle}>
      {theme === "dark" ? "Light" : "Dark"}
    </button>

    <div className={styles.layout}>

      {/* left: quiz */}
      <div className={styles.quizCol}>
        <div className={styles.hud}>
          <span className={styles.hudItem}>
            <span className={styles.hudLabel}>SCORE</span>
            <span className={styles.hudValue}>{stats.score.toLocaleString()}</span>
          </span>
          <span className={styles.hudItem}>
            <span className={styles.hudLabel}>STREAK</span>
            <span className={styles.hudValue}>{stats.streak}×</span>
          </span>
          {question && (
            <span className={styles.hudItem}>
              <span className={styles.hudLabel}>LEVEL</span>
              <span className={styles.hudValue}>{question.difficulty}</span>
            </span>
          )}
        </div>

        <div className={styles.main}>
          {phase === "loading" && <p className={styles.status}>loading question...</p>}

          {phase === "error" && (
            <div className={styles.errorBox}>
              <p className={styles.errorMsg}>{error}</p>
              <button className={styles.retryBtn} onClick={fetchQuestion}>retry</button>
            </div>
          )}

          {(phase === "answering" || phase === "submitting" || phase === "result") && question && (
            <>
              <p className={styles.prompt}>{question.prompt}</p>
              <ul className={styles.choices}>
                {question.choices.map((choice, i) => {
                  let state = "idle";
                  if (selected === choice) {
                    state = result ? (result.correct ? "correct" : "wrong") : "selected";
                  }
                  return (
                    <li key={choice}>
                      <button
                        className={`${styles.choice} ${styles[state]}`}
                        onClick={() => submitAnswer(choice)}
                        disabled={phase !== "answering"}
                      >
                        <span className={styles.choiceLetter}>{String.fromCharCode(65 + i)}</span>
                        {choice}
                      </button>
                    </li>
                  );
                })}
              </ul>
            </>
          )}

          {phase === "result" && result && (
            <div className={`${styles.result} ${result.correct ? styles.resultCorrect : styles.resultWrong}`}>
              <span className={styles.resultVerdict}>
                {result.correct ? "✓ Correct" : "✗ Wrong"}
              </span>
              {result.correct && <span className={styles.resultDelta}>+{result.scoreDelta} pts</span>}
              <div className={styles.resultStats}>
                <span>streak {result.newStreak}×</span>
                <span>rank #{result.leaderboardRankScore}</span>
                <span>lvl {result.newDifficulty}</span>
              </div>
              <button className={styles.nextBtn} onClick={fetchQuestion}>Next →</button>
            </div>
          )}
        </div>
      </div>

      {/* right: leaderboard */}
      <Leaderboard currentUsername={user?.username} />

    </div>
  </div>
);}
