import { createContext, useContext, useState, useEffect, useCallback } from 'react'

const AppContext = createContext()

export const useApp = () => {
  const context = useContext(AppContext)
  if (!context) {
    throw new Error('useApp must be used within AppProvider')
  }
  return context
}

export const AppProvider = ({ children }) => {
  // Dark mode state
  const [darkMode, setDarkMode] = useState(() => {
    const saved = localStorage.getItem('darkMode')
    return saved ? JSON.parse(saved) : false
  })

  // Font size state
  const [fontSize, setFontSize] = useState(() => {
    const saved = localStorage.getItem('fontSize')
    return saved ? parseInt(saved) : 22
  })

  // Show Arabic state
  const [showArabic, setShowArabic] = useState(() => {
    const saved = localStorage.getItem('showArabic')
    return saved ? JSON.parse(saved) : true
  })

  // Bookmarks state
  const [bookmarks, setBookmarks] = useState(() => {
    const saved = localStorage.getItem('bookmarks')
    return saved ? JSON.parse(saved) : []
  })

  // Daily hadith
  const [dailyHadith, setDailyHadith] = useState(null)

  // Update dark mode
  useEffect(() => {
    localStorage.setItem('darkMode', JSON.stringify(darkMode))
    if (darkMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [darkMode])

  // Update font size
  useEffect(() => {
    localStorage.setItem('fontSize', fontSize.toString())
  }, [fontSize])

  // Update showArabic
  useEffect(() => {
    localStorage.setItem('showArabic', JSON.stringify(showArabic))
  }, [showArabic])

  // Update bookmarks
  useEffect(() => {
    localStorage.setItem('bookmarks', JSON.stringify(bookmarks))
  }, [bookmarks])

  // Toggle dark mode
  const toggleDarkMode = useCallback(() => {
    setDarkMode(prev => !prev)
  }, [])

  // Increase font size
  const increaseFontSize = useCallback(() => {
    setFontSize(prev => Math.min(prev + 2, 32))
  }, [])

  // Decrease font size
  const decreaseFontSize = useCallback(() => {
    setFontSize(prev => Math.max(prev - 2, 12))
  }, [])

  // Toggle Arabic visibility
  const toggleArabic = useCallback(() => {
    setShowArabic(prev => !prev)
  }, [])

  // Add bookmark
  const addBookmark = useCallback((hadith) => {
    setBookmarks(prev => {
      const exists = prev.some(b => b.id === hadith.id)
      if (exists) return prev
      return [...prev, hadith]
    })
  }, [])

  // Remove bookmark
  const removeBookmark = useCallback((hadithId) => {
    setBookmarks(prev => prev.filter(b => b.id !== hadithId))
  }, [])

  // Check if hadith is bookmarked
  const isBookmarked = useCallback((hadithId) => {
    return bookmarks.some(b => b.id === hadithId)
  }, [bookmarks])

  // Toggle bookmark
  const toggleBookmark = useCallback((hadith) => {
    if (isBookmarked(hadith.id)) {
      removeBookmark(hadith.id)
    } else {
      addBookmark(hadith)
    }
  }, [isBookmarked, addBookmark, removeBookmark])

  const value = {
    darkMode,
    toggleDarkMode,
    fontSize,
    setFontSize,
    increaseFontSize,
    decreaseFontSize,
    showArabic,
    toggleArabic,
    bookmarks,
    addBookmark,
    removeBookmark,
    isBookmarked,
    toggleBookmark,
    dailyHadith,
    setDailyHadith
  }

  return (
    <AppContext.Provider value={value}>
      {children}
    </AppContext.Provider>
  )
}

export default AppContext

