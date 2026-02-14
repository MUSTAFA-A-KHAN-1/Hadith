# Hadith Portal - Technical Specification

## 1. Project Overview

**Project Name:** Hadith Portal  
**Type:** Single Page Application (SPA) - Islamic Hadith Reference Web Application  
**Core Functionality:** A premium, elegant Islamic Hadith portal that allows users to browse, search, and read Hadiths from various collections (Sahih al-Bukhari, Sahih Muslim, etc.) using the sunnah.com API.  
**Target Users:** Muslims seeking to read and study authentic Hadiths, researchers, students of Islamic studies.

---

## 2. UI/UX Specification

### 2.1 Layout Structure

#### Page Sections
- **Navbar:** Fixed/sticky top navigation with logo, navigation links, dark mode toggle, font size adjuster
- **Hero Section (Home):** Full-width immersive hero with search functionality
- **Content Areas:** Centered container with max-width 1280px
- **Footer:** Minimal footer with credits

#### Grid/Flex Layout Specifications
- **Container:** Max-width 1280px, centered, padding 24px (mobile), 48px (desktop)
- **Grid:** 12-column grid system
- **Cards:** CSS Grid with auto-fit, minmax(320px, 1fr) for responsive cards
- **Gap:** 24px between cards, 16px between elements

#### Responsive Breakpoints
- **Mobile:** < 640px (single column)
- **Tablet:** 640px - 1024px (2 columns)
- **Desktop:** > 1024px (3-4 columns)

### 2.2 Visual Design

#### Color Palette

**Light Mode:**
- Primary: `#065F46` (Emerald 800 - Deep Green)
- Primary Light: `#047857` (Emerald 700)
- Primary Dark: `#064E3B` (Emerald 900)
- Accent/Gold: `#D4A574` (Warm Gold)
- Accent Light: `#E8C9A0` (Light Gold)
- Background: `#FAF8F5` (Warm Beige)
- Background Card: `#FFFFFF`
- Text Primary: `#1F2937` (Gray 800)
- Text Secondary: `#6B7280` (Gray 500)
- Border: `#E5E0D8`

**Dark Mode:**
- Primary: `#10B981` (Emerald 500)
- Primary Light: `#34D399` (Emerald 400)
- Background: `#0F172A` (Slate 900)
- Background Card: `#1E293B` (Slate 800)
- Text Primary: `#F1F5F9` (Slate 100)
- Text Secondary: `#94A3B8` (Slate 400)
- Border: `#334155` (Slate 700)

#### Typography

**Arabic Font:** 
- Font Family: `'Amiri', 'Scheherazade', serif`
- Usage: Hadith Arabic text, chapter titles
- Size: 
  - Large: 28px (mobile), 36px (desktop)
  - Medium: 24px
  - Small: 18px
  - Base: 16px

**English Font:**
- Font Family: `'Inter', system-ui, sans-serif`
- Usage: UI elements, translations, metadata
- Sizes:
  - XL: 24px
  - LG: 18px
  - Base: 16px
  - SM: 14px
  - XS: 12px

**Font Weights:**
- Bold: 700
- Semibold: 600
- Medium: 500
- Regular: 400

#### Spacing System
- Base unit: 4px
- Spacing scale: 4, 8, 12, 16, 24, 32, 48, 64, 96px
- Card padding: 24px
- Section padding: 64px vertical
- Container padding: 24px (mobile), 48px (desktop)

#### Visual Effects
- **Card Shadow (Light):** `0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03)`
- **Card Shadow Hover:** `0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)`
- **Card Border Radius:** 16px
- **Button Border Radius:** 8px
- **Gold Divider:** 2px solid with gradient `#D4A574` to `#E8C9A0`
- **Backdrop Blur:** 12px for navbar

#### Animations

**Transitions:**
- Default: `all 0.3s cubic-bezier(0.4, 0, 0.2, 1)`
- Fast: `all 0.15s ease-out`
- Slow: `all 0.5s cubic-bezier(0.4, 0, 0.2, 1)`

**Keyframe Animations:**
- Fade In: `opacity 0 → 1, translateY 20px → 0` over 0.5s
- Slide In: `opacity 0 → 1, translateX -20px → 0` over 0.4s
- Pulse: Scale 1 → 1.05 → 1 for loading skeletons

**Staggered Animation Delays:**
- Cards: 50ms delay between each
- List items: 100ms delay between each

### 2.3 Components

#### Navbar
- Fixed position, backdrop blur
- Logo (text-based with Arabic calligraphy accent)
- Navigation links: Home, Collections, Search, Bookmarks
- Right side: Dark mode toggle, Font size controls
- Height: 72px
- Mobile: Hamburger menu

#### Hero Section
- Full viewport height minus navbar
- Centered content
- Large title with Arabic accent
- Subtitle description
- Prominent search bar
- Decorative Islamic geometric pattern (subtle SVG background)

#### Search Bar
- Large input with icon
- Placeholder: "Search hadiths..."
- Debounce: 300ms
- Loading state with skeleton
- Results dropdown or navigate to search page

#### Collection Card
- Collection name (English & Arabic)
- Book count badge
- Hadith count
- Subtle hover elevation
- Click to view collection

#### Hadith Card
- Arabic text (large, centered)
- Gold divider line
- English translation
- Narration chain (Isnad)
- Grade badge (Sahih, Hasan, etc.)
- Reference details
- Copy and Share buttons
- Bookmark icon

#### Buttons
- Primary: Emerald background, white text
- Secondary: Transparent, emerald border
- Ghost: No border, subtle hover
- Icon buttons: Circular, subtle background

