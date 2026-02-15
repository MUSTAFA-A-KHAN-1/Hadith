import { Link } from 'react-router-dom'
import Card from '../common/Card'
import GradeBadge from './GradeBadge'
import HadithActions from './HadithActions'
import { useApp } from '../../context/AppContext'
import { getCollectionDisplayName } from '../../utils/constants'
import { getArabicText, getEnglishText, getGrade } from '../../utils/helpers'

const HadithDisplay = ({ 
  hadith, 
  collection, 
  bookNumber, 
  hadithNumber,
  collectionId,
  previousHadith,
  nextHadith
}) => {
  const { showArabic, fontSize, isBookmarked, toggleBookmark } = useApp()
  
  const arabicText = getArabicText(hadith)
  const englishText = getEnglishText(hadith)
  const grade = getGrade(hadith)
  const collectionName = getCollectionDisplayName(collection)
  const bookmarked = isBookmarked(hadithNumber.toString())

  const hadithData = {
    id: hadithNumber.toString(),
    collection: collectionId,
    bookNumber: bookNumber,
    hadithNumber: hadithNumber,
    arabic: arabicText,
    english: englishText,
    grade: grade
  }

  return (
    <div className="max-w-4xl mx-auto">
      {/* Navigation */}
      {(previousHadith || nextHadith) && (
        <div className="flex items-center justify-between mb-6">
          {previousHadith ? (
            <Link
              to={`/collections/${collectionId}/books/${bookNumber}/hadith/${previousHadith}`}
              className="flex items-center gap-2 text-primary dark:text-primary-400 hover:underline"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              Previous
            </Link>
          ) : <div />}
          
          {nextHadith && (
            <Link
              to={`/collections/${collectionId}/books/${bookNumber}/hadith/${nextHadith}`}
              className="flex items-center gap-2 text-primary dark:text-primary-400 hover:underline"
            >
              Next
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </Link>
          )}
        </div>
      )}

      {/* Hadith Card */}
      <Card className="mb-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div className="flex items-center gap-3">
            <Link
              to={`/collections/${collectionId}/books/${bookNumber}`}
              className="text-sm text-primary dark:text-primary-400 hover:underline"
            >
              {collectionName}, Book {bookNumber}, Hadith {hadithNumber}
            </Link>
            <GradeBadge grade={grade} />
          </div>
          
          <div className="flex items-center gap-3">
            <button
              onClick={() => toggleBookmark(hadithData)}
              className={`
                p-2 rounded-lg transition-all duration-200
                ${bookmarked 
                  ? 'bg-accent/10 text-accent' 
                  : 'bg-gray-100 dark:bg-gray-800 text-gray-500 hover:bg-gray-200 dark:hover:bg-gray-700'
                }
              `}
              aria-label={bookmarked ? "Remove bookmark" : "Add bookmark"}
            >
              <svg 
                className="w-5 h-5" 
                fill={bookmarked ? "currentColor" : "none"} 
                viewBox="0 0 24 24" 
                stroke="currentColor"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
              </svg>
            </button>
            
            <HadithActions 
              hadith={hadith}
              collection={collectionName}
              bookNumber={bookNumber}
              hadithNumber={hadithNumber}
            />
          </div>
        </div>

        {/* Arabic Text */}
        {showArabic && arabicText && (
          <div className="mb-6">
            <p 
              className="font-arabic text-gray-900 dark:text-white leading-arabic text-right arabic-text"
              style={{ fontSize: `${fontSize}px` }}
            >
              {arabicText}
            </p>
            <div className="gold-divider my-6" />
          </div>
        )}

        {/* English Text */}
        {englishText && (
          <div className="mb-6">
            <p 
              className="text-gray-700 dark:text-gray-300 leading-relaxed"
              style={{ fontSize: `${fontSize - 2}px` }}
            >
              {englishText}
            </p>
          </div>
        )}

        {/* Narrator Chain (Isnad) */}
        {hadith?.attribution && (
          <div className="pt-4 border-t border-gray-100 dark:border-gray-700">
            <p className="text-sm text-gray-500 dark:text-gray-400">
              <span className="font-medium text-gray-700 dark:text-gray-300">Narrated by:</span> {hadith.attribution}
            </p>
          </div>
        )}

        {/* Grades */}
        {hadith?.grades && hadith.grades.length > 1 && (
          <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
            <p className="text-sm text-gray-500 dark:text-gray-400 mb-2">
              <span className="font-medium text-gray-700 dark:text-gray-300">Other grades:</span>
            </p>
            <div className="flex flex-wrap gap-2">
              {hadith.grades.slice(1).map((g, idx) => (
                <GradeBadge key={idx} grade={g.grade} size="sm" />
              ))}
            </div>
          </div>
        )}
      </Card>

      {/* Collection Info */}
      <div className="bg-white dark:bg-background-card-dark rounded-2xl p-6 shadow-card">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          About this Hadith
        </h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <p className="text-sm text-gray-500 dark:text-gray-400">Collection</p>
            <Link
              to={`/collections/${collectionId}`}
              className="text-primary dark:text-primary-400 hover:underline font-medium"
            >
              {collectionName}
            </Link>
          </div>
          <div>
            <p className="text-sm text-gray-500 dark:text-gray-400">Book</p>
            <Link
              to={`/collections/${collectionId}/books/${bookNumber}`}
              className="text-primary dark:text-primary-400 hover:underline font-medium"
            >
              Book {bookNumber}
            </Link>
          </div>
          <div>
            <p className="text-sm text-gray-500 dark:text-gray-400">Hadith Number</p>
            <p className="font-medium text-gray-900 dark:text-white">{hadithNumber}</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default HadithDisplay

