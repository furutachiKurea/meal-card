import { Layout, Menu, Typography, Button } from 'antd'
import { useNavigate, useLocation, Outlet } from 'react-router-dom'
import {
  CreditCardOutlined,
  WalletOutlined,
  WarningOutlined,
  StopOutlined,
  BarChartOutlined,
  AppstoreOutlined,
  HomeOutlined,
} from '@ant-design/icons'

const { Sider, Header, Content } = Layout
const { Title } = Typography

const menuItems = [
  { key: '/admin/issue',      icon: <CreditCardOutlined />, label: '发卡' },
  { key: '/admin/deposit',    icon: <WalletOutlined />,     label: '存款' },
  { key: '/admin/loss',       icon: <WarningOutlined />,    label: '挂失管理' },
  { key: '/admin/cancel',     icon: <StopOutlined />,       label: '注销' },
  { key: '/admin/statistics', icon: <BarChartOutlined />,   label: '汇总统计' },
  { key: '/admin/windows',    icon: <AppstoreOutlined />,   label: '窗口管理' },
]

export default function AdminLayout() {
  const navigate = useNavigate()
  const location = useLocation()

  const selectedKey = menuItems.find(item => location.pathname.startsWith(item.key))?.key
    ?? '/admin/issue'

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        theme="dark"
        width={180}
        style={{ position: 'fixed', height: '100vh', left: 0, top: 0, zIndex: 100 }}
      >
        <div style={{ padding: '20px 16px 12px', borderBottom: '1px solid rgba(255,255,255,0.1)' }}>
          <Title level={5} style={{ color: '#fff', margin: 0, lineHeight: 1.4 }}>
            食堂饭卡
            <br />
            <span style={{ fontSize: 11, fontWeight: 400, opacity: 0.7 }}>管理端</span>
          </Title>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{ marginTop: 8 }}
        />
        <div style={{ position: 'absolute', bottom: 16, left: 0, right: 0, padding: '0 12px' }}>
          <Button
            icon={<HomeOutlined />}
            onClick={() => navigate('/')}
            style={{ width: '100%' }}
            size="small"
          >
            返回首页
          </Button>
        </div>
      </Sider>

      <Layout style={{ marginLeft: 180 }}>
        <Header style={{ background: '#fff', padding: '0 24px', borderBottom: '1px solid #f0f0f0', lineHeight: '64px' }}>
          <Title level={4} style={{ margin: 0, color: '#1677ff' }}>
            {menuItems.find(item => location.pathname.startsWith(item.key))?.label ?? '管理端'}
          </Title>
        </Header>
        <Content style={{ padding: '24px', background: '#f5f5f5', minHeight: 'calc(100vh - 64px)' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
