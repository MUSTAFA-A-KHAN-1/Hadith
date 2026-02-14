import Badge from '../common/Badge'
import { getGradeColor } from '../../utils/constants'
import { getGrade } from '../../utils/helpers'

const GradeBadge = ({ grade, size = 'md' }) => {
  const gradeText = getGrade(grade)
  const colorClass = getGradeColor(gradeText)

  return (
    <Badge className={colorClass} size={size}>
      {gradeText}
    </Badge>
  )
}

export default GradeBadge

