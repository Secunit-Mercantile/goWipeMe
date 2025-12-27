<script>
  import { onMount } from 'svelte'
  import { GetBackupPreview, CreateBackup } from '../../wailsjs/go/gui/App'

  let { onBack } = $props()

  let loading = $state(true)
  let preview = $state(null)
  let creating = $state(false)
  let complete = $state(false)
  let backupInfo = $state(null)
  let error = $state(null)

  onMount(async () => {
    await loadPreview()
  })

  async function loadPreview() {
    try {
      loading = true
      preview = await GetBackupPreview()
      loading = false
    } catch (err) {
      error = err.message
      loading = false
    }
  }

  async function handleCreateBackup() {
    try {
      creating = true
      error = null
      backupInfo = await CreateBackup()
      creating = false
      complete = true
    } catch (err) {
      error = err.message
      creating = false
    }
  }

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  function handleBack() {
    onBack()
  }
</script>

<div class="backup">
  <div class="header">
    <button class="back-btn" onclick={handleBack}>← Back</button>
    <h1>Backup Data</h1>
  </div>

  <div class="content">
    {#if loading}
      <div class="loading">
        <div class="spinner"></div>
        <p>Scanning for items to backup...</p>
      </div>
    {:else if error}
      <div class="error">
        <p>Error: {error}</p>
        <button onclick={loadPreview}>Retry</button>
      </div>
    {:else if complete}
      <div class="success">
        <div class="success-icon">✓</div>
        <h2>Backup Created!</h2>
        <p>Your data has been saved successfully.</p>
        <div class="backup-details">
          <p><strong>Backup ID:</strong> {backupInfo.id}</p>
          <p><strong>Size:</strong> {formatBytes(backupInfo.size)}</p>
          <p><strong>Items:</strong> {backupInfo.items.length} files</p>
        </div>
        <button class="primary-btn" onclick={handleBack}>Back to Home</button>
      </div>
    {:else if !preview || preview.items.length === 0}
      <div class="empty">
        <p>No items found to backup.</p>
        <button onclick={handleBack}>Back to Home</button>
      </div>
    {:else}
      <div class="preview">
        <h2>Items to be backed up:</h2>
        <p class="subtitle">The following files will be saved to ~/.gowipeme/backups/</p>

        <div class="items-list">
          {#each preview.items as item}
            <div class="item">
              <span class="item-icon">📄</span>
              <span class="item-name">{item}</span>
            </div>
          {/each}
        </div>

        <div class="info-box">
          <p>Backups are stored locally and can be restored at any time.</p>
        </div>

        <div class="actions">
          <button class="secondary-btn" onclick={handleBack}>Cancel</button>
          <button
            class="primary-btn"
            onclick={handleCreateBackup}
            disabled={creating}
          >
            {creating ? 'Creating Backup...' : 'Create Backup'}
          </button>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .backup {
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
    font-weight: 700;
  }

  .success p {
    color: var(--text-secondary);
  }

  .backup-details {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    padding: 20px;
    margin: 20px 0;
    text-align: left;
  }

  .backup-details p {
    margin: 8px 0;
    color: var(--text-secondary);
  }

  .backup-details strong {
    color: var(--text-primary);
  }

  .preview h2 {
    margin-bottom: 10px;
    font-size: 1.5rem;
    color: var(--text-primary);
    font-weight: 700;
  }

  .subtitle {
    color: var(--text-secondary);
    margin-bottom: 30px;
  }

  .items-list {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 30px;
    max-height: 300px;
    overflow-y: auto;
  }

  .item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 0;
    border-bottom: 1px solid var(--border-subtle);
  }

  .item:last-child {
    border-bottom: none;
  }

  .item-icon {
    font-size: 1.2rem;
  }

  .item-name {
    color: var(--text-secondary);
  }

  .info-box {
    background: rgba(32, 227, 178, 0.08);
    border: 1px solid var(--accent-secondary);
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 30px;
    text-align: center;
  }

  .info-box p {
    color: var(--accent-primary);
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
