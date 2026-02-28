import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import { api } from '../api'

export default function AuthCallback() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const { login } = useAuth()

  useEffect(() => {
    const token = params.get('token')
    if (token) {
      localStorage.setItem('token', token)
      api.auth.me().then((user) => {
        login(token, user)
        navigate('/')
      }).catch(() => navigate('/auth'))
    } else {
      navigate('/auth')
    }
  }, [params, login, navigate])

  return (
    <div className="min-h-[60vh] flex items-center justify-center">
      <div className="animate-spin w-8 h-8 border-4 border-gray-200 border-t-gray-900 rounded-full" />
    </div>
  )
}
