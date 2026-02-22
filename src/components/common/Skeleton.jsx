const Skeleton = ({ 
  variant = 'rect', 
  width, 
  height, 
  className = '',
  count = 1 
}) => {
  const variants = {
    rect: 'rounded-lg',
    circle: 'rounded-full',
    text: 'rounded h-4'
  }

  const elements = Array.from({ length: count }, (_, i) => (
    <div
      key={i}
      className={`skeleton ${variants[variant]} ${className}`}
      style={{
        width: width || (variant === 'text' ? '100%' : undefined),
        height: height || (variant === 'text' ? '1rem' : undefined)
      }}
    />
  ))

  return <>{elements}</>
}

// Skeleton for cards
export const CardSkeleton = ({ className = '' }) => (
  <div className={`bg-white dark:bg-background-card-dark rounded-2xl p-6 shadow-card ${className}`}>
    <Skeleton height="1.5rem" width="60%" className="mb-4" />
    <Skeleton height="1rem" className="mb-2" />
    <Skeleton height="1rem" width="80%" className="mb-4" />
    <Skeleton height="8rem" />
  </div>
)

// Skeleton for hadith
export const HadithSkeleton = ({ className = '' }) => (
  <div className={`bg-white dark:bg-background-card-dark rounded-2xl p-6 shadow-card ${className}`}>
    <Skeleton height="2rem" width="40%" className="mb-4" />
    <Skeleton height="1rem" className="mb-2" />
    <Skeleton height="1rem" width="90%" className="mb-2" />
    <Skeleton height="1rem" width="75%" className="mb-4" />
    <div className="gold-divider my-4" />
    <Skeleton height="1rem" className="mb-2" />
    <Skeleton height="1rem" width="85%" className="mb-2" />
    <Skeleton height="1rem" width="70%" />
  </div>
)

// Skeleton for collection card
export const CollectionSkeleton = ({ className = '' }) => (
  <div className={`bg-white dark:bg-background-card-dark rounded-2xl p-6 shadow-card ${className}`}>
    <div className="flex items-center gap-4 mb-4">
      <Skeleton variant="circle" width="60px" height="60px" />
      <div className="flex-1">
        <Skeleton height="1.5rem" width="70%" className="mb-2" />
        <Skeleton height="1rem" width="40%" />
      </div>
    </div>
    <Skeleton height="3rem" />
  </div>
)

export default Skeleton

