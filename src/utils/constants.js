// Collection names mapping
export const COLLECTION_NAMES = {
  bukhari: 'Sahih al-Bukhari',
  muslim: 'Sahih Muslim',
  abudawud: 'Sunan Abu Dawood',
  tirmidhi: 'Jamiʿ at-Tirmidhi',
  nasai: 'Sunan an-Nasai',
  ibnmajah: 'Sunan Ibn Majah',
  muwatta: 'Muwatta Imam Malik',
  riyadussaliheen: 'Riyad as-Salihin',
  adab: 'Al-Adab al-Mufrad',
  'shamaa-il': 'Shamaa\'il Tirmidhi',
  mishkat: 'Mishkat al-Masabih'
}

// Grade colors
export const GRADE_COLORS = {
  'Sahih': 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
  'Hasan': 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
  'Daif': 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
  'Mawdu': 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
  'Munkar': 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200',
  'Mudallas': 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
  ' Hasan li-ghayrihi': 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900 dark:text-cyan-200',
  'Daif li-ghayrihi': 'bg-pink-100 text-pink-800 dark:bg-pink-900 dark:text-pink-200'
}

// Default grade if not found
export const DEFAULT_GRADE_COLOR = 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200'

// Get grade color
export const getGradeColor = (grade) => {
  if (!grade) return DEFAULT_GRADE_COLOR
  
  // Check exact match first
  if (GRADE_COLORS[grade]) {
    return GRADE_COLORS[grade]
  }
  
  // Check partial match
  for (const [key, value] of Object.entries(GRADE_COLORS)) {
    if (grade.includes(key)) {
      return value
    }
  }
  
  return DEFAULT_GRADE_COLOR
}

// Get collection display name
export const getCollectionDisplayName = (collection) => {
  if (!collection) return ''
  
  // If it's already a string, treat it as collectionId
  if (typeof collection === 'string') {
    return COLLECTION_NAMES[collection.toLowerCase()] || collection
  }
  
  // If it's an object with name property
  if (collection.name && COLLECTION_NAMES[collection.name.toLowerCase()]) {
    return COLLECTION_NAMES[collection.name.toLowerCase()]
  }
  
  return collection.name || collection.title || ''
}

// Format hadith reference
export const formatHadithReference = (collection, bookNumber, hadithNumber) => {
  const collectionName = getCollectionDisplayName(collection)
  return `${collectionName}, Book ${bookNumber}, Hadith ${hadithNumber}`
}

// Debounce function
export const debounce = (func, wait) => {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      func(...args)
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

// Format number with commas
export const formatNumber = (num) => {
  if (!num) return '0'
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

