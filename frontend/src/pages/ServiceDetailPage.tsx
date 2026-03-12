import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { ArrowLeft, ThumbsUp, Share2, Download, Calendar, Tag, User, Eye, Check, Copy } from 'lucide-react'
import { api, Service, Comment } from '../api'
import { useAuth } from '../hooks/useAuth'

export default function ServiceDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { t, i18n } = useTranslation()
  const { user } = useAuth()
  const isZH = i18n.language === 'zh'

  const [service, setService] = useState<Service | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [newComment, setNewComment] = useState('')
  const [liked, setLiked] = useState(false)
  const [loading, setLoading] = useState(true)
  const [shared, setShared] = useState(false)
  const [configCopied, setConfigCopied] = useState(false)

  useEffect(() => {
    if (!id) return
    setLoading(true)
    Promise.all([
      api.services.get(id).then(setService),
      api.services.comments(id).then(setComments),
    ])
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [id])

  const handleLike = async () => {
    if (!id || !user) return
    const res = await api.services.like(id)
    setLiked(res.liked)
    if (service) {
      setService({ ...service, likes: service.likes + (res.liked ? 1 : -1) })
    }
  }

  const handleInstall = async () => {
    if (!id) return
    const res = await api.services.install(id)
    if (service && res.counted) {
      setService({ ...service, install_count: service.install_count + 1 })
    }
  }

  const handleComment = async () => {
    if (!id || !newComment.trim()) return
    await api.services.addComment(id, newComment)
    setNewComment('')
    const updated = await api.services.comments(id)
    setComments(updated)
  }

  const handleShare = async () => {
    const url = window.location.href
    if (navigator.share) {
      try {
        await navigator.share({ title: service?.title, url })
      } catch { /* user cancelled */ }
    } else {
      await navigator.clipboard.writeText(url)
      setShared(true)
      setTimeout(() => setShared(false), 2000)
    }
  }

  const handleCopyConfig = async () => {
    if (!service?.config_json) return
    await navigator.clipboard.writeText(service.config_json)
    setConfigCopied(true)
    setTimeout(() => setConfigCopied(false), 2000)
    handleInstall()
  }

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-12">
        <div className="animate-pulse space-y-6">
          <div className="h-8 bg-gray-100 rounded w-48" />
          <div className="h-12 bg-gray-100 rounded w-3/4" />
          <div className="h-64 bg-gray-100 rounded" />
        </div>
      </div>
    )
  }

  if (!service) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-20 text-center">
        <p className="text-gray-400 text-lg">MCP server not found</p>
      </div>
    )
  }

  const title = isZH && service.title_zh ? service.title_zh : service.title
  const desc = isZH && service.desc_zh ? service.desc_zh : service.description
  const content = isZH && service.content_zh ? service.content_zh : service.content
  const tags = service.tags ? service.tags.split(',').map(t => t.trim()).filter(Boolean) : []

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700 mb-6 transition-colors"
      >
        <ArrowLeft className="w-4 h-4" /> {t('service.back')}
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2">
          <div className="mb-6">
            <h1 className="text-3xl font-bold text-gray-900 mb-3">
              <span className="mr-3">{service.icon || '🔌'}</span>{title}
            </h1>
            <p className="text-gray-500 text-lg leading-relaxed">{desc}</p>
          </div>

          <div className="flex items-center gap-3 mb-8">
            <button
              onClick={handleLike}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg border text-sm font-medium transition-colors ${liked ? 'border-primary-200 bg-primary-50 text-primary-600' : 'border-gray-200 text-gray-600 hover:bg-gray-50'}`}
            >
              <ThumbsUp className="w-4 h-4" /> {service.likes}
            </button>
            <button
              onClick={handleShare}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg border text-sm font-medium transition-colors ${shared ? 'border-green-200 bg-green-50 text-green-600' : 'border-gray-200 text-gray-600 hover:bg-gray-50'}`}
            >
              {shared ? <Check className="w-4 h-4" /> : <Share2 className="w-4 h-4" />}
              {shared ? (isZH ? '已复制链接' : 'Link Copied!') : t('service.share')}
            </button>
          </div>

          {service.config_json && (
            <div className="mb-8">
              <div className="flex items-center justify-between mb-3">
                <h2 className="text-lg font-bold text-gray-900">{t('service.configJson')}</h2>
                <button
                  onClick={handleCopyConfig}
                  className={`flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${configCopied ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}`}
                >
                  {configCopied ? <Check className="w-3.5 h-3.5" /> : <Copy className="w-3.5 h-3.5" />}
                  {configCopied ? t('service.configCopied') : t('service.copyConfig')}
                </button>
              </div>
              <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto text-sm">
                <code>{service.config_json}</code>
              </pre>
            </div>
          )}

          {content && (
            <div className="prose prose-gray max-w-none mb-12">
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                components={{
                  code: ({ children, className, ...props }) => {
                    const isBlock = className?.includes('language-')
                    return isBlock ? (
                      <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto">
                        <code className={className} {...props}>{children}</code>
                      </pre>
                    ) : (
                      <code className="bg-gray-100 text-gray-800 px-1.5 py-0.5 rounded text-sm" {...props}>{children}</code>
                    )
                  },
                  table: ({ children }) => (
                    <div className="overflow-x-auto my-4">
                      <table className="min-w-full border border-gray-200 rounded-lg text-sm">{children}</table>
                    </div>
                  ),
                  thead: ({ children }) => <thead className="bg-gray-50">{children}</thead>,
                  th: ({ children }) => <th className="px-4 py-2 text-left font-semibold text-gray-700 border-b border-gray-200">{children}</th>,
                  td: ({ children }) => <td className="px-4 py-2 text-gray-600 border-b border-gray-100">{children}</td>,
                }}
              >
                {content}
              </ReactMarkdown>
            </div>
          )}

          <section className="border-t border-gray-100 pt-8">
            <h2 className="text-lg font-bold text-gray-900 mb-6">
              {t('service.comments')} ({comments.length})
            </h2>

            {user ? (
              <div className="mb-6">
                <textarea
                  value={newComment}
                  onChange={(e) => setNewComment(e.target.value)}
                  placeholder={t('service.addComment')}
                  className="w-full p-3 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
                  rows={3}
                />
                <button
                  onClick={handleComment}
                  disabled={!newComment.trim()}
                  className="mt-2 px-4 py-2 bg-gray-900 text-white text-sm font-medium rounded-lg hover:bg-gray-800 disabled:opacity-30 transition-colors"
                >
                  {t('service.submitComment')}
                </button>
              </div>
            ) : (
              <div className="mb-6 p-4 bg-gray-50 rounded-lg text-center">
                <p className="text-sm text-gray-500">{t('service.loginToComment')}</p>
              </div>
            )}

            <div className="space-y-4">
              {comments.map((c) => (
                <div key={c.id} className="flex gap-3 p-4 rounded-lg bg-gray-50">
                  <div className="w-8 h-8 rounded-full bg-primary-100 text-primary-600 flex items-center justify-center text-sm font-medium flex-shrink-0">
                    {c.avatar ? (
                      <img src={c.avatar} alt="" className="w-8 h-8 rounded-full" />
                    ) : (
                      c.user_name[0]
                    )}
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-sm font-medium text-gray-900">{c.user_name}</span>
                      <span className="text-xs text-gray-400">
                        {new Date(c.created_at).toLocaleDateString()}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600">{c.content}</p>
                  </div>
                </div>
              ))}
            </div>
          </section>
        </div>

        <div className="lg:col-span-1">
          <div className="sticky top-24 space-y-6">
            <div className="bg-gray-50 rounded-xl p-6">
              <div className="w-full aspect-square bg-gradient-to-br from-gray-100 to-gray-200 rounded-lg flex items-center justify-center mb-6">
                <span className="text-6xl">{service.icon || '🔌'}</span>
              </div>

              <dl className="space-y-4">
                <div className="flex justify-between">
                  <dt className="text-sm text-gray-500 flex items-center gap-1.5"><Tag className="w-3.5 h-3.5" /> {t('service.category')}</dt>
                  <dd className="text-sm font-medium text-gray-900">{service.category}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-sm text-gray-500 flex items-center gap-1.5"><Calendar className="w-3.5 h-3.5" /> {t('service.created')}</dt>
                  <dd className="text-sm font-medium text-gray-900">{new Date(service.created_at).toLocaleDateString()}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-sm text-gray-500 flex items-center gap-1.5"><User className="w-3.5 h-3.5" /> {t('service.author')}</dt>
                  <dd className="text-sm font-medium text-gray-900">{service.author_name}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-sm text-gray-500 flex items-center gap-1.5"><Eye className="w-3.5 h-3.5" /> {t('home.views')}</dt>
                  <dd className="text-sm font-medium text-gray-900">{service.views}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-sm text-gray-500 flex items-center gap-1.5"><Download className="w-3.5 h-3.5" /> {t('service.installs')}</dt>
                  <dd className="text-sm font-medium text-gray-900">{service.install_count}</dd>
                </div>
              </dl>

              {tags.length > 0 && (
                <div className="mt-4 pt-4 border-t border-gray-200">
                  <dt className="text-sm text-gray-500 mb-2">{t('service.tags')}</dt>
                  <div className="flex flex-wrap gap-1.5">
                    {tags.map((tag) => (
                      <span key={tag} className="px-2 py-0.5 bg-gray-200 text-gray-600 rounded text-xs">
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {service.config_json && (
              <button
                onClick={handleCopyConfig}
                className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-gray-900 text-white rounded-lg font-medium hover:bg-gray-800 transition-colors"
              >
                {configCopied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
                {configCopied ? t('service.configCopied') : t('service.copyConfig')}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
