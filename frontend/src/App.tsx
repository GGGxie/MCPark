import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useEffect, useState } from "react"

interface MCPService {
  id: number
  name: string
  display_name: string
  description: string
  config_json: string
  category_id?: number
  category?: {
    id: number
    name: string
    display_name: string
    description: string
    icon: string
  }
  tags: string
  author: string
  repository_url: string
  documentation_url: string
  version: string
  install_count: number
  is_verified: boolean
  is_active: boolean
  created_at: string
  updated_at: string
}



function App() {
  const [currentView, setCurrentView] = useState<'home' | 'add-service'>('home')
  const [services, setServices] = useState<MCPService[]>([])
  const [loading, setLoading] = useState(true)

  // 从API获取服务列表
  const fetchServices = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/v1/services')
      const data = await response.json()
      if (data.code === 0) {
        setServices(data.data || [])
      } else {
        console.error('Failed to fetch services:', data.message)
      }
    } catch (error) {
      console.error('Error fetching services:', error)
    } finally {
      setLoading(false)
    }
  }

  // 初始化时获取数据
  useEffect(() => {
    fetchServices()
  }, [])

  // 处理浏览器后退按钮
  useEffect(() => {
    const handlePopState = () => {
      const hash = window.location.hash.slice(1) // 去掉 # 符号
      if (hash === 'add-service') {
        setCurrentView('add-service')
      } else {
        setCurrentView('home')
      }
    }

    // 监听浏览器后退/前进按钮
    window.addEventListener('popstate', handlePopState)
    
    // 初始化时检查当前 hash
    handlePopState()

    return () => window.removeEventListener('popstate', handlePopState)
  }, [])

  // 导航函数，同时更新状态和浏览器历史
  const navigateTo = (view: 'home' | 'add-service') => {
    setCurrentView(view)
    if (view === 'home') {
      window.history.pushState(null, '', '/')
    } else {
      window.history.pushState(null, '', `/#${view}`)
    }
  }

  // 生成Cursor Deep Link
  const generateCursorLink = (service: MCPService) => {
    try {
      const config = JSON.parse(service.config_json)
      // 使用 base64 编码配置，符合官方文档要求
      const encodedConfig = btoa(JSON.stringify(config))
      // 使用官方文档的正确URL格式
      return `cursor://anysphere.cursor-deeplink/mcp/install?name=${encodeURIComponent(service.name)}&config=${encodedConfig}`
    } catch (error) {
      console.error('Invalid JSON config:', error)
      return null
    }
  }

  // 渲染不同的视图
  const renderContent = () => {
    switch (currentView) {
      case 'add-service':
        return <AddServiceView onAdd={handleAddService} />
      default:
        return <HomeView services={services} onAddService={() => navigateTo('add-service')} onGenerateLink={generateCursorLink} />
    }
  }

  // 处理添加服务
  const handleAddService = async (formData: {
    name: string
    description: string
    category: string
    author: string
    configJson: string
  }) => {
    try {
      const response = await fetch('/api/v1/services', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: formData.name,
          description: formData.description,
          config_json: formData.configJson,
          category: formData.category,
          author: formData.author,
        }),
      })

      const data = await response.json()
      if (data.code === 0) {
        // 创建成功，重新获取服务列表
        await fetchServices()
        navigateTo('home')
        alert('服务添加成功！')
      } else {
        alert('添加失败：' + data.message)
      }
    } catch (error) {
      console.error('Error adding service:', error)
      alert('添加失败：网络错误')
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <div className="text-center mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-4 cursor-pointer hover:text-blue-600 transition-colors" onClick={() => navigateTo('home')}>
              🚀 MCPark
            </h1>
          </div>

          {/* 主导航按钮（仅在首页显示） */}
          {currentView === 'home' && (
            <div className="flex justify-center">
              <Button size="lg" onClick={() => navigateTo('add-service')}>+ 添加 MCP 服务</Button>
            </div>
          )}
        </div>

        {/* 渲染内容 */}
        {loading ? (
          <div className="text-center py-16">
            <div className="text-lg text-gray-600">⏳ 加载中...</div>
          </div>
        ) : (
          renderContent()
        )}
      </div>
    </div>
  )
}

