import { useState, useEffect, useCallback } from 'react'
import { getCollections, getCollection } from '../services/api'

export const useCollections = () => {
  const [collections, setCollections] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const fetchCollections = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await getCollections()
      setCollections(data)
    } catch (err) {
      setError(err.message || 'Failed to fetch collections')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchCollections()
  }, [fetchCollections])

  return { collections, loading, error, refetch: fetchCollections }
}

export const useCollection = (collectionId) => {
  const [collection, setCollection] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const fetchCollection = useCallback(async () => {
    if (!collectionId) return
    try {
      setLoading(true)
      setError(null)
      const data = await getCollection(collectionId)
      setCollection(data)
    } catch (err) {
      setError(err.message || 'Failed to fetch collection')
    } finally {
      setLoading(false)
    }
  }, [collectionId])

  useEffect(() => {
    fetchCollection()
  }, [fetchCollection])

  return { collection, loading, error, refetch: fetchCollection }
}

export default useCollections

