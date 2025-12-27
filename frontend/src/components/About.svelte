<script>
  import logo from '../assets/brand/stacked.png'

  let { isOpen = $bindable(false) } = $props()

  function closeDialog() {
    isOpen = false
  }

  function handleBackdropClick(e) {
    if (e.target === e.currentTarget) {
      closeDialog()
    }
  }

  const currentYear = new Date().getFullYear()
</script>

{#if isOpen}
  <div class="backdrop" onclick={handleBackdropClick} role="presentation">
    <div class="dialog" role="dialog" aria-modal="true" aria-labelledby="about-title">
      <button class="close-btn" onclick={closeDialog} aria-label="Close">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M1 1l12 12M13 1L1 13"/>
        </svg>
      </button>

      <div class="content">
        <img src={logo} alt="goWipeMe Logo" class="logo" />

        <h1 id="about-title" class="title gradient-text">goWipeMe</h1>
        <div class="version-badge">PRO</div>

        <div class="info">
          <p class="version">Version 1.0.0</p>
          <p class="company">secunit mercantile</p>
          <p class="copyright">© {currentYear} All rights reserved</p>
        </div>

        <div class="description">
          <p>Professional Privacy & Security Suite</p>
          <p>Secure data wiping and history management for macOS</p>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10000;
    animation: fade-in 0.2s ease-out;
  }

  @keyframes fade-in {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .dialog {
    position: relative;
    background: var(--bg-elevated);
    border: 1px solid var(--border-medium);
    border-radius: 20px;
    padding: 48px 40px;
    max-width: 500px;
    width: 90%;
    box-shadow: var(--shadow-lg);
    animation: slide-up 0.3s ease-out;
  }

  @keyframes slide-up {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .close-btn {
    position: absolute;
    top: 16px;
    right: 16px;
    background: transparent;
    border: none;
    color: var(--text-tertiary);
    cursor: pointer;
    padding: 8px;
    border-radius: 6px;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .close-btn:hover {
    background: rgba(255, 255, 255, 0.05);
    color: var(--text-primary);
  }

  .content {
    text-align: center;
  }

  .logo {
    max-width: 120px;
    max-height: 120px;
    width: auto;
    height: auto;
    object-fit: contain;
    margin: 0 auto 24px;
    display: block;
    filter: drop-shadow(0 8px 24px rgba(32, 227, 178, 0.2));
  }

  .title {
    font-size: 42px;
    font-weight: 900;
    letter-spacing: -0.04em;
    margin: 0 0 12px 0;
  }

  .gradient-text {
    background: linear-gradient(135deg, var(--accent-primary) 0%, var(--accent-hover) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .version-badge {
    display: inline-block;
    background: var(--accent-primary);
    color: var(--bg-primary);
    font-size: 10px;
    font-weight: 800;
    padding: 4px 8px;
    border-radius: 6px;
    letter-spacing: 0.5px;
    margin-bottom: 24px;
  }

  .info {
    margin-bottom: 24px;
    padding-bottom: 24px;
    border-bottom: 1px solid var(--border-subtle);
  }

  .info p {
    margin: 6px 0;
    font-size: 14px;
  }

  .version {
    color: var(--text-primary);
    font-weight: 600;
  }

  .company {
    color: var(--accent-primary);
    font-weight: 600;
    font-size: 16px;
  }

  .copyright {
    color: var(--text-tertiary);
    font-size: 12px;
  }

  .description {
    color: var(--text-secondary);
  }

  .description p {
    margin: 8px 0;
    font-size: 13px;
    line-height: 1.6;
  }
</style>
