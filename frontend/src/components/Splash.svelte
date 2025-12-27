<script>
  import { onMount } from 'svelte'
  import { theme } from '../stores/theme.js'
  import logo from '../assets/gui-logo.png'

  let { onComplete } = $props()
  let visible = $state(true)
  let fadeOut = $state(false)
  let dotsCount = $state(0)
  let dots = $derived('.'.repeat(dotsCount))

  onMount(() => {
    // Animated dots for loading
    const dotsInterval = setInterval(() => {
      dotsCount = (dotsCount + 1) % 4
    }, 300)

    // Show splash for 2 seconds, then fade out
    setTimeout(() => {
      clearInterval(dotsInterval)
      fadeOut = true
      // Wait for fade animation to complete
      setTimeout(() => {
        visible = false
        onComplete()
      }, 500)
    }, 2000)

    return () => clearInterval(dotsInterval)
  })
</script>

{#if visible}
  <div class="splash" class:fade-out={fadeOut}>
    <div class="splash-content">
      <div class="logo-wrapper">
        <img src={logo} alt="goWipeMe" class="logo" />
        <div class="glow-ring"></div>
      </div>

      <h1 class="app-title">
        <span class="gradient-text">goWipeMe</span>
        <span class="pro-badge">PRO</span>
      </h1>

      <div class="loading-container">
        <div class="loading-bar">
          <div class="loading-progress"></div>
        </div>
        <p class="loading-text">Initializing{dots}</p>
      </div>
    </div>
  </div>
{/if}

<style>
  .splash {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    background: var(--bg-primary);
    opacity: 1;
    transition: opacity 0.5s ease-out;
  }

  .splash.fade-out {
    opacity: 0;
  }

  .splash-content {
    text-align: center;
    animation: splash-fade-in 0.6s ease-out;
  }

  @keyframes splash-fade-in {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .logo-wrapper {
    position: relative;
    display: inline-block;
    margin-bottom: 32px;
  }

  .logo {
    max-width: 160px;
    max-height: 160px;
    width: auto;
    height: auto;
    object-fit: contain;
    position: relative;
    z-index: 2;
    filter: drop-shadow(0 8px 24px rgba(32, 227, 178, 0.2));
    animation: logo-float 3s ease-in-out infinite;
  }

  @keyframes logo-float {
    0%, 100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(-10px);
    }
  }

  .glow-ring {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 200px;
    height: 200px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(32, 227, 178, 0.1) 0%, transparent 70%);
    animation: glow-pulse 2s ease-in-out infinite;
  }

  @keyframes glow-pulse {
    0%, 100% {
      opacity: 0.5;
      transform: translate(-50%, -50%) scale(1);
    }
    50% {
      opacity: 0.8;
      transform: translate(-50%, -50%) scale(1.1);
    }
  }

  .app-title {
    font-size: 2.5rem;
    font-weight: 800;
    margin: 0 0 48px 0;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    font-family: var(--font-sans);
  }

  .gradient-text {
    background: linear-gradient(135deg, var(--accent-primary) 0%, var(--accent-hover) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    animation: gradient-shift 3s ease-in-out infinite;
  }

  @keyframes gradient-shift {
    0%, 100% {
      filter: brightness(1);
    }
    50% {
      filter: brightness(1.2);
    }
  }

  .pro-badge {
    font-size: 0.5em;
    font-weight: 900;
    color: var(--bg-primary);
    background: linear-gradient(135deg, var(--accent-primary), var(--accent-hover));
    padding: 4px 10px;
    border-radius: 6px;
    letter-spacing: 1px;
    box-shadow: var(--shadow-accent);
  }

  .loading-container {
    max-width: 300px;
    margin: 0 auto;
  }

  .loading-bar {
    width: 100%;
    height: 3px;
    background: var(--bg-tertiary);
    border-radius: 2px;
    overflow: hidden;
    margin-bottom: 16px;
  }

  .loading-progress {
    height: 100%;
    background: linear-gradient(90deg, var(--accent-primary), var(--accent-hover));
    box-shadow: 0 0 12px rgba(32, 227, 178, 0.4);
    animation: loading-progress 2s ease-out forwards;
  }

  @keyframes loading-progress {
    from {
      width: 0%;
    }
    to {
      width: 100%;
    }
  }

  .loading-text {
    font-size: 0.875rem;
    color: var(--text-secondary);
    font-weight: 500;
    letter-spacing: 0.5px;
    font-family: var(--font-sans);
    min-width: 120px;
    text-align: center;
  }
</style>
