import type { PlanRequest, GenerateResult } from '../types/itinerary'

const API_BASE = '/api'

async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...init?.headers,
    },
  })
  if (!res.ok) {
    const body = await res.text()
    throw new Error(`API error ${res.status}: ${body}`)
  }
  return res.json()
}

export async function generatePlan(req: PlanRequest): Promise<GenerateResult> {
  return fetchJSON<GenerateResult>('/plans', {
    method: 'POST',
    body: JSON.stringify(req),
  })
}

export async function healthCheck(): Promise<{ status: string }> {
  return fetchJSON<{ status: string }>('/health')
}

export interface RegisteredUser {
  id: string
  email: string
  created_at: string
}

export async function registerUser(email: string): Promise<RegisteredUser> {
  return fetchJSON<RegisteredUser>('/users/register', {
    method: 'POST',
    body: JSON.stringify({ email }),
  })
}
