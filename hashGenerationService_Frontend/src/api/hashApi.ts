import type { HashResponse } from '../types'

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

export async function generateHash(input: string): Promise<HashResponse> {
  const res = await fetch(`${BASE_URL}/hash`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ input }),
  })

  if (!res.ok) {
    let message = 'Failed to generate hash'
    try {
      const data = await res.json()
      message = data.error ?? message
    } catch {
      // non-JSON error response (e.g. 413 from MaxBytesReader)
    }
    throw new Error(message)
  }

  return res.json() as Promise<HashResponse>
}
