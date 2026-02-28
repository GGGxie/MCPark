import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Github, Twitter, Mail } from 'lucide-react'

export default function Footer() {
  const { t } = useTranslation()

  return (
    <footer className="bg-gray-50 border-t border-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div className="md:col-span-1">
            <Link to="/" className="flex items-center gap-2 font-bold text-lg text-gray-900">
              <span className="text-xl">🔌</span>
              <span>MCP Park</span>
            </Link>
            <p className="mt-3 text-sm text-gray-500 leading-relaxed">
              {t('footer.description')}
            </p>
          </div>

          <div>
            <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('footer.product')}</h3>
            <ul className="space-y-3">
              <li><Link to="/" className="text-sm text-gray-500 hover:text-gray-700">{t('nav.discover')}</Link></li>
              <li><Link to="/learn" className="text-sm text-gray-500 hover:text-gray-700">{t('nav.learn')}</Link></li>
              <li><Link to="/help" className="text-sm text-gray-500 hover:text-gray-700">{t('nav.help')}</Link></li>
              <li><Link to="/submit" className="text-sm text-gray-500 hover:text-gray-700">{t('nav.submit')}</Link></li>
            </ul>
          </div>

          <div>
            <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('footer.legal')}</h3>
            <ul className="space-y-3">
              <li><span className="text-sm text-gray-400 cursor-default">{t('footer.terms')}</span></li>
              <li><span className="text-sm text-gray-400 cursor-default">{t('footer.privacy')}</span></li>
            </ul>
          </div>

          <div>
            <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('footer.connect')}</h3>
            <ul className="space-y-3">
              <li>
                <a href="https://github.com" target="_blank" rel="noopener noreferrer" className="flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700">
                  <Github className="w-4 h-4" /> GitHub
                </a>
              </li>
              <li>
                <a href="https://twitter.com" target="_blank" rel="noopener noreferrer" className="flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700">
                  <Twitter className="w-4 h-4" /> Twitter
                </a>
              </li>
              <li>
                <a href="mailto:contact@mcpark.cc" className="flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700">
                  <Mail className="w-4 h-4" /> Email
                </a>
              </li>
            </ul>
          </div>
        </div>

        <div className="mt-10 pt-6 border-t border-gray-200 flex flex-col sm:flex-row items-center justify-between gap-4">
          <p className="text-sm text-gray-400">
            &copy; {new Date().getFullYear()} MCP Park. {t('footer.rights')}
          </p>
        </div>
      </div>
    </footer>
  )
}
