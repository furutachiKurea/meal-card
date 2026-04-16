import { useState, useEffect } from 'react'
import { Form, Input, InputNumber, Button, Alert, Descriptions, Card, Typography, Select, Space, Steps } from 'antd'
import { getCard, createTransaction, listWindows } from '../api.js'

const { Title } = Typography

export default function MealPage() {
  const [cardId, setCardId] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [alarm, setAlarm] = useState('')
  const [txResult, setTxResult] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [windows, setWindows] = useState([])
  const [selectedWindowId, setSelectedWindowId] = useState(null)
  const [settleForm] = Form.useForm()

  const currentStep = txResult ? 2 : cardInfo ? 1 : 0

  useEffect(() => {
    listWindows().then(res => {
      const wins = res.windows || []
      setWindows(wins)
      if (wins.length > 0) {
        setSelectedWindowId(wins[0].id)
      }
    }).catch(() => {})
  }, [])

  async function handleQueryCard() {
    if (!cardId.trim()) return
    setError('')
    setCardInfo(null)
    setAlarm('')
    setTxResult(null)
    setLoading(true)
    try {
      const res = await getCard(cardId.trim())
      if (res.status === 'cancelled') {
        setAlarm('警报：此卡已注销，禁止就餐！')
        return
      }
      if (res.status === 'lost') {
        setAlarm('警报：此卡已挂失，禁止就餐！')
        return
      }
      setCardInfo(res)
    } catch (err) {
      if (err.status === 404) {
        setAlarm('警报：此卡非本单位所发，禁止就餐！')
      } else {
        setError(err.message || '查询失败')
      }
    } finally {
      setLoading(false)
    }
  }

  async function handleSettle(values) {
    setError('')
    setLoading(true)
    try {
      const amountFen = Math.round(values.amount * 100)
      const res = await createTransaction(cardId.trim(), {
        windowId: selectedWindowId,
        amount: amountFen,
      })
      setTxResult(res)
      setCardInfo(prev => prev ? { ...prev, balance: res.newBalance } : null)
      settleForm.resetFields()
    } catch (err) {
      setError(err.message || '结算失败')
    } finally {
      setLoading(false)
    }
  }

  function handleReset() {
    setCardId('')
    setCardInfo(null)
    setAlarm('')
    setTxResult(null)
    setError('')
    settleForm.resetFields()
  }

  return (
    <Card style={{ maxWidth: 560, margin: '0 auto' }}>
      <Title level={4} style={{ marginTop: 0 }}>就餐消费（窗口机）</Title>

      {windows.length > 0 && (
        <Form.Item label="当前窗口" style={{ marginBottom: 16 }}>
          <Select
            value={selectedWindowId}
            onChange={setSelectedWindowId}
            options={windows.map(w => ({ value: w.id, label: w.name }))}
            style={{ width: '100%' }}
          />
        </Form.Item>
      )}

      <Steps
        current={currentStep}
        items={[
          { title: '刷卡验证' },
          { title: '输入金额' },
          { title: '结算完成' },
        ]}
        style={{ marginBottom: 24 }}
      />

      {currentStep === 0 && (
        <Space.Compact style={{ width: '100%' }}>
          <Input
            placeholder="请输入卡号（模拟刷卡）"
            value={cardId}
            onChange={e => { setCardId(e.target.value); setAlarm(''); setError('') }}
            onPressEnter={handleQueryCard}
          />
          <Button type="primary" loading={loading} onClick={handleQueryCard}>
            刷卡
          </Button>
        </Space.Compact>
      )}

      {alarm && (
        <Alert
          type="error"
          message={alarm}
          style={{ marginTop: 16, fontSize: 16, fontWeight: 'bold' }}
          showIcon
        />
      )}

      {cardInfo && currentStep >= 1 && (
        <Descriptions column={1} size="small" bordered style={{ marginBottom: 16 }}>
          <Descriptions.Item label="持卡人">{cardInfo.cardHolder.name}</Descriptions.Item>
          <Descriptions.Item label="卡号">{cardInfo.id}</Descriptions.Item>
          <Descriptions.Item label="余额">
            <span style={{ fontSize: 18, fontWeight: 'bold', color: '#1677ff' }}>
              {(cardInfo.balance / 100).toFixed(2)} 元
            </span>
          </Descriptions.Item>
        </Descriptions>
      )}

      {cardInfo && !txResult && (
        <Form form={settleForm} layout="vertical" onFinish={handleSettle}>
          <Form.Item label="本次消费金额（元）" name="amount" rules={[{ required: true, message: '请输入消费金额' }]}>
            <InputNumber
              min={0.01}
              step={0.01}
              precision={2}
              placeholder="0.00"
              style={{ width: '100%', fontSize: 18 }}
              size="large"
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} size="large">
              确认结算
            </Button>
          </Form.Item>
        </Form>
      )}

      {error && (
        <Alert type="error" message={error} style={{ marginTop: 8 }} showIcon />
      )}

      {txResult && (
        <Card
          size="small"
          title="结算成功"
          style={{ marginTop: 16, background: '#f6ffed', borderColor: '#b7eb8f' }}
          extra={
            <Button onClick={handleReset}>下一位</Button>
          }
        >
          <Descriptions column={1} size="small">
            <Descriptions.Item label="消费金额">{(txResult.amount / 100).toFixed(2)} 元</Descriptions.Item>
            <Descriptions.Item label="扣款后余额">
              <span style={{ fontSize: 16, fontWeight: 'bold', color: '#52c41a' }}>
                {(txResult.newBalance / 100).toFixed(2)} 元
              </span>
            </Descriptions.Item>
          </Descriptions>
        </Card>
      )}
    </Card>
  )
}
