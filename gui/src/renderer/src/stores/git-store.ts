import { create } from 'zustand'
import type { GitInfo } from '../generated/watchfire_pb'
import { getProjectClient } from '../lib/grpc-client'

interface GitState {
  gitInfo: Record<string, GitInfo>
  fetchGitInfo: (projectId: string) => Promise<void>
}

export const useGitStore = create<GitState>((set) => ({
  gitInfo: {},

  fetchGitInfo: async (projectId) => {
    try {
      const client = getProjectClient()
      const info = await client.getGitInfo({ projectId })
      set((s) => ({ gitInfo: { ...s.gitInfo, [projectId]: info } }))
    } catch {
      // ignore — project may not be a git repo
    }
  }
}))
