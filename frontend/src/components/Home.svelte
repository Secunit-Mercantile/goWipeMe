<script>
  let { onNavigate } = $props()
  let warningDismissed = $state(false)

  function dismissWarning() {
    warningDismissed = true
  }

  const features = [
    {
      id: 'backup',
      title: 'Backup',
      description: 'Save browser and shell history before cleaning',
      details: '6 browsers • 3 shells • Encrypted storage',
      icon: 'backup',
      color: 'blue',
      action: () => onNavigate('backup')
    },
    {
      id: 'restore',
      title: 'Restore',
      description: 'Recover data from previous backups',
      details: 'Version control • Selective restore',
      icon: 'restore',
      color: 'teal',
      action: () => onNavigate('restore')
    },
    {
      id: 'cleaner',
      title: 'Clear History',
      description: 'Permanently delete browsing and shell history',
      details: 'Browsers • Shells • Caches • Clipboard',
      icon: 'cleaner',
      color: 'orange',
      action: () => onNavigate('cleaner')
    },
    {
      id: 'wiper',
      title: 'Secure Wipe',
      description: 'Cryptographically wipe free disk space',
      details: 'DoD 5220.22-M • Gutmann • 1-35 passes',
      icon: 'wiper',
      color: 'red',
      action: () => onNavigate('wiper')
    }
  ]
</script>

