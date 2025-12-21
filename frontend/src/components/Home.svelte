<script>
  let { onNavigate } = $props()

  function handleBackup() {
    onNavigate('backup')
  }

  function handleRestore() {
    onNavigate('restore')
  }

  function handleCleaner() {
    onNavigate('cleaner')
  }

  function handleWiper() {
    onNavigate('wiper')
  }
</script>

<div class="home">
  <div class="header">
    <h1>goWipeMe</h1>
    <p class="tagline">Privacy Tool for macOS</p>
  </div>

  <div class="cards">
    <button class="card backup" onclick={handleBackup}>
      <div class="card-icon">💾</div>
      <h2>Backup</h2>
      <p>Save your browser history, shell history, and other data before cleaning</p>
      <div class="features">
        <span class="badge">Browser History</span>
        <span class="badge">Shell History</span>
      </div>
    </button>

    <button class="card restore" onclick={handleRestore}>
      <div class="card-icon">🔄</div>
      <h2>Restore</h2>
      <p>Restore previously backed up data from a saved backup</p>
      <div class="features">
        <span class="badge">View Backups</span>
        <span class="badge">Restore Data</span>
      </div>
    </button>

    <button class="card cleaner" onclick={handleCleaner}>
      <div class="card-icon">🧹</div>
      <h2>Clear All History</h2>
      <p>Remove browser history, shell history, caches, recent files, and clipboard</p>
      <div class="features">
        <span class="badge">6 Browsers</span>
        <span class="badge">Shell History</span>
        <span class="badge">App Caches</span>
      </div>
    </button>

    <button class="card wiper" onclick={handleWiper}>
      <div class="card-icon">🔒</div>
      <h2>Secure Wipe Free Space</h2>
      <p>Securely wipe free disk space with multiple security levels</p>
      <div class="features">
        <span class="badge">Single Pass</span>
        <span class="badge">DoD 5220.22-M</span>
        <span class="badge">Gutmann</span>
      </div>
    </button>
  </div>

  <div class="footer">
    <p>All cleaning operations are irreversible. Create a backup first!</p>
  </div>
</div>

<style>
  .home {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px;
    background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
    overflow-y: auto;
  }

  .header {
    text-align: center;
    margin-bottom: 40px;
  }

  h1 {
    font-size: 3rem;
    font-weight: 700;
    margin-bottom: 10px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .tagline {
    font-size: 1.1rem;
    color: #999;
  }

  .cards {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
    max-width: 900px;
    width: 100%;
    margin-bottom: 30px;
  }

  .card {
    background: #2a2a2a;
    border: 2px solid #3a3a3a;
    border-radius: 16px;
    padding: 30px 25px;
    cursor: pointer;
    transition: all 0.3s ease;
    text-align: center;
  }

  .card:hover {
    transform: translateY(-3px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
  }

  .card.backup:hover {
    border-color: #4ade80;
    box-shadow: 0 8px 30px rgba(74, 222, 128, 0.2);
  }

  .card.restore:hover {
    border-color: #60a5fa;
    box-shadow: 0 8px 30px rgba(96, 165, 250, 0.2);
  }

  .card.cleaner:hover {
    border-color: #f59e0b;
    box-shadow: 0 8px 30px rgba(245, 158, 11, 0.2);
  }

  .card.wiper:hover {
    border-color: #ef4444;
    box-shadow: 0 8px 30px rgba(239, 68, 68, 0.2);
  }

  .card-icon {
    font-size: 3rem;
    margin-bottom: 15px;
  }

  .card h2 {
    font-size: 1.3rem;
    margin-bottom: 10px;
    color: #fff;
  }

  .card p {
    color: #aaa;
    margin-bottom: 15px;
    line-height: 1.5;
    font-size: 0.9rem;
  }

  .features {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    justify-content: center;
  }

  .badge {
    background: #3a3a3a;
    padding: 4px 10px;
    border-radius: 15px;
    font-size: 0.75rem;
    color: #999;
  }

  .footer {
    text-align: center;
    color: #666;
    font-size: 0.85rem;
    max-width: 600px;
  }

  @media (max-width: 700px) {
    .cards {
      grid-template-columns: 1fr;
    }
  }
</style>
