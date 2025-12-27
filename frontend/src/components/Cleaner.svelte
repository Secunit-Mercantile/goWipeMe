<script>
  import { onMount } from 'svelte'
  import { GetCleanerStatus, RunCleaner } from '../../wailsjs/go/gui/App'

  let { onBack } = $props()

  let loading = $state(true)
  let cleaners = $state([])
  let cleaning = $state(false)
  let complete = $state(false)
  let error = $state(null)

  onMount(async () => {
    await loadCleaners()
  })

  async function loadCleaners() {
    try {
      loading = true
      const data = await GetCleanerStatus()
      cleaners = data || []
      loading = false
    } catch (err) {
      error = err.message
      loading = false
    }
  }

  async function handleClean() {
    if (!confirm('This will permanently delete the listed items. Are you sure?')) {
      return
    }

    try {
      cleaning = true
      error = null
      await RunCleaner()
      cleaning = false
      complete = true
    } catch (err) {
      error = err.message
      cleaning = false
    }
  }

  function handleBack() {
    onBack()
  }
</script>

<div class="cleaner">
  <div class="header">
    <button class="back-btn" onclick={handleBack}>← Back</button>
    <h1>Clear All History</h1>
  </div>

  <div class="content">
    {#if loading}
      <div class="loading">
        <div class="spinner"></div>
        <p>Scanning for items to clean...</p>
      </div>
    {:else if error}
      <div class="error">
        <p>Error: {error}</p>
        <button onclick={loadCleaners}>Retry</button>
      </div>
    {:else if complete}
      <div class="success">
        <div class="success-icon">✓</div>
        <h2>Cleaning Complete!</h2>
        <p>All selected items have been successfully removed.</p>
        <button class="primary-btn" onclick={handleBack}>Back to Home</button>
      </div>
    {:else if cleaners.length === 0}
      <div class="empty">
        <p>Nothing to clean! Your system is already clean.</p>
        <button onclick={handleBack}>Back to Home</button>
      </div>
    {:else}
      <div class="preview">
        <h2>Items to be cleaned:</h2>

        <div class="cleaners-list">
          {#each cleaners as cleaner}
            <div class="cleaner-card">
              <h3>{cleaner.name} ({cleaner.count} items)</h3>
              <ul>
                {#each cleaner.items as item}
                  <li>{item}</li>
                {/each}
              </ul>
            </div>
          {/each}
        </div>

        <div class="warning">
          <p>WARNING: This action cannot be undone!</p>
        </div>

        <div class="actions">
          <button class="secondary-btn" onclick={handleBack}>Cancel</button>
          <button
            class="primary-btn danger"
            onclick={handleClean}
            disabled={cleaning}
          >
            {cleaning ? 'Cleaning...' : 'Confirm & Clean'}
          </button>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .cleaner {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    background: var(--bg-primary);
    font-family: var(--font-sans);
  }

  .header {
    padding: 30px 40px;
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .header h1 {
    font-size: 2rem;
    font-weight: 700;
    margin: 0;
    color: var(--text-primary);
  }

  .back-btn {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    color: var(--text-primary);
    padding: 10px 20px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 1rem;
    font-weight: 500;
    transition: all 0.3s ease;
  }

  .back-btn:hover {
    background: var(--bg-tertiary);
    border-color: var(--border-medium);
  }

  .content {
    flex: 1;
    overflow-y: auto;
    padding: 40px;
  }

  .loading, .error, .success, .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    text-align: center;
  }

  .spinner {
    width: 50px;
    height: 50px;
    border: 4px solid var(--bg-tertiary);
    border-top-color: var(--accent-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .success-icon {
    font-size: 5rem;
    margin-bottom: 20px;
    color: var(--accent-primary);
    text-shadow: 0 0 20px rgba(32, 227, 178, 0.4);
  }

  .success h2 {
    color: var(--accent-primary);
    margin-bottom: 10px;
  }

  .success p {
    color: var(--text-secondary);
  }

  .preview h2 {
    margin-bottom: 30px;
    font-size: 1.5rem;
    color: var(--text-primary);
    font-weight: 700;
  }

  .cleaners-list {
    display: grid;
    gap: 20px;
    margin-bottom: 30px;
  }

  .cleaner-card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    padding: 24px;
    transition: all 0.3s ease;
  }

  .cleaner-card:hover {
    border-color: var(--border-medium);
    box-shadow: var(--shadow-md);
  }

  .cleaner-card h3 {
    margin-bottom: 15px;
    color: var(--accent-primary);
    font-size: 1.1rem;
    font-weight: 600;
  }

  .cleaner-card ul {
    list-style: none;
    padding: 0;
  }

  .cleaner-card li {
    padding: 8px 0;
    border-bottom: 1px solid var(--border-subtle);
    color: var(--text-secondary);
    font-size: 0.9rem;
  }

  .cleaner-card li:last-child {
    border-bottom: none;
  }

  .warning {
    background: rgba(255, 107, 107, 0.08);
    border: 1px solid var(--accent-danger);
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 30px;
    text-align: center;
  }

  .warning p {
    color: var(--accent-danger);
    font-weight: 600;
    margin: 0;
  }

  .actions {
    display: flex;
    gap: 15px;
    justify-content: center;
  }

  .primary-btn, .secondary-btn {
    padding: 14px 40px;
    border-radius: 8px;
    font-size: 1rem;
    cursor: pointer;
    border: none;
    font-weight: 600;
    transition: all 0.3s ease;
    font-family: var(--font-sans);
  }

  .primary-btn {
    background: linear-gradient(135deg, var(--accent-primary), var(--accent-hover));
    color: var(--bg-primary);
    box-shadow: var(--shadow-accent);
  }

  .primary-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 32px rgba(32, 227, 178, 0.25);
  }

  .primary-btn.danger {
    background: var(--accent-danger);
    box-shadow: 0 8px 24px rgba(255, 107, 107, 0.15);
  }

  .primary-btn.danger:hover {
    background: #ff5252;
    box-shadow: 0 12px 32px rgba(255, 107, 107, 0.25);
  }

  .primary-btn:disabled {
    background: var(--bg-tertiary);
    color: var(--text-tertiary);
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }

  .secondary-btn {
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border-subtle);
  }

  .secondary-btn:hover {
    background: var(--bg-tertiary);
    border-color: var(--border-medium);
  }

  .error {
    color: var(--accent-danger);
  }

  .error button {
    margin-top: 20px;
    padding: 12px 32px;
    background: linear-gradient(135deg, var(--accent-primary), var(--accent-hover));
    color: var(--bg-primary);
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    font-family: var(--font-sans);
  }

  .empty {
    color: var(--text-secondary);
  }

  .empty button {
    margin-top: 20px;
    padding: 12px 32px;
    background: linear-gradient(135deg, var(--accent-primary), var(--accent-hover));
    color: var(--bg-primary);
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    font-family: var(--font-sans);
  }
</style>
