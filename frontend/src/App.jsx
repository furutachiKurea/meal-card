import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom'
import { Card, Row, Col, Typography, Space } from 'antd'
import { SettingOutlined, DesktopOutlined, UserOutlined } from '@ant-design/icons'
import AdminLayout from './layouts/AdminLayout.jsx'
import IssuePage from './pages/IssuePage.jsx'
import DepositPage from './pages/DepositPage.jsx'
import MealPage from './pages/MealPage.jsx'
import CustomerScreen from './pages/CustomerScreen.jsx'
import LossPage from './pages/LossPage.jsx'
import CancelPage from './pages/CancelPage.jsx'
import StatisticsPage from './pages/StatisticsPage.jsx'
import WindowsPage from './pages/WindowsPage.jsx'
import NotFoundPage from './pages/NotFoundPage.jsx'

const { Title, Text } = Typography

function HomePage() {
  const navigate = useNavigate()

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #e6f4ff 0%, #f0f5ff 100%)',
        padding: 24,
      }}
    >
      <Space direction="vertical" align="center" style={{ marginBottom: 48 }}>
        <Title style={{ margin: 0, fontSize: 32, color: '#1677ff' }}>食堂饭卡管理系统</Title>
        <Text type="secondary">请选择使用模式</Text>
      </Space>

      <Row gutter={[24, 24]} justify="center">
        <Col>
          <Card
            hoverable
            style={{ width: 200, textAlign: 'center', borderRadius: 16, boxShadow: '0 4px 20px rgba(22,119,255,0.12)' }}
            onClick={() => navigate('/admin/issue')}
          >
            <Space direction="vertical" size={12} style={{ padding: '16px 0' }}>
              <SettingOutlined style={{ fontSize: 48, color: '#1677ff' }} />
              <Title level={4} style={{ margin: 0 }}>管理端</Title>
              <Text type="secondary" style={{ fontSize: 12 }}>发卡/存款/挂失/统计</Text>
            </Space>
          </Card>
        </Col>

        <Col>
          <Card
            hoverable
            style={{ width: 200, textAlign: 'center', borderRadius: 16, boxShadow: '0 4px 20px rgba(82,196,26,0.12)' }}
            onClick={() => navigate('/window')}
          >
            <Space direction="vertical" size={12} style={{ padding: '16px 0' }}>
              <DesktopOutlined style={{ fontSize: 48, color: '#52c41a' }} />
              <Title level={4} style={{ margin: 0 }}>窗口操作端</Title>
              <Text type="secondary" style={{ fontSize: 12 }}>工作人员刷卡结算</Text>
            </Space>
          </Card>
        </Col>

        <Col>
          <Card
            hoverable
            style={{ width: 200, textAlign: 'center', borderRadius: 16, boxShadow: '0 4px 20px rgba(250,140,22,0.12)' }}
            onClick={() => navigate('/window/customer')}
          >
            <Space direction="vertical" size={12} style={{ padding: '16px 0' }}>
              <UserOutlined style={{ fontSize: 48, color: '#fa8c16' }} />
              <Title level={4} style={{ margin: 0 }}>顾客屏</Title>
              <Text type="secondary" style={{ fontSize: 12 }}>学生查看余额/结果</Text>
            </Space>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />

        <Route path="/admin" element={<AdminLayout />}>
          <Route index element={<IssuePage />} />
          <Route path="issue"      element={<IssuePage />} />
          <Route path="deposit"    element={<DepositPage />} />
          <Route path="loss"       element={<LossPage />} />
          <Route path="cancel"     element={<CancelPage />} />
          <Route path="statistics" element={<StatisticsPage />} />
          <Route path="windows"    element={<WindowsPage />} />
        </Route>

        <Route path="/window" element={<MealPage />} />
        <Route path="/window/customer" element={<CustomerScreen />} />

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  )
}
