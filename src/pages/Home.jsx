import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useCollections } from '../hooks/useCollections'
import { getRandomHadith } from '../services/api'
import SearchBar from '../components/search/SearchBar'
import CollectionCard from '../components/collection/CollectionCard'
import HadithCard from '../components/hadith/HadithCard'
import Loading from '../components/common/Loading'
import Error from '../components/common/Error'
import Button from '../components/common/Button'
import Card from '../components/common/Card'

const Home = () => {
  const navigate = useNavigate()
  const { collections, loading: collectionsLoading, error: collectionsError } = useCollections()
  const [dailyHadith, setDailyHadith] = useState(null)
  const [dailyLoading, setDailyLoading] = useState(true)
  const [dailyError, setDailyError] = useState(null)

  useEffect(() => {
    const fetchDailyHadith = async () => {
      try {
        setDailyLoading(true)
        setDailyError(null)
        const data = await getRandomHadith()
        if (data && data.hadith) {
          setDailyHadith({
            hadith: data.hadith,
            collection: data.collection,
            book: data.book
          })
        }
      } catch (error) {
        console.error('Error fetching daily hadith:', error)
        setDailyError('Failed to load daily hadith')
      } finally {
        setDailyLoading(false)
      }
    }

    fetchDailyHadith()
  }, [])

  const handleSearch = (query) => {
    if (query && query.trim().length >= 2) {
      navigate(`/search?q=${encodeURIComponent(query.trim())}`)
    }
  }

  const featuredCollections = collections?.slice(0, 6) || []

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative min-h-[70vh] flex items-center justify-center islamic-pattern bg-gradient-to-b from-background-light to-white dark:from-background-dark dark:to-background-card-dark overflow-hidden">
        <div className="absolute inset-0 overflow-hidden">
          <div className="absolute -top-40 -right-40 w-80 h-80 bg-accent/10 rounded-full blur-3xl" />
          <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-primary/10 rounded-full blur-3xl" />
        </div>

        <div className="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div className="animate-fade-in">
            <div className="w-20 h-20 mx-auto mb-6 rounded-full bg-primary flex items-center justify-center shadow-lg">
              <span className="text-4xl text-white font-arabic">﷽</span>
            </div>
            
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-gray-900 dark:text-white mb-4">
              Hadith Portal
            </h1>
            
            <p className="text-lg sm:text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-2xl mx-auto">
              Explore the wisdom of the Prophet Muhammad (peace be upon him) from the most authentic sources in Islam
            </p>

            <div className="max-w-2xl mx-auto mb-8">
              <SearchBar 
                onSearch={handleSearch}
                placeholder="Search for hadiths..."
                className="w-full"
              />
            </div>

            <div className="flex flex-wrap justify-center gap-3">
              <Link to="/collections">
                <Button variant="primary" size="lg">
                  Browse Collections
                </Button>
              </Link>
              <Link to="/collections/bukhari/books/1">
                <Button variant="secondary" size="lg">
                  Start Reading
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Daily Hadith Section */}
      <section className="py-16 bg-white dark:bg-background-card-dark">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-8">
            <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white mb-2">
              Hadith of the Day
            </h2>
            <p className="text-gray-500 dark:text-gray-400">
              Reflect on the wisdom of the Prophet (peace be upon him)
            </p>
          </div>

          {dailyLoading ? (
            <Card className="max-w-3xl mx-auto">
              <Loading text="Loading daily hadith..." />
            </Card>
          ) : dailyError ? (
            <Card className="max-w-3xl mx-auto">
              <Error message={dailyError} />
            </Card>
          ) : dailyHadith ? (
            <div className="animate-fade-in">
              <HadithCard
                hadith={dailyHadith.hadith}
                collection={dailyHadith.collection}
                bookNumber={dailyHadith.book?.bookNumber || 1}
                hadithNumber={dailyHadith.hadith.hadithNumber || 1}
                collectionId={dailyHadith.collection?.name || 'bukhari'}
              />
            </div>
          ) : (
            <Card className="max-w-3xl mx-auto text-center">
              <p className="text-gray-500 dark:text-gray-400">
                Unable to load daily hadith. Please try again later.
              </p>
            </Card>
          )}
        </div>
      </section>

      {/* Collections Section */}
      <section className="py-16 bg-background-light dark:bg-background-dark">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between mb-8">
            <div>
              <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white mb-2">
                Collections
              </h2>
              <p className="text-gray-500 dark:text-gray-400">
                Explore hadith collections from renowned scholars
              </p>
            </div>
            <Link to="/collections">
              <Button variant="ghost" className="hidden sm:flex">
                View All →
              </Button>
            </Link>
          </div>

          {collectionsLoading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {[...Array(6)].map((_, i) => (
                <div key={i} className="bg-white dark:bg-background-card-dark rounded-2xl p-6 shadow-card animate-pulse">
                  <div className="h-14 w-14 bg-gray-200 dark:bg-gray-700 rounded-full mb-4" />
                  <div className="h-5 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-2" />
                  <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2" />
                </div>
              ))}
            </div>
          ) : collectionsError ? (
            <Error message={collectionsError} />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {featuredCollections.map((collection, index) => (
                <CollectionCard 
                  key={collection.name} 
                  collection={collection} 
                  index={index}
                />
              ))}
            </div>
          )}

          <div className="mt-8 text-center sm:hidden">
            <Link to="/collections">
              <Button variant="secondary">
                View All Collections →
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 bg-white dark:bg-background-card-dark">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white mb-2">
              Features
            </h2>
            <p className="text-gray-500 dark:text-gray-400">
              Everything you need to study and reflect on hadiths
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="text-center p-6">
              <div className="w-14 h-14 mx-auto mb-4 rounded-full bg-primary/10 dark:bg-primary-500/20 flex items-center justify-center">
                <svg className="w-7 h-7 text-primary dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                Powerful Search
              </h3>
              <p className="text-gray-500 dark:text-gray-400">
                Find any hadith instantly with our search feature
              </p>
            </div>

            <div className="text-center p-6">
              <div className="w-14 h-14 mx-auto mb-4 rounded-full bg-accent/20 flex items-center justify-center">
                <svg className="w-7 h-7 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                Bookmarks
              </h3>
              <p className="text-gray-500 dark:text-gray-400">
                Save your favorite hadiths for later reference
              </p>
            </div>

            <div className="text-center p-6">
              <div className="w-14 h-14 mx-auto mb-4 rounded-full bg-primary/10 dark:bg-primary-500/20 flex items-center justify-center">
                <svg className="w-7 h-7 text-primary dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                Multiple Languages
              </h3>
              <p className="text-gray-500 dark:text-gray-400">
                Read hadiths in Arabic with English translations
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

export default Home

