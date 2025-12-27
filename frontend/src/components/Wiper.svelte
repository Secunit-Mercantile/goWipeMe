<script>
  import { onMount } from 'svelte'
  import { GetWiperStatus, RunWiper } from '../../wailsjs/go/gui/App'

  let { onBack } = $props()

  let loading = $state(true)
  let wiperInfo = $state(null)
  let selectedMethod = $state(0)
  let wiping = $state(false)
  let complete = $state(false)
  let error = $state(null)
  let view = $state('select') // 'select', 'confirm', 'wiping', 'complete'

  onMount(async () => {
    await loadWiperInfo()
  })

  async function loadWiperInfo() {
    try {
      loading = true
      wiperInfo = await GetWiperStatus()
      loading = false
    } catch (err) {
      error = err.message
      loading = false
    }
  }

  function handleMethodSelect(methodId) {
    selectedMethod = methodId
    view = 'confirm'
  }

  async function handleStartWipe() {
    try {
      wiping = true
      view = 'wiping'
      await RunWiper(selectedMethod)
      wiping = false
      view = 'complete'
    } catch (err) {
      error = err.message
      wiping = false
      view = 'select'
    }
  }

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  function handleBack() {
    onBack()
  }

  function getSelectedMethod() {
    if (!wiperInfo) return null
    return wiperInfo.methods.find(m => m.id === selectedMethod)
  }
</script>

