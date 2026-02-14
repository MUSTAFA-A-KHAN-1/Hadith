import { Link } from 'react-router-dom'

const Footer = () => {
  return (
    <footer className="bg-white dark:bg-background-card-dark border-t border-gray-200 dark:border-gray-700 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center">
              <span className="text-white font-arabic text-sm">﷽</span>
            </div>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              © {new Date().getFullYear()} Hadith Portal
            </p>
          </div>
          
          <div className="flex items-center gap-6 text-sm">
            <a 
              href="https://sunnah.com" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-gray-500 dark:text-gray-400 hover:text-primary dark:hover:text-primary-400 transition-colors"
            >
              sunnah.com
            </a>
            <Link 
              to="/collections"
              className="text-gray-500 dark:text-gray-400 hover:text-primary dark:hover:text-primary-400 transition-colors"
            >
              Collections
            </Link>
          </div>
          
          <p className="text-xs text-gray-400 dark:text-gray-500">
            Seeking knowledge is an obligation upon every Muslim
          </p>
        </div>
      </div>
    </footer>
  )
}

export default Footer

