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
    border-top-color: #60a5fa;
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
    color: #4ade80;
  }

  .success h2 {
    color: #4ade80;
    margin-bottom: 10px;
  }

  .empty h2 {
    margin-bottom: 10px;
  }

  .empty p {
    color: #888;
    margin-bottom: 20px;
  }

  .backup-list h2 {
    margin-bottom: 10px;
    font-size: 1.5rem;
  }

  .subtitle {
    color: #888;
    margin-bottom: 30px;
  }

  .backups {
    display: grid;
    gap: 15px;
  }

  .backup-card {
    background: #2a2a2a;
    border: 2px solid #3a3a3a;
    border-radius: 12px;
    display: flex;
    align-items: center;
    transition: all 0.3s ease;
  }

  .backup-card:hover {
    border-color: #60a5fa;
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
  }

  .backup-meta {
    display: flex;
    gap: 20px;
    color: #888;
    font-size: 0.9rem;
  }

  .delete-btn {
    background: transparent;
    border: none;
    color: #666;
    font-size: 1.5rem;
    padding: 20px;
    cursor: pointer;
    transition: color 0.3s ease;
  }

  .delete-btn:hover {
    color: #ef4444;
  }

  .confirm h2 {
    margin-bottom: 10px;
    font-size: 1.5rem;
  }

  .backup-details {
    background: #2a2a2a;
    border: 1px solid #3a3a3a;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 20px;
  }

  .backup-details p {
    margin: 8px 0;
    color: #aaa;
  }

  .items-list {
    background: #2a2a2a;
    border: 1px solid #3a3a3a;
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
    border-bottom: 1px solid #333;
  }

  .item:last-child {
    border-bottom: none;
  }

  .item-icon {
    font-size: 1.1rem;
  }

  .item-name {
    color: #ccc;
    font-size: 0.9rem;
  }

  .warning-box {
    background: rgba(245, 158, 11, 0.1);
    border: 1px solid rgba(245, 158, 11, 0.3);
    border-radius: 8px;
    padding: 15px;
    margin-bottom: 30px;
    text-align: center;
  }

  .warning-box p {
    color: #f59e0b;
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
    background: #60a5fa;
    color: #1a1a1a;
  }

  .primary-btn:hover {
    background: #3b82f6;
  }

  .primary-btn:disabled {
    background: #444;
    color: #888;
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
    background: #60a5fa;
    color: #1a1a1a;
    border: none;
    border-radius: 8px;
    cursor: pointer;
  }
</style>
