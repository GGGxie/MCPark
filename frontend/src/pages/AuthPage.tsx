import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import { GoogleLogin, CredentialResponse } from '@react-oauth/google'
import { Lock } from 'lucide-react'
import { useAuth } from '../hooks/useAuth'
import { api } from '../api'

export default function AuthPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { login } = useAuth()
  const [error, setError] = useState('')

  const handleGoogleSuccess = async (credentialResponse: CredentialResponse) => {
    if (!credentialResponse.credential) return
    try {
      const { token, user } = await api.auth.googleToken(credentialResponse.credential)
      login(token, user)
      navigate('/')
    } catch {
      setError(t('auth.googleError') || 'Google login failed')
    }
  }

  return (
    <div className="min-h-[80vh] flex">
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-gray-900 to-gray-800 text-white flex-col justify-center items-center p-12">
        <div className="max-w-md text-center">
          <div className="text-6xl mb-6">🔌</div>
          <h2 className="text-3xl font-bold mb-4">MCP Park</h2>
          <p className="text-gray-300 text-lg leading-relaxed">
            A community-driven MCP server repository for discovering, sharing, and using Model Context Protocol servers.
          </p>
        </div>
      </div>

      <div className="flex-1 flex items-center justify-center p-8">
        <div className="w-full max-w-md">
          <div className="text-center mb-8">
            <div className="w-12 h-12 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Lock className="w-6 h-6 text-gray-600" />
            </div>
            <h1 className="text-2xl font-bold text-gray-900">{t('auth.title')}</h1>
            <p className="text-sm text-gray-500 mt-1">{t('auth.subtitle')}</p>
          </div>

          {error && (
            <div className="mb-4 p-3 bg-red-50 text-red-600 rounded-lg text-sm">{error}</div>
          )}

          <div className="flex justify-center">
            <GoogleLogin
              onSuccess={handleGoogleSuccess}
              onError={() => setError(t('auth.googleError'))}
              useOneTap
              theme="outline"
              shape="rectangular"
              size="large"
              text="signin_with"
              width="400"
            />
          </div>
        </div>
      </div>
    </div>
  )
}
