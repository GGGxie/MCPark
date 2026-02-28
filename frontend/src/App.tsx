import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import HomePage from './pages/HomePage'
import ServiceDetailPage from './pages/ServiceDetailPage'
import AuthPage from './pages/AuthPage'
import AuthCallback from './pages/AuthCallback'
import LearnPage from './pages/LearnPage'
import ArticleDetailPage from './pages/ArticleDetailPage'
import HelpPage from './pages/HelpPage'
import SubmitPage from './pages/SubmitPage'
import AdminPage from './pages/AdminPage'
import NotFoundPage from './pages/NotFoundPage'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="service/:id" element={<ServiceDetailPage />} />
        <Route path="auth" element={<AuthPage />} />
        <Route path="auth/callback" element={<AuthCallback />} />
        <Route path="learn" element={<LearnPage />} />
        <Route path="learn/:id" element={<ArticleDetailPage />} />
        <Route path="help" element={<HelpPage />} />
        <Route path="submit" element={<SubmitPage />} />
        <Route path="admin" element={<AdminPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  )
}
