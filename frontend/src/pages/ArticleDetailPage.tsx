import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { ArrowLeft, Eye, Calendar } from 'lucide-react'
import { api, Article } from '../api'

export default function ArticleDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { t, i18n } = useTranslation()
  const isZH = i18n.language === 'zh'
  const [article, setArticle] = useState<Article | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!id) return
    api.articles.get(id)
      .then(setArticle)
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [id])

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-12">
        <div className="animate-pulse space-y-6">
          <div className="h-8 bg-gray-100 rounded w-48" />
          <div className="h-12 bg-gray-100 rounded w-3/4" />
          <div className="h-64 bg-gray-100 rounded" />
        </div>
      </div>
    )
  }

  if (!article) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-20 text-center">
        <p className="text-gray-400 text-lg">Article not found</p>
      </div>
    )
  }

  const title = isZH && article.title_zh ? article.title_zh : article.title
  const content = isZH && article.content_zh ? article.content_zh : article.content

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700 mb-6 transition-colors"
      >
        <ArrowLeft className="w-4 h-4" /> {t('service.back')}
      </button>

      <header className="mb-8">
        <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-700 mb-4 inline-block">
          {article.category}
        </span>
        <h1 className="text-3xl font-bold text-gray-900 mb-4">{title}</h1>
        <div className="flex items-center gap-4 text-sm text-gray-400">
          <span className="flex items-center gap-1">
            <Calendar className="w-3.5 h-3.5" /> {new Date(article.created_at).toLocaleDateString()}
          </span>
          <span className="flex items-center gap-1">
            <Eye className="w-3.5 h-3.5" /> {article.views}
          </span>
          <span>{article.author_name}</span>
        </div>
      </header>

      <div className="prose prose-gray max-w-none">
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
    </div>
  )
}
