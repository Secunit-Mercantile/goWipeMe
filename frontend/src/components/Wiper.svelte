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
    border: 4px solid #333;
    border-top-color: #667eea;
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
    background: #2a2a2a;
    border: 1px solid #3a3a3a;
    border-radius: 12px;
    padding: 30px;
    margin-bottom: 40px;
  }

  .volume-info h2 {
    margin-bottom: 20px;
  }

  .info-grid {
    display: grid;
    gap: 15px;
  }

  .info-item {
    display: flex;
    justify-content: space-between;
    padding: 10px 0;
    border-bottom: 1px solid #333;
  }

  .info-item:last-child {
    border-bottom: none;
  }

  .label {
    color: #999;
  }

  .value {
    color: #fff;
    font-weight: 600;
  }

  .method-selection h2 {
    margin-bottom: 30px;
    font-size: 1.5rem;
  }

  .methods {
    display: grid;
    gap: 20px;
  }

  .method-card {
    background: #2a2a2a;
    border: 2px solid #3a3a3a;
    border-radius: 12px;
    padding: 25px;
    cursor: pointer;
    transition: all 0.3s ease;
    text-align: left;
    position: relative;
  }

  .method-card:hover {
    border-color: #667eea;
    transform: translateX(5px);
  }

  .method-card h3 {
    margin-bottom: 10px;
    color: #667eea;
    font-size: 1.2rem;
  }

  .method-card p {
    color: #aaa;
    line-height: 1.6;
  }

  .method-card .arrow {
    position: absolute;
    right: 25px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 1.5rem;
    color: #667eea;
  }

  .confirm h2 {
    margin-bottom: 30px;
    font-size: 1.8rem;
  }

  .confirm-details {
    background: #2a2a2a;
    border: 1px solid #3a3a3a;
    border-radius: 12px;
    padding: 25px;
    margin-bottom: 30px;
  }

  .confirm-details p {
    margin: 10px 0;
    font-size: 1.1rem;
  }

  .warning-box {
    background: rgba(239, 68, 68, 0.1);
    border: 2px solid rgba(239, 68, 68, 0.3);
    border-radius: 12px;
    padding: 25px;
    margin-bottom: 40px;
  }

  .warning-box h3 {
    color: #ef4444;
    margin-bottom: 15px;
  }

  .warning-box ul {
    list-style: none;
    padding: 0;
  }

  .warning-box li {
    padding: 8px 0;
    color: #ffcccc;
  }

  .warning-box li::before {
    content: "• ";
    color: #ef4444;
    font-weight: bold;
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

  .secondary-btn {
    background: #2a2a2a;
    color: white;
    border: 1px solid #3a3a3a;
  }

  .secondary-btn:hover {
    background: #3a3a3a;
  }

  .success-icon {
    font-size: 5rem;
    margin-bottom: 20px;
    color: #4ade80;
  }

  .complete h2 {
    color: #4ade80;
    margin-bottom: 15px;
  }

  .progress-info {
    margin-top: 30px;
    color: #999;
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
