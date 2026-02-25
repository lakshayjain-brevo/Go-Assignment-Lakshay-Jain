import type {  HashResponse } from '../types'

const BASE_URL = 'http://localhost:8080'

export async function generateHash(input: string): Promise<HashResponse> {
  const res = await fetch(`${BASE_URL}/hash`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ input }),
  })

  const data = await res.json()
  if (!res.ok) throw new Error(data.error ?? 'Failed to generate hash')
  return data as HashResponse
}