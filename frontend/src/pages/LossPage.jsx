import { useState } from 'react'
import { getCard, reportLoss, cancelLossReport } from '../api.js'

const STATUS_LABEL = { active: '正常', lost: '已挂失', cancelled: '已注销' }

export default function LossPage() {
  const [cardId, setCardId] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [successMsg, setSuccessMsg] = useState('')

  async function handleQuery(e) {
    e.preventDefault()
    setError('')
    setCardInfo(null)
    setSuccessMsg('')
    setLoading(true)
    try {
      const res = await getCard(cardId)
      setCardInfo(res)
    } catch (err) {
      setError(err.message || '查询失败')
    } finally {
      setLoading(false)
    }
  }

  async function handleReportLoss() {
    setError('')
    setSuccessMsg('')
    setLoading(true)
    try {
      const res = await reportLoss(cardId)
      setCardInfo(res)
      setSuccessMsg('挂失成功，该卡已被标记为挂失状态')
    } catch (err) {
      setError(err.message || '挂失失败')
    } finally {
      setLoading(false)
    }
  }

  async function handleCancelLoss() {
    setError('')
    setSuccessMsg('')
    setLoading(true)
    try {
      const res = await cancelLossReport(cardId)
      setCardInfo(res)
      setSuccessMsg('取消挂失成功，该卡已恢复正常')
    } catch (err) {
      setError(err.message || '取消挂失失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 480, margin: '0 auto', padding: 24 }}>
      <h2>挂失管理</h2>

      <form onSubmit={handleQuery} style={{ marginBottom: 16 }}>
        <div style={{ display: 'flex', gap: 8 }}>
          <input
            placeholder="请输入卡号"
            value={cardId}
            onChange={e => { setCardId(e.target.value); setCardInfo(null); setSuccessMsg(''); setError('') }}
            required
            style={{ flex: 1, padding: 8 }}
          />
          <button type="submit" disabled={loading} style={{ padding: '8px 16px' }}>
            查询
          </button>
        </div>
      </form>

      {cardInfo && (
        <div style={{ padding: 12, background: '#f5f5f5', borderRadius: 4, marginBottom: 16 }}>
          <p><strong>卡号：</strong>{cardInfo.id}</p>
          <p><strong>持卡人：</strong>{cardInfo.cardHolder.name}</p>
          <p><strong>证件号：</strong>{cardInfo.cardHolder.idNumber}</p>
          <p>
            <strong>状态：</strong>
            <span style={{
              color: cardInfo.status === 'active' ? '#2e7d32' : cardInfo.status === 'lost' ? '#e65100' : '#757575',
              fontWeight: 'bold',
            }}>
              {STATUS_LABEL[cardInfo.status]}
            </span>
          </p>
          <p><strong>余额：</strong>{(cardInfo.balance / 100).toFixed(2)} 元</p>

          <div style={{ marginTop: 12, display: 'flex', gap: 8 }}>
            {cardInfo.status === 'active' && (
              <button
                onClick={handleReportLoss}
                disabled={loading}
                style={{ padding: '8px 20px', background: '#e65100', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer' }}
              >
                申请挂失
              </button>
            )}
            {cardInfo.status === 'lost' && (
              <button
                onClick={handleCancelLoss}
                disabled={loading}
                style={{ padding: '8px 20px', background: '#1565c0', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer' }}
              >
                取消挂失
              </button>
            )}
            {cardInfo.status === 'cancelled' && (
              <span style={{ color: '#757575' }}>该卡已注销，无法操作</span>
            )}
          </div>
        </div>
      )}

      {successMsg && (
        <div style={{ padding: 12, background: '#e8f5e9', borderRadius: 4, color: '#2e7d32' }}>
          {successMsg}
        </div>
      )}

      {error && (
        <div style={{ padding: 12, background: '#ffe0e0', borderRadius: 4, color: '#c00' }}>
          {error}
        </div>
      )}
    </div>
  )
}
