import { useState, useEffect, useRef } from 'react'
import { Input, Button, Card, Typography, Alert } from 'antd'
import { CheckCircleOutlined, WarningFilled, CreditCardOutlined } from '@ant-design/icons'
import { useSearchParams } from 'react-router-dom'
import { getCard } from '../api.js'

const { Text } = Typography

// 顾客屏：学生刷卡入口 + 结果展示
// 学生在此输入卡号，操作端接收后结算，结果回传显示
export default function CustomerScreen() {
  const [searchParams] = useSearchParams()
  const windowId = searchParams.get('id')
  const channelRef = useRef(null)

  const [cardNo, setCardNo] = useState('')
  // 状态：idle / card_read / waiting / settled / alarm
  const [state, setState] = useState('idle')
  const [data, setData] = useState({})
  const [error, setError] = useState('')

  useEffect(() => {
    const channelName = windowId ? `window-${windowId}` : 'window-default'
    const ch = new BroadcastChannel(channelName)
    channelRef.current = ch

    ch.onmessage = (e) => {
      const msg = e.data
      switch (msg.type) {
        case 'settled':
          setState('settled')
          setData({ amount: msg.amount, newBalance: msg.newBalance })
          break
        case 'alarm':
          setState('alarm')
          setData({ message: msg.message })
          break
        case 'reset':
          setState('idle')
          setData({})
          setCardNo('')
          setError('')
          break
      }
    }

    return () => ch.close()
  }, [windowId])

  function broadcast(msg) {
    channelRef.current?.postMessage(msg)
  }

  async function handleSwipe() {
    if (!cardNo.trim()) return
    setError('')
    try {
      const res = await getCard(cardNo.trim())
      if (res.status === 'cancelled') {
        setState('alarm')
        setData({ message: '此卡已注销' })
        broadcast({ type: 'alarm', message: '此卡已注销', cardNo: cardNo.trim() })
        return
      }
      if (res.status === 'lost') {
        setState('alarm')
        setData({ message: '此卡已挂失' })
        broadcast({ type: 'alarm', message: '此卡已挂失', cardNo: cardNo.trim() })
        return
      }
      setState('waiting')
      setData({ holderName: res.cardHolder.name, balance: res.balance, cardNo: res.cardNo })
      // 通知操作端：学生已刷卡
      broadcast({ type: 'card_read', cardNo: res.cardNo, holderName: res.cardHolder.name, balance: res.balance })
    } catch (err) {
      if (err.status === 404) {
        setState('alarm')
        setData({ message: '非本单位卡' })
        broadcast({ type: 'alarm', message: '非本单位卡', cardNo: cardNo.trim() })
      } else {
        setError(err.message || '刷卡失败')
      }
    }
  }

  return (
    <div style={{
      minHeight: '100vh', background: '#0a1628',
      display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', padding: 40,
    }}>
      {/* 空闲：学生刷卡入口 */}
      {state === 'idle' && (
        <Card style={{ background: '#0d1f3c', border: '2px solid #1a3a6b', borderRadius: 16, padding: '40px', width: '100%', maxWidth: 480, textAlign: 'center' }}>
          <CreditCardOutlined style={{ fontSize: 64, color: '#4fc3f7', marginBottom: 20 }} />
          <Text style={{ fontSize: 24, color: '#e8f4fd', display: 'block', marginBottom: 24 }}>请刷卡</Text>
          <div style={{ display: 'flex', gap: 12 }}>
            <Input
              size="large"
              placeholder="请输入卡号"
              value={cardNo}
              onChange={e => setCardNo(e.target.value)}
              onPressEnter={handleSwipe}
              style={{ fontSize: 20, height: 56, background: '#061020', borderColor: '#1a3a6b', color: '#e8f4fd' }}
              styles={{ input: { background: '#061020', color: '#e8f4fd' } }}
              autoFocus
            />
            <Button type="primary" size="large" onClick={handleSwipe} style={{ height: 56, fontSize: 18, padding: '0 28px' }}>
              确认
            </Button>
          </div>
          {error && <Alert type="error" message={error} style={{ marginTop: 16 }} showIcon />}
        </Card>
      )}

      {/* 等待阿姨结算 */}
      {state === 'waiting' && (
        <Card style={{ background: '#0d1f3c', border: '2px solid #1a3a6b', borderRadius: 16, textAlign: 'center', padding: '40px 60px' }}>
          <Text style={{ color: '#8bafd4', fontSize: 20, display: 'block', marginBottom: 8 }}>{data.holderName}</Text>
          <Text style={{ color: '#8bafd4', fontSize: 18, display: 'block', marginBottom: 16 }}>卡内余额</Text>
          <span style={{ fontSize: 72, fontWeight: 'bold', color: '#4fc3f7' }}>{(data.balance / 100).toFixed(2)}</span>
          <span style={{ fontSize: 24, color: '#8bafd4', marginLeft: 8 }}>元</span>
          <Text style={{ color: '#4a6785', fontSize: 14, display: 'block', marginTop: 24 }}>等待工作人员输入金额并结算...</Text>
        </Card>
      )}

      {/* 结算完成 */}
      {state === 'settled' && (
        <Card style={{ background: '#061a0a', border: '2px solid #52c41a', borderRadius: 16, textAlign: 'center', padding: '40px 60px' }}>
          <CheckCircleOutlined style={{ fontSize: 64, color: '#52c41a', marginBottom: 16 }} />
          <Text style={{ color: '#52c41a', fontSize: 24, display: 'block', marginBottom: 32 }}>结算完成</Text>
          <div style={{ display: 'flex', gap: 60, justifyContent: 'center' }}>
            <div>
              <Text style={{ color: '#5a8a6a', display: 'block', fontSize: 16 }}>本次消费</Text>
              <span style={{ fontSize: 48, fontWeight: 'bold', color: '#ff7a45' }}>{(data.amount / 100).toFixed(2)}</span>
              <span style={{ fontSize: 20, color: '#5a8a6a', marginLeft: 4 }}>元</span>
            </div>
            <div>
              <Text style={{ color: '#5a8a6a', display: 'block', fontSize: 16 }}>剩余余额</Text>
              <span style={{ fontSize: 48, fontWeight: 'bold', color: '#52c41a' }}>{(data.newBalance / 100).toFixed(2)}</span>
              <span style={{ fontSize: 20, color: '#5a8a6a', marginLeft: 4 }}>元</span>
            </div>
          </div>
        </Card>
      )}

      {/* 报警 */}
      {state === 'alarm' && (
        <Card style={{ background: '#2d0a0a', border: '2px solid #ff4d4f', borderRadius: 16, textAlign: 'center', padding: '40px 60px' }}>
          <WarningFilled style={{ fontSize: 72, color: '#ff4d4f', marginBottom: 16 }} />
          <Text style={{ fontSize: 28, fontWeight: 'bold', color: '#ff4d4f', display: 'block' }}>{data.message}</Text>
        </Card>
      )}
    </div>
  )
}
