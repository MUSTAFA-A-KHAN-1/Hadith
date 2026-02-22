import { Link } from 'react-router-dom'
import Card from '../common/Card'
import GradeBadge from './GradeBadge'
import { useApp } from '../../context/AppContext'
import { getCollectionDisplayName } from '../../utils/constants'
import { getArabicText, getEnglishText } from '../../utils/helpers'

const HadithCard = ({ 
  hadith, 
  collection, 
  bookNumber, 
  hadithNumber,
  collectionId,
  index = 0,
  showActions = false,
  onBookmark,
  isBookmarked = false
}) => {
  const { fontSize } = useApp()
  const arabicText = getArabicText(hadith)
  const englishText = getEnglishText(hadith)
  const grade = hadith?.grades?.[0]?.grade || hadith?.grade
  const collectionName = getCollectionDisplayName(collection)

  return (
    <Link to={`/collections/${collectionId}/books/${bookNumber}/hadith/${hadithNumber}`}>
      <Card 
        className="animate-fade-in"
        style={{ animationDelay: `${index * 50}ms` }}
      >
        {/* Reference */}
        <div className="flex items-center justify-between mb-4">
          <span className="text-sm text-primary dark:text-primary-400">
            {collectionName}, Book {bookNumber}, Hadith {hadithNumber}
          </span>
          <div className="flex items-center gap-2">
            {grade && <GradeBadge grade={grade} size="sm" />}
            {showActions && onBookmark && (
              <button
                onClick={() => onBookmark(hadith)}
                className="p-1.5 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
                aria-label={isBookmarked ? "Remove bookmark" : "Add bookmark"}
              >
                <svg 
                  className={`w-5 h-5 ${isBookmarked ? 'text-accent fill-current' : 'text-gray-400'}`} 
                  fill={isBookmarked ? "currentColor" : "none"} 
                  viewBox="0 0 24 24" 
                  stroke="currentColor"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
                </svg>
              </button>
            )}
          </div>
        </div>

        {/* Arabic Text */}
        {arabicText && (
          <>
            <p 
              className="font-arabic text-gray-900 dark:text-white leading-arabic mb-4 text-right arabic-text"
              style={{ fontSize: `${fontSize}px` }}
            >
              {arabicText}
            </p>
            <div className="gold-divider my-4" />
          </>
        )}
        {/* Narrator Chain (Isnad) */}
        {(hadith?.attribution || hadith?.narrator) && (
          <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
            <p className="text-sm text-gray-500 dark:text-gray-400">
              <span className="font-medium"></span> {hadith.attribution || hadith.narrator}
            </p>
          </div>
        )}
        {/* English Text */}
        {englishText && (
          <p className="text-gray-700 dark:text-gray-300 leading-relaxed">
            {englishText}
          </p>
        )}
      </Card>
    </Link>
  )
}

export default HadithCard