#### Loading States
- Skeleton cards with pulse animation
- Spinner for inline loading
- Shimmer effect on cards

#### Error States
- Elegant error icon
- Clear error message
- Retry button

---

## 3. Functionality Specification

### 3.1 Core Features

#### A. Home Page
- Hero section with search
- Quick access to collections
- Daily Hadith feature (random hadith from API)
- Recent bookmarks section

#### B. Collections Page
- Fetch all collections from sunnah.com API
- Display in responsive grid
- Show collection info:
  - Collection name (English)
  - Number of books
  - Total hadiths
- Click to view collection books

#### C. Collection Detail Page (Books)
- List all books in a collection
- Show book number, name, hadith count
- Click to view hadiths

#### D. Hadith Display Page
- Fetch hadith by ID from API
- Display:
  - Arabic text (toggleable visibility)
  - English translation
  - Isnad (narration chain)
  - Grade (Sahih, Hasan, Da'if, etc.)
  - Reference (collection, book, hadith number)
- Copy to clipboard functionality
- Share functionality (Web Share API)
- Bookmark toggle
- Font size adjustment (affects Arabic and English)
- Previous/Next navigation

#### E. Search
- Real-time search with debounce (300ms)
- Search by keyword in hadith text
- Display results in list/grid
- Pagination or infinite scroll
- No results state
- Loading skeleton

#### F. Bookmarks
- Save hadith references to localStorage
- Display bookmarked hadiths
- Remove bookmarks
- Persist across sessions

### 3.2 API Integration

#### Base URL
`https://api.sunnah.com/v1`

#### Endpoints
- `GET /collections` - List all collections
- `GET /collections/:id/books` - Get books in collection
- `GET /collections/:collectionId/books/:bookId/hadiths` - Get hadiths in book
- `GET /hadiths/:id` - Get specific hadith
- `GET /books/:id/hadiths` - Get hadiths from book
- `GET /search` - Search hadiths

#### Authentication
- API Key in header: `x-api-key`
- API Key: Stored in `.env` file as `VITE_SUNNAH_API_KEY`

### 3.3 Data Handling
- Custom React hooks for API calls
- React Query pattern (or custom useEffect with loading/error states)
- Proper error handling and retry logic
- Response caching where appropriate
- Debounced search input

### 3.4 Edge Cases
- No internet connection: Show offline message
- API rate limiting: Handle gracefully with retry
- Empty search results: Show elegant empty state
- Very long hadith text: Proper text wrapping
- Missing translations: Show fallback message
- Invalid hadith ID: Show 404 page

---

## 4. Technical Architecture

### 4.1 Folder Structure

```
src/
├── components/
│   ├── common/
│   │   ├── Button.jsx
│   │   ├── Card.jsx
│   │   ├── Loading.jsx
│   │   ├── Error.jsx
│   │   ├── Skeleton.jsx
│   │   ├── Badge.jsx
│   │   └── ScrollToTop.jsx
│   ├── layout/
│   │   ├── Navbar.jsx
│   │   ├── Footer.jsx
│   │   └── Layout.jsx
│   ├── hadith/
│   │   ├── HadithCard.jsx
│   │   ├── HadithDisplay.jsx
│   │   ├── GradeBadge.jsx
│   │   └── HadithActions.jsx
│   ├── collection/
│   │   ├── CollectionCard.jsx
│   │   └── BookCard.jsx
│   └── search/
│       ├── SearchBar.jsx
│       └── SearchResults.jsx
├── pages/
│   ├── Home.jsx
│   ├── Collections.jsx
│   ├── CollectionDetail.jsx
│   ├── BookDetail.jsx
│   ├── HadithDetail.jsx
│   ├── Search.jsx
│   ├── Bookmarks.jsx
│   └── NotFound.jsx
├── hooks/
│   ├── useCollections.js
│   ├── useBooks.js
│   ├── useHadiths.js
│   ├── useSearch.js
│   ├── useBookmarks.js
│   └── useDarkMode.js
├── services/
│   └── api.js
├── utils/
│   ├── constants.js
│   ├── helpers.js
│   └── localStorage.js
├── context/
│   └── AppContext.jsx
├── assets/
│   └── fonts/
├── App.jsx
├── main.jsx
└── index.css
```

### 4.2 Environment Variables
```
VITE_SUNNAH_API_KEY=your_api_key_here
```

---

## 5. Acceptance Criteria

### Visual Checkpoints
- [ ] Navbar is sticky and has backdrop blur
- [ ] Hero section has Islamic geometric pattern background
- [ ] Color scheme matches specified emerald/gold/beige palette
- [ ] Arabic font (Amiri) displays correctly for hadith text
- [ ] English font (Inter) displays correctly for UI
- [ ] Dark mode toggle works and persists
- [ ] Cards have proper shadows and hover effects
- [ ] Gold divider line appears between Arabic and English
- [ ] Animations are smooth and not jarring
- [ ] Mobile layout is single column and usable

### Functional Checkpoints
- [ ] Collections load from API and display in grid
- [ ] Clicking collection shows books
- [ ] Clicking book shows hadiths
- [ ] Hadith detail page displays all required information
- [ ] Search returns results with debounce
- [ ] Bookmarks save to localStorage
- [ ] Font size adjuster works
- [ ] Copy button copies hadith text
- [ ] Share button uses Web Share API
- [ ] 404 page displays for invalid routes
- [ ] Loading states show skeletons
- [ ] Error states show retry option

### Performance Checkpoints
- [ ] Initial page load < 3 seconds
- [ ] Search debounce prevents excessive API calls
- [ ] No memory leaks from subscriptions
- [ ] Responsive on all breakpoints

