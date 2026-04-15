import { useState, useEffect } from 'react'
import { listWindows, createWindow } from '../api.js'

export default function WindowsPage() {
  const [windows, setWindows] = useState([])
  const [name, setName] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [creating, setCreating] = useState(false)

  async function fetchWindows() {
    setLoading(true)
    try {
      const res = await listWindows()
      setWindows(res.windows || [])
    } catch (err) {
      setError(err.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchWindows()
  }, [])

  async function handleCreate(e) {
    e.preventDefault()
    setError('')
    setCreating(true)
    try {
      await createWindow(name)
      setName('')
      await fetchWindows()
    } catch (err) {
      setError(err.message || '创建失败')
    } finally {
      setCreating(false)
    }
  }

  return (
    <div style={{ maxWidth: 480, margin: '0 auto', padding: 24 }}>
      <h2>窗口管理</h2>

      <h3>窗口列表</h3>
      {loading ? (
        <p>加载中...</p>
      ) : windows.length === 0 ? (
        <p style={{ color: '#999' }}>暂无窗口</p>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: 16 }}>
          <thead>
            <tr>
              <th style={{ textAlign: 'left', padding: '6px 8px', borderBottom: '2px solid #ddd' }}>ID</th>
              <th style={{ textAlign: 'left', padding: '6px 8px', borderBottom: '2px solid #ddd' }}>窗口名称</th>
            </tr>
          </thead>
          <tbody>
            {windows.map(w => (
              <tr key={w.id}>
                <td style={{ padding: '6px 8px', borderBottom: '1px solid #eee' }}>{w.id}</td>
                <td style={{ padding: '6px 8px', borderBottom: '1px solid #eee' }}>{w.name}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      <h3>新建窗口</h3>
      <form onSubmit={handleCreate}>
        <div style={{ display: 'flex', gap: 8 }}>
          <input
            placeholder="窗口名称"
            value={name}
            onChange={e => setName(e.target.value)}
            required
            style={{ flex: 1, padding: 8 }}
          />
          <button type="submit" disabled={creating} style={{ padding: '8px 16px' }}>
            {creating ? '创建中...' : '创建'}
          </button>
        </div>
      </form>

      {error && (
        <div style={{ marginTop: 12, padding: 12, background: '#ffe0e0', borderRadius: 4, color: '#c00' }}>
          {error}
        </div>
      )}
    </div>
  )
}