<div class="wiper">
  <div class="header">
    <button class="back-btn" onclick={handleBack}>← Back</button>
    <h1>Secure Wipe Free Space</h1>
  </div>

  <div class="content">
    {#if loading}
      <div class="loading">
        <div class="spinner"></div>
        <p>Loading volume information...</p>
      </div>
    {:else if error}
      <div class="error">
        <p>Error: {error}</p>
        <button onclick={loadWiperInfo}>Retry</button>
      </div>
    {:else if view === 'select'}
      <div class="method-selection">
        <div class="volume-info">
          <h2>Volume Information</h2>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">Volume:</span>
              <span class="value">{wiperInfo.volume}</span>
            </div>
            <div class="info-item">
              <span class="label">Total Space:</span>
              <span class="value">{formatBytes(wiperInfo.totalSpace)}</span>
            </div>
            <div class="info-item">
              <span class="label">Free Space:</span>
              <span class="value">{formatBytes(wiperInfo.freeSpace)}</span>
            </div>
          </div>
        </div>

        <h2>Select Wipe Method</h2>

        <div class="methods">
          {#each wiperInfo.methods as method}
            <button
              class="method-card"
              onclick={() => handleMethodSelect(method.id)}
            >
              <h3>{method.name}</h3>
              <p>{method.description}</p>
              <span class="arrow">→</span>
            </button>
          {/each}
        </div>
      </div>
    {:else if view === 'confirm'}
      <div class="confirm">
        <h2>Confirm Wipe Operation</h2>

        {#if wiperInfo}
          {@const method = getSelectedMethod()}
          <div class="confirm-details">
            <p><strong>Method:</strong> {method?.name}</p>
            <p><strong>Volume:</strong> {wiperInfo.volume}</p>
            <p><strong>Free Space to Wipe:</strong> {formatBytes(wiperInfo.freeSpace)}</p>
          </div>
        {/if}

        <div class="warning-box">
          <h3>WARNING</h3>
          <ul>
            <li>This operation will fill all free space on the volume</li>
            <li>The process will take a significant amount of time</li>
            <li>Cannot be interrupted once started</li>
            <li>All operations are irreversible</li>
          </ul>
        </div>

        <div class="actions">
          <button class="secondary-btn" onclick={() => view = 'select'}>
            Cancel
          </button>
          <button class="primary-btn danger" onclick={handleStartWipe}>
            Start Wiping
          </button>
        </div>
      </div>
    {:else if view === 'wiping'}
      <div class="wiping-view">
        <div class="spinner-large"></div>
        <h2>Wiping Free Space...</h2>
        <p>This may take a while. Please do not close the application.</p>
        <div class="progress-info">
          <p>Processing...</p>
        </div>
      </div>
    {:else if view === 'complete'}
      <div class="complete">
        <div class="success-icon">✓</div>
        <h2>Wiping Complete!</h2>
        <p>Free space has been successfully wiped.</p>
        <button class="primary-btn" onclick={handleBack}>Back to Home</button>
      </div>
    {/if}
  </div>
</div>

<style>
  .wiper {
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

  .loading, .error, .wiping-view, .complete {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    text-align: center;
  }

  .spinner, .spinner-large {
    width: 50px;
    height: 50px;
    border: 4px solid var(--bg-tertiary);
    border-top-color: var(--accent-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 20px;
  }

  .spinner-large {
    width: 80px;
    height: 80px;
    border-width: 6px;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .volume-info {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    padding: 30px;
    margin-bottom: 40px;
  }

  .volume-info h2 {
    margin-bottom: 24px;
    color: var(--text-primary);
    font-weight: 700;
  }

  .info-grid {
    display: grid;
    gap: 15px;
  }

  .info-item {
    display: flex;
    justify-content: space-between;
    padding: 12px 0;
    border-bottom: 1px solid var(--border-subtle);
  }

  .info-item:last-child {
    border-bottom: none;
  }

  .label {
    color: var(--text-secondary);
  }

  .value {
    color: var(--accent-primary);
    font-weight: 600;
  }

  .method-selection h2 {
    margin-bottom: 30px;
    font-size: 1.5rem;
    color: var(--text-primary);
    font-weight: 700;
  }

  .methods {
    display: grid;
    gap: 20px;
  }

  .method-card {
    background: var(--bg-secondary);
    border: 2px solid var(--border-subtle);
    border-radius: 12px;
    padding: 25px;
    cursor: pointer;
    transition: all 0.3s ease;
    text-align: left;
    position: relative;
  }

  .method-card:hover {
    border-color: var(--accent-primary);
    transform: translateX(5px);
    box-shadow: var(--shadow-accent);
  }

  .method-card h3 {
    margin-bottom: 10px;
    color: var(--accent-primary);
    font-size: 1.2rem;
    font-weight: 600;
  }

  .method-card p {
    color: var(--text-secondary);
    line-height: 1.6;
  }

  .method-card .arrow {
    position: absolute;
    right: 25px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 1.5rem;
    color: var(--accent-primary);
  }

  .confirm h2 {
    margin-bottom: 30px;
    font-size: 1.8rem;
    color: var(--text-primary);
    font-weight: 700;
  }

  .confirm-details {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    padding: 25px;
    margin-bottom: 30px;
  }

  .confirm-details p {
    margin: 10px 0;
    font-size: 1.1rem;
    color: var(--text-secondary);
  }

  .confirm-details strong {
    color: var(--text-primary);
  }

  .warning-box {
    background: rgba(255, 107, 107, 0.08);
    border: 2px solid var(--accent-danger);
    border-radius: 12px;
    padding: 25px;
    margin-bottom: 40px;
  }

  .warning-box h3 {
    color: var(--accent-danger);
    margin-bottom: 15px;
    font-weight: 700;
  }

  .warning-box ul {
    list-style: none;
    padding: 0;
  }

  .warning-box li {
    padding: 8px 0;
    color: var(--text-secondary);
  }

  .warning-box li::before {
    content: "• ";
    color: var(--accent-danger);
    font-weight: bold;
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

  .secondary-btn {
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border-subtle);
  }

  .secondary-btn:hover {
    background: var(--bg-tertiary);
    border-color: var(--border-medium);
  }

  .success-icon {
    font-size: 5rem;
    margin-bottom: 20px;
    color: var(--accent-primary);
    text-shadow: 0 0 20px rgba(32, 227, 178, 0.4);
  }

  .complete h2 {
    color: var(--accent-primary);
    margin-bottom: 15px;
    font-weight: 700;
  }

  .complete p {
    color: var(--text-secondary);
  }

  .progress-info {
    margin-top: 30px;
    color: var(--text-secondary);
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
