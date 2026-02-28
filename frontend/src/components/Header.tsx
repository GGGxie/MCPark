import { Link, useLocation } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Globe, Menu, X } from 'lucide-react'
import { useState } from 'react'
import { GoogleLogin, CredentialResponse } from '@react-oauth/google'
import { useAuth } from '../hooks/useAuth'
import { api } from '../api'

export default function Header() {
  const { t, i18n } = useTranslation()
  const { user, login, logout } = useAuth()
  const location = useLocation()
  const [mobileOpen, setMobileOpen] = useState(false)
  const [loginError, setLoginError] = useState('')

  const googleLocale = i18n.language === 'zh' ? 'zh_CN' : 'en'

  const handleGoogleSuccess = async (credentialResponse: CredentialResponse) => {
    if (!credentialResponse.credential) return
    try {
      const { token, user: u } = await api.auth.googleToken(credentialResponse.credential)
      login(token, u)
      setLoginError('')
    } catch {
      setLoginError(t('auth.googleError'))
    }
  }

  const toggleLang = () => {
    const next = i18n.language === 'en' ? 'zh' : 'en'
    i18n.changeLanguage(next)
    localStorage.setItem('language', next)
  }

  const navItems = [
    { to: '/', label: t('nav.discover') },
    { to: '/learn', label: t('nav.learn') },
    { to: '/help', label: t('nav.help') },
    ...(user?.is_admin ? [{ to: '/admin', label: t('nav.admin') }] : []),
  ]

  const isActive = (path: string) =>
    path === '/' ? location.pathname === '/' : location.pathname.startsWith(path)

  return (
    <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center gap-8">
            <Link to="/" className="flex items-center gap-2 font-bold text-xl text-gray-900">
              <span className="text-2xl">🔌</span>
              <span>MCP Park</span>
            </Link>

            <nav className="hidden md:flex items-center gap-1">
              {navItems.map((item) => (
                <Link
                  key={item.to}
                  to={item.to}
                  className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                    isActive(item.to)
                      ? 'bg-gray-100 text-gray-900'
                      : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                  }`}
                >
                  {item.label}
                </Link>
              ))}
            </nav>
          </div>

          <div className="hidden md:flex items-center gap-3">
            <button
              onClick={toggleLang}
              className="p-2 rounded-lg text-gray-500 hover:text-gray-700 hover:bg-gray-50 transition-colors"
              title={i18n.language === 'en' ? '切换到中文' : 'Switch to English'}
            >
              <Globe className="w-5 h-5" />
            </button>

            <Link
              to="/submit"
              className="px-4 py-2 bg-gray-900 text-white text-sm font-medium rounded-lg hover:bg-gray-800 transition-colors"
            >
              {t('nav.submit')}
            </Link>

            {user ? (
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-2">
                  {user.avatar ? (
                    <img src={user.avatar} alt="" className="w-8 h-8 rounded-full" />
                  ) : (
                    <div className="w-8 h-8 rounded-full bg-gray-200 text-gray-600 flex items-center justify-center text-sm font-medium">
                      {user.name[0]}
                    </div>
                  )}
                  <span className="text-sm font-medium text-gray-700">{user.name}</span>
                </div>
                <button
                  onClick={logout}
                  className="text-sm text-gray-500 hover:text-gray-700 transition-colors"
                >
                  {t('nav.logout')}
                </button>
              </div>
            ) : (
              <div className="flex flex-col items-end gap-1">
                <GoogleLogin
                  key={googleLocale}
                  onSuccess={handleGoogleSuccess}
                  onError={() => setLoginError(t('auth.googleError'))}
                  useOneTap={false}
                  locale={googleLocale}
                  theme="outline"
                  shape="rectangular"
                  size="large"
                  text="signin_with"
                />
                {loginError && <p className="text-xs text-red-500">{loginError}</p>}
              </div>
            )}
          </div>

          <button
            onClick={() => setMobileOpen(!mobileOpen)}
            className="md:hidden p-2 rounded-lg text-gray-500 hover:text-gray-700"
          >
            {mobileOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
          </button>
        </div>

        {mobileOpen && (
          <div className="md:hidden py-4 border-t border-gray-100">
            <nav className="flex flex-col gap-1">
              {navItems.map((item) => (
                <Link
                  key={item.to}
                  to={item.to}
                  onClick={() => setMobileOpen(false)}
                  className={`px-3 py-2 rounded-lg text-sm font-medium ${
                    isActive(item.to) ? 'bg-gray-100 text-gray-900' : 'text-gray-600'
                  }`}
                >
                  {item.label}
                </Link>
              ))}
              <Link
                to="/submit"
                onClick={() => setMobileOpen(false)}
                className="px-3 py-2 rounded-lg text-sm font-medium text-gray-600"
              >
                {t('nav.submit')}
              </Link>
              <div className="flex items-center gap-2 px-3 py-2">
                <button onClick={toggleLang} className="text-sm text-gray-500">
                  <Globe className="w-4 h-4 inline mr-1" />
                  {i18n.language === 'en' ? '中文' : 'English'}
                </button>
              </div>
              {user ? (
                <div className="flex items-center justify-between px-3 py-2">
                  <div className="flex items-center gap-2">
                    {user.avatar
                      ? <img src={user.avatar} alt="" className="w-7 h-7 rounded-full" />
                      : <div className="w-7 h-7 rounded-full bg-gray-200 text-gray-600 flex items-center justify-center text-xs font-medium">{user.name[0]}</div>
                    }
                    <span className="text-sm font-medium text-gray-700">{user.name}</span>
                  </div>
                  <button onClick={logout} className="text-sm text-gray-500">
                    {t('nav.logout')}
                  </button>
                </div>
              ) : (
                <div className="px-3 py-2">
                  <GoogleLogin
                    key={googleLocale}
                    onSuccess={handleGoogleSuccess}
                    onError={() => setLoginError(t('auth.googleError'))}
                    useOneTap={false}
                    locale={googleLocale}
                    theme="outline"
                    shape="rectangular"
                    size="large"
                    text="signin_with"
                    width="280"
                  />
                  {loginError && <p className="text-xs text-red-500 mt-1">{loginError}</p>}
                </div>
              )}
            </nav>
          </div>
        )}
      </div>
    </header>
  )
}
