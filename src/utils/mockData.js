// Mock data for the Hadith Portal
// This is used when the external API is unavailable
import bukhariData from  '../data/bukhari.json'
console.log(bukhariData);
console.log("KEYS:", Object.keys(bukhariData))
const hadithArray = bukhariData.hadiths || []
const transformedBukhari = hadithArray.map(h => ({
  hadithNumber: h.id,
  grade: 'Sahih', // all Bukhari hadith are Sahih
  arabic: h.arabic,
  english: typeof h.english === 'string' ? h.english : h.english?.text || '',
  bookNumber: h.bookId
}))

export const mockCollections = [
  {
    name: 'bukhari',
    title: 'Sahih al-Bukhari',
    hadiths: 7563,
    books: 97,
    description: 'The most authentic collection of hadiths'
  },
  {
    name: 'muslim',
    title: 'Sahih Muslim',
    hadiths: 7500,
    books: 56,
    description: 'The second most authentic collection'
  },
  {
    name: 'abudawud',
    title: 'Sunan Abu Dawood',
    hadiths: 5278,
    books: 43,
    description: 'Collection of hadiths judged by scholars to be Hassan or Sahih'
  },
  {
    name: 'tirmidhi',
    title: "Jami'at-Tirmidhi",
    hadiths: 3956,
    books: 27,
    description: 'One of the six major hadith collections'
  },
  {
    name: 'nasai',
    title: 'Sunan an-Nasai',
    hadiths: 5762,
    books: 52,
    description: 'Another of the six major hadith collections'
  },
  {
    name: 'ibnmajah',
    title: 'Sunan Ibn Majah',
    hadiths: 4341,
    books: 37,
    description: 'One of the six major Sunni hadith collections'
  },
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
  bukhari: [
    { bookNumber: 1, title: 'Revelation', hadiths: 165 },
    { bookNumber: 2, title: 'Faith', hadiths: 225 },
    { bookNumber: 3, title: 'Knowledge', hadiths: 190 },
    { bookNumber: 4, title: 'Ablutions', hadiths: 175 },
    { bookNumber: 5, title: 'Bath', hadiths: 105 },
    { bookNumber: 6, title: 'Menstruation', hadiths: 65 },
    { bookNumber: 7, title: 'Prayers', hadiths: 350 },
    { bookNumber: 8, title: 'Friday Prayer', hadiths: 110 },
    { bookNumber: 9, title: 'Eid Prayers', hadiths: 65 },
    { bookNumber: 10, title: 'Traveling', hadiths: 180 }
  ],
  muslim: [
    { bookNumber: 1, title: 'The Book of Faith', hadiths: 340 },
    { bookNumber: 2, title: 'The Book of Purification', hadiths: 280 },
    { bookNumber: 3, title: 'The Book of Menstruation', hadiths: 95 },
    { bookNumber: 4, title: 'The Book of Prayers', hadiths: 400 },
    { bookNumber: 5, title: 'The Book of Friday', hadiths: 150 },
    { bookNumber: 6, title: 'The Book of Fasting', hadiths: 220 },
    { bookNumber: 7, title: 'The Book of Hajj', hadiths: 280 },
    { bookNumber: 8, title: 'The Book of Marriage', hadiths: 210 },
    { bookNumber: 9, title: 'The Book of Divorce', hadiths: 120 },
    { bookNumber: 10, title: 'The Book of Food', hadiths: 175 }
  ],
  abudawud: [
    { bookNumber: 1, title: 'Purification', hadiths: 200 },
    { bookNumber: 2, title: 'Prayers', hadiths: 350 },
    { bookNumber: 3, title: 'Friday Prayers', hadiths: 120 },
    { bookNumber: 4, title: 'Eid Prayers', hadiths: 80 },
    { bookNumber: 5, title: 'Traveling', hadiths: 150 },
    { bookNumber: 6, title: 'Fasting', hadiths: 210 },
    { bookNumber: 7, title: 'Hajj', hadiths: 240 },
    { bookNumber: 8, title: 'Marriage', hadiths: 180 },
    { bookNumber: 9, title: 'Divorce', hadiths: 90 },
    { bookNumber: 10, title: 'Food', hadiths: 150 }
  ],
  tirmidhi: [
    { bookNumber: 1, title: 'The Book of Faith', hadiths: 180 },
    { bookNumber: 2, title: 'The Book of Purification', hadiths: 150 },
    { bookNumber: 3, title: 'The Book of Prayer', hadiths: 280 },
    { bookNumber: 4, title: 'The Book of Friday', hadiths: 95 },
    { bookNumber: 5, title: 'The Book of Fasting', hadiths: 120 },
    { bookNumber: 6, title: 'The Book of Hajj', hadiths: 165 },
    { bookNumber: 7, title: 'The Book of Marriage', hadiths: 145 },
    { bookNumber: 8, title: 'The Book of Oaths and Vows', hadiths: 85 },
    { bookNumber: 9, title: 'The Book of Blood Money', hadiths: 70 },
    { bookNumber: 10, title: 'The Book of Knowledge', hadiths: 110 }
  ],
  nasai: [
    { bookNumber: 1, title: 'The Book of Purification', hadiths: 220 },
    { bookNumber: 2, title: 'The Book of Menstruation', hadiths: 75 },
    { bookNumber: 3, title: 'The Book of Prayer', hadiths: 380 },
    { bookNumber: 4, title: 'The Book of Friday', hadiths: 125 },
    { bookNumber: 5, title: 'The Book of Eid', hadiths: 85 },
    { bookNumber: 6, title: 'The Book of Travel', hadiths: 165 },
    { bookNumber: 7, title: 'The Book of Fasting', hadiths: 195 },
    { bookNumber: 8, title: 'The Book of Hajj', hadiths: 250 },
    { bookNumber: 9, title: 'The Book of Marriage', hadiths: 175 },
    { bookNumber: 10, title: 'The Book of Divorce', hadiths: 95 }
  ],
  ibnmajah: [
    { bookNumber: 1, title: 'The Book of Purification', hadiths: 195 },
    { bookNumber: 2, title: 'The Book of Prayer', hadiths: 320 },
    { bookNumber: 3, title: 'The Book of the Call to Prayer', hadiths: 85 },
    { bookNumber: 4, title: 'The Book of Friday', hadiths: 110 },
    { bookNumber: 5, title: 'The Book of Eid', hadiths: 70 },
    { bookNumber: 6, title: 'The Book of Fasting', hadiths: 175 },
    { bookNumber: 7, title: 'The Book of Hajj', hadiths: 210 },
    { bookNumber: 8, title: 'The Book of Marriage', hadiths: 165 },
    { bookNumber: 9, title: 'The Book of Divorce', hadiths: 80 },
    { bookNumber: 10, title: 'The Book of Transactions', hadiths: 220 }
  ]
}

