import { useState } from 'react'
import { Form, Input, InputNumber, Button, Alert, Descriptions, Card, Typography, Steps, Tag, Space } from 'antd'
import { validateStudent, issueCard } from '../api.js'

const { Title } = Typography

const TYPE_LABEL = { student: '学生', staff: '教职工' }
const TYPE_COLOR = { student: 'blue', staff: 'purple' }

export default function IssuePage() {
  // 当前步骤：0 = 验证身份，1 = 录入预存款
  const [step, setStep] = useState(0)
  // 验证通过后的学籍信息
  const [studentInfo, setStudentInfo] = useState(null)
  // 发卡成功结果
  const [result, setResult] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const [idNumber, setIdNumber] = useState('')
  const [depositForm] = Form.useForm()

  // 第一步：验证证件号
  async function handleValidate() {
    if (!idNumber.trim()) return
    setError('')
    setStudentInfo(null)
    setLoading(true)
    try {
      const res = await validateStudent(idNumber.trim())
      setStudentInfo(res)
    } catch (err) {
      if (err.code === 'STUDENT_NOT_FOUND') {
        setError('证件号不存在，非本校学生/教职工')
      } else {
        setError(err.message || '验证失败')
      }
    } finally {
      setLoading(false)
    }
  }

  // 确认身份，进入第二步
  function handleConfirm() {
    setError('')
    setStep(1)
  }

  // 第二步：录入预存款并发卡
  async function handleIssue(values) {
    setError('')
    setLoading(true)
    try {
      const preDepositFen = values.preDeposit ? Math.round(values.preDeposit * 100) : 0
      const res = await issueCard({
        idNumber: studentInfo.idNumber,
        preDeposit: preDepositFen,
      })
      setResult(res)
      depositForm.resetFields()
    } catch (err) {
      setError(err.message || '发卡失败')
    } finally {
      setLoading(false)
    }
  }

  // 重置整个流程
  function handleReset() {
    setStep(0)
    setStudentInfo(null)
    setResult(null)
    setError('')
    setIdNumber('')
    depositForm.resetFields()
  }

  return (
    <Card style={{ maxWidth: 520, margin: '0 auto' }}>
      <Title level={4} style={{ marginTop: 0 }}>发卡</Title>

      <Steps
        current={step}
        size="small"
        style={{ marginBottom: 24 }}
        items={[
          { title: '验证身份' },
          { title: '录入预存款' },
        ]}
      />

      {/* 步骤 0：验证证件号 */}
      {step === 0 && !result && (
        <>
          <Space.Compact style={{ width: '100%', marginBottom: 16 }}>
            <Input
              placeholder="请输入12位证件号"
              value={idNumber}
              onChange={e => { setIdNumber(e.target.value); setStudentInfo(null); setError('') }}
              onPressEnter={handleValidate}
            />
            <Button type="primary" loading={loading} onClick={handleValidate}>
              查询
            </Button>
          </Space.Compact>

          {error && (
            <Alert type="error" message={error} style={{ marginBottom: 12 }} showIcon />
          )}

          {studentInfo && (
            <Card size="small" style={{ marginBottom: 16, background: '#f0f9ff', borderColor: '#91d5ff' }}>
              <Descriptions column={1} size="small">
                <Descriptions.Item label="姓名">{studentInfo.name}</Descriptions.Item>
                <Descriptions.Item label="证件号">{studentInfo.idNumber}</Descriptions.Item>
                <Descriptions.Item label="人员类型">
                  <Tag color={TYPE_COLOR[studentInfo.type]}>
                    {TYPE_LABEL[studentInfo.type] || studentInfo.type}
                  </Tag>
                </Descriptions.Item>
              </Descriptions>
              <Button type="primary" style={{ marginTop: 8 }} onClick={handleConfirm}>
                确认，进行发卡
              </Button>
            </Card>
          )}
        </>
      )}

      {/* 步骤 1：录入预存款 */}
      {step === 1 && !result && (
        <>
          <Card size="small" style={{ marginBottom: 16, background: '#f0f9ff', borderColor: '#91d5ff' }}>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="姓名">{studentInfo?.name}</Descriptions.Item>
              <Descriptions.Item label="证件号">{studentInfo?.idNumber}</Descriptions.Item>
              <Descriptions.Item label="人员类型">
                <Tag color={TYPE_COLOR[studentInfo?.type]}>
                  {TYPE_LABEL[studentInfo?.type] || studentInfo?.type}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="押金">20.00 元 </Descriptions.Item>
            </Descriptions>
          </Card>

          <Form form={depositForm} layout="vertical" onFinish={handleIssue}>
            <Form.Item label="预存款金额（元）" name="preDeposit">
              <InputNumber
                min={0}
                step={0.01}
                precision={2}
                placeholder="0.00（可填0）"
                style={{ width: '100%' }}
              />
            </Form.Item>

            <Form.Item style={{ marginBottom: 0 }}>
              <Space>
                <Button type="primary" htmlType="submit" loading={loading}>
                  发卡
                </Button>
                <Button onClick={handleReset}>
                  取消
                </Button>
              </Space>
            </Form.Item>
          </Form>

          {error && (
            <Alert type="error" message={error} style={{ marginTop: 12 }} showIcon />
          )}
        </>
      )}

      {/* 发卡成功结果 */}
      {result && (
        <Card
          size="small"
          title="发卡成功"
          style={{ marginTop: 16, background: '#f6ffed', borderColor: '#b7eb8f' }}
        >
          <Descriptions column={1} size="small">
            <Descriptions.Item label="卡号">{result.card.cardNo}</Descriptions.Item>
            <Descriptions.Item label="持卡人">{result.cardHolder.name}</Descriptions.Item>
            <Descriptions.Item label="证件号">{result.cardHolder.idNumber}</Descriptions.Item>
            <Descriptions.Item label="押金">{(result.card.deposit / 100).toFixed(2)} 元</Descriptions.Item>
            <Descriptions.Item label="余额">{(result.card.balance / 100).toFixed(2)} 元</Descriptions.Item>
          </Descriptions>

          {result.refund && (
            <Card
              size="small"
              title={`旧卡自动注销（卡号 ${result.refund.oldCardNo}）`}
              style={{ marginTop: 8, background: '#fffbe6', borderColor: '#ffe58f' }}
            >
              <Descriptions column={1} size="small">
                <Descriptions.Item label="退还押金">{(result.refund.deposit / 100).toFixed(2)} 元</Descriptions.Item>
                <Descriptions.Item label="退还余额">{(result.refund.balance / 100).toFixed(2)} 元</Descriptions.Item>
                <Descriptions.Item label="退还合计">{(result.refund.total / 100).toFixed(2)} 元</Descriptions.Item>
              </Descriptions>
            </Card>
          )}

          <Button style={{ marginTop: 12 }} onClick={handleReset}>
            继续办理
          </Button>
        </Card>
      )}
    </Card>
  )
}
