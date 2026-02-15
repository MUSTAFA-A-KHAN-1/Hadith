import axios from 'axios'
import { getMockCollections, getMockBooks, getMockHadiths, getMockRandomHadith, mockCollections } from '../utils/mockData'

const API_BASE_URL = 'https://api.sunnah.com/v1'
const API_KEY = import.meta.env.VITE_SUNNAH_API_KEY || 'your_api_key_here'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
    'x-api-key': API_KEY
  },
  timeout: 5000
})

// Add response interceptor for error handling
api.interceptors.response.use(
  response => response,
  error => {
    console.warn('API unavailable, using mock data:', error.message)
    return Promise.reject(error)
  }
)

// Helper to check if API is available
const isApiAvailable = async () => {
  try {
    await api.get('/collections', { timeout: 3000 })
    return true
  } catch {
    return false
  }
}

// Collections
export const getCollections = async () => {
  try {
    const available = await isApiAvailable()
    if (!available) {
      return getMockCollections()
    }
    const response = await api.get('/collections')
    return response.data
  } catch (error) {
    console.warn('Using mock data for collections')
    return getMockCollections()
  }
}

// Get collection by ID
export const getCollection = async (collectionId) => {
  try {
    const available = await isApiAvailable()
    if (!available) {
      const collection = mockCollections.find(c => c.name === collectionId)
      return collection || null
    }
    const response = await api.get(`/collections/${collectionId}`)
    return response.data
  } catch (error) {
    console.warn('Using mock data for collection')
    const collection = mockCollections.find(c => c.name === collectionId)
    return collection || null
  }
}

// Get books in a collection
export const getBooks = async (collectionId) => {
  try {
    const available = await isApiAvailable()
    if (!available) {
      return getMockBooks(collectionId)
    }
    const response = await api.get(`/collections/${collectionId}/books`)
    return response.data
  } catch (error) {
    console.warn('Using mock data for books')
    return getMockBooks(collectionId)
  }
}

// Get book by ID
export const getBook = async (collectionId, bookId) => {
  try {
    const available = await isApiAvailable()
    if (!available) {
      const books = await getMockBooks(collectionId)
      return books.find(b => b.bookNumber === parseInt(bookId)) || null
    }
    const response = await api.get(`/collections/${collectionId}/books/${bookId}`)
    return response.data
  } catch (error) {
    console.warn('Using mock data for book')
    const books = await getMockBooks(collectionId)
    return books.find(b => b.bookNumber === parseInt(bookId)) || null
  }
}

// Get hadiths in a book
export const getHadiths = async (collectionId, bookId, page = 1, limit = 20) => {
  try {
    const available = await isApiAvailable()
    if (!available) {
      return getMockHadiths(collectionId, bookId, page, limit)
    }
    const response = await api.get(`/collections/${collectionId}/books/${bookId}/hadiths`, {
      params: { page, limit }
    })
    return response.data
  } catch (error) {
    console.warn('Using mock data for hadiths')
    return getMockHadiths(collectionId, bookId, page, limit)
  }
}

// Get specific hadith
export const getSingleHadith = async (hadithId) => {
  try {
    const available = await isApiAvailable()
    if (!available) {
      return null
    }
    const response = await api.get(`/hadiths/${hadithId}`)
    return response.data
  } catch (error) {
    console.warn('Using mock data for hadith')
    return null
  }
}

// Get specific hadith from a book
export const getHadithFromBook = async (collectionId, bookId, hadithNumber) => {
  try {
    const available = await isApiAvailable()
    if (!available) {
      const data = await getMockHadiths(collectionId, bookId)
      return data.hadiths.find(h => 
        h.hadithNumber === parseInt(hadithNumber) || 
        h.hadith === parseInt(hadithNumber)
      ) || null
    }
    const response = await api.get(`/collections/${collectionId}/books/${bookId}/hadiths`, {
      params: { page: 1, limit: 100 }
    })
    const hadiths = response.data.hadiths || response.data
    return hadiths.find(h => 
      h.hadithNumber === parseInt(hadithNumber) || 
      h.hadith === parseInt(hadithNumber)
    ) || null
  } catch (error) {
    console.warn('Using mock data for hadith from book')
    const data = await getMockHadiths(collectionId, bookId)
    return data.hadiths.find(h => 
      h.hadithNumber === parseInt(hadithNumber) || 
      h.hadith === parseInt(hadithNumber)
    ) || null
  }
}

// Search hadiths
export const searchHadiths = async (query, page = 1, limit = 20) => {
  try {
    const available = await isApiAvailable()
    if (!available) {
      // Simple mock search
      const allHadiths = Object.values(mockCollections).flat ? [] : []
      // Import mock hadiths for search
      const { mockHadiths } = await import('../utils/mockData')
      const allHadithsData = Object.values(mockHadiths).flat()
      const filtered = allHadithsData.filter(h => 
        h.english?.toLowerCase().includes(query.toLowerCase()) ||
        h.arabic?.includes(query)
      ).slice(0, limit)
      return { hadiths: filtered, total: filtered.length, page, totalPages: 1 }
    }
    const response = await api.get('/search', {
      params: { query, page, limit }
    })
    return response.data
  } catch (error) {
    console.warn('Using mock data for search')
    const { mockHadiths } = await import('../utils/mockData')
    const allHadiths = Object.values(mockHadiths).flat()
    const filtered = allHadiths.filter(h => 
      h.english?.toLowerCase().includes(query.toLowerCase()) ||
      h.arabic?.includes(query)
    ).slice(0, limit)
    return { hadiths: filtered, total: filtered.length, page, totalPages: 1 }
  }
}

// Get random hadith (for daily hadith)
export const getRandomHadith = async () => {
  try {
    const available = await isApiAvailable()
    if (!available) {
      return getMockRandomHadith()
    }
    
    // First get collections
    const collections = await getCollections()
    if (!collections || collections.length === 0) return getMockRandomHadith()

    // Pick a random collection (prioritize major ones)
    const majorCollections = collections.filter(c => 
      ['bukhari', 'muslim', 'abudawud', 'tirmidhi', 'nasai', 'ibnmajah'].includes(c.name?.toLowerCase())
    )
    const collectionList = majorCollections.length > 0 ? majorCollections : collections
    const randomCollection = collectionList[Math.floor(Math.random() * collectionList.length)]

    // Get books in the collection
    const books = await getBooks(randomCollection.name)
    if (!books || books.length === 0) return getMockRandomHadith()

    // Pick a random book
    const randomBook = books[Math.floor(Math.random() * books.length)]

    // Get hadiths from the book
    const hadithsData = await getHadiths(randomCollection.name, randomBook.bookNumber, 1, 50)
    if (!hadithsData || !hadithsData.hadiths || hadithsData.hadiths.length === 0) return getMockRandomHadith()

    // Pick a random hadith
    const randomHadith = hadithsData.hadiths[Math.floor(Math.random() * hadithsData.hadiths.length)]

    return {
      hadith: randomHadith,
      collection: randomCollection,
      book: randomBook
    }
  } catch (error) {
    console.warn('Using mock data for random hadith')
    return getMockRandomHadith()
  }
}

export default api

