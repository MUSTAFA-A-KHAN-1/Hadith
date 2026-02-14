// Copy text to clipboard
export const copyToClipboard = async (text) => {
  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(text)
      return true
    }
    // Fallback for older browsers
    const textArea = document.createElement('textarea')
    textArea.value = text
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    textArea.style.top = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()
    const result = document.execCommand('copy')
    document.body.removeChild(textArea)
    return result
  } catch (error) {
    console.error('Failed to copy:', error)
    return false
  }
}

// Share content using Web Share API
export const shareContent = async (title, text, url) => {
  try {
    if (navigator.share && navigator.canShare) {
      const shareData = { title, text, url }
      if (navigator.canShare(shareData)) {
        await navigator.share(shareData)
        return true
      }
    }
    // Fallback: copy to clipboard
    return await copyToClipboard(`${title}\n\n${text}\n\n${url}`)
  } catch (error) {
    // User cancelled or error
    if (error.name !== 'AbortError') {
      console.error('Failed to share:', error)
    }
    return false
  }
}

// Format date
export const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// Truncate text
export const truncateText = (text, maxLength = 150) => {
  if (!text || text.length <= maxLength) return text
  return text.substring(0, maxLength).trim() + '...'
}

// Validate hadith ID format
export const isValidHadithId = (id) => {
  return /^\d+$/.test(id)
}

// Generate hadith ID from params
export const generateHadithId = (collection, bookNumber, hadithNumber) => {
  return `${collection}-${bookNumber}-${hadithNumber}`
}

// Parse hadith ID
export const parseHadithId = (id) => {
  const parts = id.split('-')
  if (parts.length === 3) {
    return {
      collection: parts[0],
      bookNumber: parts[1],
      hadithNumber: parts[2]
    }
  }
  return null
}

// Scroll to top
export const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

// Get arabic text from hadith
export const getArabicText = (hadith) => {
  return hadith?.arabic || hadith?.text || hadith?.hadith || ''
}

// Get english text from hadith
export const getEnglishText = (hadith) => {
  return hadith?.english || hadith?.translation || hadith?.text || ''
}

// Get narrator chain
export const getNarratorChain = (hadith) => {
  return hadith?.gradings?.[0]?.narrator || hadith?.attribution || ''
}

// Get grade
export const getGrade = (hadith) => {
  if (hadith?.grades?.[0]?.grade) {
    return hadith.grades[0].grade
  }
  if (hadith?.grade) {
    return hadith.grade
  }
  return 'Sahih' // Default to Sahih for major collections
}

