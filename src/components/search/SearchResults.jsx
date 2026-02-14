import { Link } from 'react-router-dom'
import Card from '../common/Card'
import GradeBadge from '../hadith/GradeBadge'
import { getCollectionDisplayName } from '../../utils/constants'
import { getArabicText, getEnglishText, truncateText } from '../../utils/helpers'

const SearchResults = ({ results, loading }) => {
  if (loading) {
    return (
      <div className="space-y-4">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="bg-white dark:bg-background-card-dark rounded-2xl p-6 shadow-card animate-pulse">
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/3 mb-4" />
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2" />
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4" />
          </div>
        ))}
      </div>
    )
  }

  if (!results || results.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
          <svg className="w-10 h-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          No Results Found
        </h3>
        <p className="text-gray-500 dark:text-gray-400">
          Try different keywords or check your spelling
        </p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
        Found {results.length} results
      </p>
      
      {results.map((result, index) => {
        const arabicText = getArabicText(result)
        const englishText = getEnglishText(result)
        const collection = result.collection || {}
        const collectionName = getCollectionDisplayName(collection)
        const grade = result.grades?.[0]?.grade || result.grade
        const bookNumber = result.bookNumber || result.book || 1
        const hadithNumber = result.hadithNumber || result.hadith || result.id

        return (
          <Link
            key={result.id || index}
            to={`/collections/${collection.name || 'bukhari'}/books/${bookNumber}/hadith/${hadithNumber}`}
          >
            <Card 
              className="animate-fade-in"
              style={{ animationDelay: `${index * 50}ms` }}
            >
              <div className="flex items-start justify-between gap-4 mb-3">
                <span className="text-sm text-primary dark:text-primary-400">
                  {collectionName}, Book {bookNumber}, Hadith {hadithNumber}
                </span>
                {grade && <GradeBadge grade={grade} size="sm" />}
              </div>

              {arabicText && (
                <p className="font-arabic text-lg text-gray-900 dark:text-white mb-3 line-clamp-2">
                  {truncateText(arabicText, 150)}
                </p>
              )}

              {englishText && (
                <p className="text-gray-600 dark:text-gray-400 line-clamp-2">
                  {truncateText(englishText, 200)}
                </p>
              )}
            </Card>
          </Link>
        )
      })}
    </div>
  )
}

export default SearchResults

