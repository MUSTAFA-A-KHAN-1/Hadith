import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useApp } from '../context/AppContext'
import { getSingleHadith } from '../services/api'
import Card from '../components/common/Card'
import Loading from '../components/common/Loading'
import Error from '../components/common/Error'
import Button from '../components/common/Button'

const Bookmarks = () => {
  const { bookmarks, removeBookmark, loading: contextLoading } = useApp()
  const [bookmarkedHadiths, setBookmarkedHadiths] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    const fetchBookmarkedHadiths = async () => {
      if (bookmarks.length === 0) {
        setBookmarkedHadiths([])
        setLoading(false)
        return
      }

      try {
        setLoading(true)
        setError(null)

        const hadiths = await Promise.all(
          bookmarks.map(async (bookmark) => {
            try {
              const data = await getSingleHadith(
                bookmark.collection,
                bookmark.bookNumber,
                bookmark.hadithNumber
              )
              return {
                ...bookmark,
                hadith: data
              }
            } catch (error) {
              console.error('Error fetching bookmarked hadith:', error)
              return null
            }
          })
        )

        setBookmarkedHadiths(hadiths.filter(Boolean))
      } catch (err) {
        console.error('Error fetching bookmarks:', err)
        setError('Failed to load bookmarks')
      } finally {
        setLoading(false)
      }
    }

    fetchBookmarkedHadiths()
  }, [bookmarks])

  const handleRemove = (bookmark) => {
    removeBookmark(bookmark.id)
    setBookmarkedHadiths(prev => prev.filter(b => b.id !== bookmark.id))
  }

  if (loading || contextLoading) {
    return (
      <div className="min-h-screen py-8">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-8">
            <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-white mb-4">
              Bookmarks
            </h1>
          </div>
          <Card>
            <Loading text="Loading bookmarks..." />
          </Card>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center mb-8 animate-fade-in">
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-white mb-4">
            Bookmarks
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Your saved hadiths for later reference
          </p>
        </div>

        {/* Error State */}
        {error && <Error message={error} />}

        {/* Empty State */}
        {!error && bookmarkedHadiths.length === 0 && (
          <div className="text-center py-12 animate-fade-in">
            <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
              <svg className="w-10 h-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
              </svg>
            </div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              No Bookmarks Yet
            </h3>
            <p className="text-gray-500 dark:text-gray-400 mb-6">
              Start saving hadiths that inspire you
            </p>
            <Link to="/collections">
              <Button variant="primary">
                Browse Collections
              </Button>
            </Link>
          </div>
        )}

        {/* Bookmarks List */}
        {!error && bookmarkedHadiths.length > 0 && (
          <div className="space-y-6">
            {bookmarkedHadiths.map((bookmark, index) => (
              <Link
                key={bookmark.id}
                to={`/collections/${bookmark.collection}/books/${bookmark.bookNumber}/hadith/${bookmark.hadithNumber}`}
              >
                <Card 
                  className="animate-fade-in relative"
                  style={{ animationDelay: `${index * 50}ms` }}
                >
                  <button
                    onClick={(e) => {
                      e.preventDefault()
                      handleRemove(bookmark)
                    }}
                    className="absolute top-4 right-4 p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
                    aria-label="Remove bookmark"
                  >
                    <svg className="w-5 h-5 text-gray-400 hover:text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>

                  <div className="pr-12">
                    <div className="flex items-center gap-2 mb-3">
                      <span className="text-sm text-primary dark:text-primary-400">
                        {bookmark.collectionName}, Book {bookmark.bookNumber}, Hadith {bookmark.hadithNumber}
                      </span>
                    </div>

                    {bookmark.hadith?.arabic && (
                      <p className="font-arabic text-lg text-gray-900 dark:text-white mb-3 line-clamp-2">
                        {bookmark.hadith.arabic}
                      </p>
                    )}

                    {bookmark.hadith?.text && (
                      <p className="text-gray-600 dark:text-gray-400 line-clamp-2">
                        {bookmark.hadith.text}
                      </p>
                    )}
                  </div>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default Bookmarks

