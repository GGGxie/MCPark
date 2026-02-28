import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Send } from 'lucide-react'
import { api } from '../api'
import { useAuth } from '../hooks/useAuth'

const categories = [
  'Development Tools',
  'Productivity',
  'Data & Analytics',
  'Communication',
  'AI & ML',
  'File Management',
  'Web & API',
  'Database',
]

export default function SubmitPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { user } = useAuth()

  const [form, setForm] = useState({
    name: '',
    title: '', title_zh: '',
    description: '', desc_zh: '',
    category: categories[0],
    icon: '',
    config_json: '',
    content: '', content_zh: '',
    tags: '',
  })
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  if (!user) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-20 text-center">
        <Send className="w-12 h-12 text-gray-300 mx-auto mb-4" />
        <p className="text-gray-500 text-lg">{t('submit.loginRequired')}</p>
      </div>
    )
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      await api.services.create(form)
      setSuccess(true)
      setTimeout(() => navigate('/'), 2000)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to submit')
    } finally {
      setSubmitting(false)
    }
  }

  const update = (field: string, value: string) =>
    setForm((prev) => ({ ...prev, [field]: value }))

  if (success) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-20 text-center">
        <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <Send className="w-8 h-8 text-green-600" />
        </div>
        <h2 className="text-xl font-bold text-gray-900 mb-2">{t('submit.success')}</h2>
        <p className="text-gray-500">{t('submit.pending')}</p>
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="text-center mb-10">
        <h1 className="text-3xl font-bold text-gray-900 mb-3">{t('submit.title')}</h1>
        <p className="text-gray-500">{t('submit.subtitle')}</p>
      </div>

      {error && (
        <div className="mb-6 p-3 bg-red-50 text-red-600 rounded-lg text-sm">{error}</div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.nameLabel')}</label>
          <input
            type="text"
            value={form.name} onChange={(e) => update('name', e.target.value)}
            placeholder="my-mcp-server"
            className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.titleLabel')} *</label>
            <input
              type="text" required
              value={form.title} onChange={(e) => update('title', e.target.value)}
              className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.titleZHLabel')}</label>
            <input
              type="text"
              value={form.title_zh} onChange={(e) => update('title_zh', e.target.value)}
              className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.descLabel')} *</label>
            <textarea
              required rows={3}
              value={form.description} onChange={(e) => update('description', e.target.value)}
              className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.descZHLabel')}</label>
            <textarea
              rows={3}
              value={form.desc_zh} onChange={(e) => update('desc_zh', e.target.value)}
              className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.categoryLabel')} *</label>
            <select
              value={form.category} onChange={(e) => update('category', e.target.value)}
              className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white"
            >
              {categories.map((c) => <option key={c} value={c}>{c}</option>)}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.iconLabel')}</label>
            <input
              type="text" placeholder="🔌"
              value={form.icon} onChange={(e) => update('icon', e.target.value)}
              className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.configLabel')}</label>
          <textarea
            rows={6}
            value={form.config_json} onChange={(e) => update('config_json', e.target.value)}
            placeholder={'{\n  "mcpServers": {\n    "my-server": {\n      "command": "npx",\n      "args": ["-y", "@scope/my-mcp-server"]\n    }\n  }\n}'}
            className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.tagsLabel')}</label>
          <input
            type="text" placeholder="files, database, api"
            value={form.tags} onChange={(e) => update('tags', e.target.value)}
            className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.contentLabel')} *</label>
          <textarea
            required rows={8}
            value={form.content} onChange={(e) => update('content', e.target.value)}
            className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
            placeholder={"# My MCP Server\n\nDescribe your MCP server in Markdown..."}
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5">{t('submit.contentZHLabel')}</label>
          <textarea
            rows={8}
            value={form.content_zh} onChange={(e) => update('content_zh', e.target.value)}
            className="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
            placeholder={"# 我的 MCP 服务\n\n用 Markdown 描述你的 MCP 服务..."}
          />
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full flex items-center justify-center gap-2 py-3 bg-gray-900 text-white rounded-lg font-medium hover:bg-gray-800 disabled:opacity-50 transition-colors"
        >
          <Send className="w-4 h-4" />
          {submitting ? '...' : t('submit.submitBtn')}
        </button>
      </form>
    </div>
  )
}
