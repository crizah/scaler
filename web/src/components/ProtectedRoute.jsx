// components/ProtectedRoute.jsx
import { Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export function ProtectedRoute({ children }) {
  const { user, loading } = useAuth();
  if (loading) return null;  // or a spinner if you want
  if (!user) return <Navigate to="/" replace />;
  return children;
}