import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

export default function NotFoundPage() {
  const { t } = useTranslation()

  return (
    <div className="min-h-[60vh] flex flex-col items-center justify-center px-4">
      <h1 className="text-8xl font-bold text-gray-200 mb-4">{t('notFound.title')}</h1>
      <p className="text-xl text-gray-500 mb-8">{t('notFound.message')}</p>
      <Link
        to="/"
        className="px-6 py-3 bg-gray-900 text-white rounded-lg font-medium hover:bg-gray-800 transition-colors"
      >
        {t('notFound.back')}
      </Link>
    </div>
  )
}
