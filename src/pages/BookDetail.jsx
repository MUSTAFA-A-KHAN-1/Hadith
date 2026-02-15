import { useParams, Link } from 'react-router-dom'
import { useHadiths } from '../hooks/useHadiths'
import { useBook } from '../hooks/useBooks'
import { getCollectionDisplayName } from '../utils/constants'
import HadithCard from '../components/hadith/HadithCard'
import Loading from '../components/common/Loading'
import Error from '../components/common/Error'
import Button from '../components/common/Button'

const BookDetail = () => {
  const { collectionId, bookNumber } = useParams()
  const { hadiths, loading, error, total } = useHadiths(collectionId, bookNumber)
  const { book } = useBook(collectionId, bookNumber)
  const collectionName = getCollectionDisplayName(collectionId)
  
  const bookTitle = book?.title || `Book ${bookNumber}`

  return (
    <div className="min-h-screen py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Breadcrumb */}
        <nav className="mb-6 animate-fade-in">
          <ol className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400 flex-wrap">
            <li>
              <Link to="/" className="hover:text-primary dark:hover:text-primary-400">Home</Link>
            </li>
            <li>/</li>
            <li>
              <Link to="/collections" className="hover:text-primary dark:hover:text-primary-400">Collections</Link>
            </li>
            <li>/</li>
            <li>
              <Link to={`/collections/${collectionId}`} className="hover:text-primary dark:hover:text-primary-400">{collectionName}</Link>
            </li>
            <li>/</li>
            <li className="text-gray-900 dark:text-white">{bookTitle}</li>
          </ol>
        </nav>

        {/* Header */}
        <div className="mb-8 animate-fade-in" style={{ animationDelay: '100ms' }}>
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-white mb-2">
            {bookTitle}
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            {collectionName} • {total ? `${total} hadiths` : 'Loading...'}
          </p>
        </div>

        {/* Loading State */}
        {loading && (
          <div className="space-y-6">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="bg-white dark:bg-background-card-dark rounded-2xl p-6 shadow-card animate-pulse">
                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4 mb-4" />
                <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2" />
                <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-2" />
                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2" />
              </div>
            ))}
          </div>
        )}

        {/* Error State */}
        {error && <Error message={error} />}

        {/* Hadiths List */}
        {!loading && !error && hadiths && hadiths.length > 0 && (
          <div className="space-y-6">
            {hadiths.map((hadithItem, index) => (
              <HadithCard
                key={hadithItem.hadithNumber || hadithItem.id || index}
                hadith={hadithItem}
                collection={collectionName}
                collectionId={collectionId}
                bookNumber={parseInt(bookNumber)}
                hadithNumber={hadithItem.hadithNumber || hadithItem.hadith || index + 1}
                index={index}
              />
            ))}
          </div>
        )}

        {/* Empty State */}
        {!loading && !error && hadiths && hadiths.length === 0 && (
          <div className="text-center py-12">
            <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
              <svg className="w-10 h-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              No Hadiths Found
            </h3>
            <p className="text-gray-500 dark:text-gray-400 mb-6">
              This book doesn't have any hadiths available.
            </p>
            <Link to={`/collections/${collectionId}`}>
              <Button variant="primary">
                Browse Other Books
              </Button>
            </Link>
          </div>
        )}
      </div>
    </div>
  )
}

export default BookDetail

