import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Search, Copy, Settings, Zap, ChevronDown, ChevronUp, Mail, Github, MessageCircle } from 'lucide-react'

export default function HelpPage() {
  const { t } = useTranslation()
  const [openFaq, setOpenFaq] = useState<number | null>(null)

  const steps = [
    { icon: <Search className="w-6 h-6" />, title: t('help.steps.browse'), desc: t('help.steps.browseDesc'), color: 'bg-blue-50 text-blue-600' },
    { icon: <Copy className="w-6 h-6" />, title: t('help.steps.download'), desc: t('help.steps.downloadDesc'), color: 'bg-green-50 text-green-600' },
    { icon: <Settings className="w-6 h-6" />, title: t('help.steps.install'), desc: t('help.steps.installDesc'), color: 'bg-purple-50 text-purple-600' },
    { icon: <Zap className="w-6 h-6" />, title: t('help.steps.use'), desc: t('help.steps.useDesc'), color: 'bg-amber-50 text-amber-600' },
  ]

  const faqs = [
    { q: t('help.faqs.q1'), a: t('help.faqs.a1') },
    { q: t('help.faqs.q2'), a: t('help.faqs.a2') },
    { q: t('help.faqs.q3'), a: t('help.faqs.a3') },
    { q: t('help.faqs.q4'), a: t('help.faqs.a4') },
    { q: t('help.faqs.q5'), a: t('help.faqs.a5') },
    { q: t('help.faqs.q6'), a: t('help.faqs.a6') },
  ]

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="text-center mb-16">
        <h1 className="text-3xl font-bold text-gray-900 mb-3">{t('help.title')}</h1>
        <p className="text-gray-500 text-lg">{t('help.subtitle')}</p>
      </div>

      <section className="mb-16">
        <h2 className="text-xl font-bold text-gray-900 mb-8 text-center">{t('help.gettingStarted')}</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {steps.map((step, i) => (
            <div key={i} className="text-center">
              <div className={`w-14 h-14 ${step.color} rounded-2xl flex items-center justify-center mx-auto mb-4`}>
                {step.icon}
              </div>
              <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-3 text-sm font-bold text-gray-500">
                {i + 1}
              </div>
              <h3 className="font-semibold text-gray-900 mb-1">{step.title}</h3>
              <p className="text-sm text-gray-500">{step.desc}</p>
            </div>
          ))}
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-xl font-bold text-gray-900 mb-8 text-center">{t('help.faq')}</h2>
        <div className="space-y-3">
          {faqs.map((faq, i) => (
            <div key={i} className="border border-gray-100 rounded-xl overflow-hidden">
              <button
                onClick={() => setOpenFaq(openFaq === i ? null : i)}
                className="w-full flex items-center justify-between p-5 text-left hover:bg-gray-50 transition-colors"
              >
                <span className="text-sm font-medium text-gray-900">{faq.q}</span>
                {openFaq === i ? (
                  <ChevronUp className="w-4 h-4 text-gray-400 flex-shrink-0" />
                ) : (
                  <ChevronDown className="w-4 h-4 text-gray-400 flex-shrink-0" />
                )}
              </button>
              {openFaq === i && (
                <div className="px-5 pb-5">
                  <p className="text-sm text-gray-500 leading-relaxed">{faq.a}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      </section>

      <section>
        <h2 className="text-xl font-bold text-gray-900 mb-8 text-center">{t('help.contact')}</h2>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
          <a
            href="mailto:contact@mcpark.cc"
            className="flex flex-col items-center p-6 rounded-xl border border-gray-100 hover:border-gray-200 hover:shadow-md transition-all"
          >
            <Mail className="w-8 h-8 text-gray-400 mb-3" />
            <span className="text-sm font-medium text-gray-900">{t('help.contactEmail')}</span>
            <span className="text-xs text-gray-400 mt-1">contact@mcpark.cc</span>
          </a>
          <a
            href="https://github.com"
            target="_blank"
            rel="noopener noreferrer"
            className="flex flex-col items-center p-6 rounded-xl border border-gray-100 hover:border-gray-200 hover:shadow-md transition-all"
          >
            <Github className="w-8 h-8 text-gray-400 mb-3" />
            <span className="text-sm font-medium text-gray-900">{t('help.contactGithub')}</span>
            <span className="text-xs text-gray-400 mt-1">github.com/mcpark</span>
          </a>
          <a
            href="#"
            className="flex flex-col items-center p-6 rounded-xl border border-gray-100 hover:border-gray-200 hover:shadow-md transition-all"
          >
            <MessageCircle className="w-8 h-8 text-gray-400 mb-3" />
            <span className="text-sm font-medium text-gray-900">{t('help.contactCommunity')}</span>
            <span className="text-xs text-gray-400 mt-1">Discord / Slack</span>
          </a>
        </div>
      </section>
    </div>
  )
}
