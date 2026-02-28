import { useState, useEffect, useRef, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { Search, ChevronLeft, ChevronRight, LayoutGrid, AlignLeft } from 'lucide-react'
import { api, Service, PaginatedResponse } from '../api'
import ServiceCard from '../components/ServiceCard'

export default function HomePage() {
  const { t } = useTranslation()
  const [featured, setFeatured] = useState<Service[]>([])
  const [services, setServices] = useState<PaginatedResponse<Service>>({ data: [], total: 0, page: 1, page_size: 12, total_pages: 0 })
  const [sort, setSort] = useState<'hottest' | 'latest'>('hottest')
  const [search, setSearch] = useState('')
  const [inputValue, setInputValue] = useState('')
  const [page, setPage] = useState(1)
  const [compact, setCompact] = useState(false)
  const [featuredIdx, setFeaturedIdx] = useState(0)
  const [loading, setLoading] = useState(true)
  const composingRef = useRef(false)
  const debounceRef = useRef<ReturnType<typeof setTimeout>>()

  useEffect(() => {
    api.services.featured().then(setFeatured).catch(() => {})
  }, [])

  useEffect(() => {
    setLoading(true)
    api.services.list({ page, sort, search: search || undefined })
      .then(setServices)
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [page, sort, search])

  const triggerSearch = useCallback((value: string) => {
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => {
      setSearch(value)
      setPage(1)
    }, 300)
  }, [])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setInputValue(value)
    if (!composingRef.current) {
      triggerSearch(value)
    }
  }

  const handleCompositionEnd = (e: React.CompositionEvent<HTMLInputElement>) => {
    composingRef.current = false
    triggerSearch((e.target as HTMLInputElement).value)
  }

  const handleSearchSubmit = () => {
    if (debounceRef.current) clearTimeout(debounceRef.current)
    setSearch(inputValue)
    setPage(1)
  }

  const featuredVisible = featured.slice(featuredIdx, featuredIdx + 3)

  return (
    <div>
      <section className="bg-gradient-to-b from-gray-50 to-white py-16 px-4">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4 tracking-tight">
            {t('hero.title')}
          </h1>
          <p className="text-lg text-gray-500 mb-8">{t('hero.subtitle')}</p>
          <div className="relative max-w-xl mx-auto">
            <button
              onClick={handleSearchSubmit}
              className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
              aria-label="Search"
            >
              <Search className="w-5 h-5" />
            </button>
            <input
              type="text"
              placeholder={t('home.search')}
              value={inputValue}
              onChange={handleInputChange}
              onCompositionStart={() => { composingRef.current = true }}
              onCompositionEnd={handleCompositionEnd}
              onKeyDown={(e) => { if (e.key === 'Enter') handleSearchSubmit() }}
              className="w-full pl-12 pr-4 py-3.5 rounded-xl border border-gray-200 bg-white text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent shadow-sm"
            />
          </div>
        </div>
      </section>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {featured.length > 0 && !search && (
          <section className="mb-12">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-bold text-gray-900">{t('home.featured')}</h2>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => setFeaturedIdx(Math.max(0, featuredIdx - 1))}
                  disabled={featuredIdx === 0}
                  className="p-1.5 rounded-lg border border-gray-200 text-gray-400 hover:text-gray-600 disabled:opacity-30"
                >
                  <ChevronLeft className="w-4 h-4" />
                </button>
                <button
                  onClick={() => setFeaturedIdx(Math.min(featured.length - 3, featuredIdx + 1))}
                  disabled={featuredIdx >= featured.length - 3}
                  className="p-1.5 rounded-lg border border-gray-200 text-gray-400 hover:text-gray-600 disabled:opacity-30"
                >
                  <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {featuredVisible.map((service) => (
                <ServiceCard key={service.id} service={service} />
              ))}
            </div>
          </section>
        )}

        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-1 bg-gray-100 rounded-lg p-1">
            <button
              onClick={() => { setSort('hottest'); setPage(1) }}
              className={`px-4 py-1.5 rounded-md text-sm font-medium transition-colors ${sort === 'hottest' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700'}`}
            >
              {t('home.hottest')}
            </button>
            <button
              onClick={() => { setSort('latest'); setPage(1) }}
              className={`px-4 py-1.5 rounded-md text-sm font-medium transition-colors ${sort === 'latest' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700'}`}
            >
              {t('home.latest')}
            </button>
          </div>

          <div className="flex items-center gap-1 bg-gray-100 rounded-lg p-1">
            <button
              onClick={() => setCompact(false)}
              className={`p-1.5 rounded-md transition-colors ${!compact ? 'bg-white shadow-sm' : 'text-gray-400'}`}
            >
              <LayoutGrid className="w-4 h-4" />
            </button>
            <button
              onClick={() => setCompact(true)}
              className={`p-1.5 rounded-md transition-colors ${compact ? 'bg-white shadow-sm' : 'text-gray-400'}`}
            >
              <AlignLeft className="w-4 h-4" />
            </button>
          </div>
        </div>

        {loading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="animate-pulse rounded-xl border border-gray-100 overflow-hidden">
                <div className="aspect-video bg-gray-100" />
                <div className="p-5 space-y-3">
                  <div className="h-4 bg-gray-100 rounded w-20" />
                  <div className="h-5 bg-gray-100 rounded w-3/4" />
                  <div className="h-4 bg-gray-100 rounded w-full" />
                </div>
              </div>
            ))}
          </div>
        ) : services.data.length === 0 ? (
          <div className="text-center py-20">
            <p className="text-gray-400 text-lg">{t('home.noResults')}</p>
          </div>
        ) : (
          <div className={`grid gap-6 ${compact ? 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3' : 'grid-cols-1 md:grid-cols-2'}`}>
            {services.data.map((service) => (
              <ServiceCard key={service.id} service={service} compact={compact} />
            ))}
          </div>
        )}

        {services.total_pages > 1 && (
          <div className="flex items-center justify-center gap-2 mt-10">
            <button
              onClick={() => setPage(page - 1)}
              disabled={page <= 1}
              className="px-4 py-2 rounded-lg border border-gray-200 text-sm font-medium text-gray-600 hover:bg-gray-50 disabled:opacity-30"
            >
              {t('home.prev')}
            </button>
            {[...Array(services.total_pages)].map((_, i) => (
              <button
                key={i}
                onClick={() => setPage(i + 1)}
                className={`w-10 h-10 rounded-lg text-sm font-medium transition-colors ${page === i + 1 ? 'bg-gray-900 text-white' : 'border border-gray-200 text-gray-600 hover:bg-gray-50'}`}
              >
                {i + 1}
              </button>
            ))}
            <button
              onClick={() => setPage(page + 1)}
              disabled={page >= services.total_pages}
              className="px-4 py-2 rounded-lg border border-gray-200 text-sm font-medium text-gray-600 hover:bg-gray-50 disabled:opacity-30"
            >
              {t('home.next')}
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
