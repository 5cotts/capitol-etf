import {
  BrowserRouter as Router,
  Routes,
  Route,
} from 'react-router-dom';

import Politicians from './components/Politicians'
import Trades from './components/Trades';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Politicians/>} />
        <Route path="/trades/:pid" element={<Trades />} />
      </Routes>
    </Router> 
  );
}

export default App;
