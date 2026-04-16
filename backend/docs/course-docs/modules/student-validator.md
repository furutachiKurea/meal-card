# StudentValidator 接口与学籍验证模块

## 作用
- 定义学籍验证的抽象接口，供 CardService 调用，屏蔽底层实现
- 发卡前校验证件号是否为本校学生/教职工，返回姓名和人员类型

## 职责边界
- 负责：定义接口契约、提供 HTTP 实现（调用 Mock 服务）
- 不负责：证件号格式校验（由上层做）、HTTP 路由/响应、业务状态流转

## 接口定义（backend/service/student_validator.go）

```go
type StudentInfo struct {
    IDNumber string
    Name     string
    Type     string // "student" | "staff"
}

type StudentValidator interface {
    Validate(idNumber string) (*StudentInfo, error)
}
```

## 实现一：HttpStudentValidator（backend/client/student_client.go）

### 输入
- idNumber：12 位证件号字符串

### 输出
- *StudentInfo（valid=true 时）
- error（STUDENT_NOT_FOUND / STUDENT_SERVICE_ERROR）

### 核心流程
1. 读取环境变量 STUDENT_SERVICE_URL（默认 http://localhost:9090）
2. GET {url}/validate?idNumber={idNumber}
3. HTTP 404 → 返回 STUDENT_NOT_FOUND 错误
4. 网络错误或非 2xx → 返回 STUDENT_SERVICE_ERROR 错误
5. 解析响应 JSON，返回 StudentInfo

### 关键实现点
- URL 通过环境变量注入，便于测试和切换环境
- 不缓存验证结果，每次调用都发请求（课设规模无需缓存）

## 实现二：FakeStudentValidator（测试用，内联在 card_service_test.go）

- 硬编码若干合法证件号及对应姓名
- 指定证件号返回 StudentInfo，未知证件号返回 STUDENT_NOT_FOUND
- 不依赖网络，单元测试直接注入

---

# mock-services/student-service/ 模块

## 作用
- 独立 Go HTTP 服务，模拟学校学籍系统的验证接口
- 供联调和演示使用，数据硬编码在 main.go

## 端口
`:9090`

## 接口
`GET /validate?idNumber=xxx`

### 响应（合法证件号）
```json
{ "valid": true, "idNumber": "202305133513", "name": "张三", "type": "student" }
```

### 响应（不存在）
```json
{ "valid": false }
```

## 内置学生/教职工名单（硬编码，按需修改 main.go）

| 证件号 | 姓名 | 类型 |
|---|---|---|
| 202305133513 | 张三 | student |
| 202305133514 | 李四 | student |
| 202305133515 | 王五 | student |
| 202201000001 | 陈老师 | staff |
| 202201000002 | 刘老师 | staff |
