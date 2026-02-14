import { useState, useEffect, useCallback } from 'react'
import { useSearchParams, Link } from 'react-router-dom'
import { searchHadiths } from '../services/api'
import SearchBar from '../components/search/SearchBar'
import SearchResults from '../components/search/SearchResults'
import Loading from '../components/common/Loading'
import Error from '../components/common/Error'
import Button from '../components/common/Button'

const Search = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const initialQuery = searchParams.get('q') || ''
  
  const [query, setQuery] = useState(initialQuery)
  const [results, setResults] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [hasSearched, setHasSearched] = useState(false)

  const performSearch = useCallback(async (searchQuery) => {
    if (!searchQuery || searchQuery.trim().length < 2) {
      setResults([])
      setHasSearched(false)
      return
    }

    try {
      setLoading(true)
      setError(null)
      setHasSearched(true)
      
      const data = await searchHadiths(searchQuery.trim())
      setResults(data || [])
    } catch (err) {
      console.error('Search error:', err)
      setError('Failed to search hadiths. Please try again.')
      setResults([])
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    const debounceTimer = setTimeout(() => {
      if (initialQuery) {
        performSearch(initialQuery)
      }
    }, 300)

    return () => clearTimeout(debounceTimer)
  }, [initialQuery, performSearch])

  const handleSearch = (newQuery) => {
    setQuery(newQuery)
    if (newQuery.trim().length >= 2) {
      setSearchParams({ q: newQuery.trim() })
    } else {
      setSearchParams({})
    }
  }

  const handleSearchSubmit = (e) => {
    e.preventDefault()
    performSearch(query)
  }

  return (
    <div className="min-h-screen py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center mb-8 animate-fade-in">
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-white mb-4">
            Search Hadiths
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Search across all collections for specific hadiths
          </p>
        </div>

        {/* Search Form */}
        <div className="mb-8 animate-fade-in" style={{ animationDelay: '100ms' }}>
          <form onSubmit={handleSearchSubmit}>
            <SearchBar
              onSearch={handleSearch}
              initialValue={query}
              placeholder="Search for hadiths..."
              autoFocus
            />
          </form>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-2 text-center">
            Type at least 2 characters to search
          </p>
        </div>

        {/* Results */}
        <div className="animate-fade-in" style={{ animationDelay: '200ms' }}>
          {loading ? (
            <SearchResults results={[]} loading={true} />
          ) : error ? (
            <Error message={error} />
          ) : hasSearched ? (
            <SearchResults results={results} loading={false} />
          ) : (
            <div className="text-center py-12">
              <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
                <svg className="w-10 h-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                Start Your Search
              </h3>
              <p className="text-gray-500 dark:text-gray-400 mb-6">
                Enter keywords to find hadiths from all collections
              </p>
              <div className="flex flex-wrap justify-center gap-2">
                <span className="text-sm text-gray-500 dark:text-gray-400">Try:</span>
                <button 
                  onClick={() => handleSearch('faith')} 
                  className="text-sm text-primary dark:text-primary-400 hover:underline"
                >
                  faith
                </button>
                <button 
                  onClick={() => handleSearch('prayer')} 
                  className="text-sm text-primary dark:text-primary-400 hover:underline"
                >
                  prayer
                </button>
                <button 
                  onClick={() => handleSearch('charity')} 
                  className="text-sm text-primary dark:text-primary-400 hover:underline"
                >
                  charity
                </button>
                <button 
                  onClick={() => handleSearch('patience')} 
                  className="text-sm text-primary dark:text-primary-400 hover:underline"
                >
                  patience
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Search

