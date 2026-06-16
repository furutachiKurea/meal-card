import { useState, useEffect, useRef } from 'react'
import { Typography, Card } from 'antd'
import { CheckCircleOutlined, WarningFilled, CreditCardOutlined } from '@ant-design/icons'
import { useSearchParams } from 'react-router-dom'

const { Text } = Typography

// 顾客屏：学生面对的大字只读界面
// 通过 BroadcastChannel 接收操作端推送的状态
export default function CustomerScreen() {
  const [searchParams] = useSearchParams()
  const windowId = searchParams.get('id')
  const channelRef = useRef(null)

  // 状态：idle / card_read / settled / alarm
  const [state, setState] = useState('idle')
  const [data, setData] = useState({})

  useEffect(() => {
    const channelName = windowId ? `window-${windowId}` : 'window-default'
    const ch = new BroadcastChannel(channelName)
    channelRef.current = ch

    ch.onmessage = (e) => {
      const msg = e.data
      switch (msg.type) {
        case 'card_read':
          setState('card_read')
          setData({ holderName: msg.holderName, balance: msg.balance })
          break
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
          break
      }
    }

    return () => ch.close()
  }, [windowId])

  return (
    <div style={{
      minHeight: '100vh',
      background: '#0a1628',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      padding: 40,
    }}>
      {/* 空闲状态 */}
      {state === 'idle' && (
        <div style={{ textAlign: 'center' }}>
          <CreditCardOutlined style={{ fontSize: 80, color: '#1a3a6b', marginBottom: 24 }} />
          <Text style={{ fontSize: 28, color: '#4a6785', display: 'block' }}>
            请刷卡
          </Text>
        </div>
      )}

      {/* 读卡成功，显示余额 */}
      {state === 'card_read' && (
        <Card style={{
          background: '#0d1f3c',
          border: '2px solid #1a3a6b',
          borderRadius: 16,
          textAlign: 'center',
          padding: '40px 60px',
        }}>
          <Text style={{ color: '#8bafd4', fontSize: 20, display: 'block', marginBottom: 8 }}>
            {data.holderName}
          </Text>
          <Text style={{ color: '#8bafd4', fontSize: 18, display: 'block', marginBottom: 16 }}>
            卡内余额
          </Text>
          <span style={{ fontSize: 80, fontWeight: 'bold', color: '#4fc3f7' }}>
            {(data.balance / 100).toFixed(2)}
          </span>
          <span style={{ fontSize: 28, color: '#8bafd4', marginLeft: 8 }}>元</span>
        </Card>
      )}

      {/* 结算完成 */}
      {state === 'settled' && (
        <Card style={{
          background: '#061a0a',
          border: '2px solid #52c41a',
          borderRadius: 16,
          textAlign: 'center',
          padding: '40px 60px',
        }}>
          <CheckCircleOutlined style={{ fontSize: 64, color: '#52c41a', marginBottom: 16 }} />
          <Text style={{ color: '#52c41a', fontSize: 24, display: 'block', marginBottom: 32 }}>
            结算完成
          </Text>
          <div style={{ display: 'flex', gap: 60, justifyContent: 'center' }}>
            <div>
              <Text style={{ color: '#5a8a6a', display: 'block', fontSize: 16 }}>本次消费</Text>
              <span style={{ fontSize: 48, fontWeight: 'bold', color: '#ff7a45' }}>
                {(data.amount / 100).toFixed(2)}
              </span>
              <span style={{ fontSize: 20, color: '#5a8a6a', marginLeft: 4 }}>元</span>
            </div>
            <div>
              <Text style={{ color: '#5a8a6a', display: 'block', fontSize: 16 }}>剩余余额</Text>
              <span style={{ fontSize: 48, fontWeight: 'bold', color: '#52c41a' }}>
                {(data.newBalance / 100).toFixed(2)}
              </span>
              <span style={{ fontSize: 20, color: '#5a8a6a', marginLeft: 4 }}>元</span>
            </div>
          </div>
        </Card>
      )}

      {/* 报警 */}
      {state === 'alarm' && (
        <Card style={{
          background: '#2d0a0a',
          border: '2px solid #ff4d4f',
          borderRadius: 16,
          textAlign: 'center',
          padding: '40px 60px',
        }}>
          <WarningFilled style={{ fontSize: 72, color: '#ff4d4f', marginBottom: 16 }} />
          <Text style={{ fontSize: 28, fontWeight: 'bold', color: '#ff4d4f', display: 'block' }}>
            {data.message}
          </Text>
        </Card>
      )}
    </div>
  )
}
