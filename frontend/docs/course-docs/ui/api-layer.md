# api.js 封装

## 作用
- 统一封装前端所有 HTTP 请求，屏蔽 fetch 细节，对外暴露语义化函数

## 核心实现

- BASE_URL：`http://localhost:8080`
- 内部 `request(path, options)` 函数统一处理：
  - 默认设置 `Content-Type: application/json`
  - 统一解析响应 JSON
  - 非 2xx 状态码：构造 Error 对象，附带 `err.code`、`err.status`，throw 出去
  - 2xx：直接返回 data

## 导出函数列表

### 卡管理
| 函数 | 方法 | 路径 |
|------|------|------|
| issueCard({ name, idNumber, deposit, preDeposit }) | POST | /api/cards |
| getCard(id) | GET | /api/cards/{id} |

### 存款
| 函数 | 方法 | 路径 |
|------|------|------|
| deposit(id, amount) | POST | /api/cards/{id}/deposits |

### 就餐消费
| 函数 | 方法 | 路径 |
|------|------|------|
| createTransaction(id, { windowId, amount }) | POST | /api/cards/{id}/transactions |

### 挂失
| 函数 | 方法 | 路径 |
|------|------|------|
| reportLoss(id) | PUT | /api/cards/{id}/loss-report |
| cancelLossReport(id) | DELETE | /api/cards/{id}/loss-report |

### 注销
| 函数 | 方法 | 路径 |
|------|------|------|
| cancelCard(id) | POST | /api/cards/{id}/cancellation |

### 汇总统计
| 函数 | 方法 | 路径 |
|------|------|------|
| getMealRevenue({ startTime, endTime }) | GET | /api/statistics/meal-revenue |
| getWindowRevenue({ startTime, endTime }) | GET | /api/statistics/window-revenue |
| getDepositDetails({ startTime?, endTime? }) | GET | /api/statistics/deposit-details |
| getDepositSummary() | GET | /api/statistics/deposit-summary |
| getActiveBalance() | GET | /api/statistics/active-balance |
| getDailyReport(date) | GET | /api/statistics/daily-report |
| getYearlyReport(year) | GET | /api/statistics/yearly-report |

### 窗口管理
| 函数 | 方法 | 路径 |
|------|------|------|
| listWindows() | GET | /api/windows |
| createWindow(name) | POST | /api/windows |

## 金额单位约定
- 所有金额接口传参和响应均以**分（整数）**为单位
- 页面显示时 ÷ 100 转元：`(fen / 100).toFixed(2)`
- 页面提交时 × 100 转分：`Math.round(yuan * 100)`
- 此约定在所有页面中统一执行，api.js 本身不做转换

## 错误处理
- HTTP 非 2xx → throw Error，页面 catch 后通过 `err.message` 展示 Alert
- 特殊场景：MealPage 中 catch 块判断 `err.status === 404` 识别「非本单位卡」警报
