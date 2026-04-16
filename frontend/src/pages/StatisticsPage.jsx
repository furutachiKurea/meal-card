import { useState } from 'react'
import { Button, Alert, Card, Typography, Table, DatePicker, InputNumber, Space, Descriptions, Row, Col } from 'antd'
import dayjs from 'dayjs'
import {
  getMealRevenue,
  getWindowRevenue,
  getDepositDetails,
  getDepositSummary,
  getActiveBalance,
  getDailyReport,
  getYearlyReport,
} from '../api.js'

const { Title } = Typography
const { RangePicker } = DatePicker

function formatYuan(fen) {
  return (fen / 100).toFixed(2) + ' 元'
}

export default function StatisticsPage() {
  // 本餐售饭总收入
  const [revenueRange, setRevenueRange] = useState(null)
  const [mealRevenue, setMealRevenue] = useState(null)

  // 各窗口收入
  const [winRevenueRange, setWinRevenueRange] = useState(null)
  const [windowRevenue, setWindowRevenue] = useState(null)

  // 存款明细
  const [depositRange, setDepositRange] = useState(null)
  const [depositDetails, setDepositDetails] = useState(null)
  const [depositPage, setDepositPage] = useState(1)
  const depositPageSize = 5

  // 本日/本月存款
  const [depositSummary, setDepositSummary] = useState(null)

  // 流动资金
  const [activeBalance, setActiveBalance] = useState(null)

  // 日餐报表
  const [dailyDate, setDailyDate] = useState(dayjs())
  const [dailyReport, setDailyReport] = useState(null)

  // 年餐报表
  const [yearlyYear, setYearlyYear] = useState(new Date().getFullYear())
  const [yearlyReport, setYearlyReport] = useState(null)

  const [errors, setErrors] = useState({})
  const [loading, setLoading] = useState({})

  async function wrap(key, fn) {
    setLoading(l => ({ ...l, [key]: true }))
    setErrors(e => ({ ...e, [key]: '' }))
    try {
      await fn()
    } catch (err) {
      setErrors(e => ({ ...e, [key]: err.message || '查询失败' }))
    } finally {
      setLoading(l => ({ ...l, [key]: false }))
    }
  }

  const windowRevenueColumns = [
    { title: '窗口', dataIndex: 'windowName', key: 'windowName' },
    { title: '收入', dataIndex: 'revenue', key: 'revenue', align: 'right', render: v => formatYuan(v) },
  ]

  const dailyReportColumns = [
    { title: '窗口', dataIndex: 'windowName', key: 'windowName' },
    { title: '收入', dataIndex: 'revenue', key: 'revenue', align: 'right', render: v => formatYuan(v) },
    { title: '笔数', dataIndex: 'transactionCount', key: 'transactionCount', align: 'right' },
  ]

  const yearlyReportColumns = [
    { title: '月份', dataIndex: 'month', key: 'month', render: v => `${v} 月` },
    { title: '收入', dataIndex: 'revenue', key: 'revenue', align: 'right', render: v => formatYuan(v) },
    { title: '笔数', dataIndex: 'transactionCount', key: 'transactionCount', align: 'right' },
  ]

  const depositDetailColumns = [
    { title: '卡号', dataIndex: 'cardNo', key: 'cardNo' },
    { title: '金额', dataIndex: 'amount', key: 'amount', align: 'right', render: v => formatYuan(v) },
    { title: '时间', dataIndex: 'createdAt', key: 'createdAt', render: v => new Date(v).toLocaleString('zh-CN') },
  ]

  return (
    <div style={{ maxWidth: 860, margin: '0 auto' }}>
      <Title level={4}>汇总统计</Title>

      <Row gutter={[16, 16]}>
        {/* 本餐售饭总收入 */}
        <Col span={24}>
          <Card title="本餐售饭总收入" size="small">
            <Space wrap>
              <RangePicker
                showTime
                value={revenueRange}
                onChange={setRevenueRange}
              />
              <Button
                type="primary"
                loading={loading.mealRevenue}
                disabled={!revenueRange}
                onClick={() => wrap('mealRevenue', async () => {
                  const res = await getMealRevenue({
                    startTime: revenueRange[0].toISOString(),
                    endTime: revenueRange[1].toISOString(),
                  })
                  setMealRevenue(res)
                })}
              >
                查询
              </Button>
            </Space>
            {errors.mealRevenue && <Alert type="error" message={errors.mealRevenue} style={{ marginTop: 8 }} showIcon />}
            {mealRevenue && (
              <p style={{ marginTop: 8, fontSize: 16 }}>
                总收入：<strong>{formatYuan(mealRevenue.totalRevenue)}</strong>
              </p>
            )}
          </Card>
        </Col>

        {/* 各窗口收入 */}
        <Col span={24}>
          <Card title="各窗口收入" size="small">
            <Space wrap>
              <RangePicker
                showTime
                value={winRevenueRange}
                onChange={setWinRevenueRange}
              />
              <Button
                type="primary"
                loading={loading.windowRevenue}
                disabled={!winRevenueRange}
                onClick={() => wrap('windowRevenue', async () => {
                  const res = await getWindowRevenue({
                    startTime: winRevenueRange[0].toISOString(),
                    endTime: winRevenueRange[1].toISOString(),
                  })
                  setWindowRevenue(res)
                })}
              >
                查询
              </Button>
            </Space>
            {errors.windowRevenue && <Alert type="error" message={errors.windowRevenue} style={{ marginTop: 8 }} showIcon />}
            {windowRevenue && (
              <Table
                style={{ marginTop: 8 }}
                size="small"
                rowKey="windowId"
                dataSource={windowRevenue.windows || []}
                columns={windowRevenueColumns}
                pagination={{ pageSize: 10 }}
              />
            )}
          </Card>
        </Col>

        {/* 各持卡人存款明细 */}
        <Col span={24}>
          <Card title="各持卡人存款明细" size="small">
            <Space wrap>
              <RangePicker
                showTime
                value={depositRange}
                onChange={setDepositRange}
              />
              <Button
                type="primary"
                loading={loading.depositDetails}
                onClick={() => {
                  setDepositPage(1)
                  wrap('depositDetails', async () => {
                    const params = { page: 1, pageSize: depositPageSize }
                    if (depositRange) {
                      params.startTime = depositRange[0].toISOString()
                      params.endTime = depositRange[1].toISOString()
                    }
                    const res = await getDepositDetails(params)
                    setDepositDetails(res)
                  })
                }}
              >
                查询
              </Button>
            </Space>
            {errors.depositDetails && <Alert type="error" message={errors.depositDetails} style={{ marginTop: 8 }} showIcon />}
            {depositDetails && (
              <>
                {(depositDetails.holders || []).map(h => (
                  <Card
                    key={h.holderId}
                    size="small"
                    style={{ marginTop: 8 }}
                    title={`${h.holderName}（${h.idNumber}）— 合计：${formatYuan(h.totalAmount)}`}
                  >
                    <Table
                      size="small"
                      rowKey="id"
                      dataSource={h.deposits || []}
                      columns={depositDetailColumns}
                      pagination={{ pageSize: 10 }}
                    />
                  </Card>
                ))}
                {depositDetails.total > depositPageSize && (
                  <div style={{ marginTop: 12, textAlign: 'right' }}>
                    <Button.Group>
                      {Array.from({ length: Math.ceil(depositDetails.total / depositPageSize) }, (_, i) => i + 1).map(p => (
                        <Button
                          key={p}
                          size="small"
                          type={p === depositPage ? 'primary' : 'default'}
                          loading={loading.depositDetails && p === depositPage}
                          onClick={() => {
                            setDepositPage(p)
                            wrap('depositDetails', async () => {
                              const params = { page: p, pageSize: depositPageSize }
                              if (depositRange) {
                                params.startTime = depositRange[0].toISOString()
                                params.endTime = depositRange[1].toISOString()
                              }
                              const res = await getDepositDetails(params)
                              setDepositDetails(res)
                            })
                          }}
                        >
                          {p}
                        </Button>
                      ))}
                    </Button.Group>
                    <span style={{ marginLeft: 8, fontSize: 13, color: '#666' }}>
                      共 {depositDetails.total} 位持卡人
                    </span>
                  </div>
                )}
              </>
            )}
          </Card>
        </Col>

        {/* 本日/本月存款 */}
        <Col xs={24} md={12}>
          <Card title="本日 / 本月存款金额" size="small">
            <Button
              type="primary"
              loading={loading.depositSummary}
              onClick={() => wrap('depositSummary', async () => {
                const res = await getDepositSummary()
                setDepositSummary(res)
              })}
            >
              查询
            </Button>
            {errors.depositSummary && <Alert type="error" message={errors.depositSummary} style={{ marginTop: 8 }} showIcon />}
            {depositSummary && (
              <Descriptions column={1} size="small" style={{ marginTop: 8 }}>
                <Descriptions.Item label="今日存款">{formatYuan(depositSummary.todayTotal)}</Descriptions.Item>
                <Descriptions.Item label="本月存款">{formatYuan(depositSummary.monthTotal)}</Descriptions.Item>
              </Descriptions>
            )}
          </Card>
        </Col>

        {/* 流动资金总额 */}
        <Col xs={24} md={12}>
          <Card title="卡中流动资金总额" size="small">
            <Button
              type="primary"
              loading={loading.activeBalance}
              onClick={() => wrap('activeBalance', async () => {
                const res = await getActiveBalance()
                setActiveBalance(res)
              })}
            >
              查询
            </Button>
            {errors.activeBalance && <Alert type="error" message={errors.activeBalance} style={{ marginTop: 8 }} showIcon />}
            {activeBalance && (
              <p style={{ marginTop: 8, fontSize: 16 }}>
                流动资金总额：<strong>{formatYuan(activeBalance.totalBalance)}</strong>
              </p>
            )}
          </Card>
        </Col>

        {/* 日餐报表 */}
        <Col span={24}>
          <Card title="日餐报表" size="small">
            <Space wrap>
              <DatePicker value={dailyDate} onChange={setDailyDate} />
              <Button
                type="primary"
                loading={loading.dailyReport}
                disabled={!dailyDate}
                onClick={() => wrap('dailyReport', async () => {
                  const res = await getDailyReport(dailyDate.format('YYYY-MM-DD'))
                  setDailyReport(res)
                })}
              >
                查询
              </Button>
            </Space>
            {errors.dailyReport && <Alert type="error" message={errors.dailyReport} style={{ marginTop: 8 }} showIcon />}
            {dailyReport && (
              <>
                <Descriptions column={2} size="small" style={{ marginTop: 8, marginBottom: 8 }}>
                  <Descriptions.Item label="日期">{dailyReport.date}</Descriptions.Item>
                  <Descriptions.Item label="消费笔数">{dailyReport.transactionCount}</Descriptions.Item>
                  <Descriptions.Item label="总收入">{formatYuan(dailyReport.totalRevenue)}</Descriptions.Item>
                </Descriptions>
                <Table
                  size="small"
                  rowKey="windowId"
                  dataSource={dailyReport.windows || []}
                  columns={dailyReportColumns}
                  pagination={{ pageSize: 10 }}
                />
              </>
            )}
          </Card>
        </Col>

        {/* 年餐报表 */}
        <Col span={24}>
          <Card title="年餐报表" size="small">
            <Space wrap>
              <InputNumber
                min={2000}
                max={2099}
                value={yearlyYear}
                onChange={setYearlyYear}
                style={{ width: 100 }}
              />
              <Button
                type="primary"
                loading={loading.yearlyReport}
                disabled={!yearlyYear}
                onClick={() => wrap('yearlyReport', async () => {
                  const res = await getYearlyReport(parseInt(yearlyYear))
                  setYearlyReport(res)
                })}
              >
                查询
              </Button>
            </Space>
            {errors.yearlyReport && <Alert type="error" message={errors.yearlyReport} style={{ marginTop: 8 }} showIcon />}
            {yearlyReport && (
              <>
                <Descriptions column={2} size="small" style={{ marginTop: 8, marginBottom: 8 }}>
                  <Descriptions.Item label="年份">{yearlyReport.year} 年</Descriptions.Item>
                  <Descriptions.Item label="消费笔数">{yearlyReport.transactionCount}</Descriptions.Item>
                  <Descriptions.Item label="总收入">{formatYuan(yearlyReport.totalRevenue)}</Descriptions.Item>
                </Descriptions>
                <Table
                  size="small"
                  rowKey="month"
                  dataSource={yearlyReport.months || []}
                  columns={yearlyReportColumns}
                  pagination={{ pageSize: 10 }}
                />
              </>
            )}
          </Card>
        </Col>
      </Row>
    </div>
  )
}
