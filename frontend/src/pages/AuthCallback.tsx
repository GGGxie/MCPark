import { useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import { api } from '../api'

export default function AuthCallback() {
  const navigate = useNavigate()
  const { login } = useAuth()
  // Captured once: StrictMode runs effects twice in dev, and the fragment is cleared on first read.
  const tokenRef = useRef<string | null | undefined>(undefined)

  useEffect(() => {
    // The backend returns the token in the URL fragment (#token=...) so it is
    // never sent to a server or written to access logs.
    if (tokenRef.current === undefined) {
      tokenRef.current = new URLSearchParams(window.location.hash.replace(/^#/, '')).get('token')
      // Drop the fragment from the address bar and history right away.
      window.history.replaceState(null, '', window.location.pathname + window.location.search)
    }
    const token = tokenRef.current

    if (!token) {
      navigate('/auth')
      return
    }

    localStorage.setItem('token', token)
    api.auth.me().then((user) => {
      login(token, user)
      navigate('/')
    }).catch(() => {
      localStorage.removeItem('token')
      navigate('/auth')
    })
  }, [login, navigate])

  return (
    <div className="min-h-[60vh] flex items-center justify-center">
      <div className="animate-spin w-8 h-8 border-4 border-gray-200 border-t-gray-900 rounded-full" />
    </div>
  )
}
