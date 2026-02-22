import { Link } from 'react-router-dom'
import Button from '../components/common/Button'

const NotFound = () => {
  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4">
      <div className="text-center animate-fade-in">
        {/* Decorative elements */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 -right-40 w-80 h-80 bg-accent/5 rounded-full blur-3xl" />
          <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-primary/5 rounded-full blur-3xl" />
        </div>

        <div className="relative">
          {/* Large 404 text */}
          <div className="text-[150px] sm:text-[200px] font-bold text-primary/10 dark:text-primary-500/10 leading-none select-none">
            404
          </div>

          {/* Content */}
          <div className="absolute inset-0 flex flex-col items-center justify-center">
            <div className="w-20 h-20 mx-auto mb-6 rounded-full bg-primary/10 dark:bg-primary-500/20 flex items-center justify-center">
              <svg className="w-10 h-10 text-primary dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>

            <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-white mb-4">
              Page Not Found
            </h1>

            <p className="text-lg text-gray-600 dark:text-gray-400 mb-8 max-w-md">
              The page you're looking for doesn't exist or has been moved. Let us guide you back.
            </p>

            <div className="flex flex-wrap justify-center gap-4">
              <Link to="/">
                <Button variant="primary" size="lg">
                  Go Home
                </Button>
              </Link>
              <Link to="/collections">
                <Button variant="secondary" size="lg">
                  Browse Collections
                </Button>
              </Link>
            </div>
          </div>
        </div>

        {/* Quran verse decoration */}
        <div className="mt-16 pt-8 border-t border-gray-200 dark:border-gray-700">
          <p className="font-arabic text-2xl text-gray-700 dark:text-gray-300 mb-2">
            وَعَسَى أَن تَكۡرَهُوا۟ شَيۡئًا وَهُوَ خَيۡرٞ لَّكُمۡ
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            "And it may be that you dislike a thing which is good for you..."
          </p>
          <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">
            — Surah Al-Baqarah 2:216
          </p>
        </div>
      </div>
    </div>
  )
}

export default NotFound

