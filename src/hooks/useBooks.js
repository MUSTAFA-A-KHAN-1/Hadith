import { useState, useEffect, useCallback } from 'react'
import { getBooks, getBook } from '../services/api'

export const useBooks = (collectionId) => {
  const [books, setBooks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const fetchBooks = useCallback(async () => {
    if (!collectionId) return
    try {
      setLoading(true)
      setError(null)
      const data = await getBooks(collectionId)
      setBooks(data)
    } catch (err) {
      setError(err.message || 'Failed to fetch books')
    } finally {
      setLoading(false)
    }
  }, [collectionId])

  useEffect(() => {
    fetchBooks()
  }, [fetchBooks])

  return { books, loading, error, refetch: fetchBooks }
}

export const useBook = (collectionId, bookId) => {
  const [book, setBook] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const fetchBook = useCallback(async () => {
    if (!collectionId || !bookId) return
    try {
      setLoading(true)
      setError(null)
      const data = await getBook(collectionId, bookId)
      setBook(data)
    } catch (err) {
      setError(err.message || 'Failed to fetch book')
    } finally {
      setLoading(false)
    }
  }, [collectionId, bookId])

  useEffect(() => {
    fetchBook()
  }, [fetchBook])

  return { book, loading, error, refetch: fetchBook }
}

export default useBooks

