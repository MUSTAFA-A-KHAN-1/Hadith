import { useParams, Link } from 'react-router-dom'
import { useSingleHadith, useHadiths } from '../hooks/useHadiths'
import { getCollectionDisplayName } from '../utils/constants'
import HadithDisplay from '../components/hadith/HadithDisplay'
import Loading from '../components/common/Loading'
import Error from '../components/common/Error'

const HadithDetail = () => {
  const { collectionId, bookNumber, hadithNumber } = useParams()
  const { hadith, loading, error } = useSingleHadith(collectionId, bookNumber, hadithNumber)
  const { hadiths: allHadiths } = useHadiths(collectionId, bookNumber, 1, 1000)
  const collectionName = getCollectionDisplayName(collectionId)

  const currentHadithNum = parseInt(hadithNumber)
  
  // Get valid previous/next hadith numbers based on actual available hadiths
  const getValidPreviousHadith = () => {
    if (!allHadiths || allHadiths.length === 0) return currentHadithNum > 1 ? currentHadithNum - 1 : null
    const hadithNumbers = allHadiths
      .map(h => h.hadithNumber || h.hadith || h.id)
      .filter(n => n != null)
      .sort((a, b) => a - b)
    const currentIndex = hadithNumbers.indexOf(currentHadithNum)
    if (currentIndex > 0) return hadithNumbers[currentIndex - 1]
    return null
  }
  
  const getValidNextHadith = () => {
    if (!allHadiths || allHadiths.length === 0) return currentHadithNum + 1
    const hadithNumbers = allHadiths
      .map(h => h.hadithNumber || h.hadith || h.id)
      .filter(n => n != null)
      .sort((a, b) => a - b)
    const currentIndex = hadithNumbers.indexOf(currentHadithNum)
    if (currentIndex >= 0 && currentIndex < hadithNumbers.length - 1) return hadithNumbers[currentIndex + 1]
    return null
  }
  
  const previousHadith = getValidPreviousHadith()
  const nextHadith = getValidNextHadith()

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
            <li>
              <Link to={`/collections/${collectionId}/books/${bookNumber}`} className="hover:text-primary dark:hover:text-primary-400">Book {bookNumber}</Link>
            </li>
            <li>/</li>
            <li className="text-gray-900 dark:text-white">Hadith {hadithNumber}</li>
          </ol>
        </nav>

        {/* Loading State */}
        {loading && (
          <div className="bg-white dark:bg-background-card-dark rounded-2xl p-8 shadow-card">
            <div className="animate-pulse space-y-4">
              <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4" />
              <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-full" />
              <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-3/4" />
              <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-1/2" />
            </div>
          </div>
        )}

        {/* Error State */}
        {error && <Error message={error} />}

        {/* Hadith Display */}
        {!loading && !error && hadith && (
          <HadithDisplay
            hadith={hadith}
            collection={collectionName}
            collectionId={collectionId}
            bookNumber={parseInt(bookNumber)}
            hadithNumber={currentHadithNum}
            previousHadith={previousHadith}
            nextHadith={nextHadith}
          />
        )}
      </div>
    </div>
  )
}

export default HadithDetail

