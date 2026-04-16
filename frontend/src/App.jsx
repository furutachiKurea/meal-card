import { BrowserRouter, Routes, Route, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Typography } from 'antd'
import IssuePage from './pages/IssuePage.jsx'
import DepositPage from './pages/DepositPage.jsx'
import MealPage from './pages/MealPage.jsx'
import LossPage from './pages/LossPage.jsx'
import CancelPage from './pages/CancelPage.jsx'
import StatisticsPage from './pages/StatisticsPage.jsx'
import WindowsPage from './pages/WindowsPage.jsx'

const { Header, Content } = Layout
const { Title } = Typography

const navItems = [
  { key: '/issue', label: '发卡' },
  { key: '/deposit', label: '存款' },
  { key: '/meal', label: '就餐消费' },
  { key: '/loss', label: '挂失管理' },
  { key: '/cancel', label: '注销' },
  { key: '/statistics', label: '汇总统计' },
  { key: '/windows', label: '窗口管理' },
]

function AppLayout() {
  const navigate = useNavigate()
  const location = useLocation()

  const selectedKey = location.pathname === '/' ? '/issue' : location.pathname

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', padding: '0 24px', gap: 24 }}>
        <Title level={4} style={{ color: '#fff', margin: 0, whiteSpace: 'nowrap' }}>
          食堂饭卡管理系统
        </Title>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[selectedKey]}
          items={navItems}
          onClick={({ key }) => navigate(key)}
          style={{ flex: 1, minWidth: 0 }}
        />
      </Header>
      <Content style={{ padding: '24px', background: '#f5f5f5' }}>
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
      </Content>
    </Layout>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AppLayout />
    </BrowserRouter>
  )
}
