const BASE_URL = ''

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

// ==================== 学籍验证 ====================

/** 校验证件号 GET /api/validate-student?idNumber=xxx */
export function validateStudent(idNumber) {
  return request(`/api/validate-student?idNumber=${encodeURIComponent(idNumber)}`)
}

// ==================== 卡管理 ====================

/** 发卡 POST /api/cards，请求体只含 idNumber + preDeposit */
export function issueCard({ idNumber, preDeposit = 0 }) {
  return request('/api/cards', {
    method: 'POST',
    body: JSON.stringify({ idNumber, preDeposit }),
  })
}

/** 按证件号查询当前有效卡 GET /api/cards?idNumber=xxx */
export function getCardByIDNumber(idNumber) {
  return request(`/api/cards?idNumber=${encodeURIComponent(idNumber)}`)
}

/** 按卡号查询卡信息 GET /api/cards/{cardNo} */
export function getCard(cardNo) {
  return request(`/api/cards/${cardNo}`)
}

// ==================== 存款 ====================

/** 存款 POST /api/cards/{cardNo}/deposits */
export function deposit(cardNo, amount) {
  return request(`/api/cards/${cardNo}/deposits`, {
    method: 'POST',
    body: JSON.stringify({ amount }),
  })
}

// ==================== 就餐消费 ====================

/** 就餐消费 POST /api/cards/{cardNo}/transactions */
export function createTransaction(cardNo, { windowId, amount }) {
  return request(`/api/cards/${cardNo}/transactions`, {
    method: 'POST',
    body: JSON.stringify({ windowId, amount }),
  })
}

/** 查询消费历史 GET /api/cards/{cardNo}/transactions */
export function getCardTransactions(cardNo, { page, pageSize } = {}) {
  const params = new URLSearchParams()
  if (page != null) params.set('page', page)
  if (pageSize != null) params.set('pageSize', pageSize)
  const qs = params.toString()
  return request(`/api/cards/${cardNo}/transactions${qs ? '?' + qs : ''}`)
}

/** 查询存款历史 GET /api/cards/{cardNo}/deposits */
export function getCardDeposits(cardNo, { page, pageSize } = {}) {
  const params = new URLSearchParams()
  if (page != null) params.set('page', page)
  if (pageSize != null) params.set('pageSize', pageSize)
  const qs = params.toString()
  return request(`/api/cards/${cardNo}/deposits${qs ? '?' + qs : ''}`)
}

// ==================== 挂失 ====================

/** 挂失 PUT /api/cards/{cardNo}/loss-report */
export function reportLoss(cardNo) {
  return request(`/api/cards/${cardNo}/loss-report`, { method: 'PUT' })
}

/** 取消挂失 DELETE /api/cards/{cardNo}/loss-report */
export function cancelLossReport(cardNo) {
  return request(`/api/cards/${cardNo}/loss-report`, { method: 'DELETE' })
}

// ==================== 注销 ====================

/** 注销 POST /api/cards/{cardNo}/cancellation */
export function cancelCard(cardNo) {
  return request(`/api/cards/${cardNo}/cancellation`, { method: 'POST' })
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
export function getDepositDetails({ startTime, endTime, page, pageSize } = {}) {
  const params = new URLSearchParams()
  if (startTime) params.set('startTime', startTime)
  if (endTime) params.set('endTime', endTime)
  if (page != null) params.set('page', page)
  if (pageSize != null) params.set('pageSize', pageSize)
  const qs = params.toString()
  return request(`/api/statistics/deposit-details${qs ? '?' + qs : ''}`)
}

/** 单个持卡人存款记录（分页） GET /api/statistics/holder-deposits */
export function getHolderDeposits({ holderId, page, pageSize, startTime, endTime } = {}) {
  const params = new URLSearchParams()
  if (holderId != null) params.set('holderId', holderId)
  if (page != null) params.set('page', page)
  if (pageSize != null) params.set('pageSize', pageSize)
  if (startTime) params.set('startTime', startTime)
  if (endTime) params.set('endTime', endTime)
  const qs = params.toString()
  return request(`/api/statistics/holder-deposits${qs ? '?' + qs : ''}`)
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
