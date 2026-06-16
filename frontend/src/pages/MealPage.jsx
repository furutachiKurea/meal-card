import { useState, useEffect, useRef } from 'react'
import { Form, InputNumber, Button, Alert, Card, Typography, Select, Steps } from 'antd'
import { Input } from 'antd'
import { getCard, createTransaction, listWindows } from '../api.js'
import { HomeOutlined, CheckCircleOutlined, WarningFilled } from '@ant-design/icons'
import { useNavigate, useSearchParams } from 'react-router-dom'

const { Title, Text } = Typography

export default function MealPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const lockedWindowId = searchParams.get('id') ? Number(searchParams.get('id')) : null
  const [cardNo, setCardNo] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [alarm, setAlarm] = useState('')
  const [txResult, setTxResult] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [windows, setWindows] = useState([])
  const [selectedWindowId, setSelectedWindowId] = useState(lockedWindowId)
  const [settleForm] = Form.useForm()
  const channelRef = useRef(null)

  const currentStep = txResult ? 2 : cardInfo ? 1 : 0

  // 初始化 BroadcastChannel 向顾客屏推送状态
  useEffect(() => {
    const channelName = lockedWindowId ? `window-${lockedWindowId}` : 'window-default'
    channelRef.current = new BroadcastChannel(channelName)
    return () => channelRef.current?.close()
  }, [lockedWindowId])

  function broadcast(msg) {
    channelRef.current?.postMessage(msg)
  }

  useEffect(() => {
    listWindows().then(res => {
      const wins = res.windows || []
      setWindows(wins)
      if (!lockedWindowId && wins.length > 0) {
        setSelectedWindowId(wins[0].id)
      }
    }).catch(() => {})
  }, [])

  async function handleQueryCard() {
    if (!cardNo.trim()) return
    setError('')
    setCardInfo(null)
    setAlarm('')
    setTxResult(null)
    setLoading(true)
    try {
      const res = await getCard(cardNo.trim())
      if (res.status === 'cancelled') {
        setAlarm('警报：此卡已注销，禁止就餐！')
        broadcast({ type: 'alarm', message: '此卡已注销' })
        return
      }
      if (res.status === 'lost') {
        setAlarm('警报：此卡已挂失，禁止就餐！')
        broadcast({ type: 'alarm', message: '此卡已挂失' })
        return
      }
      setCardInfo(res)
      broadcast({ type: 'card_read', holderName: res.cardHolder.name, balance: res.balance })
    } catch (err) {
      if (err.status === 404) {
        setAlarm('警报：此卡非本单位所发，禁止就餐！')
        broadcast({ type: 'alarm', message: '非本单位卡' })
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
      const res = await createTransaction(cardNo.trim(), {
        windowId: selectedWindowId,
        amount: amountFen,
      })
      setTxResult(res)
      setCardInfo(prev => prev ? { ...prev, balance: res.newBalance } : null)
      settleForm.resetFields()
      broadcast({ type: 'settled', amount: res.amount, newBalance: res.newBalance })
    } catch (err) {
      setError(err.message || '结算失败')
    } finally {
      setLoading(false)
    }
  }

  function handleReset() {
    setCardNo('')
    setCardInfo(null)
    setAlarm('')
    setTxResult(null)
    setError('')
    settleForm.resetFields()
    broadcast({ type: 'reset' })
  }

  return (
    <div style={{ minHeight: '100vh', background: '#0a1628', display: 'flex', flexDirection: 'column' }}>
      <style>{'.meal-amount-input .ant-input-number-input { color: #e8f4fd !important; background: #061020 !important; }'}</style>

      {/* 顶栏 */}
      <div style={{
        display: 'flex', alignItems: 'center', justifyContent: 'space-between',
        padding: '12px 24px', background: '#0d1f3c', borderBottom: '1px solid #1a3a6b',
      }}>
        <Title level={4} style={{ color: '#4fc3f7', margin: 0 }}>窗口操作端</Title>
        {windows.length > 0 && !lockedWindowId && (
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <Text style={{ color: '#8bafd4', fontSize: 13 }}>当前窗口：</Text>
            <Select
              value={selectedWindowId}
              onChange={setSelectedWindowId}
              options={windows.map(w => ({ value: w.id, label: w.name }))}
              style={{ width: 140 }}
              size="small"
            />
          </div>
        )}
        {lockedWindowId && (
          <Text style={{ color: '#8bafd4', fontSize: 13 }}>
            窗口：{windows.find(w => w.id === lockedWindowId)?.name || `#${lockedWindowId}`}
          </Text>
        )}
        <Button
          icon={<HomeOutlined />}
          size="small"
          onClick={() => navigate('/')}
          style={{ background: 'transparent', borderColor: '#1a3a6b', color: '#8bafd4' }}
        >
          首页
        </Button>
      </div>

      {/* 步骤条 */}
      <div style={{ padding: '16px 40px 0', background: '#0d1f3c' }}>
        <Steps
          current={currentStep}
          size="small"
          items={[
            { title: <span style={{ color: currentStep >= 0 ? '#4fc3f7' : '#4a6785' }}>刷卡验证</span> },
            { title: <span style={{ color: currentStep >= 1 ? '#4fc3f7' : '#4a6785' }}>输入金额</span> },
            { title: <span style={{ color: currentStep >= 2 ? '#52c41a' : '#4a6785' }}>结算完成</span> },
          ]}
          style={{ maxWidth: 500 }}
        />
      </div>

      {/* 主内容区 */}
      <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '32px 24px' }}>
        <div style={{ width: '100%', maxWidth: 560 }}>

          {/* 警报区 */}
          {alarm && (
            <Card style={{ background: '#2d0a0a', border: '2px solid #ff4d4f', borderRadius: 12, marginBottom: 24, textAlign: 'center' }}>
              <WarningFilled style={{ fontSize: 48, color: '#ff4d4f', display: 'block', marginBottom: 12 }} />
              <Text style={{ fontSize: 22, fontWeight: 'bold', color: '#ff4d4f', display: 'block' }}>{alarm}</Text>
              <Button size="large" onClick={handleReset} style={{ marginTop: 16, background: '#1a1a1a', borderColor: '#4a4a4a', color: '#ccc' }}>
                重新刷卡
              </Button>
            </Card>
          )}

          {/* 步骤 0：刷卡输入 */}
          {currentStep === 0 && !alarm && (
            <Card style={{ background: '#0d1f3c', border: '1px solid #1a3a6b', borderRadius: 12 }}>
              <Text style={{ color: '#8bafd4', display: 'block', marginBottom: 16, fontSize: 15 }}>
                请刷卡或输入卡号（16位）
              </Text>
              <div style={{ display: 'flex', gap: 12 }}>
                <Input
                  size="large"
                  placeholder="卡号"
                  value={cardNo}
                  onChange={e => { setCardNo(e.target.value); setAlarm(''); setError('') }}
                  onPressEnter={handleQueryCard}
                  style={{ flex: 1, fontSize: 22, height: 56, background: '#061020', borderColor: '#1a3a6b', color: '#e8f4fd' }}
                  styles={{ input: { background: '#061020', color: '#e8f4fd' } }}
                  autoFocus
                />
                <Button type="primary" size="large" loading={loading} onClick={handleQueryCard} style={{ height: 56, fontSize: 18, padding: '0 28px' }}>
                  刷卡
                </Button>
              </div>
              {error && <Alert type="error" message={error} style={{ marginTop: 16 }} showIcon />}
            </Card>
          )}

          {/* 步骤 1：显示余额 + 输入金额 */}
          {cardInfo && !txResult && (
            <>
              <Card style={{ background: '#0d1f3c', border: '1px solid #1a3a6b', borderRadius: 12, marginBottom: 16, textAlign: 'center' }}>
                <Text style={{ color: '#8bafd4', fontSize: 15, display: 'block', marginBottom: 4 }}>
                  {cardInfo.cardHolder.name}（{cardInfo.cardNo}）
                </Text>
                <div style={{ margin: '16px 0' }}>
                  <Text style={{ color: '#8bafd4', fontSize: 16 }}>卡内余额</Text>
                  <div>
                    <span style={{ fontSize: 64, fontWeight: 'bold', color: '#4fc3f7', lineHeight: 1.1 }}>
                      {(cardInfo.balance / 100).toFixed(2)}
                    </span>
                    <span style={{ fontSize: 22, color: '#8bafd4', marginLeft: 6 }}>元</span>
                  </div>
                </div>
              </Card>

              <Card style={{ background: '#0d1f3c', border: '1px solid #1a3a6b', borderRadius: 12 }}>
                <Form form={settleForm} layout="vertical" onFinish={handleSettle}>
                  <Form.Item
                    label={<span style={{ color: '#8bafd4', fontSize: 16 }}>本次消费金额（元）</span>}
                    name="amount"
                    rules={[{ required: true, message: '请输入消费金额' }]}
                    style={{ marginBottom: 16 }}
                  >
                    <InputNumber
                      min={0.01} step={0.01} precision={2} placeholder="0.00" size="large"
                      className="meal-amount-input"
                      style={{ width: '100%', fontSize: 28, height: 64, background: '#061020', borderColor: '#1a3a6b', color: '#e8f4fd' }}
                      autoFocus
                    />
                  </Form.Item>
                  <Button type="primary" htmlType="submit" loading={loading} size="large" block style={{ height: 56, fontSize: 20 }}>
                    确认结算
                  </Button>
                </Form>
                {error && <Alert type="error" message={error} style={{ marginTop: 12 }} showIcon />}
              </Card>
            </>
          )}

          {/* 步骤 2：结算完成 */}
          {txResult && (
            <Card style={{ background: '#061a0a', border: '2px solid #52c41a', borderRadius: 12, textAlign: 'center' }}>
              <CheckCircleOutlined style={{ fontSize: 56, color: '#52c41a', display: 'block', marginBottom: 16 }} />
              <Text style={{ color: '#52c41a', fontSize: 22, fontWeight: 'bold', display: 'block', marginBottom: 24 }}>结算成功</Text>
              <div style={{ display: 'flex', justifyContent: 'space-around', marginBottom: 32 }}>
                <div>
                  <Text style={{ color: '#5a8a6a', display: 'block', fontSize: 14 }}>本次消费</Text>
                  <span style={{ fontSize: 36, fontWeight: 'bold', color: '#ff7a45' }}>{(txResult.amount / 100).toFixed(2)}</span>
                  <span style={{ fontSize: 16, color: '#5a8a6a', marginLeft: 4 }}>元</span>
                </div>
                <div>
                  <Text style={{ color: '#5a8a6a', display: 'block', fontSize: 14 }}>扣款后余额</Text>
                  <span style={{ fontSize: 36, fontWeight: 'bold', color: '#52c41a' }}>{(txResult.newBalance / 100).toFixed(2)}</span>
                  <span style={{ fontSize: 16, color: '#5a8a6a', marginLeft: 4 }}>元</span>
                </div>
              </div>
              <Button type="primary" size="large" onClick={handleReset} style={{ height: 56, fontSize: 18, padding: '0 48px' }}>
                下一位
              </Button>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}
