import { useState } from 'react'
import { getCard, deposit } from '../api.js'

export default function DepositPage() {
  const [cardId, setCardId] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [amount, setAmount] = useState('')
  const [receipt, setReceipt] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleQueryCard(e) {
    e.preventDefault()
    setError('')
    setCardInfo(null)
    setReceipt(null)
    setLoading(true)
    try {
      const res = await getCard(cardId)
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

  async function handleDeposit(e) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const amountFen = Math.round(parseFloat(amount) * 100)
      const res = await deposit(cardId, amountFen)
      setReceipt(res)
      setCardInfo(prev => prev ? { ...prev, balance: res.newBalance } : null)
      setAmount('')
    } catch (err) {
      setError(err.message || '存款失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 480, margin: '0 auto', padding: 24 }}>
      <h2>存款</h2>

      <form onSubmit={handleQueryCard} style={{ marginBottom: 16 }}>
        <div style={{ display: 'flex', gap: 8 }}>
          <input
            placeholder="请输入卡号"
            value={cardId}
            onChange={e => setCardId(e.target.value)}
            required
            style={{ flex: 1, padding: 8 }}
          />
          <button type="submit" disabled={loading} style={{ padding: '8px 16px' }}>
            查询
          </button>
        </div>
      </form>

      {cardInfo && (
        <div style={{ padding: 12, background: '#e3f2fd', borderRadius: 4, marginBottom: 16 }}>
          <p><strong>持卡人：</strong>{cardInfo.cardHolder.name}</p>
          <p><strong>卡号：</strong>{cardInfo.id}</p>
          <p><strong>当前余额：</strong>{(cardInfo.balance / 100).toFixed(2)} 元</p>
        </div>
      )}

      {cardInfo && (
        <form onSubmit={handleDeposit}>
          <div style={{ marginBottom: 12 }}>
            <label>存款金额（元）</label>
            <input
              type="number"
              step="0.01"
              min="0.01"
              value={amount}
              onChange={e => setAmount(e.target.value)}
              required
              style={{ display: 'block', width: '100%', padding: 8, marginTop: 4 }}
            />
          </div>
          <button type="submit" disabled={loading} style={{ padding: '8px 24px' }}>
            {loading ? '处理中...' : '确认存款'}
          </button>
        </form>
      )}

      {error && (
        <div style={{ marginTop: 16, padding: 12, background: '#ffe0e0', borderRadius: 4, color: '#c00' }}>
          {error}
        </div>
      )}

      {receipt && (
        <div style={{ marginTop: 16, padding: 16, background: '#e8f5e9', borderRadius: 4 }}>
          <h3>存款收据</h3>
          <hr />
          <p><strong>卡号：</strong>{receipt.cardId}</p>
          <p><strong>持卡人：</strong>{receipt.holderName}</p>
          <p><strong>充值金额：</strong>{(receipt.amount / 100).toFixed(2)} 元</p>
          <p><strong>充值后余额：</strong>{(receipt.newBalance / 100).toFixed(2)} 元</p>
          <p><strong>充值时间：</strong>{new Date(receipt.createdAt).toLocaleString('zh-CN')}</p>
        </div>
      )}
    </div>
  )
}
