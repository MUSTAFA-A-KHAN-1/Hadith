import { useParams, Link } from 'react-router-dom'
import { useBooks } from '../hooks/useBooks'
import { getCollectionDisplayName } from '../utils/constants'
import BookCard from '../components/collection/BookCard'
import Loading from '../components/common/Loading'
import Error from '../components/common/Error'
import Button from '../components/common/Button'

const CollectionDetail = () => {
  const { collectionId } = useParams()
  const { books, loading, error } = useBooks(collectionId)
  const collectionName = getCollectionDisplayName(collectionId)

  const getArabicName = (id) => {
    const arabicNames = {
      'bukhari': 'صحيح البخاري',
      'muslim': 'صحيح مسلم',
      'abudawud': 'سنن أبي داود',
      'tirmidhi': 'جامع الترمذي',
      'nasai': 'سنن النسائي',
      'ibnmajah': 'سنن ابن ماجه',
      'muwatta': 'موطأ مالك',
      'riyadussaliheen': 'رياض الصالحين',
      'adab': 'الأدب المفرد',
      'shamaa-il': 'شمائل الترمذي',
      'mishkat': 'مشكاة المصابيح'
    }
    return arabicNames[id?.toLowerCase()] || ''
  }

  return (
    <div className="min-h-screen py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Breadcrumb */}
        <nav className="mb-6 animate-fade-in">
          <ol className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
            <li>
              <Link to="/" className="hover:text-primary dark:hover:text-primary-400">Home</Link>
            </li>
            <li>/</li>
            <li>
              <Link to="/collections" className="hover:text-primary dark:hover:text-primary-400">Collections</Link>
            </li>
            <li>/</li>
            <li className="text-gray-900 dark:text-white">{collectionName}</li>
          </ol>
        </nav>

        {/* Header */}
        <div className="mb-8 animate-fade-in" style={{ animationDelay: '100ms' }}>
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-white mb-2">
            {collectionName}
          </h1>
          {getArabicName(collectionId) && (
            <p className="text-2xl font-arabic text-gray-600 dark:text-gray-400">
              {getArabicName(collectionId)}
            </p>
          )}
        </div>

        {/* Loading State */}
        {loading && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="bg-white dark:bg-background-card-dark rounded-2xl p-6 shadow-card animate-pulse">
                <div className="h-12 w-12 bg-gray-200 dark:bg-gray-700 rounded-xl mb-4" />
                <div className="h-5 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-2" />
                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2" />
              </div>
            ))}
          </div>
        )}

        {/* Error State */}
        {error && <Error message={error} />}

        {/* Books Grid */}
        {!loading && !error && books && books.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {books.map((book, index) => (
              <BookCard 
                key={book.bookNumber || book.id} 
                book={book}
                collectionId={collectionId}
                index={index}
              />
            ))}
          </div>
        )}

        {/* Empty State */}
        {!loading && !error && books && books.length === 0 && (
          <div className="text-center py-12">
            <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
              <svg className="w-10 h-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
              </svg>
            </div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              No Books Found
            </h3>
            <p className="text-gray-500 dark:text-gray-400 mb-6">
              This collection doesn't have any books available.
            </p>
            <Link to="/collections">
              <Button variant="primary">
                Browse Other Collections
              </Button>
            </Link>
          </div>
        )}
      </div>
    </div>
  )
}

export default CollectionDetail