<div class="home">
  <div class="container">
    <!-- Header -->
    <div class="header">
      <div class="brand">
        <h1 class="title gradient-text">goWipeMe</h1>
        <div class="badge">PRO</div>
      </div>
      <p class="subtitle">Professional Privacy & Security Suite</p>
    </div>

    <!-- Stats Bar -->
    <div class="stats">
      <div class="stat-item">
        <div class="stat-value">6</div>
        <div class="stat-label">Browsers</div>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <div class="stat-value">3</div>
        <div class="stat-label">Shells</div>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <div class="stat-value">35</div>
        <div class="stat-label">Max Passes</div>
      </div>
    </div>

    <!-- Warning Banner -->
    {#if !warningDismissed}
      <div class="warning">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
          <path d="M10 0L0 18h20L10 0zm0 14a1 1 0 110-2 1 1 0 010 2zm0-3a1 1 0 01-1-1V7a1 1 0 012 0v3a1 1 0 01-1 1z"/>
        </svg>
        <span>Please perform a system backup in addition to the backup features offered below.</span>
        <button class="warning-close" onclick={dismissWarning} aria-label="Dismiss warning">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M1 1l12 12M13 1L1 13"/>
          </svg>
        </button>
      </div>
    {/if}

    <!-- Feature Cards -->
    <div class="grid">
      {#each features as feature}
        <button class="card" onclick={feature.action}>
          <div class="card-glow"></div>
          <div class="card-header">
            <div class="icon {feature.color}">
              {#if feature.icon === 'backup'}
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3"/>
                </svg>
              {:else if feature.icon === 'restore'}
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 12a9 9 0 009 9 9.75 9.75 0 006.74-2.74M21 12A9 9 0 003.26 5.26M3 12V7m0 0h5m13 5v5m0 0h-5"/>
                </svg>
              {:else if feature.icon === 'cleaner'}
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2m3 0v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6h14z"/>
                </svg>
              {:else if feature.icon === 'wiper'}
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                  <path d="M9 12l2 2 4-4"/>
                </svg>
              {/if}
            </div>
            <div class="card-badge">{feature.color}</div>
          </div>
          <h3 class="card-title">{feature.title}</h3>
          <p class="card-description">{feature.description}</p>
          <div class="card-footer">
            <span class="card-details">{feature.details}</span>
            <svg class="card-arrow" width="16" height="16" viewBox="0 0 16 16" fill="none">
              <path d="M6 12l4-4-4-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </div>
        </button>
      {/each}
    </div>
  </div>
</div>

<style>
  .home {
    width: 100%;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background: var(--bg-primary);
    padding: 40px 20px;
    overflow-y: auto;
  }

  .container {
    max-width: 1100px;
    margin: 0 auto;
  }

  .header {
    text-align: center;
    margin-bottom: 32px;
  }

  .brand {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    margin-bottom: 8px;
  }

  .title {
    font-size: 56px;
    font-weight: 900;
    letter-spacing: -0.04em;
  }

  .gradient-text {
    background: linear-gradient(135deg, var(--accent-primary) 0%, var(--accent-hover) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .badge {
    background: var(--accent-primary);
    color: var(--bg-primary);
    font-size: 11px;
    font-weight: 800;
    padding: 4px 8px;
    border-radius: 6px;
    letter-spacing: 0.5px;
  }

  .subtitle {
    font-size: 16px;
    color: var(--text-secondary);
    font-weight: 500;
  }

  .stats {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 40px;
    margin-bottom: 32px;
    padding: 24px;
    background: var(--bg-elevated);
    border: 1px solid var(--border-subtle);
    border-radius: 16px;
  }

  .stat-item {
    text-align: center;
  }

  .stat-value {
    font-size: 32px;
    font-weight: 800;
    color: var(--accent-primary);
    line-height: 1;
    margin-bottom: 4px;
  }

  .stat-label {
    font-size: 12px;
    color: var(--text-tertiary);
    text-transform: uppercase;
    letter-spacing: 1px;
    font-weight: 600;
  }

  .stat-divider {
    width: 1px;
    height: 40px;
    background: var(--border-subtle);
  }

  .warning {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 18px 24px;
    padding-right: 48px;
    background: rgba(255, 107, 107, 0.05);
    border: 1px solid rgba(255, 107, 107, 0.2);
    border-radius: 12px;
    margin-bottom: 32px;
  }

  .warning svg {
    flex-shrink: 0;
    color: var(--accent-danger);
  }

  .warning span {
    font-size: 13px;
    color: var(--text-secondary);
    font-weight: 500;
    text-align: center;
  }

  .warning-close {
    position: absolute;
    top: 50%;
    right: 16px;
    transform: translateY(-50%);
    background: transparent;
    border: none;
    color: var(--text-tertiary);
    cursor: pointer;
    padding: 6px;
    border-radius: 6px;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .warning-close:hover {
    background: rgba(255, 107, 107, 0.1);
    color: var(--accent-danger);
  }

  .warning-close svg {
    color: inherit;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
    margin-bottom: 32px;
  }

  .card {
    position: relative;
    padding: 28px;
    background: var(--bg-elevated);
    border: 1px solid var(--border-subtle);
    border-radius: 20px;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    text-align: left;
    font-family: inherit;
    color: inherit;
    overflow: hidden;
  }

  .card-glow {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg,
      transparent,
      var(--accent-primary),
      transparent
    );
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .card:hover {
    transform: translateY(-4px);
    border-color: var(--border-medium);
    box-shadow: var(--shadow-accent);
  }

  .card:hover .card-glow {
    opacity: 1;
  }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
  }

  .icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg,
      rgba(32, 227, 178, 0.1),
      rgba(32, 227, 178, 0.05)
    );
    border: 1px solid rgba(32, 227, 178, 0.2);
  }

  .icon svg {
    color: var(--accent-primary);
  }

  .card-badge {
    font-size: 10px;
    font-weight: 700;
    color: var(--text-tertiary);
    text-transform: uppercase;
    letter-spacing: 1px;
  }

  .card-title {
    font-size: 22px;
    font-weight: 700;
    margin-bottom: 8px;
    color: var(--text-primary);
  }

  .card-description {
    font-size: 14px;
    color: var(--text-secondary);
    line-height: 1.6;
    margin-bottom: 20px;
  }

  .card-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .card-details {
    font-size: 12px;
    font-family: var(--font-mono);
    color: var(--text-tertiary);
  }

  .card-arrow {
    color: var(--accent-primary);
    opacity: 0;
    transform: translateX(-8px);
    transition: all 0.3s ease;
  }

  .card:hover .card-arrow {
    opacity: 1;
    transform: translateX(0);
  }

  /* Responsive */
  @media (max-width: 768px) {
    .grid {
      grid-template-columns: 1fr;
    }

    .title {
      font-size: 36px;
    }

    .stats {
      gap: 20px;
    }
  }
</style>
