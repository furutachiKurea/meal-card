import { useState } from 'react'
import { getCard, cancelCard } from '../api.js'

const STATUS_LABEL = { active: '正常', lost: '已挂失', cancelled: '已注销' }

export default function CancelPage() {
  const [cardId, setCardId] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [result, setResult] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [confirmed, setConfirmed] = useState(false)

  async function handleQuery(e) {
    e.preventDefault()
    setError('')
    setCardInfo(null)
    setResult(null)
    setConfirmed(false)
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

  async function handleCancel() {
    setError('')
    setLoading(true)
    try {
      const res = await cancelCard(cardId)
      setResult(res)
      setCardInfo(null)
    } catch (err) {
      setError(err.message || '注销失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 480, margin: '0 auto', padding: 24 }}>
      <h2>注销</h2>

      <form onSubmit={handleQuery} style={{ marginBottom: 16 }}>
        <div style={{ display: 'flex', gap: 8 }}>
          <input
            placeholder="请输入卡号"
            value={cardId}
            onChange={e => { setCardId(e.target.value); setCardInfo(null); setResult(null); setConfirmed(false); setError('') }}
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
            <span style={{ color: cardInfo.status === 'active' ? '#2e7d32' : '#e65100', fontWeight: 'bold' }}>
              {STATUS_LABEL[cardInfo.status]}
            </span>
          </p>
          <p><strong>余额：</strong>{(cardInfo.balance / 100).toFixed(2)} 元</p>
          <p><strong>押金：</strong>{(cardInfo.deposit / 100).toFixed(2)} 元</p>
          <p><strong>预计退款：</strong>{((cardInfo.balance + cardInfo.deposit) / 100).toFixed(2)} 元</p>

          {cardInfo.status === 'cancelled' ? (
            <p style={{ color: '#757575' }}>该卡已注销</p>
          ) : (
            <>
              {!confirmed ? (
                <button
                  onClick={() => setConfirmed(true)}
                  style={{ marginTop: 12, padding: '8px 20px', background: '#c62828', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer' }}
                >
                  申请注销
                </button>
              ) : (
                <div style={{ marginTop: 12 }}>
                  <p style={{ color: '#c62828', fontWeight: 'bold' }}>确认要注销该卡吗？此操作不可撤销！</p>
                  <div style={{ display: 'flex', gap: 8 }}>
                    <button
                      onClick={handleCancel}
                      disabled={loading}
                      style={{ padding: '8px 20px', background: '#c62828', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer' }}
                    >
                      {loading ? '处理中...' : '确认注销'}
                    </button>
                    <button
                      onClick={() => setConfirmed(false)}
                      style={{ padding: '8px 16px' }}
                    >
                      取消
                    </button>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      )}

      {error && (
        <div style={{ padding: 12, background: '#ffe0e0', borderRadius: 4, color: '#c00' }}>
          {error}
        </div>
      )}

      {result && (
        <div style={{ padding: 16, background: '#e8f5e9', borderRadius: 4 }}>
          <h3>注销成功</h3>
          <hr />
          <p><strong>卡号：</strong>{result.card.id}</p>
          <h4>退款明细</h4>
          <p><strong>退还押金：</strong>{(result.refund.deposit / 100).toFixed(2)} 元</p>
          <p><strong>退还余额：</strong>{(result.refund.balance / 100).toFixed(2)} 元</p>
          <p style={{ fontSize: 18, color: '#2e7d32' }}>
            <strong>应退合计：</strong>{(result.refund.total / 100).toFixed(2)} 元
          </p>
        </div>
      )}
    </div>
  )
}