export const mockHadiths = {
  bukhari: transformedBukhari,
  muslim: [
    { hadithNumber: 1, grade: 'Sahih', arabic: 'إِنَّمَا الأَعْمَالُ بِالنِّيَّاتِ', english: 'Indeed, actions are but by intentions.', bookNumber: 1 },
    { hadithNumber: 2, grade: 'Sahih', arabic: 'الإِيمَانُ بِضْعٌ وَسِتُّونَ أَوْ سَبْعُونَ شُعْبَةً', english: 'Faith has sixty or seventy branches.', bookNumber: 1 },
    { hadithNumber: 3, grade: 'Sahih', arabic: 'الطُّهُرُ شَطْرُ الإِيمَانِ', english: 'Purification is half of faith.', bookNumber: 1 },
    { hadithNumber: 4, grade: 'Sahih', arabic: 'مَنْ أَحْدَثَ فِي أَمْرِنَا مَا لَيْسَ مِنْهُ فَهُوَ رَدٌّ', english: 'Whoever introduces something into our affair that is not part of it will be rejected.', bookNumber: 1 },
    { hadithNumber: 5, grade: 'Sahih', arabic: 'حُسْنُ الْخُلُقِ نِصْفُ الإِيمَانِ', english: 'Good character is half of faith.', bookNumber: 1 },
    { hadithNumber: 6, grade: 'Sahih', arabic: 'الْمُسْلِمُ مَنْ سَلِمَ الْمُسْلِمُونَ مِنْ لِسَانِهِ وَيَدِهِ', english: 'A Muslim is one from whose tongue and hand the Muslims are safe.', bookNumber: 1 },
    { hadithNumber: 7, grade: 'Sahih', arabic: 'لَا يُؤْمِنُ أَحَدُكُمْ حَتَّى يُحِبَّ لِأَخِيهِ مَا يُحِبُّ لِنَفْسِهِ', english: 'None of you has faith until he loves for his brother what he loves for himself.', bookNumber: 1 },
    { hadithNumber: 8, grade: 'Sahih', arabic: 'خَيْرُكُمْ مَنْ تَعَلَّمَ الْقُرْآنَ وَعَلَّمَهُ', english: 'The best among you are those who learn the Quran and teach it.', bookNumber: 1 },
    { hadithNumber: 9, grade: 'Sahih', arabic: 'صِفُوا الْمَلَائِكَةَ وَصِفُوا الْجِنَّ', english: 'Describe the angels and describe the jinn.', bookNumber: 1 },
    { hadithNumber: 10, grade: 'Sahih', arabic: 'الْعِلْمُ قَبْلَ الْقَوْلِ وَالْعَمَلِ', english: 'Knowledge comes before speech and action.', bookNumber: 1 }
  ],
  abudawud: [
    { hadithNumber: 1, grade: 'Sahih', arabic: 'إِنَّمَا الأَعْمَالُ بِالنِّيَّاتِ', english: 'Indeed, actions are but by intentions.', bookNumber: 1 },
    { hadithNumber: 2, grade: 'Sahih', arabic: 'مَنْ أَحْدَثَ فِي أَمْرِنَا مَا لَيْسَ مِنْهُ فَهُوَ رَدٌّ', english: 'Whoever introduces something into our affair that is not part of it will be rejected.', bookNumber: 1 },
    { hadithNumber: 3, grade: 'Hasan', arabic: 'الطُّهُرُ شَطْرُ الإِيمَانِ', english: 'Purification is half of faith.', bookNumber: 1 },
    { hadithNumber: 4, grade: 'Sahih', arabic: 'صِفُوا الْمَلَائِكَةَ وَصِفُوا الْجِنَّ', english: 'Describe the angels and describe the jinn.', bookNumber: 1 },
    { hadithNumber: 5, grade: 'Sahih', arabic: 'الْعِلْمُ قَبْلَ الْقَوْلِ وَالْعَمَلِ', english: 'Knowledge comes before speech and action.', bookNumber: 1 },
    { hadithNumber: 6, grade: 'Hasan', arabic: 'مَنْ سَلِكَ طَرِيقًا يَلْتَمِسُ فِيهِ عِلْمًا سَهَّلَ اللَّهُ لَهُ طَرِيقًا إِلَى الْجَنَّةِ', english: 'Whoever follows a path seeking knowledge, Allah will make easy for him a path to Paradise.', bookNumber: 1 },
    { hadithNumber: 7, grade: 'Sahih', arabic: 'الْمُؤْمِنُ يَأْلَفُ وَلَا يُوحِشُ', english: 'The believer is friendly and not harsh.', bookNumber: 1 },
    { hadithNumber: 8, grade: 'Sahih', arabic: 'لَا يُؤْمِنُ أَحَدُكُمْ حَتَّى يُحِبَّ لِأَخِيهِ مَا يُحِبُّ لِنَفْسِهِ', english: 'None of you has faith until he loves for his brother what he loves for himself.', bookNumber: 1 },
    { hadithNumber: 9, grade: 'Sahih', arabic: 'خَيْرُ النَّاسِ أَعْمَلُهُمْ لِلنَّاسِ', english: 'The best of people are those most beneficial to others.', bookNumber: 1 },
    { hadithNumber: 10, grade: 'Hasan', arabic: 'حُسْنُ الْخُلُقِ نِصْفُ الإِيمَانِ', english: 'Good character is half of faith.', bookNumber: 1 }
  ],
  tirmidhi: [
    { hadithNumber: 1, grade: 'Sahih', arabic: 'إِنَّمَا الأَعْمَالُ بِالنِّيَّاتِ', english: 'Indeed, actions are but by intentions.', bookNumber: 1 },
    { hadithNumber: 2, grade: 'Sahih', arabic: 'الْإِيمَانُ بِضْعٌ وَسِتُّونَ شُعْبَةً', english: 'Faith has over sixty branches.', bookNumber: 1 },
    { hadithNumber: 3, grade: 'Sahih', arabic: 'الطُّهُرُ شَطْرُ الإِيمَانِ', english: 'Purification is half of faith.', bookNumber: 1 },
    { hadithNumber: 4, grade: 'Sahih', arabic: 'مَنْ أَحْدَثَ فِي أَمْرِنَا مَا لَيْسَ مِنْهُ فَهُوَ رَدٌّ', english: 'Whoever introduces something into our affair that is not part of it will be rejected.', bookNumber: 1 },
    { hadithNumber: 5, grade: 'Hasan', arabic: 'حُسْنُ الْخُلُقِ نِصْفُ الإِيمَانِ', english: 'Good character is half of faith.', bookNumber: 1 },
    { hadithNumber: 6, grade: 'Sahih', arabic: 'الْمُسْلِمُ مَنْ سَلِمَ الْمُسْلِمُونَ مِنْ لِسَانِهِ وَيَدِهِ', english: 'A Muslim is one from whose tongue and hand the Muslims are safe.', bookNumber: 1 },
    { hadithNumber: 7, grade: 'Sahih', arabic: 'لَا يُؤْمِنُ أَحَدُكُمْ حَتَّى يُحِبَّ لِأَخِيهِ مَا يُحِبُّ لِنَفْسِهِ', english: 'None of you has faith until he loves for his brother what he loves for himself.', bookNumber: 1 },
    { hadithNumber: 8, grade: 'Sahih', arabic: 'خَيْرُكُمْ مَنْ تَعَلَّمَ الْقُرْآنَ وَعَلَّمَهُ', english: 'The best among you are those who learn the Quran and teach it.', bookNumber: 1 },
    { hadithNumber: 9, grade: 'Sahih', arabic: 'الْعِلْمُ قَبْلَ الْقَوْلِ وَالْعَمَلِ', english: 'Knowledge comes before speech and action.', bookNumber: 1 },
    { hadithNumber: 10, grade: 'Hasan', arabic: 'مَنْ سَلِكَ طَرِيقًا يَلْتَمِسُ فِيهِ عِلْمًا سَهَّلَ اللَّهُ لَهُ طَرِيقًا إِلَى الْجَنَّةِ', english: 'Whoever follows a path seeking knowledge, Allah will make easy for him a path to Paradise.', bookNumber: 1 }
  ],
  nasai: [
    { hadithNumber: 1, grade: 'Sahih', arabic: 'إِنَّمَا الأَعْمَالُ بِالنِّيَّاتِ', english: 'Indeed, actions are but by intentions.', bookNumber: 1 },
    { hadithNumber: 2, grade: 'Sahih', arabic: 'الطُّهُرُ شَطْرُ الإِيمَانِ', english: 'Purification is half of faith.', bookNumber: 1 },
    { hadithNumber: 3, grade: 'Sahih', arabic: 'الْإِيمَانُ بِضْعٌ وَسِتُّونَ شُعْبَةً', english: 'Faith has over sixty branches.', bookNumber: 1 },
    { hadithNumber: 4, grade: 'Sahih', arabic: 'مَنْ أَحْدَثَ فِي أَمْرِنَا مَا لَيْسَ مِنْهُ فَهُوَ رَدٌّ', english: 'Whoever introduces something into our affair that is not part of it will be rejected.', bookNumber: 1 },
    { hadithNumber: 5, grade: 'Sahih', arabic: 'حُسْنُ الْخُلُقِ نِصْفُ الإِيمَانِ', english: 'Good character is half of faith.', bookNumber: 1 },
    { hadithNumber: 6, grade: 'Sahih', arabic: 'الْمُسْلِمُ مَنْ سَلِمَ الْمُسْلِمُونَ مِنْ لِسَانِهِ وَيَدِهِ', english: 'A Muslim is one from whose tongue and hand the Muslims are safe.', bookNumber: 1 },
    { hadithNumber: 7, grade: 'Sahih', arabic: 'لَا يُؤْمِنُ أَحَدُكُمْ حَتَّى يُحِبَّ لِأَخِيهِ مَا يُحِبُّ لِنَفْسِهِ', english: 'None of you has faith until he loves for his brother what he loves for himself.', bookNumber: 1 },
    { hadithNumber: 8, grade: 'Sahih', arabic: 'خَيْرُ النَّاسِ أَعْمَلُهُمْ لِلنَّاسِ', english: 'The best of people are those most beneficial to others.', bookNumber: 1 },
    { hadithNumber: 9, grade: 'Sahih', arabic: 'الْمُؤْمِنُ يَأْلَفُ وَلَا يُوحِشُ', english: 'The believer is friendly and not harsh.', bookNumber: 1 },
    { hadithNumber: 10, grade: 'Sahih', arabic: 'خَيْرُكُمْ مَنْ تَعَلَّمَ الْقُرْآنَ وَعَلَّمَهُ', english: 'The best among you are those who learn the Quran and teach it.', bookNumber: 1 }
  ],
  ibnmajah: [
    { hadithNumber: 1, grade: 'Sahih', arabic: 'إِنَّمَا الأَعْمَالُ بِالنِّيَّاتِ', english: 'Indeed, actions are but by intentions.', bookNumber: 1 },
    { hadithNumber: 2, grade: 'Sahih', arabic: 'الطُّهُرُ شَطْرُ الإِيمَانِ', english: 'Purification is half of faith.', bookNumber: 1 },
    { hadithNumber: 3, grade: 'Hasan', arabic: 'الْإِيمَانُ بِضْعٌ وَسِتُّونَ شُعْبَةً', english: 'Faith has over sixty branches.', bookNumber: 1 },
    { hadithNumber: 4, grade: 'Sahih', arabic: 'مَنْ أَحْدَثَ فِي أَمْرِنَا مَا لَيْسَ مِنْهُ فَهُوَ رَدٌّ', english: 'Whoever introduces something into our affair that is not part of it will be rejected.', bookNumber: 1 },
    { hadithNumber: 5, grade: 'Hasan', arabic: 'حُسْنُ الْخُلُقِ نِصْفُ الإِيمَانِ', english: 'Good character is half of faith.', bookNumber: 1 },
    { hadithNumber: 6, grade: 'Sahih', arabic: 'الْمُسْلِمُ مَنْ سَلِمَ الْمُسْلِمُونَ مِنْ لِسَانِهِ وَيَدِهِ', english: 'A Muslim is one from whose tongue and hand the Muslims are safe.', bookNumber: 1 },
    { hadithNumber: 7, grade: 'Sahih', arabic: 'لَا يُؤْمِنُ أَحَدُكُمْ حَتَّى يُحِبَّ لِأَخِيهِ مَا يُحِبُّ لِنَفْسِهِ', english: 'None of you has faith until he loves for his brother what he loves for himself.', bookNumber: 1 },
    { hadithNumber: 8, grade: 'Sahih', arabic: 'خَيْرُ النَّاسِ أَعْمَلُهُمْ لِلنَّاسِ', english: 'The best of people are those most beneficial to others.', bookNumber: 1 },
    { hadithNumber: 9, grade: 'Hasan', arabic: 'الْمُؤْمِنُ يَأْلَفُ وَلَا يُوحِشُ', english: 'The believer is friendly and not harsh.', bookNumber: 1 },
    { hadithNumber: 10, grade: 'Sahih', arabic: 'خَيْرُكُمْ مَنْ تَعَلَّمَ الْقُرْآنَ وَعَلَّمَهُ', english: 'The best among you are those who learn the Quran and teach it.', bookNumber: 1 }
  ]
}

// Helper function to get mock data
export const getMockCollections = () => Promise.resolve(mockCollections)

export const getMockBooks = (collection) => {
  const books = mockBooks[collection] || mockBooks.bukhari
  return Promise.resolve(books)
}

export const getMockHadiths = (collection, bookNumber, page = 1, limit = 20) => {
  const collectionHadiths = mockHadiths[collection] || mockHadiths.bukhari
  const bookHadiths = collectionHadiths.slice(0, 10)
  return Promise.resolve({
    hadiths: bookHadiths,
    total: bookHadiths.length,
    page,
    totalPages: 1
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

