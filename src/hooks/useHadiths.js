import { useState, useEffect, useCallback } from 'react'
import { getHadiths, getHadithFromBook } from '../services/api'

export const useHadiths = (collectionId, bookId, page = 1, limit = 20) => {
  const [hadiths, setHadiths] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [total, setTotal] = useState(0)

  const fetchHadiths = useCallback(async () => {
    if (!collectionId || !bookId) return
    try {
      setLoading(true)
      setError(null)
      const data = await getHadiths(collectionId, bookId, page, limit)
      
      // Handle different response formats
      if (Array.isArray(data)) {
        setHadiths(data)
        setTotal(data.length)
      } else if (data && data.hadiths) {
        setHadiths(data.hadiths)
        setTotal(data.hadiths.length)
      } else {
        setHadiths([])
        setTotal(0)
      }
    } catch (err) {
      setError(err.message || 'Failed to fetch hadiths')
      setHadiths([])
    } finally {
      setLoading(false)
    }
  }, [collectionId, bookId, page, limit])

  useEffect(() => {
    fetchHadiths()
  }, [fetchHadiths])

  return { hadiths, loading, error, total, refetch: fetchHadiths }
}

export const useSingleHadith = (collectionId, bookId, hadithNumber) => {
  const [hadith, setHadith] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const fetchHadith = useCallback(async () => {
    if (!collectionId || !bookId || !hadithNumber) return
    try {
      setLoading(true)
      setError(null)
      
      // Fetch hadiths from the book and find the specific one
      const data = await getHadiths(collectionId, bookId, 1, 100)
      
      // Handle different response formats
      let hadithList = []
      if (Array.isArray(data)) {
        hadithList = data
      } else if (data && data.hadiths) {
        hadithList = data.hadiths
      }
      
      // Find the specific hadith
      const foundHadith = hadithList.find(h => 
        h.hadithNumber === parseInt(hadithNumber) || 
        h.hadith === parseInt(hadithNumber) ||
        h.id === parseInt(hadithNumber)
      )
      
      if (foundHadith) {
        setHadith(foundHadith)
      } else {
      // If not found in the list, try to get it directly
        const directData = await getHadithFromBook(collectionId, bookId, hadithNumber)
        setHadith(directData)
      }
    } catch (err) {
      setError(err.message || 'Failed to fetch hadith')
      setHadith(null)
    } finally {
      setLoading(false)
    }
  }, [collectionId, bookId, hadithNumber])

  useEffect(() => {
    fetchHadith()
  }, [fetchHadith])

  return { hadith, loading, error, refetch: fetchHadith }
}

export default useHadiths

