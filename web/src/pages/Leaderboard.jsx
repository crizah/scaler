
import { useState, useEffect, useRef } from "react";
import axios from "axios";
import styles from "./Leaderboard.module.css";


const BASE_URL = window.RUNTIME_CONFIG.BACKEND_URL;

function getAuthHeader() {
  return { Authorization: `Bearer ${localStorage.getItem("sessionToken")}` };
}

export function Leaderboard({ currentUsername }) {
  const [tab, setTab]       = useState("score");
  const [scores, setScores] = useState(null);
  const [streaks, setStreaks] = useState(null);
  const intervalRef = useRef(null);

  async function fetchBoth() {
    try {
      const [s, st] = await Promise.all([
        axios.get(`${BASE_URL}/v1/leaderboard/score`,  { headers: getAuthHeader() }),
        axios.get(`${BASE_URL}/v1/leaderboard/streak`, { headers: getAuthHeader() }),
      ]);
      setScores(s.data);
      setStreaks(st.data);
    } catch (_) {}
  }

  useEffect(() => {
    fetchBoth();
    intervalRef.current = setInterval(fetchBoth, 5000);
    return () => clearInterval(intervalRef.current);
  }, []);

  const data = tab === "score" ? scores : streaks;

  return (
    <aside className={styles.panel}>
      <div className={styles.header}>
        <span className={styles.title}>LEADERBOARD</span>
        <span className={styles.live}><span className={styles.dot} />LIVE</span>
      </div>

      <div className={styles.tabs}>
        <button
          className={`${styles.tab} ${tab === "score" ? styles.active : ""}`}
          onClick={() => setTab("score")}
        >Score</button>
        <button
          className={`${styles.tab} ${tab === "streak" ? styles.active : ""}`}
          onClick={() => setTab("streak")}
        >Streak</button>
      </div>

      <ul className={styles.list}>
        {!data && <li className={styles.empty}>loading...</li>}

        {data?.entries.map((entry) => (
          <li
            key={entry.username}
            className={`${styles.row} ${entry.username === currentUsername ? styles.you : ""}`}
          >
            <span className={styles.rank}>
              {entry.rank <= 3 ? ["1st yay","2nd yay","3rd yay"][entry.rank - 1] : `#${entry.rank}`}
            </span>
            <span className={styles.name}>{entry.username}</span>
            <span className={styles.val}>
              {tab === "score" ? entry.value.toLocaleString() : `${entry.value}×`}
            </span>
          </li>
        ))}

        
        {data && !data.entries.find(e => e.username === currentUsername) && (
          <>
            <li className={styles.separator}>···</li>
            <li className={`${styles.row} ${styles.you}`}>
              <span className={styles.rank}>#{data.currentUser.rank}</span>
              <span className={styles.name}>{data.currentUser.username}</span>
              <span className={styles.val}>
                {tab === "score"
                  ? data.currentUser.value.toLocaleString()
                  : `${data.currentUser.value}×`}
              </span>
            </li>
          </>
        )}
      </ul>
    </aside>
  );
}