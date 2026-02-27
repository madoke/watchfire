import { useEffect, type ReactNode } from 'react'
import { X } from 'lucide-react'

interface SlidePanelProps {
  open: boolean
  onClose: () => void
  title: string
  children: ReactNode
  footer?: ReactNode
}

export function SlidePanel({ open, onClose, title, children, footer }: SlidePanelProps) {
  useEffect(() => {
    if (!open) return
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [open, onClose])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-[200] flex justify-end">
      <div className="absolute inset-0 bg-black/60" onClick={onClose} />
      <div
        className="relative w-[560px] max-w-full h-full bg-[var(--wf-bg-secondary)] border-l border-[var(--wf-border)] shadow-wf-lg flex flex-col"
        style={{ animation: 'slideInRight 0.2s ease-out' }}
      >
        <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--wf-border)] shrink-0">
          <h3 className="text-base font-semibold">{title}</h3>
          <button onClick={onClose} className="text-[var(--wf-text-muted)] hover:text-[var(--wf-text-primary)] transition-colors">
            <X size={18} />
          </button>
        </div>
        <div className="flex-1 overflow-y-auto px-5 py-4">{children}</div>
        {footer && (
          <div className="flex items-center justify-end gap-2 px-5 py-3 border-t border-[var(--wf-border)] shrink-0">
            {footer}
          </div>
        )}
      </div>
    </div>
  )
}
