import { useState, useEffect } from 'react'
import { getCard, createTransaction, listWindows } from '../api.js'

export default function MealPage() {
  const [cardId, setCardId] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [alarm, setAlarm] = useState('')
  const [amount, setAmount] = useState('')
  const [txResult, setTxResult] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [windows, setWindows] = useState([])
  const [selectedWindowId, setSelectedWindowId] = useState('')

  useEffect(() => {
    listWindows().then(res => {
      setWindows(res.windows || [])
      if (res.windows && res.windows.length > 0) {
        setSelectedWindowId(String(res.windows[0].id))
      }
    }).catch(() => {})
  }, [])

  async function handleQueryCard(e) {
    e.preventDefault()
    setError('')
    setCardInfo(null)
    setAlarm('')
    setTxResult(null)
    setLoading(true)
    try {
      const res = await getCard(cardId)
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

  async function handleSettle(e) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const amountFen = Math.round(parseFloat(amount) * 100)
      const res = await createTransaction(cardId, {
        windowId: parseInt(selectedWindowId),
        amount: amountFen,
      })
      setTxResult(res)
      setCardInfo(prev => prev ? { ...prev, balance: res.newBalance } : null)
      setAmount('')
    } catch (err) {
      setError(err.message || '结算失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 480, margin: '0 auto', padding: 24 }}>
      <h2>就餐消费（窗口机）</h2>

      {windows.length > 0 && (
        <div style={{ marginBottom: 12 }}>
          <label>当前窗口</label>
          <select
            value={selectedWindowId}
            onChange={e => setSelectedWindowId(e.target.value)}
            style={{ display: 'block', width: '100%', padding: 8, marginTop: 4 }}
          >
            {windows.map(w => (
              <option key={w.id} value={String(w.id)}>{w.name}</option>
            ))}
          </select>
        </div>
      )}

      <form onSubmit={handleQueryCard} style={{ marginBottom: 16 }}>
        <div style={{ display: 'flex', gap: 8 }}>
          <input
            placeholder="请输入卡号（模拟刷卡）"
            value={cardId}
            onChange={e => { setCardId(e.target.value); setCardInfo(null); setAlarm(''); setTxResult(null); setError('') }}
            required
            style={{ flex: 1, padding: 8 }}
          />
          <button type="submit" disabled={loading} style={{ padding: '8px 16px' }}>
            刷卡
          </button>
        </div>
      </form>

      {alarm && (
        <div style={{
          padding: 16,
          background: '#ff1744',
          color: '#fff',
          borderRadius: 4,
          marginBottom: 16,
          fontWeight: 'bold',
          fontSize: 18,
          textAlign: 'center',
        }}>
          ⚠ {alarm}
        </div>
      )}

      {cardInfo && (
        <div style={{ padding: 12, background: '#e3f2fd', borderRadius: 4, marginBottom: 16 }}>
          <p><strong>持卡人：</strong>{cardInfo.cardHolder.name}</p>
          <p><strong>卡号：</strong>{cardInfo.id}</p>
          <p style={{ fontSize: 20, color: '#1565c0' }}>
            <strong>余额：</strong>{(cardInfo.balance / 100).toFixed(2)} 元
          </p>
        </div>
      )}

      {cardInfo && !txResult && (
        <form onSubmit={handleSettle}>
          <div style={{ marginBottom: 12 }}>
            <label>本次消费金额（元）</label>
            <input
              type="number"
              step="0.01"
              min="0.01"
              value={amount}
              onChange={e => setAmount(e.target.value)}
              required
              style={{ display: 'block', width: '100%', padding: 8, marginTop: 4, fontSize: 18 }}
            />
          </div>
          <button type="submit" disabled={loading} style={{ padding: '10px 32px', fontSize: 16 }}>
            {loading ? '结算中...' : '确认结算'}
          </button>
        </form>
      )}

      {error && (
        <div style={{ marginTop: 16, padding: 12, background: '#ffe0e0', borderRadius: 4, color: '#c00' }}>
          {error}
        </div>
      )}

      {txResult && (
        <div style={{ marginTop: 16, padding: 16, background: '#e8f5e9', borderRadius: 4 }}>
          <h3>结算成功</h3>
          <p><strong>消费金额：</strong>{(txResult.amount / 100).toFixed(2)} 元</p>
          <p style={{ fontSize: 18, color: '#2e7d32' }}>
            <strong>扣款后余额：</strong>{(txResult.newBalance / 100).toFixed(2)} 元
          </p>
          <button
            onClick={() => { setCardInfo(null); setCardId(''); setTxResult(null); setAmount('') }}
            style={{ marginTop: 8, padding: '6px 16px' }}
          >
            下一位
          </button>
        </div>
      )}
    </div>
  )
}
