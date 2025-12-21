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
    border-top-color: #4ade80;
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

  .backup-details {
    background: #2a2a2a;
    border: 1px solid #3a3a3a;
    border-radius: 12px;
    padding: 20px;
    margin: 20px 0;
    text-align: left;
  }

  .backup-details p {
    margin: 8px 0;
    color: #aaa;
  }

  .preview h2 {
    margin-bottom: 10px;
    font-size: 1.5rem;
  }

  .subtitle {
    color: #888;
    margin-bottom: 30px;
  }

  .items-list {
    background: #2a2a2a;
    border: 1px solid #3a3a3a;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 30px;
  }

  .item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 0;
    border-bottom: 1px solid #333;
  }

  .item:last-child {
    border-bottom: none;
  }

  .item-icon {
    font-size: 1.2rem;
  }

  .item-name {
    color: #ccc;
  }

  .info-box {
    background: rgba(74, 222, 128, 0.1);
    border: 1px solid rgba(74, 222, 128, 0.3);
    border-radius: 8px;
    padding: 15px;
    margin-bottom: 30px;
    text-align: center;
  }

  .info-box p {
    color: #4ade80;
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
    background: #4ade80;
    color: #1a1a1a;
  }

  .primary-btn:hover {
    background: #22c55e;
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
    background: #4ade80;
    color: #1a1a1a;
    border: none;
    border-radius: 8px;
    cursor: pointer;
  }
</style>
