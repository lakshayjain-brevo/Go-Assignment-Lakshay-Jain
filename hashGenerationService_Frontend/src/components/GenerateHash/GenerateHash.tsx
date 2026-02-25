import { useState } from 'react'
import { generateHash } from '../../api/hashApi'
import type { HashResponse } from '../../types'
import styles from './GenerateHash.module.css'

export default function GenerateHash() {
  const [input, setInput] = useState('')
  const [result, setResult] = useState<HashResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [copied, setCopied] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!input.trim()) return

    setError(null)
    setResult(null)
    setLoading(true)

    try {
      const data = await generateHash(input.trim())
      setResult(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  async function handleCopy() {
    if (!result) return
    try {
      await navigator.clipboard.writeText(result.hash)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      // clipboard not available in non-secure context — silently ignore
    }
  }

  function handleReset() {
    setInput('')
    setResult(null)
    setError(null)
  }

  return (
    <div className={styles.card}>
      <div className={styles.cardHeader}>
        <div className={styles.iconWrapper}>⚡</div>
        <div>
          <h2 className={styles.title}>Generate Hash</h2>
          <p className={styles.subtitle}>Alphanumeric input → 10-char unique hash</p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.field}>
          <label htmlFor="gen-input" className={styles.label}>
            Input
          </label>
          <input
            id="gen-input"
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="e.g. Hello123"
            className={styles.input}
            disabled={loading}
            autoComplete="off"
          />
          <span className={styles.hint}>Only letters and numbers (a–z, A–Z, 0–9)</span>
        </div>

        <div className={styles.actions}>
          <button
            type="submit"
            className={styles.primaryBtn}
            disabled={loading || !input.trim()}
          >
            {loading ? (
              <>
                <span className={styles.spinner} />
                Generating...
              </>
            ) : (
              'Generate Hash'
            )}
          </button>

          {(result !== null || error !== null) && (
            <button type="button" className={styles.secondaryBtn} onClick={handleReset}>
              Reset
            </button>
          )}
        </div>
      </form>

      {error !== null && (
        <div className={styles.error}>
          <span className={styles.errorIcon}>✕</span>
          {error}
        </div>
      )}

      {result !== null && (
        <div className={styles.result}>
          <div className={styles.resultHeader}>Result</div>
          <div className={styles.resultRow}>
            <span className={styles.resultLabel}>Input</span>
            <span className={styles.resultValue}>{result.input}</span>
          </div>
          <div className={styles.resultRow}>
            <span className={styles.resultLabel}>Hash</span>
            <div className={styles.hashRow}>
              <code className={styles.hash}>{result.hash}</code>
              <button
                className={`${styles.copyBtn} ${copied ? styles.copiedBtn : ''}`}
                onClick={handleCopy}
              >
                {copied ? '✓ Copied' : 'Copy'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
