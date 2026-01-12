import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import './App.css'
import { Counter } from './pages/Counter.jsx'
import { Admin } from './pages/Admin.jsx'
import { Home } from './pages/Home.jsx'

function App() {
  return (
    <div>
      <BrowserRouter>
        <Routes>
          <Route index element={<Navigate to="/home" replace />} />
          <Route path='/home' element={<Home />} />
          <Route path='/:festivalCode' element={<Counter />} />
          <Route path='/:festivalCode/admin' element={<Admin />} />
          <Route path='/:festivalCode/*' element={<Navigate to="/home" replace />} />
        </Routes>
      </BrowserRouter>
    </div>
  )
}

export default App
