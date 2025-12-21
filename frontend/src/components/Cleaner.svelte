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
    background: #1a1a1a;
  }

  .header {
    padding: 30px 40px;
    border-bottom: 1px solid #333;
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .header h1 {
    font-size: 2rem;
    margin: 0;
  }

  .back-btn {
    background: #2a2a2a;
    border: 1px solid #3a3a3a;
    color: #fff;
    padding: 10px 20px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 1rem;
  }

  .back-btn:hover {
    background: #3a3a3a;
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
    border: 4px solid #333;
    border-top-color: #667eea;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .success-icon {
    font-size: 5rem;
    margin-bottom: 20px;
    color: #4ade80;
  }

  .success h2 {
    color: #4ade80;
    margin-bottom: 10px;
  }

  .preview h2 {
    margin-bottom: 30px;
    font-size: 1.5rem;
  }

  .cleaners-list {
    display: grid;
    gap: 20px;
    margin-bottom: 30px;
  }

  .cleaner-card {
    background: #2a2a2a;
    border: 1px solid #3a3a3a;
    border-radius: 12px;
    padding: 20px;
  }

  .cleaner-card h3 {
    margin-bottom: 15px;
    color: #667eea;
    font-size: 1.1rem;
  }

  .cleaner-card ul {
    list-style: none;
    padding: 0;
  }

  .cleaner-card li {
    padding: 8px 0;
    border-bottom: 1px solid #333;
    color: #aaa;
    font-size: 0.9rem;
  }

  .cleaner-card li:last-child {
    border-bottom: none;
  }

  .warning {
    background: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.3);
    border-radius: 8px;
    padding: 15px;
    margin-bottom: 30px;
    text-align: center;
  }

  .warning p {
    color: #ef4444;
    font-weight: 600;
  }

  .actions {
    display: flex;
    gap: 15px;
    justify-content: center;
  }

  .primary-btn, .secondary-btn {
    padding: 15px 40px;
    border-radius: 8px;
    font-size: 1rem;
    cursor: pointer;
    border: none;
    font-weight: 600;
    transition: all 0.3s ease;
  }

  .primary-btn {
    background: #667eea;
    color: white;
  }

  .primary-btn:hover {
    background: #5568d3;
  }

  .primary-btn.danger {
    background: #ef4444;
  }

  .primary-btn.danger:hover {
    background: #dc2626;
  }

  .primary-btn:disabled {
    background: #444;
    cursor: not-allowed;
  }

  .secondary-btn {
    background: #2a2a2a;
    color: white;
    border: 1px solid #3a3a3a;
  }

  .secondary-btn:hover {
    background: #3a3a3a;
  }

  .error {
    color: #ef4444;
  }

  .error button {
    margin-top: 20px;
    padding: 10px 30px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
  }
</style>
