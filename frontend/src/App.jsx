import { BrowserRouter, Routes, Route, NavLink } from 'react-router-dom'
import IssuePage from './pages/IssuePage.jsx'
import DepositPage from './pages/DepositPage.jsx'
import MealPage from './pages/MealPage.jsx'
import LossPage from './pages/LossPage.jsx'
import CancelPage from './pages/CancelPage.jsx'
import StatisticsPage from './pages/StatisticsPage.jsx'
import WindowsPage from './pages/WindowsPage.jsx'

const navItems = [
  { to: '/issue', label: '发卡' },
  { to: '/deposit', label: '存款' },
  { to: '/meal', label: '就餐消费' },
  { to: '/loss', label: '挂失管理' },
  { to: '/cancel', label: '注销' },
  { to: '/statistics', label: '汇总统计' },
  { to: '/windows', label: '窗口管理' },
]

const navStyle = {
  display: 'flex',
  gap: 0,
  background: '#1565c0',
  padding: '0 16px',
  flexWrap: 'wrap',
}

const linkStyle = ({ isActive }) => ({
  color: isActive ? '#fff' : '#bbdefb',
  padding: '12px 16px',
  textDecoration: 'none',
  fontWeight: isActive ? 'bold' : 'normal',
  borderBottom: isActive ? '3px solid #fff' : '3px solid transparent',
  display: 'inline-block',
})

export default function App() {
  return (
    <BrowserRouter>
      <header>
        <div style={{ background: '#0d47a1', color: '#fff', padding: '12px 24px' }}>
          <h1 style={{ margin: 0, fontSize: 20 }}>食堂饭卡管理系统</h1>
        </div>
        <nav style={navStyle}>
          {navItems.map(item => (
            <NavLink key={item.to} to={item.to} style={linkStyle}>
              {item.label}
            </NavLink>
          ))}
        </nav>
      </header>
      <main style={{ padding: '16px 0' }}>
        <Routes>
          <Route path="/" element={<IssuePage />} />
          <Route path="/issue" element={<IssuePage />} />
          <Route path="/deposit" element={<DepositPage />} />
          <Route path="/meal" element={<MealPage />} />
          <Route path="/loss" element={<LossPage />} />
          <Route path="/cancel" element={<CancelPage />} />
          <Route path="/statistics" element={<StatisticsPage />} />
          <Route path="/windows" element={<WindowsPage />} />
        </Routes>
      </main>
    </BrowserRouter>
  )
}
