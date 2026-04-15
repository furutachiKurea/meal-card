import { useState } from 'react'
import { issueCard } from '../api.js'

export default function IssuePage() {
  const [form, setForm] = useState({ name: '', idNumber: '', deposit: '', preDeposit: '' })
  const [result, setResult] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  function handleChange(e) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    setResult(null)
    setLoading(true)
    try {
      const depositFen = Math.round(parseFloat(form.deposit) * 100)
      const preDepositFen = form.preDeposit ? Math.round(parseFloat(form.preDeposit) * 100) : 0
      const res = await issueCard({
        name: form.name,
        idNumber: form.idNumber,
        deposit: depositFen,
        preDeposit: preDepositFen,
      })
      setResult(res)
    } catch (err) {
      setError(err.message || '发卡失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 480, margin: '0 auto', padding: 24 }}>
      <h2>发卡</h2>
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: 12 }}>
          <label>持卡人姓名</label>
          <input
            name="name"
            value={form.name}
            onChange={handleChange}
            required
            style={{ display: 'block', width: '100%', padding: 8, marginTop: 4 }}
          />
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>证件号</label>
          <input
            name="idNumber"
            value={form.idNumber}
            onChange={handleChange}
            required
            style={{ display: 'block', width: '100%', padding: 8, marginTop: 4 }}
          />
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>押金（元）</label>
          <input
            name="deposit"
            type="number"
            step="0.01"
            min="0.01"
            value={form.deposit}
            onChange={handleChange}
            required
            style={{ display: 'block', width: '100%', padding: 8, marginTop: 4 }}
          />
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>预存金额（元）</label>
          <input
            name="preDeposit"
            type="number"
            step="0.01"
            min="0"
            value={form.preDeposit}
            onChange={handleChange}
            style={{ display: 'block', width: '100%', padding: 8, marginTop: 4 }}
          />
        </div>
        <button type="submit" disabled={loading} style={{ padding: '8px 24px' }}>
          {loading ? '处理中...' : '办理发卡'}
        </button>
      </form>

      {error && (
        <div style={{ marginTop: 16, padding: 12, background: '#ffe0e0', borderRadius: 4, color: '#c00' }}>
          {error}
        </div>
      )}

      {result && (
        <div style={{ marginTop: 16, padding: 16, background: '#e8f5e9', borderRadius: 4 }}>
          <h3>发卡成功</h3>
          <p><strong>卡号：</strong>{result.card.id}</p>
          <p><strong>持卡人：</strong>{result.cardHolder.name}</p>
          <p><strong>证件号：</strong>{result.cardHolder.idNumber}</p>
          <p><strong>押金：</strong>{(result.card.deposit / 100).toFixed(2)} 元</p>
          <p><strong>余额：</strong>{(result.card.balance / 100).toFixed(2)} 元</p>
          {result.refund && (
            <div style={{ marginTop: 8, padding: 8, background: '#fff3e0', borderRadius: 4 }}>
              <p><strong>旧卡自动注销（卡号 {result.refund.oldCardId}）</strong></p>
              <p>退还押金：{(result.refund.deposit / 100).toFixed(2)} 元</p>
              <p>退还余额：{(result.refund.balance / 100).toFixed(2)} 元</p>
              <p>退还合计：{(result.refund.total / 100).toFixed(2)} 元</p>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
