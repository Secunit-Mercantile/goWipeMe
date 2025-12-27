<script>
  import { onMount } from 'svelte'
  import { ListBackups, RestoreBackup, DeleteBackup } from '../../wailsjs/go/gui/App'

  let { onBack } = $props()

  let loading = $state(true)
  let backups = $state([])
  let selectedBackup = $state(null)
  let restoring = $state(false)
  let complete = $state(false)
  let error = $state(null)
  let view = $state('list') // 'list', 'confirm', 'complete'

  onMount(async () => {
    await loadBackups()
  })

  async function loadBackups() {
    try {
      loading = true
      backups = await ListBackups() || []
      loading = false
    } catch (err) {
      error = err.message
      loading = false
    }
  }

  function handleSelectBackup(backup) {
    selectedBackup = backup
    view = 'confirm'
  }

  async function handleRestore() {
    try {
      restoring = true
      error = null
      await RestoreBackup(selectedBackup.id)
      restoring = false
      view = 'complete'
    } catch (err) {
      error = err.message
      restoring = false
    }
  }

  async function handleDelete(backup) {
    if (!confirm(`Delete backup from ${formatDate(backup.timestamp)}?`)) {
      return
    }

    try {
      await DeleteBackup(backup.id)
      await loadBackups()
    } catch (err) {
      error = err.message
    }
  }

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  function formatDate(timestamp) {
    return new Date(timestamp).toLocaleString()
  }

  function handleBack() {
    if (view === 'confirm') {
      view = 'list'
      selectedBackup = null
    } else {
      onBack()
    }
  }
</script>

<div class="restore">
  <div class="header">
    <button class="back-btn" onclick={handleBack}>← Back</button>
    <h1>Restore Backup</h1>
  </div>

  <div class="content">
    {#if loading}
      <div class="loading">
        <div class="spinner"></div>
        <p>Loading backups...</p>
      </div>
    {:else if error}
      <div class="error">
        <p>Error: {error}</p>
        <button onclick={loadBackups}>Retry</button>
      </div>
    {:else if view === 'complete'}
      <div class="success">
        <div class="success-icon">✓</div>
        <h2>Restore Complete!</h2>
        <p>Your data has been restored successfully.</p>
        <button class="primary-btn" onclick={() => onBack()}>Back to Home</button>
      </div>
    {:else if view === 'confirm'}
      <div class="confirm">
        <h2>Confirm Restore</h2>
        <p class="subtitle">The following items will be restored:</p>

        <div class="backup-details">
          <p><strong>Backup Date:</strong> {formatDate(selectedBackup.timestamp)}</p>
          <p><strong>Size:</strong> {formatBytes(selectedBackup.size)}</p>
        </div>

        <div class="items-list">
          {#each selectedBackup.items as item}
            <div class="item">
              <span class="item-icon">📄</span>
              <span class="item-name">{item}</span>
            </div>
          {/each}
        </div>

        <div class="warning-box">
          <p>This will overwrite any existing files at the original locations.</p>
        </div>

        <div class="actions">
          <button class="secondary-btn" onclick={() => view = 'list'}>Cancel</button>
          <button
            class="primary-btn"
            onclick={handleRestore}
            disabled={restoring}
          >
            {restoring ? 'Restoring...' : 'Restore Backup'}
          </button>
        </div>
      </div>
    {:else if backups.length === 0}
      <div class="empty">
        <div class="empty-icon">📁</div>
        <h2>No Backups Found</h2>
        <p>Create a backup first before you can restore.</p>
        <button class="primary-btn" onclick={() => onBack()}>Back to Home</button>
      </div>
    {:else}
      <div class="backup-list">
        <h2>Available Backups</h2>
        <p class="subtitle">Select a backup to restore</p>

        <div class="backups">
          {#each backups as backup}
            <div class="backup-card">
              <div class="backup-info" onclick={() => handleSelectBackup(backup)}>
                <div class="backup-date">{formatDate(backup.timestamp)}</div>
                <div class="backup-meta">
                  <span>{backup.items.length} items</span>
                  <span>{formatBytes(backup.size)}</span>
                </div>
              </div>
              <button class="delete-btn" onclick={() => handleDelete(backup)} title="Delete backup">
                ×
              </button>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .restore {
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

  .success-icon, .empty-icon {
    font-size: 5rem;
    margin-bottom: 20px;
  }

  .success-icon {
    color: var(--accent-primary);
    text-shadow: 0 0 20px rgba(32, 227, 178, 0.4);
  }

  .empty-icon {
    opacity: 0.5;
  }

  .success h2 {
    color: var(--accent-primary);
    margin-bottom: 10px;
    font-weight: 700;
  }

  .success p {
    color: var(--text-secondary);
  }

  .empty h2 {
    margin-bottom: 10px;
    color: var(--text-primary);
    font-weight: 700;
  }

  .empty p {
    color: var(--text-secondary);
    margin-bottom: 20px;
  }

  .backup-list h2 {
    margin-bottom: 10px;
    font-size: 1.5rem;
    color: var(--text-primary);
    font-weight: 700;
  }

  .subtitle {
    color: var(--text-secondary);
    margin-bottom: 30px;
  }

  .backups {
    display: grid;
    gap: 15px;
  }

  .backup-card {
    background: var(--bg-secondary);
    border: 2px solid var(--border-subtle);
    border-radius: 12px;
    display: flex;
    align-items: center;
    transition: all 0.3s ease;
  }

  .backup-card:hover {
    border-color: var(--accent-primary);
    box-shadow: var(--shadow-accent);
  }

  .backup-info {
    flex: 1;
    padding: 20px;
    cursor: pointer;
  }

  .backup-date {
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 8px;
    color: var(--text-primary);
  }

  .backup-meta {
    display: flex;
    gap: 20px;
    color: var(--text-secondary);
    font-size: 0.9rem;
  }

  .delete-btn {
    background: transparent;
    border: none;
    color: var(--text-tertiary);
    font-size: 1.5rem;
    padding: 20px;
    cursor: pointer;
    transition: color 0.3s ease;
  }

  .delete-btn:hover {
    color: var(--accent-danger);
  }

  .confirm h2 {
    margin-bottom: 10px;
    font-size: 1.5rem;
    color: var(--text-primary);
    font-weight: 700;
  }

  .backup-details {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 20px;
  }

  .backup-details p {
    margin: 8px 0;
    color: var(--text-secondary);
  }

  .backup-details strong {
    color: var(--text-primary);
  }

  .items-list {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 20px;
    max-height: 200px;
    overflow-y: auto;
  }

  .item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 0;
    border-bottom: 1px solid var(--border-subtle);
  }

  .item:last-child {
    border-bottom: none;
  }

  .item-icon {
    font-size: 1.1rem;
  }

  .item-name {
    color: var(--text-secondary);
    font-size: 0.9rem;
  }

  .warning-box {
    background: rgba(255, 178, 32, 0.08);
    border: 1px solid var(--accent-warning);
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 30px;
    text-align: center;
  }

  .warning-box p {
    color: var(--accent-warning);
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
</style>
