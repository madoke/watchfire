import { useState, useEffect, useRef, useCallback } from 'react'
import type { Project } from '../../generated/watchfire_pb'
import { getProjectClient } from '../../lib/grpc-client'
import { useToast } from '../../components/ui/Toast'

interface Props {
  projectId: string
  project: Project
}

export function SecretsTab({ projectId, project }: Props) {
  const [value, setValue] = useState(project.secretsInstructions || '')
  const [saved, setSaved] = useState(true)
  const { toast } = useToast()
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    setValue(project.secretsInstructions || '')
    setSaved(true)
  }, [project.secretsInstructions, projectId])

  const save = useCallback(async (text: string) => {
    try {
      const client = getProjectClient()
      await client.updateProject({ projectId, secretsInstructions: text })
      setSaved(true)
    } catch (err) {
      toast('Failed to save secrets instructions', 'error')
    }
  }, [projectId])

  const handleChange = (text: string) => {
    setValue(text)
    setSaved(false)
    if (timerRef.current) clearTimeout(timerRef.current)
    timerRef.current = setTimeout(() => save(text), 1000)
  }

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between px-4 py-2 border-b border-[var(--wf-border)]">
        <span className="text-xs text-[var(--wf-text-muted)]">
          Secrets & setup instructions — Markdown
        </span>
        <span className="text-xs text-[var(--wf-text-muted)]">
          {saved ? 'Saved' : 'Saving...'}
        </span>
      </div>
      <textarea
        value={value}
        onChange={(e) => handleChange(e.target.value)}
        className="flex-1 w-full px-4 py-3 bg-[var(--wf-bg-primary)] text-sm font-mono leading-relaxed text-[var(--wf-text-primary)] placeholder-[var(--wf-text-muted)] focus:outline-none resize-none"
        placeholder="Tell agents how to access external services, API keys, CLI tools..."
      />
    </div>
  )
}
