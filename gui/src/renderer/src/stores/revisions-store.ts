import { create } from 'zustand'
import type { Revision } from '../generated/watchfire_pb'
import { getRevisionClient } from '../lib/grpc-client'

interface RevisionsState {
  revisions: Record<string, Revision[]>
  loading: boolean
  error: string | null

  fetchRevisions: (projectId: string) => Promise<void>
  createRevision: (projectId: string, title: string, content: string) => Promise<Revision>
  updateRevision: (projectId: string, revisionNumber: number, updates: {
    title?: string
    content?: string
    complete?: boolean
  }) => Promise<void>
  deleteRevision: (projectId: string, revisionNumber: number) => Promise<void>
}

export const useRevisionsStore = create<RevisionsState>((set, get) => ({
  revisions: {},
  loading: false,
  error: null,

  fetchRevisions: async (projectId) => {
    set({ loading: true, error: null })
    try {
      const client = getRevisionClient()
      const resp = await client.listRevisions({ projectId })
      set((s) => ({
        revisions: { ...s.revisions, [projectId]: resp.revisions },
        loading: false
      }))
    } catch (err) {
      set({ error: String(err), loading: false })
    }
  },

  createRevision: async (projectId, title, content) => {
    const client = getRevisionClient()
    const revision = await client.createRevision({ projectId, title, content })
    get().fetchRevisions(projectId)
    return revision
  },

  updateRevision: async (projectId, revisionNumber, updates) => {
    const client = getRevisionClient()
    await client.updateRevision({ projectId, revisionNumber, ...updates })
    get().fetchRevisions(projectId)
  },

  deleteRevision: async (projectId, revisionNumber) => {
    const client = getRevisionClient()
    await client.deleteRevision({ projectId, revisionNumber })
    get().fetchRevisions(projectId)
  }
}))
