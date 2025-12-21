<script>
  import { onMount } from 'svelte'
  import logo from '../assets/gui-logo.png'

  let { onComplete } = $props()

  let visible = $state(true)
  let fadeOut = $state(false)

  onMount(() => {
    // Show splash for 2 seconds, then fade out
    setTimeout(() => {
      fadeOut = true
      // Wait for fade animation to complete
      setTimeout(() => {
        visible = false
        onComplete()
      }, 500)
    }, 2000)
  })
</script>

{#if visible}
  <div class="splash" class:fade-out={fadeOut}>
    <div class="logo-container">
      <img src={logo} alt="goWipeMe" class="logo" />
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
    background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    opacity: 1;
    transition: opacity 0.5s ease-out;
  }

  .splash.fade-out {
    opacity: 0;
  }

  .logo-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    animation: pulse 2s ease-in-out infinite;
  }

  .logo {
    max-width: 400px;
    max-height: 400px;
    width: auto;
    height: auto;
    object-fit: contain;
    filter: drop-shadow(0 0 30px rgba(102, 126, 234, 0.3));
  }

  @keyframes pulse {
    0%, 100% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.02);
    }
  }
</style>
