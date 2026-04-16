import { useState } from 'react'
import { Input, Button, Alert, Descriptions, Card, Typography, Space, Tag } from 'antd'
import { getCardByIDNumber, reportLoss, cancelLossReport } from '../api.js'

const { Title } = Typography

const STATUS_LABEL = { active: '正常', lost: '已挂失', cancelled: '已注销' }
const STATUS_COLOR = { active: 'success', lost: 'warning', cancelled: 'default' }

export default function LossPage() {
  const [idNumber, setIdNumber] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [successMsg, setSuccessMsg] = useState('')

  async function handleQuery() {
    if (!idNumber.trim()) return
    setError('')
    setCardInfo(null)
    setSuccessMsg('')
    setLoading(true)
    try {
      const res = await getCardByIDNumber(idNumber.trim())
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
      const res = await reportLoss(cardInfo.cardNo)
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
      const res = await cancelLossReport(cardInfo.cardNo)
      setCardInfo(res)
      setSuccessMsg('取消挂失成功，该卡已恢复正常')
    } catch (err) {
      setError(err.message || '取消挂失失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card style={{ maxWidth: 520, margin: '0 auto' }}>
      <Title level={4} style={{ marginTop: 0 }}>挂失管理</Title>

      <Space.Compact style={{ width: '100%', marginBottom: 16 }}>
        <Input
          placeholder="请输入证件号（12位）"
          value={idNumber}
          onChange={e => { setIdNumber(e.target.value); setCardInfo(null); setSuccessMsg(''); setError('') }}
          onPressEnter={handleQuery}
        />
        <Button type="primary" loading={loading} onClick={handleQuery}>
          查询
        </Button>
      </Space.Compact>

      {cardInfo && (
        <Card size="small" style={{ marginBottom: 16 }}>
          <Descriptions column={1} size="small">
            <Descriptions.Item label="卡号">{cardInfo.cardNo}</Descriptions.Item>
            <Descriptions.Item label="持卡人">{cardInfo.cardHolder.name}</Descriptions.Item>
            <Descriptions.Item label="证件号">{cardInfo.cardHolder.idNumber}</Descriptions.Item>
            <Descriptions.Item label="状态">
              <Tag color={STATUS_COLOR[cardInfo.status]}>{STATUS_LABEL[cardInfo.status]}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="余额">{(cardInfo.balance / 100).toFixed(2)} 元</Descriptions.Item>
          </Descriptions>

          <div style={{ marginTop: 12 }}>
            {cardInfo.status === 'active' && (
              <Button danger onClick={handleReportLoss} loading={loading}>
                申请挂失
              </Button>
            )}
            {cardInfo.status === 'lost' && (
              <Button type="primary" onClick={handleCancelLoss} loading={loading}>
                取消挂失
              </Button>
            )}
            {cardInfo.status === 'cancelled' && (
              <span style={{ color: '#999' }}>该卡已注销，无法操作</span>
            )}
          </div>
        </Card>
      )}

      {successMsg && (
        <Alert type="success" message={successMsg} style={{ marginBottom: 8 }} showIcon />
      )}

      {error && (
        <Alert type="error" message={error} style={{ marginBottom: 8 }} showIcon />
      )}
    </Card>
  )
}
