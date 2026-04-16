import { useState } from 'react'
import { Form, Input, InputNumber, Button, Alert, Descriptions, Card, Typography, Space } from 'antd'
import { getCard, deposit } from '../api.js'

const { Title } = Typography

export default function DepositPage() {
  const [cardId, setCardId] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [receipt, setReceipt] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [depositForm] = Form.useForm()

  async function handleQueryCard() {
    if (!cardId.trim()) return
    setError('')
    setCardInfo(null)
    setReceipt(null)
    setLoading(true)
    try {
      const res = await getCard(cardId.trim())
      if (res.status !== 'active') {
        const statusText = res.status === 'lost' ? '已挂失' : '已注销'
        setError(`该卡${statusText}，无法充值`)
        return
      }
      setCardInfo(res)
    } catch (err) {
      setError(err.message || '查询失败')
    } finally {
      setLoading(false)
    }
  }

  async function handleDeposit(values) {
    setError('')
    setLoading(true)
    try {
      const amountFen = Math.round(values.amount * 100)
      const res = await deposit(cardId.trim(), amountFen)
      setReceipt(res)
      setCardInfo(prev => prev ? { ...prev, balance: res.newBalance } : null)
      depositForm.resetFields()
    } catch (err) {
      setError(err.message || '存款失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card style={{ maxWidth: 520, margin: '0 auto' }}>
      <Title level={4} style={{ marginTop: 0 }}>存款</Title>

      <Space.Compact style={{ width: '100%', marginBottom: 16 }}>
        <Input
          placeholder="请输入卡号"
          value={cardId}
          onChange={e => { setCardId(e.target.value); setCardInfo(null); setReceipt(null); setError('') }}
          onPressEnter={handleQueryCard}
        />
        <Button type="primary" loading={loading} onClick={handleQueryCard}>
          查询
        </Button>
      </Space.Compact>

      {cardInfo && (
        <Descriptions
          column={1}
          size="small"
          bordered
          style={{ marginBottom: 16 }}
        >
          <Descriptions.Item label="持卡人">{cardInfo.cardHolder.name}</Descriptions.Item>
          <Descriptions.Item label="卡号">{cardInfo.id}</Descriptions.Item>
          <Descriptions.Item label="当前余额">{(cardInfo.balance / 100).toFixed(2)} 元</Descriptions.Item>
        </Descriptions>
      )}

      {cardInfo && (
        <Form form={depositForm} layout="vertical" onFinish={handleDeposit}>
          <Form.Item label="存款金额（元）" name="amount" rules={[{ required: true, message: '请输入存款金额' }]}>
            <InputNumber
              min={0.01}
              step={0.01}
              precision={2}
              placeholder="0.00"
              style={{ width: '100%' }}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              确认存款
            </Button>
          </Form.Item>
        </Form>
      )}

      {error && (
        <Alert type="error" message={error} style={{ marginTop: 8 }} showIcon />
      )}

      {receipt && (
        <Card
          size="small"
          title="存款收据"
          style={{ marginTop: 16, background: '#f6ffed', borderColor: '#b7eb8f' }}
        >
          <Descriptions column={1} size="small">
            <Descriptions.Item label="卡号">{receipt.cardId}</Descriptions.Item>
            <Descriptions.Item label="持卡人">{receipt.holderName}</Descriptions.Item>
            <Descriptions.Item label="充值金额">{(receipt.amount / 100).toFixed(2)} 元</Descriptions.Item>
            <Descriptions.Item label="充值后余额">{(receipt.newBalance / 100).toFixed(2)} 元</Descriptions.Item>
            <Descriptions.Item label="充值时间">{new Date(receipt.createdAt).toLocaleString('zh-CN')}</Descriptions.Item>
          </Descriptions>
        </Card>
      )}
    </Card>
  )
}
