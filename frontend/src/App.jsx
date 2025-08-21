import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Navbar from './Navbar';
import Home from './pages/Home';
import Game from './pages/Game';

export default function App() {
  return (
    <Router>
      <div className="app" style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
        <Navbar />
        <main style={{ flex: 1, padding: '20px', overflow: 'auto' }}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/game/:gameID" element={<Game />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}
