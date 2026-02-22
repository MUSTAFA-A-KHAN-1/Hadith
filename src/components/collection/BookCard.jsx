import { Link } from 'react-router-dom'
import Card from '../common/Card'
import { formatNumber } from '../../utils/constants'

const BookCard = ({ book, collectionId, index = 0 }) => {
  const bookNumber = book.bookNumber || book.number || 1
  const hadithCount = book.hadithCount || book.hadithsCount || 0
  const title = book.title || book.name || `Book ${bookNumber}`
  const englishTitle = book.englishTitle || book.title || ''

  return (
    <Link to={`/collections/${collectionId}/books/${bookNumber}`}>
      <Card 
        className="h-full animate-fade-in"
        style={{ animationDelay: `${index * 50}ms` }}
      >
        <div className="flex items-start gap-4">
          <div className="w-12 h-12 rounded-xl bg-accent/10 dark:bg-accent/20 flex items-center justify-center flex-shrink-0">
            <span className="text-xl font-bold text-accent dark:text-accent-light">
              {bookNumber}
            </span>
          </div>
          <div className="flex-1 min-w-0">
            <h3 className="text-base font-semibold text-gray-900 dark:text-white mb-1 truncate">
              {title}
            </h3>
            {englishTitle && (
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-2 truncate">
                {englishTitle}
              </p>
            )}
            <p className="text-xs text-gray-500 dark:text-gray-500">
              {formatNumber(hadithCount)} hadiths
            </p>
          </div>
        </div>
      </Card>
    </Link>
  )
}

export default BookCard

