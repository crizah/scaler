import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from "./context/AuthContext";
// import { Msg } from './pages/GetMsg';
import {Home} from './pages/Home';
import { Quiz } from './pages/Quiz';
// import { Notif } from './pages/GetNotif';
// import { SignUp } from './pages/Signup';  
// import { Login } from './pages/LogIn';
import { ProtectedRoute } from './components/ProtectedRoute';
// import { Verification } from './pages/Verification';





import './App.css';


function App() {
 
  return (
<AuthProvider>
 
    <Router>
      <Routes>
        

        <Route
          path={"/"}
          element={<Home />}
        />

        <Route
            path="/quiz"
            element={
              <ProtectedRoute>
                <Quiz />
              </ProtectedRoute>
            }
          />

     
      
      </Routes>
    </Router>

</AuthProvider>
  );
}

export default App;



