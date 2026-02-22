const Card = ({
  children,
  className = '',
  hover = true,
  padding = 'md',
  onClick,
  ...props
}) => {
  const paddingClasses = {
    none: '',
    sm: 'p-4',
    md: 'p-6',
    lg: 'p-8'
  }

  return (
    <div
      className={`
        bg-white dark:bg-background-card-dark rounded-2xl
        shadow-card dark:shadow-gray-900/50
        ${hover ? 'card-hover cursor-pointer' : ''}
        ${paddingClasses[padding]}
        ${className}
      `}
      onClick={onClick}
      {...props}
    >
      {children}
    </div>
  )
}

export default Card

