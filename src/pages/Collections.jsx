import { useCollections } from '../hooks/useCollections'
import CollectionCard from '../components/collection/CollectionCard'
import Loading from '../components/common/Loading'
import Error from '../components/common/Error'

const Collections = () => {
  const { collections, loading, error } = useCollections()

  return (
    <div className="min-h-screen py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center mb-12 animate-fade-in">
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-white mb-4">
            Hadith Collections
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">
            Explore the six major hadith collections (Kutub al-Sittah) and other famous compilations from renowned Islamic scholars
          </p>
        </div>

        {/* Loading State */}
        {loading && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="bg-white dark:bg-background-card-dark rounded-2xl p-6 shadow-card animate-pulse">
                <div className="h-14 w-14 bg-gray-200 dark:bg-gray-700 rounded-full mb-4" />
                <div className="h-5 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-2" />
                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2" />
              </div>
            ))}
          </div>
        )}

        {/* Error State */}
        {error && <Error message={error} />}

        {/* Collections Grid */}
        {!loading && !error && collections && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {collections.map((collection, index) => (
              <CollectionCard 
                key={collection.name} 
                collection={collection} 
                index={index}
              />
            ))}
          </div>
        )}

        {/* Info Section */}
        {!loading && !error && (
          <div className="mt-12 p-6 bg-white dark:bg-background-card-dark rounded-2xl shadow-card">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
              About Hadith Collections
            </h2>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              The six major hadith collections (Kutub al-Sittah) are considered the most authentic after the Quran:
            </p>
            <ul className="grid grid-cols-1 md:grid-cols-2 gap-3 text-gray-600 dark:text-gray-400">
              <li>• <strong>Sahih al-Bukhari</strong> - Compiled by Imam Bukhari</li>
              <li>• <strong>Sahih Muslim</strong> - Compiled by Imam Muslim</li>
              <li>• <strong>Sunan Abu Dawood</strong> - Compiled by Abu Dawood</li>
              <li>• <strong>Jami' at-Tirmidhi</strong> - Compiled by Imam Tirmidhi</li>
              <li>• <strong>Sunan an-Nasai</strong> - Compiled by Imam Nasai</li>
              <li>• <strong>Sunan Ibn Majah</strong> - Compiled by Ibn Majah</li>
            </ul>
          </div>
        )}
      </div>
    </div>
  )
}

export default Collections

