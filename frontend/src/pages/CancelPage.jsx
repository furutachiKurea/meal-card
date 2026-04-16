import { useState } from 'react'
import { Input, Button, Alert, Descriptions, Card, Typography, Space, Tag, Modal } from 'antd'
import { ExclamationCircleFilled } from '@ant-design/icons'
import { getCardByIDNumber, cancelCard } from '../api.js'

const { Title } = Typography
const STATUS_LABEL = { active: '正常', lost: '已挂失', cancelled: '已注销' }
const STATUS_COLOR = { active: 'success', lost: 'warning', cancelled: 'default' }

export default function CancelPage() {
  const [idNumber, setIdNumber] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [result, setResult] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleQuery() {
    if (!idNumber.trim()) return
    setError('')
    setCardInfo(null)
    setResult(null)
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

  function handleCancelConfirm() {
    Modal.confirm({
      title: '确认注销',
      icon: <ExclamationCircleFilled />,
      content: '确认要注销该卡吗？此操作不可撤销！',
      okText: '确认注销',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        setError('')
        setLoading(true)
        try {
          const res = await cancelCard(cardInfo.cardNo)
          setResult(res)
          setCardInfo(null)
        } catch (err) {
          setError(err.message || '注销失败')
        } finally {
          setLoading(false)
        }
      },
    })
  }

  return (
    <Card style={{ maxWidth: 520, margin: '0 auto' }}>
      <Title level={4} style={{ marginTop: 0 }}>注销</Title>

      <Space.Compact style={{ width: '100%', marginBottom: 16 }}>
        <Input
          placeholder="请输入证件号（12位）"
          value={idNumber}
          onChange={e => { setIdNumber(e.target.value); setCardInfo(null); setResult(null); setError('') }}
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
            <Descriptions.Item label="当前余额">{(cardInfo.balance / 100).toFixed(2)} 元</Descriptions.Item>
            <Descriptions.Item label="押金">{(cardInfo.deposit / 100).toFixed(2)} 元</Descriptions.Item>
            <Descriptions.Item label="预计退款">
              <span style={{ fontWeight: 'bold' }}>
                {((cardInfo.balance + cardInfo.deposit) / 100).toFixed(2)} 元
              </span>
            </Descriptions.Item>
          </Descriptions>

          <div style={{ marginTop: 12 }}>
            {cardInfo.status === 'cancelled' ? (
              <span style={{ color: '#999' }}>该卡已注销</span>
            ) : (
              <Button danger type="primary" onClick={handleCancelConfirm} loading={loading}>
                申请注销
              </Button>
            )}
          </div>
        </Card>
      )}

      {error && (
        <Alert type="error" message={error} style={{ marginBottom: 8 }} showIcon />
      )}

      {result && (
        <Card
          size="small"
          title="注销成功"
          style={{ marginTop: 16, background: '#f6ffed', borderColor: '#b7eb8f' }}
        >
          <Descriptions column={1} size="small" title="退款明细">
            <Descriptions.Item label="卡号">{result.card.cardNo}</Descriptions.Item>
            <Descriptions.Item label="退还押金">{(result.refund.deposit / 100).toFixed(2)} 元</Descriptions.Item>
            <Descriptions.Item label="退还余额">{(result.refund.balance / 100).toFixed(2)} 元</Descriptions.Item>
            <Descriptions.Item label="应退合计">
              <span style={{ fontSize: 16, fontWeight: 'bold', color: '#52c41a' }}>
                {(result.refund.total / 100).toFixed(2)} 元
              </span>
            </Descriptions.Item>
          </Descriptions>
        </Card>
      )}
    </Card>
  )
}
