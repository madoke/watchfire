import { statusLabel, statusColor, cn } from '../lib/utils'

interface TaskStatusBadgeProps {
  status: string
  success?: boolean
  className?: string
}

export function TaskStatusBadge({ status, success, className }: TaskStatusBadgeProps) {
  const isFailed = status === 'done' && success !== true

  const bgMap: Record<string, string> = {
    draft: 'bg-[var(--wf-bg-elevated)]',
    ready: 'bg-amber-900/30',
    done: 'bg-emerald-900/30'
  }

  return (
    <span
      className={cn(
        'inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium',
        isFailed ? 'bg-red-900/30 text-red-400' : bgMap[status] || 'bg-[var(--wf-bg-elevated)]',
        !isFailed && statusColor(status),
        className
      )}
    >
      {isFailed ? 'Failed' : statusLabel(status)}
    </span>
  )
}
