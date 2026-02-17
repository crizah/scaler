
import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useTheme } from "../hooks/useTheme";

import styles from './Home.module.css'

export function Home() {
  const { login, isAuthenticated } = useAuth()
  const navigate = useNavigate()
  const { theme, toggleTheme } = useTheme();


  const [username, setUsername]   = useState('')
  const [mode, setMode]           = useState('register')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError]         = useState('')

  useEffect(() => {
    if (isAuthenticated) navigate('/quiz', { replace: true })
  }, [isAuthenticated, navigate])

  async function handleSubmit(e) {
    e.preventDefault()
    const name = username.trim()
    if (!name) return

    setIsLoading(true)
    setError('')
    try {
      await login(name, mode === 'register')
      navigate('/quiz')
    } catch (err) {
      setError(err.message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      <div className={styles.card}>
        <button
  type="button"
  onClick={toggleTheme}
  className={styles.themeToggle}
>
  {theme === "dark" ? "Light" : "Dark"}
</button>


        <h1 className={styles.title}>QUIZOS</h1>

        <div className={styles.tabs}>
          <button
            type="button"
            className={`${styles.tab} ${mode === 'register' ? styles.active : ''}`}
            onClick={() => { setMode('register'); setError('') }}
          >
            Register
          </button>
          <button
            type="button"
            className={`${styles.tab} ${mode === 'login' ? styles.active : ''}`}
            onClick={() => { setMode('login'); setError('') }}
          >
            Login
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label htmlFor="username" className={styles.label}>Username</label>
            <input
              id="username"
              className={`${styles.input} ${error ? styles.inputError : ''}`}
              type="text"
              value={username}
              onChange={(e) => { setUsername(e.target.value); setError('') }}
              placeholder="enter username"
              maxLength={10}
              autoFocus
              disabled={isLoading}
              spellCheck={false}
            />
            {error && <span className={styles.error} role="alert">{error}</span>}
          </div>

          <button type="submit" className={styles.btn} disabled={isLoading || !username.trim()}>
            {isLoading ? <span className={styles.spinner} /> : mode === 'register' ? 'Create Account' : 'Login'}
          </button>
        </form>

      </div>
    </div>
  )
}