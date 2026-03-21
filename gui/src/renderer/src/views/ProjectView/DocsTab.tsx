import { useState, useEffect, useRef, useCallback, useMemo } from 'react'
import { FileText, Plus } from 'lucide-react'
import type { Revision } from '../../generated/watchfire_pb'
import { useRevisionsStore } from '../../stores/revisions-store'
import { useTasksStore } from '../../stores/tasks-store'
import { useProjectsStore } from '../../stores/projects-store'
import { getProjectClient } from '../../lib/grpc-client'
import { useToast } from '../../components/ui/Toast'
import { MarkdownEditor } from '../../components/MarkdownEditor'
import { cn } from '../../lib/utils'

const EMPTY_REVISIONS: Revision[] = []

type RevisionStatus = 'done' | 'in-dev' | 'todo'
type Selection = { type: 'definition' } | { type: 'revision'; revisionNumber: number }

interface Props {
  projectId: string
}

export function DocsTab({ projectId }: Props) {
  const project = useProjectsStore((s) => s.projects.find((p) => p.projectId === projectId))
  const fetchProjects = useProjectsStore((s) => s.fetchProjects)
  const updateProjectLocal = useProjectsStore((s) => s.updateProjectLocal)

  const revisions = useRevisionsStore((s) => s.revisions[projectId]) ?? EMPTY_REVISIONS
  const fetchRevisions = useRevisionsStore((s) => s.fetchRevisions)
  const tasks = useTasksStore((s) => s.tasks[projectId]) ?? []

  const revisionStatuses = useMemo(() => {
    const statuses: Record<number, RevisionStatus> = {}
    for (const rev of revisions) {
      if (rev.complete) {
        statuses[rev.revisionNumber] = 'done'
        continue
      }
      const revTasks = tasks.filter((t) => t.revisionNumber === rev.revisionNumber)
      const hasInDev = revTasks.some((t) => t.status === 'ready')
      statuses[rev.revisionNumber] = hasInDev ? 'in-dev' : 'todo'
    }
    return statuses
  }, [revisions, tasks])
  const createRevision = useRevisionsStore((s) => s.createRevision)
  const updateRevision = useRevisionsStore((s) => s.updateRevision)
  const deleteRevision = useRevisionsStore((s) => s.deleteRevision)
  const { toast } = useToast()

  const [selection, setSelection] = useState<Selection>({ type: 'definition' })
  const [saved, setSaved] = useState(true)
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const dirtyRef = useRef(false)

  // Fetch revisions on mount and poll
  useEffect(() => {
    fetchRevisions(projectId)
  }, [projectId])

  useEffect(() => {
    const interval = setInterval(() => {
      if (!dirtyRef.current) fetchRevisions(projectId)
    }, 5000)
    return () => clearInterval(interval)
  }, [projectId])

  const selectedRevision = selection.type === 'revision'
    ? revisions.find((r) => r.revisionNumber === selection.revisionNumber)
    : null

  // Resolve current editor value and key
  const editorKey = selection.type === 'definition' ? `def-${projectId}` : `rev-${(selection as any).revisionNumber}`
  const editorValue = selection.type === 'definition'
    ? (project?.definition ?? '')
    : (selectedRevision?.content ?? '')

  const handleEditorChange = useCallback((content: string) => {
    setSaved(false)
    dirtyRef.current = true
    if (timerRef.current) clearTimeout(timerRef.current)

    if (selection.type === 'definition') {
      timerRef.current = setTimeout(async () => {
        updateProjectLocal(projectId, { definition: content })
        try {
          const client = getProjectClient()
          await client.updateProject({ projectId, definition: content })
          setSaved(true)
          dirtyRef.current = false
          await fetchProjects()
        } catch {
          toast('Failed to save definition', 'error')
        }
      }, 1000)
    } else {
      const revNum = selection.revisionNumber
      const match = content.match(/^#\s+(.+)$/m)
      const autoTitle = match ? match[1].trim() : null

      timerRef.current = setTimeout(async () => {
        const updates: { content: string; title?: string } = { content }
        if (autoTitle) updates.title = autoTitle
        try {
          await updateRevision(projectId, revNum, updates)
          setSaved(true)
          dirtyRef.current = false
        } catch {
          toast('Failed to save revision', 'error')
        }
      }, 1000)
    }
  }, [projectId, selection, updateProjectLocal, fetchProjects, updateRevision])

  const titleTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const handleTitleChange = (revNum: number, title: string) => {
    if (titleTimerRef.current) clearTimeout(titleTimerRef.current)
    titleTimerRef.current = setTimeout(async () => {
      try {
        await updateRevision(projectId, revNum, { title })
      } catch {
        toast('Failed to update title', 'error')
      }
    }, 500)
  }

  const handleCreate = async () => {
    try {
      const rev = await createRevision(projectId, 'New Revision', '')
      setSelection({ type: 'revision', revisionNumber: rev.revisionNumber })
    } catch {
      toast('Failed to create revision', 'error')
    }
  }

  const handleDelete = async (revNum: number) => {
    try {
      await deleteRevision(projectId, revNum)
      if (selection.type === 'revision' && selection.revisionNumber === revNum) {
        setSelection({ type: 'definition' })
      }
    } catch {
      toast('Failed to delete revision', 'error')
    }
  }

  const switchTo = (sel: Selection) => {
    dirtyRef.current = false
    setSaved(true)
    setSelection(sel)
  }

  return (
    <div className="flex h-full">
      {/* Left sidebar */}
      <div className="w-56 shrink-0 border-r border-[var(--wf-border)] flex flex-col">
        {/* Project Definition entry */}
        <div className="border-b border-[var(--wf-border)]">
          <button
            onClick={() => switchTo({ type: 'definition' })}
            className={cn(
              'w-full flex items-center gap-2 px-3 py-2 text-left text-xs transition-colors',
              selection.type === 'definition'
                ? 'bg-[var(--wf-bg-elevated)] text-[var(--wf-text-primary)]'
                : 'text-[var(--wf-text-secondary)] hover:bg-[var(--wf-bg-elevated)]'
            )}
          >
            <FileText size={12} className="shrink-0" />
            <span className="flex-1 truncate">Project Definition</span>
          </button>
        </div>

        {/* Revisions header */}
        <div className="flex items-center justify-between px-3 py-2 border-b border-[var(--wf-border)]">
          <span className="text-xs font-medium text-[var(--wf-text-muted)]">Revisions</span>
          <button
            onClick={handleCreate}
            className="p-1 text-[var(--wf-text-muted)] hover:text-[var(--wf-text-primary)] transition-colors"
            title="New revision"
          >
            <Plus size={14} />
          </button>
        </div>
        <div className="flex-1 overflow-y-auto">
          {revisions.map((rev) => (
            <RevisionItem
              key={rev.revisionNumber}
              revision={rev}
              status={revisionStatuses[rev.revisionNumber] ?? 'todo'}
              selected={selection.type === 'revision' && selection.revisionNumber === rev.revisionNumber}
              onSelect={() => switchTo({ type: 'revision', revisionNumber: rev.revisionNumber })}
              onTitleChange={(title) => handleTitleChange(rev.revisionNumber, title)}
              onDelete={() => handleDelete(rev.revisionNumber)}
            />
          ))}
          {revisions.length === 0 && (
            <div className="px-3 py-6 text-center text-xs text-[var(--wf-text-muted)]">
              No revisions yet
            </div>
          )}
        </div>
      </div>

      {/* Right area: shared editor */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header bar */}
        <div className="flex items-center justify-between px-4 py-2 border-b border-[var(--wf-border)] gap-3">
          <div className="flex items-center gap-2 flex-1 min-w-0">
            {selection.type === 'definition' ? (
              <span className="text-xs text-[var(--wf-text-muted)]">Project Definition</span>
            ) : selectedRevision ? (
              <>
                <span className="shrink-0 text-xs text-[var(--wf-text-muted)]">
                  #{String(selectedRevision.revisionNumber).padStart(4, '0')}
                </span>
                <input
                  value={selectedRevision.title}
                  onChange={(e) => handleTitleChange(selectedRevision.revisionNumber, e.target.value)}
                  className="flex-1 min-w-0 text-xs bg-transparent text-[var(--wf-text-primary)] outline-none border-b border-transparent focus:border-[var(--wf-accent)] transition-colors"
                  placeholder="Revision title"
                />
              </>
            ) : null}
          </div>
          <span className="shrink-0 text-xs text-[var(--wf-text-muted)]">
            {saved ? 'Saved' : 'Saving...'}
          </span>
        </div>

        {/* Editor */}
        <div className="flex-1 overflow-hidden">
          <MarkdownEditor
            key={editorKey}
            value={editorValue}
            onChange={handleEditorChange}
            placeholder={selection.type === 'definition'
              ? 'Describe your project...'
              : 'Describe the work for this revision...'}
          />
        </div>
      </div>
    </div>
  )
}

function RevisionItem({
  revision,
  status,
  selected,
  onSelect,
  onTitleChange,
  onDelete
}: {
  revision: Revision
  status: RevisionStatus
  selected: boolean
  onSelect: () => void
  onTitleChange: (title: string) => void
  onDelete: () => void
}) {
  const [editing, setEditing] = useState(false)
  const [title, setTitle] = useState(revision.title)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    setTitle(revision.title)
  }, [revision.title])

  useEffect(() => {
    if (editing && inputRef.current) {
      inputRef.current.focus()
      inputRef.current.select()
    }
  }, [editing])

  const commitTitle = () => {
    setEditing(false)
    if (title.trim() && title !== revision.title) {
      onTitleChange(title.trim())
    } else {
      setTitle(revision.title)
    }
  }

  return (
    <button
      onClick={onSelect}
      onDoubleClick={() => setEditing(true)}
      onContextMenu={(e) => {
        e.preventDefault()
        if (confirm(`Delete revision "${revision.title}"?`)) {
          onDelete()
        }
      }}
      className={cn(
        'w-full flex items-center gap-2 px-3 py-2 text-left text-xs transition-colors',
        selected
          ? 'bg-[var(--wf-bg-elevated)] text-[var(--wf-text-primary)]'
          : 'text-[var(--wf-text-secondary)] hover:bg-[var(--wf-bg-elevated)]'
      )}
    >
      {editing ? (
        <input
          ref={inputRef}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          onBlur={commitTitle}
          onKeyDown={(e) => {
            if (e.key === 'Enter') commitTitle()
            if (e.key === 'Escape') { setTitle(revision.title); setEditing(false) }
          }}
          className="flex-1 bg-transparent text-xs outline-none border-b border-[var(--wf-accent)]"
          onClick={(e) => e.stopPropagation()}
        />
      ) : (
        <span className="flex-1 truncate">{revision.title}</span>
      )}
      <span
        className={cn(
          'shrink-0 inline-flex items-center px-1.5 py-0.5 rounded-full text-[10px] font-medium',
          status === 'done' && 'bg-emerald-900/30 text-[var(--wf-success)]',
          status === 'in-dev' && 'bg-amber-900/30 text-[var(--wf-warning)]',
          status === 'todo' && 'bg-[var(--wf-bg-elevated)] text-[var(--wf-text-muted)]'
        )}
      >
        {status === 'done' ? 'Done' : status === 'in-dev' ? 'In Dev' : 'Todo'}
      </span>
    </button>
  )
}
