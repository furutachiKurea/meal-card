import { useState } from 'react'
import {
  getMealRevenue,
  getWindowRevenue,
  getDepositDetails,
  getDepositSummary,
  getActiveBalance,
  getDailyReport,
  getYearlyReport,
} from '../api.js'

function formatYuan(fen) {
  return (fen / 100).toFixed(2) + ' 元'
}

function today() {
  return new Date().toISOString().slice(0, 10)
}

function thisYear() {
  return new Date().getFullYear()
}

export default function StatisticsPage() {
  // 本餐售饭总收入
  const [revenueStart, setRevenueStart] = useState('')
  const [revenueEnd, setRevenueEnd] = useState('')
  const [mealRevenue, setMealRevenue] = useState(null)

  // 各窗口收入
  const [winRevenueStart, setWinRevenueStart] = useState('')
  const [winRevenueEnd, setWinRevenueEnd] = useState('')
  const [windowRevenue, setWindowRevenue] = useState(null)

  // 存款明细
  const [depositStart, setDepositStart] = useState('')
  const [depositEnd, setDepositEnd] = useState('')
  const [depositDetails, setDepositDetails] = useState(null)

  // 本日/本月存款
  const [depositSummary, setDepositSummary] = useState(null)

  // 流动资金
  const [activeBalance, setActiveBalance] = useState(null)

  // 日餐报表
  const [dailyDate, setDailyDate] = useState(today())
  const [dailyReport, setDailyReport] = useState(null)

  // 年餐报表
  const [yearlyYear, setYearlyYear] = useState(String(thisYear()))
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

  return (
    <div style={{ maxWidth: 800, margin: '0 auto', padding: 24 }}>
      <h2>汇总统计</h2>

      {/* 本餐售饭总收入 */}
      <section style={{ marginBottom: 24, padding: 16, border: '1px solid #ddd', borderRadius: 4 }}>
        <h3>本餐售饭总收入</h3>
        <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center' }}>
          <input type="datetime-local" value={revenueStart} onChange={e => setRevenueStart(e.target.value)} style={{ padding: 6 }} />
          <span>至</span>
          <input type="datetime-local" value={revenueEnd} onChange={e => setRevenueEnd(e.target.value)} style={{ padding: 6 }} />
          <button
            onClick={() => wrap('mealRevenue', async () => {
              const res = await getMealRevenue({ startTime: new Date(revenueStart).toISOString(), endTime: new Date(revenueEnd).toISOString() })
              setMealRevenue(res)
            })}
            disabled={loading.mealRevenue || !revenueStart || !revenueEnd}
            style={{ padding: '6px 16px' }}
          >
            查询
          </button>
        </div>
        {errors.mealRevenue && <p style={{ color: '#c00' }}>{errors.mealRevenue}</p>}
        {mealRevenue && (
          <p style={{ marginTop: 8, fontSize: 18 }}>
            总收入：<strong>{formatYuan(mealRevenue.totalRevenue)}</strong>
          </p>
        )}
      </section>

      {/* 各窗口收入 */}
      <section style={{ marginBottom: 24, padding: 16, border: '1px solid #ddd', borderRadius: 4 }}>
        <h3>各窗口收入</h3>
        <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center' }}>
          <input type="datetime-local" value={winRevenueStart} onChange={e => setWinRevenueStart(e.target.value)} style={{ padding: 6 }} />
          <span>至</span>
          <input type="datetime-local" value={winRevenueEnd} onChange={e => setWinRevenueEnd(e.target.value)} style={{ padding: 6 }} />
          <button
            onClick={() => wrap('windowRevenue', async () => {
              const res = await getWindowRevenue({ startTime: new Date(winRevenueStart).toISOString(), endTime: new Date(winRevenueEnd).toISOString() })
              setWindowRevenue(res)
            })}
            disabled={loading.windowRevenue || !winRevenueStart || !winRevenueEnd}
            style={{ padding: '6px 16px' }}
          >
            查询
          </button>
        </div>
        {errors.windowRevenue && <p style={{ color: '#c00' }}>{errors.windowRevenue}</p>}
        {windowRevenue && (
          <table style={{ marginTop: 8, width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '4px 8px', borderBottom: '1px solid #ddd' }}>窗口</th>
                <th style={{ textAlign: 'right', padding: '4px 8px', borderBottom: '1px solid #ddd' }}>收入</th>
              </tr>
            </thead>
            <tbody>
              {(windowRevenue.windows || []).map(w => (
                <tr key={w.windowId}>
                  <td style={{ padding: '4px 8px' }}>{w.windowName}</td>
                  <td style={{ textAlign: 'right', padding: '4px 8px' }}>{formatYuan(w.revenue)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

      {/* 各持卡人存款明细 */}
      <section style={{ marginBottom: 24, padding: 16, border: '1px solid #ddd', borderRadius: 4 }}>
        <h3>各持卡人存款明细</h3>
        <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center' }}>
          <input type="datetime-local" value={depositStart} onChange={e => setDepositStart(e.target.value)} style={{ padding: 6 }} />
          <span>至</span>
          <input type="datetime-local" value={depositEnd} onChange={e => setDepositEnd(e.target.value)} style={{ padding: 6 }} />
          <button
            onClick={() => wrap('depositDetails', async () => {
              const params = {}
              if (depositStart) params.startTime = new Date(depositStart).toISOString()
              if (depositEnd) params.endTime = new Date(depositEnd).toISOString()
              const res = await getDepositDetails(params)
              setDepositDetails(res)
            })}
            disabled={loading.depositDetails}
            style={{ padding: '6px 16px' }}
          >
            查询
          </button>
        </div>
        {errors.depositDetails && <p style={{ color: '#c00' }}>{errors.depositDetails}</p>}
        {depositDetails && (
          <div style={{ marginTop: 8 }}>
            {(depositDetails.holders || []).map(h => (
              <div key={h.holderId} style={{ marginBottom: 12, padding: 8, background: '#f5f5f5', borderRadius: 4 }}>
                <strong>{h.holderName}</strong>（{h.idNumber}）— 合计：{formatYuan(h.totalAmount)}
                <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: 4 }}>
                  <thead>
                    <tr>
                      <th style={{ textAlign: 'left', padding: '2px 6px', fontSize: 12, color: '#666' }}>卡号</th>
                      <th style={{ textAlign: 'right', padding: '2px 6px', fontSize: 12, color: '#666' }}>金额</th>
                      <th style={{ textAlign: 'right', padding: '2px 6px', fontSize: 12, color: '#666' }}>时间</th>
                    </tr>
                  </thead>
                  <tbody>
                    {(h.deposits || []).map(d => (
                      <tr key={d.id}>
                        <td style={{ padding: '2px 6px', fontSize: 13 }}>{d.cardId}</td>
                        <td style={{ textAlign: 'right', padding: '2px 6px', fontSize: 13 }}>{formatYuan(d.amount)}</td>
                        <td style={{ textAlign: 'right', padding: '2px 6px', fontSize: 13 }}>{new Date(d.createdAt).toLocaleString('zh-CN')}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ))}
          </div>
        )}
      </section>

      {/* 本日/本月存款 */}
      <section style={{ marginBottom: 24, padding: 16, border: '1px solid #ddd', borderRadius: 4 }}>
        <h3>本日 / 本月存款金额</h3>
        <button
          onClick={() => wrap('depositSummary', async () => {
            const res = await getDepositSummary()
            setDepositSummary(res)
          })}
          disabled={loading.depositSummary}
          style={{ padding: '6px 16px' }}
        >
          {loading.depositSummary ? '加载中...' : '查询'}
        </button>
        {errors.depositSummary && <p style={{ color: '#c00' }}>{errors.depositSummary}</p>}
        {depositSummary && (
          <div style={{ marginTop: 8, display: 'flex', gap: 24 }}>
            <p>今日存款：<strong>{formatYuan(depositSummary.todayTotal)}</strong></p>
            <p>本月存款：<strong>{formatYuan(depositSummary.monthTotal)}</strong></p>
          </div>
        )}
      </section>

      {/* 流动资金总额 */}
      <section style={{ marginBottom: 24, padding: 16, border: '1px solid #ddd', borderRadius: 4 }}>
        <h3>卡中流动资金总额</h3>
        <button
          onClick={() => wrap('activeBalance', async () => {
            const res = await getActiveBalance()
            setActiveBalance(res)
          })}
          disabled={loading.activeBalance}
          style={{ padding: '6px 16px' }}
        >
          {loading.activeBalance ? '加载中...' : '查询'}
        </button>
        {errors.activeBalance && <p style={{ color: '#c00' }}>{errors.activeBalance}</p>}
        {activeBalance && (
          <p style={{ marginTop: 8, fontSize: 18 }}>
            流动资金总额：<strong>{formatYuan(activeBalance.totalBalance)}</strong>
          </p>
        )}
      </section>

      {/* 日餐报表 */}
      <section style={{ marginBottom: 24, padding: 16, border: '1px solid #ddd', borderRadius: 4 }}>
        <h3>日餐报表</h3>
        <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
          <input type="date" value={dailyDate} onChange={e => setDailyDate(e.target.value)} style={{ padding: 6 }} />
          <button
            onClick={() => wrap('dailyReport', async () => {
              const res = await getDailyReport(dailyDate)
              setDailyReport(res)
            })}
            disabled={loading.dailyReport || !dailyDate}
            style={{ padding: '6px 16px' }}
          >
            查询
          </button>
        </div>
        {errors.dailyReport && <p style={{ color: '#c00' }}>{errors.dailyReport}</p>}
        {dailyReport && (
          <div style={{ marginTop: 8 }}>
            <p>日期：{dailyReport.date}</p>
            <p>总收入：<strong>{formatYuan(dailyReport.totalRevenue)}</strong> | 消费笔数：{dailyReport.transactionCount}</p>
            <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: 4 }}>
              <thead>
                <tr>
                  <th style={{ textAlign: 'left', padding: '4px 8px', borderBottom: '1px solid #ddd' }}>窗口</th>
                  <th style={{ textAlign: 'right', padding: '4px 8px', borderBottom: '1px solid #ddd' }}>收入</th>
                  <th style={{ textAlign: 'right', padding: '4px 8px', borderBottom: '1px solid #ddd' }}>笔数</th>
                </tr>
              </thead>
              <tbody>
                {(dailyReport.windows || []).map(w => (
                  <tr key={w.windowId}>
                    <td style={{ padding: '4px 8px' }}>{w.windowName}</td>
                    <td style={{ textAlign: 'right', padding: '4px 8px' }}>{formatYuan(w.revenue)}</td>
                    <td style={{ textAlign: 'right', padding: '4px 8px' }}>{w.transactionCount}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* 年餐报表 */}
      <section style={{ marginBottom: 24, padding: 16, border: '1px solid #ddd', borderRadius: 4 }}>
        <h3>年餐报表</h3>
        <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
          <input
            type="number"
            value={yearlyYear}
            onChange={e => setYearlyYear(e.target.value)}
            min="2000"
            max="2099"
            style={{ padding: 6, width: 100 }}
          />
          <button
            onClick={() => wrap('yearlyReport', async () => {
              const res = await getYearlyReport(parseInt(yearlyYear))
              setYearlyReport(res)
            })}
            disabled={loading.yearlyReport || !yearlyYear}
            style={{ padding: '6px 16px' }}
          >
            查询
          </button>
        </div>
        {errors.yearlyReport && <p style={{ color: '#c00' }}>{errors.yearlyReport}</p>}
        {yearlyReport && (
          <div style={{ marginTop: 8 }}>
            <p>{yearlyReport.year} 年 | 总收入：<strong>{formatYuan(yearlyReport.totalRevenue)}</strong> | 消费笔数：{yearlyReport.transactionCount}</p>
            <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: 4 }}>
              <thead>
                <tr>
                  <th style={{ textAlign: 'left', padding: '4px 8px', borderBottom: '1px solid #ddd' }}>月份</th>
                  <th style={{ textAlign: 'right', padding: '4px 8px', borderBottom: '1px solid #ddd' }}>收入</th>
                  <th style={{ textAlign: 'right', padding: '4px 8px', borderBottom: '1px solid #ddd' }}>笔数</th>
                </tr>
              </thead>
              <tbody>
                {(yearlyReport.months || []).map(m => (
                  <tr key={m.month}>
                    <td style={{ padding: '4px 8px' }}>{m.month} 月</td>
                    <td style={{ textAlign: 'right', padding: '4px 8px' }}>{formatYuan(m.revenue)}</td>
                    <td style={{ textAlign: 'right', padding: '4px 8px' }}>{m.transactionCount}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  )
}
