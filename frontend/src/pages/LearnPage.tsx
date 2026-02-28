import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { BookOpen, Eye } from 'lucide-react'
import { api, Article } from '../api'

export default function LearnPage() {
  const navigate = useNavigate()
  const { t, i18n } = useTranslation()
  const isZH = i18n.language === 'zh'
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.articles.list()
      .then(setArticles)
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const categoryColors: Record<string, string> = {
    'Getting Started': 'bg-green-100 text-green-700',
    'Development': 'bg-blue-100 text-blue-700',
    'Concepts': 'bg-purple-100 text-purple-700',
    'Integration': 'bg-orange-100 text-orange-700',
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="text-center mb-12">
        <h1 className="text-3xl font-bold text-gray-900 mb-3">{t('learn.title')}</h1>
        <p className="text-gray-500 text-lg">{t('learn.subtitle')}</p>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="animate-pulse rounded-xl border border-gray-100 p-6 space-y-3">
              <div className="h-4 bg-gray-100 rounded w-20" />
              <div className="h-5 bg-gray-100 rounded w-3/4" />
              <div className="h-4 bg-gray-100 rounded w-full" />
            </div>
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {articles.map((article) => {
            const title = isZH && article.title_zh ? article.title_zh : article.title
            const desc = isZH && article.desc_zh ? article.desc_zh : article.description
            const colorClass = categoryColors[article.category] || 'bg-gray-100 text-gray-700'

            return (
              <article
                key={article.id}
                onClick={() => navigate(`/learn/${article.id}`)}
                className="group bg-white rounded-xl border border-gray-100 hover:border-gray-200 hover:shadow-lg transition-all duration-300 p-6 cursor-pointer"
              >
                <div className="flex items-center gap-2 mb-3">
                  <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${colorClass}`}>
                    {article.category}
                  </span>
                </div>

                <h3 className="text-base font-semibold text-gray-900 group-hover:text-primary-600 transition-colors mb-2">
                  <BookOpen className="w-4 h-4 inline mr-2 text-gray-400" />
                  {title}
                </h3>
                <p className="text-sm text-gray-500 line-clamp-2 leading-relaxed mb-4">{desc}</p>

                <div className="flex items-center justify-between text-sm text-gray-400">
                  <span className="flex items-center gap-1">
                    <Eye className="w-3.5 h-3.5" /> {article.views}
                  </span>
                  <span className="text-xs">
                    {new Date(article.created_at).toLocaleDateString()}
                  </span>
                </div>
              </article>
            )
          })}
        </div>
      )}
    </div>
  )
}
