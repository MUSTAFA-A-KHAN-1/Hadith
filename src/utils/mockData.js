// Mock data for the Hadith Portal
// This is used when the external API is unavailable
import bukhariData from '../data/bukhari.json'
import muslimData from '../data/muslim.json'
import abudawudData from '../data/abudawud.json'
import tirmidhiData from '../data/tirmidhi.json'
import nasaiData from '../data/nasai.json'
import ibnmajahData from '../data/ibnmajah.json'
import darimiData from '../data/darimi.json'

// Transform hadith data from JSON files to the format used by the app
const transformHadiths = (data, defaultGrade = 'Sahih') => {
  const hadithArray = data.hadiths || []
  return hadithArray.map(h => ({
    hadithNumber: h.idInBook || h.id,
    grade: h.grade || defaultGrade,
    arabic: h.arabic,
    english: typeof h.english === 'string' ? h.english : (h.english?.text || h.english?.narrator || ''),
    chapterId: h.chapterId,
    bookId: h.bookId
  }))
}

// Transform book data from JSON files
const transformBooks = (data) => {
  const chaptersArray = data.chapters || []
  return chaptersArray.map(ch => ({
    bookNumber: ch.id,
    title: ch.english || ch.arabic || '',
    english: ch.english || '',
    arabic: ch.arabic || '',
    hadiths: 0 // Will be calculated when loading hadiths
  }))
}

// Transform collection metadata from JSON files
const transformCollection = (data, name, defaultGrade = 'Sahih') => {
  return {
    name,
    title: data.metadata?.english?.title || data.metadata?.arabic?.title || data.metadata?.english || data.metadata?.arabic || '',
    author: data.metadata?.english?.author || data.metadata?.arabic?.author || data.metadata?.english || data.metadata?.arabic || '',
    hadiths: data.hadiths?.length || 0,
    books: data.chapters?.length || 0,
    description: data.metadata?.english?.title || data.metadata?.arabic?.title || '',
    grade: defaultGrade
  }
}

// Pre-transform all data
const bukhariHadiths = transformHadiths(bukhariData, 'Sahih')
const bukhariBooks = transformBooks(bukhariData)
const bukhariCollection = transformCollection(bukhariData, 'bukhari', 'Sahih')

const muslimHadiths = transformHadiths(muslimData, 'Sahih')
const muslimBooks = transformBooks(muslimData)
const muslimCollection = transformCollection(muslimData, 'muslim', 'Sahih')

const abudawudHadiths = transformHadiths(abudawudData, 'Sahih')
const abudawudBooks = transformBooks(abudawudData)
const abudawudCollection = transformCollection(abudawudData, 'abudawud', 'Sahih')

const tirmidhiHadiths = transformHadiths(tirmidhiData, 'Sahih')
const tirmidhiBooks = transformBooks(tirmidhiData)
const tirmidhiCollection = transformCollection(tirmidhiData, 'tirmidhi', 'Sahih')

const nasaiHadiths = transformHadiths(nasaiData, 'Sahih')
const nasaiBooks = transformBooks(nasaiData)
const nasaiCollection = transformCollection(nasaiData, 'nasai', 'Sahih')

const ibnmajahHadiths = transformHadiths(ibnmajahData, 'Sahih')
const ibnmajahBooks = transformBooks(ibnmajahData)
const ibnmajahCollection = transformCollection(ibnmajahData, 'ibnmajah', 'Sahih')

const darimiHadiths = transformHadiths(darimiData, 'Sahih')
const darimiBooks = transformBooks(darimiData)
const darimiCollection = transformCollection(darimiData, 'darimi', 'Sahih')

export const mockCollections = [
  bukhariCollection,
  muslimCollection,
  abudawudCollection,
  tirmidhiCollection,
  nasaiCollection,
  ibnmajahCollection,
  darimiCollection,
  {
    name: 'muwatta',
    title: 'Muwatta Imam Malik',
    hadiths: 1900,
    books: 61,
    description: 'The earliest collection of hadith'
  },
  {
    name: 'riyadussaliheen',
    title: 'Riyad as-Salihin',
    hadiths: 1906,
    books: 19,
    description: 'Gardens of the Righteous'
  },
  {
    name: 'adab',
    title: 'Al-Adab al-Mufrad',
    hadiths: 1321,
    books: 16,
    description: 'The Book of Manners'
  },
  {
    name: 'shamaa-il',
    title: "Shamaa'il Tirmidhi",
    hadiths: 395,
    books: 1,
    description: 'The Beauty of the Prophet Muhammad'
  },
  {
    name: 'mishkat',
    title: 'Mishkat al-Masabih',
    hadiths: 6000,
    books: 29,
    description: 'The Hollow Place of Light'
  }
]

export const mockBooks = {
  bukhari: bukhariBooks,
  muslim: muslimBooks,
  abudawud: abudawudBooks,
  tirmidhi: tirmidhiBooks,
  nasai: nasaiBooks,
  ibnmajah: ibnmajahBooks,
  darimi: darimiBooks
}

export const mockHadiths = {
  bukhari: bukhariHadiths,
  muslim: muslimHadiths,
  abudawud: abudawudHadiths,
  tirmidhi: tirmidhiHadiths,
  nasai: nasaiHadiths,
  ibnmajah: ibnmajahHadiths,
  darimi: darimiHadiths
}

// Helper function to get mock data
export const getMockCollections = () => Promise.resolve(mockCollections)

export const getMockBooks = (collection) => {
  const books = mockBooks[collection] || mockBooks.bukhari
  return Promise.resolve(books)
}

export const getMockHadiths = (collection, chapterId, page = 1, limit = 20) => {
  let collectionHadiths = mockHadiths[collection] || mockHadiths.bukhari
  
  // Filter hadiths by chapterId (which is the book/chapter number)
  if (chapterId) {
    collectionHadiths = collectionHadiths.filter(h => h.chapterId === parseInt(chapterId))
  }
  
  // If no filtered results, return all from collection
  if (collectionHadiths.length === 0) {
    collectionHadiths = mockHadiths[collection] || mockHadiths.bukhari
  }
  
  // Apply pagination
  const startIndex = (page - 1) * limit
  const endIndex = startIndex + limit
  const bookHadiths = collectionHadiths.slice(startIndex, endIndex)
  
  return Promise.resolve({
    hadiths: bookHadiths,
    total: collectionHadiths.length,
    page,
    totalPages: Math.ceil(collectionHadiths.length / limit)
  })
}

export const getMockRandomHadith = () => {
  const collections = Object.keys(mockHadiths)
  const randomCollection = collections[Math.floor(Math.random() * collections.length)]
  const hadiths = mockHadiths[randomCollection]
  const randomHadith = hadiths[Math.floor(Math.random() * hadiths.length)]
  
  const collection = mockCollections.find(c => c.name === randomCollection)
  const books = mockBooks[randomCollection] || mockBooks.bukhari
  const book = books[0]
  
  return Promise.resolve({
    hadith: randomHadith,
    collection,
    book
  })
}

