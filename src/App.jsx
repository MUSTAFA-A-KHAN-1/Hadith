import { Routes, Route } from 'react-router-dom'
import Layout from './components/layout/Layout'
import Home from './pages/Home'
import Collections from './pages/Collections'
import CollectionDetail from './pages/CollectionDetail'
import BookDetail from './pages/BookDetail'
import HadithDetail from './pages/HadithDetail'
import Search from './pages/Search'
import Bookmarks from './pages/Bookmarks'
import NotFound from './pages/NotFound'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="collections" element={<Collections />} />
        <Route path="collections/:collectionId" element={<CollectionDetail />} />
        <Route path="collections/:collectionId/books/:bookNumber" element={<BookDetail />} />
        <Route path="collections/:collectionId/books/:bookNumber/hadith/:hadithNumber" element={<HadithDetail />} />
        <Route path="search" element={<Search />} />
        <Route path="bookmarks" element={<Bookmarks />} />
        <Route path="*" element={<NotFound />} />
      </Route>
    </Routes>
  )
}

export default App

