import { useState, useEffect } from 'react'
import { Form, Input, Button, Alert, Card, Typography, Table } from 'antd'
import { listWindows, createWindow } from '../api.js'

const { Title } = Typography

export default function WindowsPage() {
  const [windows, setWindows] = useState([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [creating, setCreating] = useState(false)
  const [form] = Form.useForm()

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '窗口名称', dataIndex: 'name', key: 'name' },
  ]

  async function fetchWindows() {
    setLoading(true)
    try {
      const res = await listWindows()
      setWindows(res.windows || [])
    } catch (err) {
      setError(err.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchWindows()
  }, [])

  async function handleCreate(values) {
    setError('')
    setCreating(true)
    try {
      await createWindow(values.name)
      form.resetFields()
      await fetchWindows()
    } catch (err) {
      setError(err.message || '创建失败')
    } finally {
      setCreating(false)
    }
  }

  return (
    <Card style={{ maxWidth: 520, margin: '0 auto' }}>
      <Title level={4} style={{ marginTop: 0 }}>窗口管理</Title>

      <Title level={5}>窗口列表</Title>
      <Table
        rowKey="id"
        dataSource={windows}
        columns={columns}
        loading={loading}
        pagination={false}
        size="small"
        locale={{ emptyText: '暂无窗口' }}
        style={{ marginBottom: 24 }}
      />

      <Title level={5}>新建窗口</Title>
      <Form form={form} layout="inline" onFinish={handleCreate}>
        <Form.Item name="name" rules={[{ required: true, message: '请输入窗口名称' }]} style={{ flex: 1 }}>
          <Input placeholder="窗口名称" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={creating}>
            创建
          </Button>
        </Form.Item>
      </Form>

      {error && (
        <Alert type="error" message={error} style={{ marginTop: 12 }} showIcon />
      )}
    </Card>
  )
}
