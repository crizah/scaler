// contexts/AuthContext.jsx
import { createContext, useContext, useState, useEffect } from "react";
import axios from "axios";

const AuthContext = createContext();
// const BASE_URL = process.env.REACT_APP_BACKEND_URL;

const BASE_URL = window.RUNTIME_CONFIG.BACKEND_URL;
export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Rehydrate from localStorage on mount
  useEffect(() => {
    const username = localStorage.getItem("username");
    const token    = localStorage.getItem("sessionToken");
    if (username && token) setUser({ username });
    setLoading(false);
  }, []);

  const login = async (username, isNew) => {
    const endpoint = isNew ? "/v1/auth/register" : "/v1/auth/session";
    const res = await axios.post(`${BASE_URL}${endpoint}`, { 
      username: username });
    localStorage.setItem("sessionToken", res.data.sessionToken);
    localStorage.setItem("username", res.data.username);
    setUser({ username: res.data.username });
  };

  const logout = () => {
    localStorage.removeItem("sessionToken");
    localStorage.removeItem("username");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}