// 首页视图组件
function HomeView({ services, onAddService, onGenerateLink }: { 
  services: MCPService[], 
  onAddService: () => void,
  onGenerateLink: (service: MCPService) => string | null
}) {
  const handleCursorImport = (service: MCPService) => {
    const link = onGenerateLink(service)
    if (link) {
      window.open(link, '_blank')
    } else {
      alert('配置格式错误，无法生成 Cursor 链接')
    }
  }

  const copyJsonConfig = (service: MCPService) => {
    navigator.clipboard.writeText(service.config_json).then(() => {
      alert('JSON 配置已复制到剪贴板！')
    }).catch(() => {
      alert('复制失败，请手动复制：' + service.config_json)
    })
  }

  return (
    <>
      {/* 服务列表 */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-900">我的 MCP 服务</h2>
          <div className="text-sm text-gray-500">共 {services.length} 个服务</div>
        </div>
        
        {services.length === 0 ? (
          <Card className="text-center py-16">
            <CardContent>
              <div className="text-gray-500 text-lg mb-4">🚀 开始添加你的第一个 MCP 服务</div>
              <div className="text-gray-400 text-sm mb-6">MCP 服务可以让 Cursor IDE 具备更强大的功能扩展</div>
              <Button onClick={onAddService} size="lg">添加 MCP 服务</Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {services.map((service) => (
              <Card key={service.id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <CardTitle className="text-lg">{service.name}</CardTitle>
                      <CardDescription className="mt-1">{service.description || '无描述'}</CardDescription>
                    </div>
                    {service.category && (
                      <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full">
                        {service.category.display_name || service.category.name}
                      </span>
                    )}
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {service.author && (
                      <div className="text-sm text-gray-600">👤 {service.author}</div>
                    )}
                    <div className="flex gap-2">
                      <Button 
                        size="sm" 
                        className="flex-1"
                        onClick={() => handleCursorImport(service)}
                      >
                        🚀 导入 Cursor
                      </Button>
                      <Button 
                        size="sm" 
                        variant="outline"
                        onClick={() => copyJsonConfig(service)}
                        title="复制 JSON 配置"
                      >
                        📋
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>

      {/* 快速提示 */}
      {services.length > 0 && (
        <Card className="bg-blue-50 border-blue-200">
          <CardContent className="pt-6">
            <div className="flex items-start gap-3">
              <div className="text-blue-500 text-xl">💡</div>
              <div>
                <div className="font-medium text-blue-900 mb-1">使用提示</div>
                <div className="text-blue-700 text-sm">
                  点击"🚀 导入 Cursor"会自动打开 Cursor IDE 并安装对应的 MCP 服务。点击"📋"可以复制 JSON 配置到剪贴板。
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </>
  )
}



// 添加服务视图组件  
function AddServiceView({ onAdd }: { onAdd: (formData: {
  name: string
  description: string
  category: string
  author: string
  configJson: string
}) => Promise<void> }) {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    category: '',
    author: '',
    configJson: ''
  })
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (formData.name && formData.configJson) {
      setSubmitting(true)
      try {
        await onAdd(formData)
      } finally {
        setSubmitting(false)
      }
    } else {
      alert('请填写必要信息')
    }
  }

  return (
    <>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-gray-900 text-center">添加 MCP 服务</h2>
      </div>

      <Card className="max-w-2xl mx-auto">
        <CardHeader>
          <CardTitle>新建服务配置</CardTitle>
          <CardDescription>填写下面的信息来添加一个新的 MCP 服务</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">服务名称 *</label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="例如: postgres-server"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">描述</label>
              <input
                type="text"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="简短描述这个服务的功能"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">分类</label>
              <select
                value={formData.category}
                onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">选择分类</option>
                <option value="database">数据库</option>
                <option value="api">API服务</option>
                <option value="file">文件系统</option>
                <option value="tool">开发工具</option>
                <option value="other">其他</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">作者</label>
              <input
                type="text"
                value={formData.author}
                onChange={(e) => setFormData({ ...formData, author: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="作者名称"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">JSON 配置 *</label>
              <textarea
                value={formData.configJson}
                onChange={(e) => setFormData({ ...formData, configJson: e.target.value })}
                rows={10}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
                placeholder='{"command": "npx", "args": ["-y", "@modelcontextprotocol/server-postgres", "postgresql://localhost/mydb"]}'
              />
            </div>

                                    <div className="pt-4">
                          <Button type="submit" className="w-full" disabled={submitting}>
                            {submitting ? '添加中...' : '添加服务'}
                          </Button>
                        </div>
          </form>
        </CardContent>
      </Card>
    </>
  )
}

export default App
