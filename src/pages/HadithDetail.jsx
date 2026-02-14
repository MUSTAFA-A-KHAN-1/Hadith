import { useParams, Link } from 'react-router-dom'
import { useSingleHadith } from '../hooks/useHadiths'
import { getCollectionDisplayName } from '../utils/constants'
import HadithDisplay from '../components/hadith/HadithDisplay'
import Loading from '../components/common/Loading'
import Error from '../components/common/Error'

const HadithDetail = () => {
  const { collectionId, bookNumber, hadithNumber } = useParams()
  const { hadith, loading, error } = useSingleHadith(collectionId, bookNumber, hadithNumber)
  const collectionName = getCollectionDisplayName(collectionId)

  const currentHadithNum = parseInt(hadithNumber)
  const previousHadith = currentHadithNum > 1 ? currentHadithNum - 1 : null
  const nextHadith = currentHadithNum + 1

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

