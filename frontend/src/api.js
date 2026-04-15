const BASE_URL = 'http://localhost:8080'

async function request(path, options = {}) {
  const url = `${BASE_URL}${path}`
  const res = await fetch(url, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  const data = await res.json()
  if (!res.ok) {
    const err = new Error(data.message || '请求失败')
    err.code = data.code
    err.status = res.status
    throw err
  }
  return data
}

// ==================== 卡管理 ====================

/** 发卡 POST /api/cards */
export function issueCard({ name, idNumber, deposit, preDeposit = 0 }) {
  return request('/api/cards', {
    method: 'POST',
    body: JSON.stringify({ name, idNumber, deposit, preDeposit }),
  })
}

/** 查询卡信息 GET /api/cards/{id} */
export function getCard(id) {
  return request(`/api/cards/${id}`)
}

// ==================== 存款 ====================

/** 存款 POST /api/cards/{id}/deposits */
export function deposit(id, amount) {
  return request(`/api/cards/${id}/deposits`, {
    method: 'POST',
    body: JSON.stringify({ amount }),
  })
}

// ==================== 就餐消费 ====================

/** 就餐消费 POST /api/cards/{id}/transactions */
export function createTransaction(id, { windowId, amount }) {
  return request(`/api/cards/${id}/transactions`, {
    method: 'POST',
    body: JSON.stringify({ windowId, amount }),
  })
}

// ==================== 挂失 ====================

/** 挂失 PUT /api/cards/{id}/loss-report */
export function reportLoss(id) {
  return request(`/api/cards/${id}/loss-report`, { method: 'PUT' })
}

/** 取消挂失 DELETE /api/cards/{id}/loss-report */
export function cancelLossReport(id) {
  return request(`/api/cards/${id}/loss-report`, { method: 'DELETE' })
}

// ==================== 注销 ====================

/** 注销 POST /api/cards/{id}/cancellation */
export function cancelCard(id) {
  return request(`/api/cards/${id}/cancellation`, { method: 'POST' })
}

// ==================== 汇总统计 ====================

/** 本餐售饭总收入 GET /api/statistics/meal-revenue */
export function getMealRevenue({ startTime, endTime }) {
  return request(`/api/statistics/meal-revenue?startTime=${encodeURIComponent(startTime)}&endTime=${encodeURIComponent(endTime)}`)
}

/** 各窗口收入 GET /api/statistics/window-revenue */
export function getWindowRevenue({ startTime, endTime }) {
  return request(`/api/statistics/window-revenue?startTime=${encodeURIComponent(startTime)}&endTime=${encodeURIComponent(endTime)}`)
}

/** 各持卡人存款明细 GET /api/statistics/deposit-details */
export function getDepositDetails({ startTime, endTime } = {}) {
  const params = new URLSearchParams()
  if (startTime) params.set('startTime', startTime)
  if (endTime) params.set('endTime', endTime)
  const qs = params.toString()
  return request(`/api/statistics/deposit-details${qs ? '?' + qs : ''}`)
}

/** 本日/本月存款金额 GET /api/statistics/deposit-summary */
export function getDepositSummary() {
  return request('/api/statistics/deposit-summary')
}

/** 卡中流动资金总额 GET /api/statistics/active-balance */
export function getActiveBalance() {
  return request('/api/statistics/active-balance')
}

/** 日餐报表 GET /api/statistics/daily-report */
export function getDailyReport(date) {
  return request(`/api/statistics/daily-report?date=${encodeURIComponent(date)}`)
}

/** 年餐报表 GET /api/statistics/yearly-report */
export function getYearlyReport(year) {
  return request(`/api/statistics/yearly-report?year=${encodeURIComponent(year)}`)
}

// ==================== 窗口管理 ====================

/** 获取窗口列表 GET /api/windows */
export function listWindows() {
  return request('/api/windows')
}

/** 创建窗口 POST /api/windows */
export function createWindow(name) {
  return request('/api/windows', {
    method: 'POST',
    body: JSON.stringify({ name }),
  })
}
