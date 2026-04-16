import { useState } from 'react'
import { Form, Input, InputNumber, Button, Alert, Descriptions, Card, Typography } from 'antd'
import { issueCard } from '../api.js'

const { Title } = Typography

export default function IssuePage() {
  const [form] = Form.useForm()
  const [result, setResult] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(values) {
    setError('')
    setResult(null)
    setLoading(true)
    try {
      const depositFen = Math.round(values.deposit * 100)
      const preDepositFen = values.preDeposit ? Math.round(values.preDeposit * 100) : 0
      const res = await issueCard({
        name: values.name,
        idNumber: values.idNumber,
        deposit: depositFen,
        preDeposit: preDepositFen,
      })
      setResult(res)
      form.resetFields()
    } catch (err) {
      setError(err.message || '发卡失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card style={{ maxWidth: 520, margin: '0 auto' }}>
      <Title level={4} style={{ marginTop: 0 }}>发卡</Title>

      <Form form={form} layout="vertical" onFinish={handleSubmit}>
        <Form.Item label="持卡人姓名" name="name" rules={[{ required: true, message: '请输入持卡人姓名' }]}>
          <Input placeholder="请输入姓名" />
        </Form.Item>

        <Form.Item label="证件号" name="idNumber" rules={[{ required: true, message: '请输入证件号' }]}>
          <Input placeholder="请输入证件号" />
        </Form.Item>

        <Form.Item label="押金（元）" name="deposit" rules={[{ required: true, message: '请输入押金金额' }]}>
          <InputNumber
            min={0.01}
            step={0.01}
            precision={2}
            placeholder="0.00"
            style={{ width: '100%' }}
          />
        </Form.Item>

        <Form.Item label="预存金额（元）" name="preDeposit">
          <InputNumber
            min={0}
            step={0.01}
            precision={2}
            placeholder="0.00"
            style={{ width: '100%' }}
          />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading}>
            办理发卡
          </Button>
        </Form.Item>
      </Form>

      {error && (
        <Alert type="error" message={error} style={{ marginTop: 8 }} showIcon />
      )}

      {result && (
        <Card
          size="small"
          title="发卡成功"
          style={{ marginTop: 16, background: '#f6ffed', borderColor: '#b7eb8f' }}
        >
          <Descriptions column={1} size="small">
            <Descriptions.Item label="卡号">{result.card.id}</Descriptions.Item>
            <Descriptions.Item label="持卡人">{result.cardHolder.name}</Descriptions.Item>
            <Descriptions.Item label="证件号">{result.cardHolder.idNumber}</Descriptions.Item>
            <Descriptions.Item label="押金">{(result.card.deposit / 100).toFixed(2)} 元</Descriptions.Item>
            <Descriptions.Item label="余额">{(result.card.balance / 100).toFixed(2)} 元</Descriptions.Item>
          </Descriptions>

          {result.refund && (
            <Card
              size="small"
              title={`旧卡自动注销（卡号 ${result.refund.oldCardId}）`}
              style={{ marginTop: 8, background: '#fffbe6', borderColor: '#ffe58f' }}
            >
              <Descriptions column={1} size="small">
                <Descriptions.Item label="退还押金">{(result.refund.deposit / 100).toFixed(2)} 元</Descriptions.Item>
                <Descriptions.Item label="退还余额">{(result.refund.balance / 100).toFixed(2)} 元</Descriptions.Item>
                <Descriptions.Item label="退还合计">{(result.refund.total / 100).toFixed(2)} 元</Descriptions.Item>
              </Descriptions>
            </Card>
          )}
        </Card>
      )}
    </Card>
  )
}
