// Local storage keys
const STORAGE_KEYS = {
  BOOKMARKS: 'hadith_bookmarks',
  DARK_MODE: 'dark_mode',
  FONT_SIZE: 'font_size',
  SHOW_ARABIC: 'show_arabic',
  DAILY_HADITH_DATE: 'daily_hadith_date',
  DAILY_HADITH: 'daily_hadith'
}

// Get item from local storage
export const getItem = (key, defaultValue = null) => {
  try {
    const item = localStorage.getItem(key)
    return item ? JSON.parse(item) : defaultValue
  } catch (error) {
    console.error('Error reading from localStorage:', error)
    return defaultValue
  }
}

// Set item in local storage
export const setItem = (key, value) => {
  try {
    localStorage.setItem(key, JSON.stringify(value))
    return true
  } catch (error) {
    console.error('Error writing to localStorage:', error)
    return false
  }
}

// Remove item from local storage
export const removeItem = (key) => {
  try {
    localStorage.removeItem(key)
    return true
  } catch (error) {
    console.error('Error removing from localStorage:', error)
    return false
  }
}

// Get bookmarks
export const getBookmarks = () => {
  return getItem(STORAGE_KEYS.BOOKMARKS, [])
}

// Save bookmark
export const saveBookmark = (hadith) => {
  const bookmarks = getBookmarks()
  const exists = bookmarks.some(b => b.id === hadith.id)
  if (!exists) {
    bookmarks.push(hadith)
    setItem(STORAGE_KEYS.BOOKMARKS, bookmarks)
  }
  return bookmarks
}

// Remove bookmark
export const removeBookmark = (hadithId) => {
  const bookmarks = getBookmarks()
  const filtered = bookmarks.filter(b => b.id !== hadithId)
  setItem(STORAGE_KEYS.BOOKMARKS, filtered)
  return filtered
}

// Check if hadith is bookmarked
export const isBookmarked = (hadithId) => {
  const bookmarks = getBookmarks()
  return bookmarks.some(b => b.id === hadithId)
}

// Get daily hadith (cached for a day)
export const getDailyHadith = () => {
  const today = new Date().toDateString()
  const savedDate = getItem(STORAGE_KEYS.DAILY_HADITH_DATE)
  
  if (savedDate !== today) {
    // Clear old daily hadith
    removeItem(STORAGE_KEYS.DAILY_HADITH)
    return null
  }
  
  return getItem(STORAGE_KEYS.DAILY_HADITH)
}

// Save daily hadith
export const saveDailyHadith = (hadith) => {
  const today = new Date().toDateString()
  setItem(STORAGE_KEYS.DAILY_HADITH_DATE, today)
  setItem(STORAGE_KEYS.DAILY_HADITH, hadith)
}

export { STORAGE_KEYS }

