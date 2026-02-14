import { Link } from 'react-router-dom'
import Card from '../common/Card'
import Badge from '../common/Badge'
import { getCollectionDisplayName, formatNumber } from '../../utils/constants'

const CollectionCard = ({ collection, index = 0 }) => {
  const name = getCollectionDisplayName(collection)
  const bookCount = collection.bookCount || collection.books?.length || 0
  const hadithCount = collection.hadithCount || collection.hadithsCount || 0

  // Get Arabic name based on collection
  const getArabicName = (collectionName) => {
    const arabicNames = {
      'bukhari': 'صحيح البخاري',
      'muslim': 'صحيح مسلم',
      'abudawud': 'سنن أبي داود',
      'tirmidhi': 'جامع الترمذي',
      'nasai': 'سنن النسائي',
      'ibnmajah': 'سنن ابن ماجه',
      'muwatta': 'موطأ مالك',
      'riyadussaliheen': 'رياض الصالحين',
      'adab': ' الأدب المفرد',
      'shamaa-il': 'شمائل الترمذي',
      'mishkat': 'مشكاة المصابيح'
    }
    return arabicNames[collectionName?.toLowerCase()] || ''
  }

  return (
    <Link to={`/collections/${collection.name}`}>
      <Card 
        className="h-full animate-fade-in"
        style={{ animationDelay: `${index * 50}ms` }}
      >
        <div className="flex items-start gap-4">
          <div className="w-14 h-14 rounded-full bg-primary/10 dark:bg-primary-500/20 flex items-center justify-center flex-shrink-0">
            <span className="text-2xl">📚</span>
          </div>
          <div className="flex-1 min-w-0">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-1 truncate">
              {name}
            </h3>
            {getArabicName(collection.name) && (
              <p className="font-arabic text-lg text-gray-600 dark:text-gray-400 mb-2 truncate">
                {getArabicName(collection.name)}
              </p>
            )}
            <div className="flex flex-wrap gap-2 mt-3">
              <Badge variant="primary" size="sm">
                {formatNumber(bookCount)} Books
              </Badge>
              <Badge variant="accent" size="sm">
                {formatNumber(hadithCount)} Hadiths
              </Badge>
            </div>
          </div>
        </div>
      </Card>
    </Link>
  )
}

export default CollectionCard

