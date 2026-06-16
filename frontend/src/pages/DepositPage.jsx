import { useState } from 'react'
import { Form, Input, InputNumber, Button, Alert, Descriptions, Card, Typography, Space, Table } from 'antd'
import { getCard, deposit, getCardDeposits } from '../api.js'

const { Title } = Typography

export default function DepositPage() {
  const [cardNo, setCardNo] = useState('')
  const [cardInfo, setCardInfo] = useState(null)
  const [receipt, setReceipt] = useState(null)
  const [error, setError] = useState('')
  const [queryLoading, setQueryLoading] = useState(false)
  const [depositLoading, setDepositLoading] = useState(false)
  const [depositForm] = Form.useForm()

  // 存款历史
  const [history, setHistory] = useState(null)
  const [historyPage, setHistoryPage] = useState(1)

  async function handleQueryCard() {
    if (!cardNo.trim()) return
    setError('')
    setCardInfo(null)
    setReceipt(null)
    setHistory(null)
    setQueryLoading(true)
    try {
      const res = await getCard(cardNo.trim())
      if (res.status !== 'active') {
        const statusText = res.status === 'lost' ? '已挂失' : '已注销'
        setError(`该卡${statusText}，无法充值`)
        return
      }
      setCardInfo(res)
      // 同时加载存款历史
      loadHistory(1)
    } catch (err) {
      setError(err.message || '查询失败')
    } finally {
      setQueryLoading(false)
    }
  }

  async function loadHistory(page) {
    try {
      const res = await getCardDeposits(cardNo.trim(), { page, pageSize: 5 })
      setHistory(res)
      setHistoryPage(page)
    } catch {}
  }

  async function handleDeposit(values) {
    setError('')
    setDepositLoading(true)
    try {
      const amountFen = Math.round(values.amount * 100)
      const res = await deposit(cardNo.trim(), amountFen)
      setReceipt(res)
      setCardInfo(prev => prev ? { ...prev, balance: res.newBalance } : null)
      depositForm.resetFields()
      // 刷新历史
      loadHistory(1)
    } catch (err) {
      setError(err.message || '存款失败')
    } finally {
      setDepositLoading(false)
    }
  }

  const historyColumns = [
    { title: '金额', dataIndex: 'amount', key: 'amount', render: v => `${(v / 100).toFixed(2)} 元` },
    { title: '时间', dataIndex: 'createdAt', key: 'createdAt', render: v => new Date(v).toLocaleString('zh-CN') },
  ]

  return (
    <Card style={{ maxWidth: 520, margin: '0 auto' }}>
      <Title level={4} style={{ marginTop: 0 }}>存款</Title>

      <Space.Compact style={{ width: '100%', marginBottom: 16 }}>
        <Input
          placeholder="请输入16位卡号"
          value={cardNo}
          onChange={e => { setCardNo(e.target.value); setCardInfo(null); setReceipt(null); setHistory(null); setError('') }}
          onPressEnter={handleQueryCard}
        />
        <Button type="primary" loading={queryLoading} onClick={handleQueryCard}>
          查询
        </Button>
      </Space.Compact>

      {cardInfo && (
        <Descriptions column={1} size="small" bordered style={{ marginBottom: 16 }}>
          <Descriptions.Item label="持卡人">{cardInfo.cardHolder.name}</Descriptions.Item>
          <Descriptions.Item label="卡号">{cardInfo.cardNo}</Descriptions.Item>
          <Descriptions.Item label="当前余额">{(cardInfo.balance / 100).toFixed(2)} 元</Descriptions.Item>
        </Descriptions>
      )}

      {cardInfo && (
        <Form form={depositForm} layout="vertical" onFinish={handleDeposit}>
          <Form.Item label="存款金额（元）" name="amount" rules={[{ required: true, message: '请输入存款金额' }]}>
            <InputNumber min={0.01} step={0.01} precision={2} placeholder="0.00" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={depositLoading}>确认存款</Button>
          </Form.Item>
        </Form>
      )}

      {error && <Alert type="error" message={error} style={{ marginTop: 8 }} showIcon />}

      {receipt && (
        <>
          <style>{`
            @media print {
              body * { visibility: hidden; }
              .deposit-receipt-print, .deposit-receipt-print * { visibility: visible; }
              .deposit-receipt-print { position: fixed; top: 0; left: 0; width: 100%; }
            }
          `}</style>
          <Card
            size="small"
            title="存款收据"
            className="deposit-receipt-print"
            style={{ marginTop: 16, background: '#f6ffed', borderColor: '#b7eb8f' }}
            extra={<Button size="small" onClick={() => window.print()}>打印收据</Button>}
          >
            <Descriptions column={1} size="small">
              <Descriptions.Item label="单号">{receipt.id}</Descriptions.Item>
              <Descriptions.Item label="卡号">{receipt.cardNo}</Descriptions.Item>
              <Descriptions.Item label="持卡人">{receipt.holderName}</Descriptions.Item>
              <Descriptions.Item label="充值金额">{(receipt.amount / 100).toFixed(2)} 元</Descriptions.Item>
              <Descriptions.Item label="充值后余额">{(receipt.newBalance / 100).toFixed(2)} 元</Descriptions.Item>
              <Descriptions.Item label="充值时间">{new Date(receipt.createdAt).toLocaleString('zh-CN')}</Descriptions.Item>
            </Descriptions>
          </Card>
        </>
      )}

      {history && history.total > 0 && (
        <Card size="small" title="充值记录" style={{ marginTop: 16 }}>
          <Table
            size="small"
            rowKey="id"
            dataSource={history.records}
            columns={historyColumns}
            pagination={{
              current: historyPage,
              total: history.total,
              pageSize: 5,
              size: 'small',
              onChange: p => loadHistory(p),
            }}
          />
        </Card>
      )}
    </Card>
  )
}
