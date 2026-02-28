import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Check, X, Trash2, ExternalLink } from 'lucide-react'
import { api, Service } from '../api'
import { useAuth } from '../hooks/useAuth'

type StatusFilter = 'all' | 'pending' | 'approved' | 'rejected'

export default function AdminPage() {
  const { t, i18n } = useTranslation()
  const navigate = useNavigate()
  const { user, isLoading } = useAuth()
  const isZH = i18n.language === 'zh'

  const [services, setServices] = useState<Service[]>([])
  const [loading, setLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('pending')
  const [page] = useState(1)
  const [total, setTotal] = useState(0)
  const [acting, setActing] = useState<number | null>(null)

  useEffect(() => {
    if (!isLoading && (!user || !user.is_admin)) {
      navigate('/')
    }
  }, [user, isLoading, navigate])

  const loadServices = () => {
    if (!user?.is_admin) return
    setLoading(true)
    api.admin.listServices({ page, status: statusFilter === 'all' ? undefined : statusFilter })
      .then((res) => {
        setServices(res.data)
        setTotal(res.total)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadServices()
  }, [user, statusFilter, page])

  const handleApprove = async (id: number) => {
    setActing(id)
    try {
      await api.admin.approve(id)
      loadServices()
    } finally {
      setActing(null)
    }
  }

  const handleReject = async (id: number) => {
    setActing(id)
    try {
      await api.admin.reject(id)
      loadServices()
    } finally {
      setActing(null)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this MCP server?')) return
    setActing(id)
    try {
      await api.admin.delete(id)
      loadServices()
    } finally {
      setActing(null)
    }
  }

  if (isLoading || !user?.is_admin) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="animate-spin w-8 h-8 border-4 border-gray-200 border-t-gray-900 rounded-full" />
      </div>
    )
  }

  const statusBadge = (status: string) => {
    const cls = ({
      pending: 'bg-yellow-100 text-yellow-700',
      approved: 'bg-green-100 text-green-700',
      rejected: 'bg-red-100 text-red-700',
    } as Record<string, string>)[status] || 'bg-gray-100 text-gray-700'
    return <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${cls}`}>{status}</span>
  }

  const filters: StatusFilter[] = ['all', 'pending', 'approved', 'rejected']

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">{t('admin.title')}</h1>

      <div className="flex items-center gap-2 mb-6">
        {filters.map((s) => (
          <button
            key={s}
            onClick={() => setStatusFilter(s)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${statusFilter === s ? 'bg-gray-900 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
          >
            {t(`admin.${s}`)}
          </button>
        ))}
        <span className="ml-auto text-sm text-gray-400">{total} total</span>
      </div>

      {loading ? (
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="animate-pulse h-20 bg-gray-100 rounded-xl" />
          ))}
        </div>
      ) : services.length === 0 ? (
        <div className="text-center py-20">
          <p className="text-gray-400">{t('admin.noServices')}</p>
        </div>
      ) : (
        <div className="space-y-3">
          {services.map((s) => {
            const title = isZH && s.title_zh ? s.title_zh : s.title
            return (
              <div key={s.id} className="flex items-center gap-4 p-4 bg-white border border-gray-100 rounded-xl hover:border-gray-200 transition-colors">
                <span className="text-2xl flex-shrink-0">{s.icon || '🔌'}</span>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <span className="text-sm font-semibold text-gray-900 truncate">{title}</span>
                    {statusBadge(s.status)}
                  </div>
                  <p className="text-xs text-gray-400">{s.category} · {s.author_name} · {new Date(s.created_at).toLocaleDateString()}</p>
                </div>
                <div className="flex items-center gap-2 flex-shrink-0">
                  <a
                    href={`/service/${s.id}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors"
                  >
                    <ExternalLink className="w-4 h-4" />
                  </a>
                  {s.status !== 'approved' && (
                    <button
                      onClick={() => handleApprove(s.id)}
                      disabled={acting === s.id}
                      className="flex items-center gap-1 px-3 py-1.5 bg-green-100 text-green-700 rounded-lg text-sm font-medium hover:bg-green-200 disabled:opacity-50 transition-colors"
                    >
                      <Check className="w-3.5 h-3.5" /> {t('admin.approve')}
                    </button>
                  )}
                  {s.status !== 'rejected' && (
                    <button
                      onClick={() => handleReject(s.id)}
                      disabled={acting === s.id}
                      className="flex items-center gap-1 px-3 py-1.5 bg-red-100 text-red-700 rounded-lg text-sm font-medium hover:bg-red-200 disabled:opacity-50 transition-colors"
                    >
                      <X className="w-3.5 h-3.5" /> {t('admin.reject')}
                    </button>
                  )}
                  <button
                    onClick={() => handleDelete(s.id)}
                    disabled={acting === s.id}
                    className="flex items-center gap-1 px-3 py-1.5 bg-gray-100 text-gray-600 rounded-lg text-sm font-medium hover:bg-gray-200 disabled:opacity-50 transition-colors"
                  >
                    <Trash2 className="w-3.5 h-3.5" /> {t('admin.delete')}
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
