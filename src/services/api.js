import axios from 'axios'

const API_BASE_URL = 'https://api.sunnah.com/v1'
const API_KEY = import.meta.env.VITE_SUNNAH_API_KEY || 'your_api_key_here'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
    'x-api-key': API_KEY
  }
})

// Add response interceptor for error handling
api.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error.response?.data || error.message)
    return Promise.reject(error)
  }
)

// Collections
export const getCollections = async () => {
  try {
    const response = await api.get('/collections')
    return response.data
  } catch (error) {
    console.error('Error fetching collections:', error)
    throw error
  }
}

// Get collection by ID
export const getCollection = async (collectionId) => {
  try {
    const response = await api.get(`/collections/${collectionId}`)
    return response.data
  } catch (error) {
    console.error('Error fetching collection:', error)
    throw error
  }
}

// Get books in a collection
export const getBooks = async (collectionId) => {
  try {
    const response = await api.get(`/collections/${collectionId}/books`)
    return response.data
  } catch (error) {
    console.error('Error fetching books:', error)
    throw error
  }
}

// Get book by ID
export const getBook = async (collectionId, bookId) => {
  try {
    const response = await api.get(`/collections/${collectionId}/books/${bookId}`)
    return response.data
  } catch (error) {
    console.error('Error fetching book:', error)
    throw error
  }
}

// Get hadiths in a book
export const getHadiths = async (collectionId, bookId, page = 1, limit = 20) => {
  try {
    const response = await api.get(`/collections/${collectionId}/books/${bookId}/hadiths`, {
      params: { page, limit }
    })
    return response.data
  } catch (error) {
    console.error('Error fetching hadiths:', error)
    throw error
  }
}

// Get specific hadith
export const getSingleHadith = async (hadithId) => {
  try {
    const response = await api.get(`/hadiths/${hadithId}`)
    return response.data
  } catch (error) {
    console.error('Error fetching hadith:', error)
    throw error
  }
}

// Get specific hadith from a book
export const getHadithFromBook = async (collectionId, bookId, hadithNumber) => {
  try {
    const response = await api.get(`/collections/${collectionId}/books/${bookId}/hadiths`, {
      params: { page: 1, limit: 100 }
    })
    const hadiths = response.data.hadiths || response.data
    return hadiths.find(h => 
      h.hadithNumber === parseInt(hadithNumber) || 
      h.hadith === parseInt(hadithNumber)
    ) || null
  } catch (error) {
    console.error('Error fetching hadith from book:', error)
    throw error
  }
}

// Search hadiths
export const searchHadiths = async (query, page = 1, limit = 20) => {
  try {
    const response = await api.get('/search', {
      params: { query, page, limit }
    })
    return response.data
  } catch (error) {
    console.error('Error searching hadiths:', error)
    throw error
  }
}

// Get random hadith (for daily hadith)
export const getRandomHadith = async () => {
  try {
    // First get collections
    const collections = await getCollections()
    if (!collections || collections.length === 0) return null

    // Pick a random collection (prioritize major ones)
    const majorCollections = collections.filter(c => 
      ['bukhari', 'muslim', 'abudawud', 'tirmidhi', 'nasai', 'ibnmajah'].includes(c.name?.toLowerCase())
    )
    const collectionList = majorCollections.length > 0 ? majorCollections : collections
    const randomCollection = collectionList[Math.floor(Math.random() * collectionList.length)]

    // Get books in the collection
    const books = await getBooks(randomCollection.name)
    if (!books || books.length === 0) return null

    // Pick a random book
    const randomBook = books[Math.floor(Math.random() * books.length)]

    // Get hadiths from the book
    const hadithsData = await getHadiths(randomCollection.name, randomBook.bookNumber, 1, 50)
    if (!hadithsData || !hadithsData.hadiths || hadithsData.hadiths.length === 0) return null

    // Pick a random hadith
    const randomHadith = hadithsData.hadiths[Math.floor(Math.random() * hadithsData.hadiths.length)]

    return {
      hadith: randomHadith,
      collection: randomCollection,
      book: randomBook
    }
  } catch (error) {
    console.error('Error fetching random hadith:', error)
    return null
  }
}

export default api

