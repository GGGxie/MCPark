import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { ThumbsUp, Eye, Download } from 'lucide-react'
import { Service } from '../api'

const categoryColors: Record<string, string> = {
  'Development Tools': 'bg-blue-100 text-blue-700',
  'Productivity': 'bg-green-100 text-green-700',
  'Data & Analytics': 'bg-purple-100 text-purple-700',
  'Communication': 'bg-orange-100 text-orange-700',
  'AI & ML': 'bg-pink-100 text-pink-700',
  'File Management': 'bg-yellow-100 text-yellow-700',
  'Web & API': 'bg-cyan-100 text-cyan-700',
  'Database': 'bg-indigo-100 text-indigo-700',
}

interface Props {
  service: Service
  compact?: boolean
}

export default function ServiceCard({ service, compact }: Props) {
  const { i18n } = useTranslation()
  const isZH = i18n.language === 'zh'
  const title = isZH && service.title_zh ? service.title_zh : service.title
  const desc = isZH && service.desc_zh ? service.desc_zh : service.description
  const colorClass = categoryColors[service.category] || 'bg-gray-100 text-gray-700'

  return (
    <Link
      to={`/service/${service.id}`}
      className="group block bg-white rounded-xl border border-gray-100 hover:border-gray-200 hover:shadow-lg transition-all duration-300 overflow-hidden"
    >
      {!compact && service.image && (
        <div className="aspect-video bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center overflow-hidden">
          <img src={service.image} alt={title} className="w-full h-full object-cover" />
        </div>
      )}
      {!compact && !service.image && (
        <div className="aspect-video bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center">
          <span className="text-5xl">{service.icon || '🔌'}</span>
        </div>
      )}

      <div className="p-5">
        <div className="flex items-center gap-2 mb-3">
          <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${colorClass}`}>
            {service.category}
          </span>
          {service.featured && (
            <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-700">
              Featured
            </span>
          )}
        </div>

        <h3 className="text-base font-semibold text-gray-900 group-hover:text-primary-600 transition-colors mb-2 line-clamp-1">
          {compact && <span className="mr-2">{service.icon || '🔌'}</span>}
          {title}
        </h3>
        <p className="text-sm text-gray-500 line-clamp-2 leading-relaxed mb-4">{desc}</p>

        <div className="flex items-center justify-between text-sm text-gray-400">
          <div className="flex items-center gap-4">
            <span className="flex items-center gap-1">
              <ThumbsUp className="w-3.5 h-3.5" /> {service.likes}
            </span>
            <span className="flex items-center gap-1">
              <Eye className="w-3.5 h-3.5" /> {service.views}
            </span>
            <span className="flex items-center gap-1">
              <Download className="w-3.5 h-3.5" /> {service.install_count}
            </span>
          </div>
          {service.author_name && (
            <span className="text-xs text-gray-400">{service.author_name}</span>
          )}
        </div>
      </div>
    </Link>
  )
}
