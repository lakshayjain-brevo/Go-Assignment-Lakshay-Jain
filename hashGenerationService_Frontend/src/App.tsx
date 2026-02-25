import GenerateHash from './components/GenerateHash/GenerateHash'
import styles from './App.module.css'

export default function App() {
  return (
    <div className={styles.app}>
      <header className={styles.header}>
        <div className={styles.headerInner}>
          <div className={styles.logo}>
            <span className={styles.logoIcon}>#</span>
            <span>Hash Generation Service</span>
          </div>
          <span className={styles.badge}>v1.0</span>
        </div>
      </header>

      <main className={styles.main}>
        <div className={styles.hero}>
          <h1>Generate Hashes</h1>
          <p>
            Enter an alphanumeric string to get a unique 10-character hash.
          </p>
        </div>

        <div className={styles.grid}>
          <GenerateHash />
        </div>
      </main>
    </div>
  )
}
